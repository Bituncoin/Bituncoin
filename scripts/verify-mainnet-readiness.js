#!/usr/bin/env node

const { default: fetch } = require('node-fetch')

function parseArgs(argv) {
    const args = {}
    for (let index = 0; index < argv.length; index += 1) {
        const token = argv[index]
        if (!token.startsWith('--')) {
            continue
        }

        const key = token.slice(2)
        const next = argv[index + 1]
        if (!next || next.startsWith('--')) {
            args[key] = true
            continue
        }

        args[key] = next
        index += 1
    }
    return args
}

function toNumber(value, fallback) {
    const num = Number(value)
    return Number.isFinite(num) ? num : fallback
}

async function verifyMainnetReadiness() {
    const args = parseArgs(process.argv.slice(2))
    const quiet = Boolean(args.quiet) || String(process.env.BTNG_READINESS_QUIET || 'false').toLowerCase() === 'true'
    const watchtowerBaseUrl = String(
        args.watchtowerUrl || process.env.BTNG_WATCHTOWER_URL || 'http://localhost:3001'
    ).replace(/\/$/, '')

    const minOnlineNodes = toNumber(args.minOnlineNodes || process.env.BTNG_MIN_ONLINE_NODES, 2)
    const requireSignedSnapshot = String(
        args.requireSignedSnapshot || process.env.BTNG_REQUIRE_SIGNED_SNAPSHOT || 'true'
    ).toLowerCase() !== 'false'
    const requireAnchoredSnapshot = String(
        args.requireAnchoredSnapshot || process.env.BTNG_REQUIRE_ANCHORED_SNAPSHOT || 'true'
    ).toLowerCase() !== 'false'
    const requireDocumentationIdentity = String(
        args.requireDocumentationIdentity || process.env.BTNG_REQUIRE_DOCUMENTATION_IDENTITY || 'true'
    ).toLowerCase() !== 'false'

    const result = {
        ok: false,
        watchtowerBaseUrl,
        checks: [],
        timestamp: new Date().toISOString()
    }

    const addCheck = (name, ok, details) => {
        result.checks.push({ name, ok, details })
    }

    try {
        const snapshotResponse = await fetch(`${watchtowerBaseUrl}/api/watchtower/nodes`, { cache: 'no-store' })
        addCheck('watchtower_api_reachable', snapshotResponse.ok, `status=${snapshotResponse.status}`)
        if (!snapshotResponse.ok) {
            throw new Error(`Watchtower endpoint returned HTTP ${snapshotResponse.status}`)
        }

        const snapshotBody = await snapshotResponse.json()
        const metrics = snapshotBody?.metrics || {}
        const snapshot = snapshotBody?.snapshot || {}
        const documentation = snapshotBody?.documentation || {}

        const onlineNodes = toNumber(metrics.onlineNodes, 0)
        const totalNodes = toNumber(metrics.totalNodes, 0)
        const snapshotSequence = toNumber(snapshot.snapshotSequence, 0)
        const hasPrevious = Boolean(snapshot.previousSnapshotHash)
        const snapshotSigned = Boolean(snapshotBody?.snapshotSigned)

        addCheck(
            'minimum_online_nodes',
            onlineNodes >= minOnlineNodes,
            `online=${onlineNodes}, required=${minOnlineNodes}, total=${totalNodes}`
        )
        addCheck('snapshot_sequence_progressed', snapshotSequence >= 2, `sequence=${snapshotSequence}`)

        if (requireAnchoredSnapshot) {
            addCheck('snapshot_anchored', hasPrevious, `previousHashPresent=${hasPrevious}`)
        }

        if (requireSignedSnapshot) {
            addCheck('snapshot_signed', snapshotSigned, `snapshotSigned=${snapshotSigned}`)
        }

        if (requireDocumentationIdentity) {
            addCheck(
                'documentation_identity_available',
                Boolean(documentation.available),
                `available=${Boolean(documentation.available)}`
            )
            addCheck(
                'documentation_identity_verified',
                Boolean(documentation.verified),
                `verified=${Boolean(documentation.verified)}`
            )
            addCheck(
                'documentation_version_hash_present',
                Boolean(documentation.versionHash),
                `versionHash=${documentation.versionHash || 'missing'}`
            )
        }

        const verifyResponse = await fetch(`${watchtowerBaseUrl}/api/verify/snapshot`, { cache: 'no-store' })
        addCheck('verify_snapshot_endpoint', verifyResponse.ok, `status=${verifyResponse.status}`)

        if (verifyResponse.ok) {
            const verifyBody = await verifyResponse.json()
            const signatureValid = Boolean(
                verifyBody?.signatureValid ?? verifyBody?.verification?.signatureValid
            )
            addCheck('signature_valid_from_verify_endpoint', signatureValid, `signatureValid=${signatureValid}`)
        }

        result.ok = result.checks.every((check) => check.ok)
    } catch (error) {
        addCheck('exception', false, error.message)
        result.ok = false
    }

    const printable = JSON.stringify(result, null, 2)
    if (result.ok) {
        if (!quiet) {
            console.log('MAINNET_READINESS=PASS')
            console.log(printable)
        }
        return 0
    }

    if (!quiet) {
        console.error('MAINNET_READINESS=FAIL')
        console.error(printable)
    }
    return 1
}

if (require.main === module) {
    verifyMainnetReadiness()
        .then((code) => process.exit(code))
        .catch((error) => {
            console.error(`MAINNET_READINESS=FAIL\n${error.message}`)
            process.exit(1)
        })
}

module.exports = { verifyMainnetReadiness }
