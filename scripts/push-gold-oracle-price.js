#!/usr/bin/env node

const { JsonRpcProvider, Wallet, Contract } = require('ethers')

const DEFAULT_INTERVAL_MS = 60_000
const DEFAULT_TIMEOUT_MS = 15_000
const DEFAULT_PRICE_FIELD = 'base_price_gram'
const DEFAULT_API_URL = process.env.BTNG_GOLD_PRICE_API_URL || 'http://localhost:64799/api/btng/gold/price/latest'

const ORACLE_ABI = [
    'function updatePrice(uint256 priceUSD) external',
    'function currentPriceUSD() view returns (uint256)'
]

function toPositiveNumber(value) {
    const parsed = Number(value)
    if (!Number.isFinite(parsed) || parsed <= 0) {
        return null
    }
    return parsed
}

function toCents(priceValue) {
    return Math.round(priceValue * 100)
}

function nowIso() {
    return new Date().toISOString()
}

function wait(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms))
}

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

function requiredEnv(name) {
    const value = process.env[name]
    if (!value || !String(value).trim()) {
        throw new Error(`Missing required environment variable: ${name}`)
    }
    return String(value).trim()
}

async function fetchLatestPrice(apiUrl, priceField) {
    const response = await fetch(apiUrl, {
        method: 'GET',
        headers: { Accept: 'application/json' },
        signal: AbortSignal.timeout(DEFAULT_TIMEOUT_MS)
    })

    if (!response.ok) {
        const body = await response.text()
        throw new Error(`Gold API HTTP ${response.status}: ${body}`)
    }

    const payload = await response.json()
    const candidate = payload?.[priceField]
    const numericPrice = toPositiveNumber(candidate)
    if (!numericPrice) {
        throw new Error(`Invalid or missing '${priceField}' in gold price payload`)
    }

    return {
        payload,
        numericPrice,
        cents: toCents(numericPrice)
    }
}

async function pushOnce(config) {
    const { cycleNumber, apiUrl, priceField, oracleContract } = config

    const cycleTag = `[cycle:${cycleNumber}]`
    console.log(`${nowIso()} ${cycleTag} Fetching latest gold price from ${apiUrl}`)

    const latest = await fetchLatestPrice(apiUrl, priceField)
    const onChainRaw = await oracleContract.currentPriceUSD()
    const onChainCents = Number(onChainRaw)

    console.log(
        `${nowIso()} ${cycleTag} price_field=${priceField} source=${latest.numericPrice.toFixed(4)} usd cents=${latest.cents} onchain=${onChainCents}`
    )

    if (onChainCents === latest.cents) {
        console.log(`${nowIso()} ${cycleTag} No update needed (on-chain value already current)`)
        return { updated: false, txHash: null, cents: latest.cents }
    }

    const tx = await oracleContract.updatePrice(latest.cents)
    console.log(`${nowIso()} ${cycleTag} Submitted tx ${tx.hash}`)

    const receipt = await tx.wait()
    if (!receipt || receipt.status !== 1) {
        throw new Error(`Transaction failed or reverted: ${tx.hash}`)
    }

    console.log(`${nowIso()} ${cycleTag} Oracle updated successfully at block ${receipt.blockNumber}`)
    return { updated: true, txHash: tx.hash, cents: latest.cents }
}

async function run() {
    const args = parseArgs(process.argv.slice(2))

    const rpcUrl = requiredEnv('BTNG_ORACLE_RPC_URL')
    const oracleAddress = requiredEnv('BTNG_ORACLE_ADDRESS')
    const privateKeyRaw = requiredEnv('BTNG_ORACLE_ADMIN_PRIVATE_KEY')
    const privateKey = privateKeyRaw.startsWith('0x') ? privateKeyRaw : `0x${privateKeyRaw}`

    const apiUrl = String(args.apiUrl || process.env.BTNG_GOLD_PRICE_API_URL || DEFAULT_API_URL)
    const priceField = String(args.priceField || process.env.BTNG_ORACLE_PRICE_FIELD || DEFAULT_PRICE_FIELD)
    const intervalMs = Number(args.intervalMs || process.env.BTNG_ORACLE_PUSH_INTERVAL_MS || DEFAULT_INTERVAL_MS)
    const once = Boolean(args.once)

    const provider = new JsonRpcProvider(rpcUrl)
    const wallet = new Wallet(privateKey, provider)
    const oracleContract = new Contract(oracleAddress, ORACLE_ABI, wallet)

    console.log(`${nowIso()} Starting BTNG oracle updater`)
    console.log(`${nowIso()} Oracle address: ${oracleAddress}`)
    console.log(`${nowIso()} Updater wallet: ${wallet.address}`)
    console.log(`${nowIso()} Gold API URL: ${apiUrl}`)
    console.log(`${nowIso()} Price field: ${priceField}`)
    console.log(`${nowIso()} Interval ms: ${intervalMs}`)

    let cycleNumber = 1
    let consecutiveFailures = 0

    while (true) {
        try {
            await pushOnce({
                cycleNumber,
                apiUrl,
                priceField,
                oracleContract
            })
            consecutiveFailures = 0
        } catch (error) {
            consecutiveFailures += 1
            const retryDelay = Math.min(intervalMs, 5_000 * 2 ** Math.min(consecutiveFailures - 1, 4))
            console.error(`${nowIso()} [cycle:${cycleNumber}] Update failed: ${error.message}`)
            console.error(`${nowIso()} [cycle:${cycleNumber}] Retrying in ${retryDelay}ms`)
            await wait(retryDelay)

            if (once) {
                process.exit(1)
            }

            cycleNumber += 1
            continue
        }

        if (once) {
            process.exit(0)
        }

        await wait(intervalMs)
        cycleNumber += 1
    }
}

run().catch((error) => {
    console.error(`${nowIso()} Fatal updater error: ${error.message}`)
    process.exit(1)
})
