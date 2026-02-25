#!/usr/bin/env node

const { readFileSync, writeFileSync, mkdirSync } = require('fs')
const { join, resolve } = require('path')

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

function countrySlug(country) {
    return country.trim().toLowerCase().replace(/[^a-z0-9]+/g, '-')
}

function required(args, key) {
    const value = args[key] || process.env[`BTNG_${key.toUpperCase()}`]
    if (!value) {
        throw new Error(`Missing required argument --${key}`)
    }
    return String(value).trim()
}

function optionalNumber(args, key, fallback) {
    const value = args[key]
    if (typeof value === 'undefined') {
        return fallback
    }

    const num = Number(value)
    if (!Number.isFinite(num)) {
        throw new Error(`Invalid numeric value for --${key}`)
    }
    return num
}

function generateSecondNodePackage() {
    const args = parseArgs(process.argv.slice(2))
    const country = required(args, 'country')
    const slug = countrySlug(country)

    const orgId = String(args.orgId || `btng-${slug}-member`)
    const mspid = String(args.mspid || `Btng${country.replace(/[^a-zA-Z0-9]/g, '')}MSP`)
    const peerHost = String(args.peerHost || `peer0.${orgId}.btng-fabric-network.com`)
    const peerPort = optionalNumber(args, 'peerPort', 8051)
    const caHost = String(args.caHost || `ca.${orgId}.btng-fabric-network.com`)
    const caPort = optionalNumber(args, 'caPort', 8054)
    const watchtowerBaseUrl = String(args.watchtowerUrl || process.env.BTNG_WATCHTOWER_URL || 'http://localhost:3001').replace(/\/$/, '')
    const nodeId = String(args.nodeId || `${slug}-peer-0`)

    const networkConfigPath = resolve(process.cwd(), 'btng-fabric-network.json')
    const networkConfig = JSON.parse(readFileSync(networkConfigPath, 'utf8'))

    const channelName = Object.keys(networkConfig.channels || {})[0] || 'btng-sovereign-channel'
    const orderer = Object.keys(networkConfig.orderers || {})[0] || 'orderer.btng-fabric-network.com:7050'

    const packageObject = {
        generatedAt: new Date().toISOString(),
        country,
        nodeProfile: {
            nodeId,
            organization: orgId,
            mspid,
            peerEndpoint: `${peerHost}:${peerPort}`,
            caEndpoint: `${caHost}:${caPort}`,
            sovereign: true,
            role: 'regional-validator'
        },
        fabricJoin: {
            channelName,
            orderer,
            commands: [
                `peer channel fetch 0 ${channelName}.block -c ${channelName} -o ${orderer} --tls --cafile <ORDERER_CA_FILE>`,
                `peer channel join -b ${channelName}.block`,
                `peer lifecycle chaincode queryinstalled`
            ]
        },
        watchtowerRegistration: {
            endpoint: `${watchtowerBaseUrl}/api/watchtower/nodes`,
            payloadTemplate: {
                nodeId,
                ip: '<PUBLIC_OR_PRIVATE_IP>',
                port: peerPort,
                role: 'regional-validator',
                timestamp: Math.floor(Date.now() / 1000),
                country,
                weight: 1,
                name: `${country} Sovereign Node`,
                signature: '<BASE64_ES256_SIGNATURE>'
            },
            signingStringTemplate: JSON.stringify({
                nodeId,
                ip: '<PUBLIC_OR_PRIVATE_IP>',
                port: peerPort,
                role: 'regional-validator',
                timestamp: '<UNIX_SECONDS>'
            })
        },
        postJoinChecklist: [
            'Confirm peer health endpoint responds with blockHeight',
            'Register signed node with watchtower',
            'Run npm run verify:mainnet-readiness',
            'Confirm onlineNodes >= 2 and anchored snapshot sequence increments'
        ]
    }

    const outputDir = resolve(process.cwd(), 'artifacts', 'join-packages')
    mkdirSync(outputDir, { recursive: true })
    const outputPath = join(outputDir, `${slug}-join-package.json`)
    writeFileSync(outputPath, `${JSON.stringify(packageObject, null, 2)}\n`, 'utf8')

    console.log(`JOIN_PACKAGE_CREATED=${outputPath}`)
    console.log(`NODE_ID=${nodeId}`)
    console.log(`CHANNEL=${channelName}`)
}

if (require.main === module) {
    try {
        generateSecondNodePackage()
    } catch (error) {
        console.error(`JOIN_PACKAGE_FAILED=${error.message}`)
        process.exit(1)
    }
}

module.exports = { generateSecondNodePackage }
