#!/usr/bin/env node

const { readFileSync, writeFileSync } = require('fs')
const { resolve } = require('path')

const ROOT_DIR = process.cwd()
const IDENTITY_PATH = resolve(ROOT_DIR, 'documentation.identity.json')
const FABRIC_MANIFEST_PATH = resolve(ROOT_DIR, 'btng-fabric-network.json')
const SERVICE_MANIFEST_PATH = resolve(ROOT_DIR, 'manifests/service-doc-identity.json')

function readJson(path) {
    return JSON.parse(readFileSync(path, 'utf8'))
}

function writeJson(path, payload) {
    writeFileSync(path, `${JSON.stringify(payload, null, 2)}\n`, 'utf8')
}

function propagateDocumentationIdentity() {
    const identityManifest = readJson(IDENTITY_PATH)
    const documentation = identityManifest.documentation

    if (!documentation?.versionHash) {
        throw new Error('documentation.identity.json is missing documentation.versionHash')
    }

    const fabric = readJson(FABRIC_MANIFEST_PATH)
    fabric.documentation = {
        versionHash: documentation.versionHash,
        signature: documentation.signature,
        publicKey: documentation.publicKey,
        timestamp: documentation.timestamp,
        source: documentation.source,
        signatureAlgorithm: documentation.signatureAlgorithm
    }
    writeJson(FABRIC_MANIFEST_PATH, fabric)

    const serviceManifest = {
        documentation: {
            versionHash: documentation.versionHash,
            signature: documentation.signature,
            publicKey: documentation.publicKey,
            timestamp: documentation.timestamp,
            source: documentation.source,
            signatureAlgorithm: documentation.signatureAlgorithm
        },
        services: {
            watchtower: {
                manifestPath: 'btng-api/app/api/watchtower/meta/route.ts',
                enforceOnStartup: true
            },
            goldOracle: {
                manifestPath: 'gold-oracle.config.js',
                enforceOnStartup: true
            }
        }
    }
    writeJson(SERVICE_MANIFEST_PATH, serviceManifest)

    console.log('DOC_IDENTITY_PROPAGATE=OK')
    console.log(`DOC_IDENTITY_HASH=${documentation.versionHash}`)
}

if (require.main === module) {
    try {
        propagateDocumentationIdentity()
        process.exit(0)
    } catch (error) {
        console.error(`DOC_IDENTITY_PROPAGATE=FAIL\n${error.message}`)
        process.exit(1)
    }
}

module.exports = {
    propagateDocumentationIdentity
}
