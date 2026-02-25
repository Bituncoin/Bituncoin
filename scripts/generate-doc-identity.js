#!/usr/bin/env node

const { createHash, createSign } = require('crypto')
const { readFileSync, writeFileSync } = require('fs')
const { isAbsolute, resolve } = require('path')

const ROOT_DIR = process.cwd()
const CANONICAL_CONFIG_PATH = resolve(ROOT_DIR, 'docs/canonical-docs.json')
const OUTPUT_PATH = resolve(ROOT_DIR, 'documentation.identity.json')

function loadKeyFromEnvOrFile(inlineVarName, pathVarName) {
    const inlineValue = process.env[inlineVarName]?.trim()
    if (inlineValue) {
        return inlineValue.replace(/\\n/g, '\n')
    }

    const keyPathValue = process.env[pathVarName]?.trim()
    if (!keyPathValue) {
        return null
    }

    const absolutePath = isAbsolute(keyPathValue) ? keyPathValue : resolve(ROOT_DIR, keyPathValue)
    try {
        return readFileSync(absolutePath, 'utf8')
    } catch {
        return null
    }
}

function loadDefaultKeyPair() {
    const privateKeyCandidates = [
        resolve(ROOT_DIR, 'btng-es256-private-key.pem'),
        resolve(ROOT_DIR, 'keys/btng-es256-private-key.pem')
    ]
    const publicKeyCandidates = [
        resolve(ROOT_DIR, 'btng-es256-public-key.pem'),
        resolve(ROOT_DIR, 'keys/btng-es256-public-key.pem')
    ]

    const loadFirstExisting = (candidates) => {
        for (const candidate of candidates) {
            try {
                return readFileSync(candidate, 'utf8')
            } catch {
                // keep searching
            }
        }
        return null
    }

    return {
        privateKey: loadFirstExisting(privateKeyCandidates),
        publicKey: loadFirstExisting(publicKeyCandidates)
    }
}

function sortKeysRecursively(value) {
    if (Array.isArray(value)) {
        return value.map((entry) => sortKeysRecursively(entry))
    }

    if (value && typeof value === 'object') {
        const sorted = {}
        Object.keys(value)
            .sort()
            .forEach((key) => {
                sorted[key] = sortKeysRecursively(value[key])
            })
        return sorted
    }

    return value
}

function toCanonicalJson(value) {
    return JSON.stringify(sortKeysRecursively(value))
}

function sha256Hex(value) {
    return createHash('sha256').update(value).digest('hex')
}

function signVersionHash(versionHash, privateKey) {
    const signer = createSign('SHA256')
    signer.update(versionHash)
    signer.end()
    return signer.sign(privateKey, 'base64')
}

function generateDocumentationIdentity() {
    const configRaw = readFileSync(CANONICAL_CONFIG_PATH, 'utf8')
    const config = JSON.parse(configRaw)
    const files = Array.isArray(config.files) ? config.files.slice() : []

    if (files.length === 0) {
        throw new Error('docs/canonical-docs.json contains no files to hash.')
    }

    const fileEntries = files
        .map((relativePath) => {
            const absolutePath = resolve(ROOT_DIR, relativePath)
            const content = readFileSync(absolutePath)
            return {
                path: relativePath.replace(/\\/g, '/'),
                sha256: sha256Hex(content),
                size: content.length
            }
        })
        .sort((a, b) => a.path.localeCompare(b.path))

    const identityCore = {
        source: config.source || 'BTNG Sovereign Documentation Suite',
        hashAlgorithm: 'sha256',
        files: fileEntries
    }

    const versionHash = `sha256-${sha256Hex(toCanonicalJson(identityCore))}`

    const defaultKeys = loadDefaultKeyPair()
    const privateKey = loadKeyFromEnvOrFile('BTNG_ES256_PRIVATE_KEY', 'BTNG_ES256_PRIVATE_KEY_PATH') || defaultKeys.privateKey
    const publicKey = loadKeyFromEnvOrFile('BTNG_ES256_PUBLIC_KEY', 'BTNG_ES256_PUBLIC_KEY_PATH') || defaultKeys.publicKey

    if (!privateKey) {
        throw new Error('BTNG ES256 private key not configured. Set BTNG_ES256_PRIVATE_KEY or BTNG_ES256_PRIVATE_KEY_PATH.')
    }

    if (!publicKey) {
        throw new Error('BTNG ES256 public key not configured. Set BTNG_ES256_PUBLIC_KEY or BTNG_ES256_PUBLIC_KEY_PATH.')
    }

    const signature = signVersionHash(versionHash, privateKey)

    const identityManifest = {
        documentation: {
            versionHash,
            signature,
            publicKey,
            timestamp: Math.floor(Date.now() / 1000),
            source: identityCore.source,
            hashAlgorithm: identityCore.hashAlgorithm,
            signatureAlgorithm: 'ES256/SHA256',
            files: identityCore.files
        }
    }

    writeFileSync(OUTPUT_PATH, `${JSON.stringify(identityManifest, null, 2)}\n`, 'utf8')

    console.log('DOC_IDENTITY_GENERATE=OK')
    console.log(`DOC_IDENTITY_HASH=${versionHash}`)
    console.log(`DOC_IDENTITY_FILES=${identityCore.files.length}`)

    return identityManifest
}

if (require.main === module) {
    try {
        generateDocumentationIdentity()
        process.exit(0)
    } catch (error) {
        console.error(`DOC_IDENTITY_GENERATE=FAIL\n${error.message}`)
        process.exit(1)
    }
}

module.exports = {
    generateDocumentationIdentity,
    toCanonicalJson,
    sha256Hex
}
