#!/usr/bin/env node

const { default: fetch } = require('node-fetch')
const { mkdirSync, writeFileSync } = require('fs')
const { dirname, resolve } = require('path')
const { verifyDocumentationIdentity } = require('./verify-doc-identity')
const { verifyMainnetReadiness } = require('./verify-mainnet-readiness')

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

function resolveArtifactPath(mode, explicitOutputPath) {
    if (explicitOutputPath) {
        return resolve(process.cwd(), String(explicitOutputPath))
    }

    const suffix = mode === 'soft' ? 'soft' : 'strict'
    return resolve(process.cwd(), 'cache', `verify-btng-network-${suffix}.json`)
}

function writeArtifact(filePath, payload) {
    mkdirSync(dirname(filePath), { recursive: true })
    writeFileSync(filePath, `${JSON.stringify(payload, null, 2)}\n`, 'utf8')
}

async function discoverWatchtowerUrls(startPort, endPort) {
    const urls = []
    for (let port = startPort; port <= endPort; port += 1) {
        const baseUrl = `http://localhost:${port}`
        try {
            const response = await fetch(`${baseUrl}/api/watchtower/meta`, { cache: 'no-store' })
            if (response.ok) {
                urls.push(baseUrl)
            }
        } catch {
            // continue scan
        }
    }
    return urls
}

async function checkWatchtowerBaseUrl(baseUrl) {
    const checks = []
    const addCheck = (name, ok, details) => checks.push({ name, ok, details })

    const metaResponse = await fetch(`${baseUrl}/api/watchtower/meta`, { cache: 'no-store' })
    addCheck('watchtower_meta_reachable', metaResponse.ok, `status=${metaResponse.status}`)
    if (!metaResponse.ok) {
        return checks
    }

    const metaBody = await metaResponse.json()
    addCheck('watchtower_meta_doc_verified', Boolean(metaBody?.documentation?.verified), `verified=${Boolean(metaBody?.documentation?.verified)}`)
    addCheck('watchtower_meta_doc_hash_present', Boolean(metaBody?.documentation?.versionHash), `versionHash=${metaBody?.documentation?.versionHash || 'missing'}`)

    const nodesResponse = await fetch(`${baseUrl}/api/watchtower/nodes`, { cache: 'no-store' })
    addCheck('watchtower_nodes_reachable', nodesResponse.ok, `status=${nodesResponse.status}`)
    if (!nodesResponse.ok) {
        return checks
    }

    const nodesBody = await nodesResponse.json()
    addCheck('watchtower_nodes_snapshot_signed', Boolean(nodesBody?.snapshotSigned), `snapshotSigned=${Boolean(nodesBody?.snapshotSigned)}`)
    addCheck('watchtower_nodes_doc_verified', Boolean(nodesBody?.documentation?.verified), `verified=${Boolean(nodesBody?.documentation?.verified)}`)

    const verifyResponse = await fetch(`${baseUrl}/api/verify/snapshot?url=${encodeURIComponent(baseUrl)}`, { cache: 'no-store' })
    addCheck('watchtower_verify_endpoint', verifyResponse.ok, `status=${verifyResponse.status}`)
    if (verifyResponse.ok) {
        const verifyBody = await verifyResponse.json()
        addCheck('watchtower_verify_signature_valid', Boolean(verifyBody?.signatureValid), `signatureValid=${Boolean(verifyBody?.signatureValid)}`)
        addCheck('watchtower_verify_doc_verified', Boolean(verifyBody?.documentation?.verified), `verified=${Boolean(verifyBody?.documentation?.verified)}`)
    }

    return checks
}

async function verifyBtngNetwork() {
    const args = parseArgs(process.argv.slice(2))
    const mode = String(args.mode || process.env.BTNG_NETWORK_VERIFY_MODE || 'strict').toLowerCase() === 'soft'
        ? 'soft'
        : 'strict'
    const minOnlineNodes = toNumber(args.minOnlineNodes || process.env.BTNG_MIN_ONLINE_NODES, 2)
    const startPort = toNumber(args.startPort, 3001)
    const endPort = toNumber(args.endPort, 3010)
    const outputPath = resolveArtifactPath(mode, args.output || process.env.BTNG_NETWORK_VERIFY_OUTPUT)

    const report = {
        ok: false,
        mode,
        timestamp: new Date().toISOString(),
        checks: [],
        summary: {
            totalChecks: 0,
            passedChecks: 0,
            failedChecks: 0,
            warningChecks: 0
        },
        artifactPath: outputPath
    }

    const addCheck = (name, ok, details, severity = 'error') => {
        report.checks.push({ name, ok, details, severity })
    }

    try {
        const docIdentity = verifyDocumentationIdentity()
        addCheck('documentation_identity_local_verified', Boolean(docIdentity?.ok), `versionHash=${docIdentity?.versionHash || 'unknown'}`)

        const watchtowerUrls = await discoverWatchtowerUrls(startPort, endPort)
        addCheck('watchtower_instances_discovered', watchtowerUrls.length > 0, `count=${watchtowerUrls.length}`)

        for (const baseUrl of watchtowerUrls) {
            const checks = await checkWatchtowerBaseUrl(baseUrl)
            checks.forEach((check) => {
                addCheck(`${baseUrl}:${check.name}`, check.ok, check.details)
            })

            const previousArgv = process.argv.slice()
            process.argv = ['node', 'verify-mainnet-readiness.js', '--watchtowerUrl', baseUrl, '--minOnlineNodes', String(minOnlineNodes), '--quiet']
            const readinessExitCode = await verifyMainnetReadiness()
            process.argv = previousArgv
            const readinessOk = readinessExitCode === 0
            const readinessSeverity = mode === 'soft' ? 'warn' : 'error'
            addCheck(`${baseUrl}:mainnet_readiness`, readinessOk, `exitCode=${readinessExitCode}`, readinessSeverity)
        }

        report.summary.totalChecks = report.checks.length
        report.summary.passedChecks = report.checks.filter((check) => check.ok).length
        report.summary.warningChecks = report.checks.filter((check) => !check.ok && check.severity === 'warn').length
        report.summary.failedChecks = report.checks.filter((check) => !check.ok && check.severity !== 'warn').length
        report.ok = report.summary.failedChecks === 0
    } catch (error) {
        addCheck('exception', false, error.message)
        report.ok = false
        report.summary.totalChecks = report.checks.length
        report.summary.passedChecks = report.checks.filter((check) => check.ok).length
        report.summary.warningChecks = report.checks.filter((check) => !check.ok && check.severity === 'warn').length
        report.summary.failedChecks = report.checks.filter((check) => !check.ok && check.severity !== 'warn').length
    }

    writeArtifact(outputPath, report)
    console.log(`BTNG_NETWORK_VERIFY_ARTIFACT=${outputPath}`)

    if (report.ok) {
        console.log('BTNG_NETWORK_VERIFY=PASS')
        console.log(JSON.stringify(report, null, 2))
        return 0
    }

    console.error('BTNG_NETWORK_VERIFY=FAIL')
    console.error(JSON.stringify(report, null, 2))
    return 1
}

if (require.main === module) {
    verifyBtngNetwork()
        .then((code) => process.exit(code))
        .catch((error) => {
            console.error(`BTNG_NETWORK_VERIFY=FAIL\n${error.message}`)
            process.exit(1)
        })
}

module.exports = { verifyBtngNetwork }
