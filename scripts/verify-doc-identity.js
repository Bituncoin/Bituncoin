#!/usr/bin/env node

const { createVerify, createHash } = require('crypto')
const { readFileSync } = require('fs')
const { resolve } = require('path')

const ROOT_DIR = process.cwd()
const IDENTITY_PATH = resolve(ROOT_DIR, 'documentation.identity.json')

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

function verifySignature(versionHash, signature, publicKey) {
    const verifier = createVerify('SHA256')
    verifier.update(versionHash)
    verifier.end()
    return verifier.verify(publicKey, signature, 'base64')
}

function verifyDocumentationIdentity() {
    const manifest = JSON.parse(readFileSync(IDENTITY_PATH, 'utf8'))
    const documentation = manifest.documentation || {}

    if (!documentation.versionHash || !documentation.signature || !documentation.publicKey) {
        throw new Error('documentation.identity.json missing required fields: versionHash, signature, publicKey')
    }

    const files = Array.isArray(documentation.files) ? documentation.files : []
    const core = {
        source: documentation.source || 'BTNG Sovereign Documentation Suite',
        hashAlgorithm: documentation.hashAlgorithm || 'sha256',
        files
    }

    const computedHash = `sha256-${sha256Hex(toCanonicalJson(core))}`
    const hashValid = computedHash === documentation.versionHash
    const signatureValid = verifySignature(documentation.versionHash, documentation.signature, documentation.publicKey)

    console.log(`DOC_IDENTITY_HASH=${documentation.versionHash}`)
    console.log(`DOC_IDENTITY_HASH_VALID=${hashValid}`)
    console.log(`DOC_IDENTITY_SIGNATURE_VALID=${signatureValid}`)

    if (!hashValid || !signatureValid) {
        throw new Error('Documentation identity verification failed.')
    }

    return {
        ok: true,
        versionHash: documentation.versionHash,
        hashValid,
        signatureValid,
        fileCount: files.length
    }
}

if (require.main === module) {
    try {
        verifyDocumentationIdentity()
        process.exit(0)
    } catch (error) {
        console.error(`DOC_IDENTITY_VERIFY=FAIL\n${error.message}`)
        process.exit(1)
    }
}

module.exports = {
    verifyDocumentationIdentity
}
