#!/usr/bin/env node

const { randomUUID } = require('crypto')
const { default: fetch } = require('node-fetch')

const LIVE_CONFIRMATION_TEXT = 'I UNDERSTAND REAL MONEY WILL MOVE'

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

function required(value, label) {
    if (!value) {
        throw new Error(`Missing required value: ${label}`)
    }
    return String(value).trim()
}

function asPositiveAmount(value) {
    const amount = Number(value)
    if (!Number.isFinite(amount) || amount <= 0) {
        throw new Error('Amount must be a positive number')
    }
    return amount.toFixed(2)
}

async function runMomoLive() {
    const args = parseArgs(process.argv.slice(2))

    const amount = asPositiveAmount(args.amount || process.env.BTNG_MOMO_AMOUNT || '0.10')
    const currency = String(args.currency || process.env.BTNG_MOMO_CURRENCY || 'GHS').toUpperCase()
    const msisdn = required(args.msisdn || process.env.BTNG_MOMO_MSISDN, '--msisdn or BTNG_MOMO_MSISDN')
    const payerMessage = String(args.payerMessage || process.env.BTNG_MOMO_PAYER_MESSAGE || 'BTNG sovereign penny mint').slice(0, 160)
    const payeeNote = String(args.payeeNote || process.env.BTNG_MOMO_PAYEE_NOTE || 'BTNG mint reference').slice(0, 160)
    const externalId = String(args.externalId || process.env.BTNG_MOMO_EXTERNAL_ID || `btng-live-${Date.now()}-${randomUUID().slice(0, 8)}`)

    const apiBaseUrl = String(args.baseUrl || process.env.BTNG_MOMO_API_BASE_URL || '').replace(/\/$/, '')
    const token = String(args.token || process.env.BTNG_MOMO_API_TOKEN || '')

    const isExecuteLive = Boolean(args.executeLive)
    const confirmation = String(args.confirm || '')
    const allowRealMoney = String(process.env.BTNG_ENABLE_REAL_MONEY || '').toLowerCase() === 'true'

    const payload = {
        amount,
        currency,
        externalId,
        payer: {
            partyIdType: 'MSISDN',
            partyId: msisdn
        },
        payerMessage,
        payeeNote
    }

    const summary = {
        mode: isExecuteLive ? 'execute-live' : 'dry-run',
        amount,
        currency,
        msisdn,
        externalId,
        endpoint: `${apiBaseUrl || '<BTNG_MOMO_API_BASE_URL>'}/request-to-pay`
    }

    console.log('MOMO_LIVE_SUMMARY')
    console.log(JSON.stringify(summary, null, 2))

    if (!isExecuteLive) {
        console.log('MOMO_LIVE_DRY_RUN=true')
        console.log('Payload preview:')
        console.log(JSON.stringify(payload, null, 2))
        console.log('To execute real transaction, pass --executeLive and --confirm with the exact confirmation text.')
        return 0
    }

    if (!allowRealMoney) {
        throw new Error('Real money mode blocked. Set BTNG_ENABLE_REAL_MONEY=true to allow execution.')
    }

    if (confirmation !== LIVE_CONFIRMATION_TEXT) {
        throw new Error(`Live confirmation text mismatch. Expected: ${LIVE_CONFIRMATION_TEXT}`)
    }

    if (!apiBaseUrl) {
        throw new Error('Missing BTNG_MOMO_API_BASE_URL for live execution')
    }

    if (!token) {
        throw new Error('Missing BTNG_MOMO_API_TOKEN for live execution')
    }

    const response = await fetch(`${apiBaseUrl}/request-to-pay`, {
        method: 'POST',
        headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
            'X-Reference-Id': externalId,
            'X-Target-Environment': process.env.BTNG_MOMO_TARGET_ENV || 'production'
        },
        body: JSON.stringify(payload)
    })

    const bodyText = await response.text()
    const result = {
        status: response.status,
        ok: response.ok,
        externalId,
        responseBody: bodyText
    }

    if (!response.ok) {
        console.error('MOMO_LIVE_EXECUTION=FAIL')
        console.error(JSON.stringify(result, null, 2))
        return 1
    }

    console.log('MOMO_LIVE_EXECUTION=SUCCESS')
    console.log(JSON.stringify(result, null, 2))
    return 0
}

if (require.main === module) {
    runMomoLive()
        .then((code) => process.exit(code))
        .catch((error) => {
            console.error(`MOMO_LIVE_EXECUTION=FAIL\n${error.message}`)
            process.exit(1)
        })
}

module.exports = { runMomoLive }
