#!/usr/bin/env node

const fs = require('fs')
const crypto = require('crypto')

const CANDIDATE_PORTS = Array.from({ length: 10 }, (_, index) => 3001 + index)
const privateKeyPath = process.env.BTNG_ES256_PRIVATE_KEY_PATH || './btng-es256-private-key.pem'

function serializePayload(payload) {
    return JSON.stringify({
        nodeId: payload.nodeId,
        ip: payload.ip,
        port: payload.port,
        role: payload.role,
        timestamp: payload.timestamp
    })
}

function signPayload(payload, privateKeyPem) {
    const signer = crypto.createSign('SHA256')
    signer.update(serializePayload(payload))
    signer.end()
    return signer.sign(privateKeyPem, 'base64')
}

async function findWatchtowerBaseUrl() {
    for (const port of CANDIDATE_PORTS) {
        const baseUrl = `http://localhost:${port}`
        try {
            const response = await fetch(`${baseUrl}/api/watchtower/nodes`, { cache: 'no-store' })
            if (!response.ok) {
                continue
            }

            const payload = await response.json()
            if (payload && payload.ok === true) {
                return baseUrl
            }
        } catch {
            // continue scan
        }
    }

    return null
}

async function main() {
    if (!fs.existsSync(privateKeyPath)) {
        throw new Error(`Private key not found at ${privateKeyPath}`)
    }

    const privateKey = fs.readFileSync(privateKeyPath, 'utf8')
    const baseUrl = await findWatchtowerBaseUrl()

    if (!baseUrl) {
        throw new Error('No active watchtower endpoint found on ports 3001-3010')
    }

    const registrationPayload = {
        nodeId: 'vps-1282934',
        ip: '72.62.117.125',
        port: 3003,
        role: 'btng-gold-engine',
        timestamp: Math.floor(Date.now() / 1000),
        country: 'Ghana',
        weight: 2,
        name: 'Ghana Gold Engine Node'
    }

    const signature = signPayload(registrationPayload, privateKey)
    const requestBody = {
        ...registrationPayload,
        signature
    }

    const registerResponse = await fetch(`${baseUrl}/api/watchtower/nodes`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(requestBody)
    })

    const registerBody = await registerResponse.json()
    console.log(`REGISTER_STATUS=${registerResponse.status}`)
    console.log(`REGISTER_OK=${Boolean(registerBody?.ok)}`)

    if (!registerResponse.ok) {
        console.log(`REGISTER_ERROR=${registerBody?.error || 'unknown error'}`)
        process.exit(1)
    }

    const snapshotResponse = await fetch(`${baseUrl}/api/watchtower/nodes`, { cache: 'no-store' })
    const snapshotBody = await snapshotResponse.json()
    const registeredNodeFound = Array.isArray(snapshotBody?.nodes)
        ? snapshotBody.nodes.some((node) => node.id === registrationPayload.nodeId)
        : false

    console.log(`BASE_URL=${baseUrl}`)
    console.log(`SNAPSHOT_STATUS=${snapshotResponse.status}`)
    console.log(`SNAPSHOT_SIGNED=${Boolean(snapshotBody?.snapshotSigned)}`)
    console.log(`NODE_PRESENT=${registeredNodeFound}`)

    if (!snapshotResponse.ok || !snapshotBody?.snapshotSigned || !registeredNodeFound) {
        process.exit(1)
    }
}

main().catch((error) => {
    console.error(`TEST_ERROR=${error.message}`)
    process.exit(1)
})
