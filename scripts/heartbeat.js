/* eslint-disable no-console */

const DEFAULT_GENESIS_NODE = process.env.BTNG_GENESIS_NODE || '74.118.126.72:64799'
const DEFAULT_INTERVAL_MS = Number(process.env.BTNG_HEARTBEAT_INTERVAL_MS || 15000)

const peerNodes = (process.env.BTNG_PEER_NODES || '')
    .split(',')
    .map((entry) => entry.trim())
    .filter(Boolean)

let currentLocalHeight = Number(process.env.BTNG_LOCAL_HEIGHT || 12459)

function normalizeHealthUrl(nodeAddress) {
    const trimmed = String(nodeAddress || '').trim()
    if (!trimmed) {
        return ''
    }

    if (trimmed.startsWith('http://') || trimmed.startsWith('https://')) {
        return `${trimmed.replace(/\/$/, '')}/health`
    }

    return `http://${trimmed}/health`
}

async function getNodeHealth(nodeAddress) {
    const healthUrl = normalizeHealthUrl(nodeAddress)
    if (!healthUrl) {
        return {
            ok: false,
            blockHeight: currentLocalHeight
        }
    }

    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 4000)

    try {
        const response = await fetch(healthUrl, {
            method: 'GET',
            signal: controller.signal,
            headers: { Accept: 'application/json' }
        })

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`)
        }

        const payload = await response.json()
        const blockHeight = Number.isFinite(payload?.blockHeight)
            ? payload.blockHeight
            : currentLocalHeight

        return {
            ok: true,
            blockHeight
        }
    } catch {
        return {
            ok: false,
            blockHeight: currentLocalHeight
        }
    } finally {
        clearTimeout(timeoutId)
    }
}

async function checkNetworkPulse() {
    const genesisStatus = await getNodeHealth(DEFAULT_GENESIS_NODE)

    if (genesisStatus.ok && genesisStatus.blockHeight >= currentLocalHeight) {
        currentLocalHeight = genesisStatus.blockHeight
        console.log('✅ Sync Verified: We are holding the Ancient Treasure together.')
        return
    }

    console.log('⚠️ Lighthouse Offline: Initiating Local Consensus Mode.')

    for (const node of peerNodes) {
        const peerStatus = await getNodeHealth(node)
        if (peerStatus.ok && peerStatus.blockHeight > currentLocalHeight) {
            currentLocalHeight = peerStatus.blockHeight
            console.log(`🔄 Recovered newer height from peer ${node}: ${currentLocalHeight}`)
        }
    }
}

async function startHeartbeat() {
    console.log(`🫀 BTNG Heartbeat active. Genesis node: ${DEFAULT_GENESIS_NODE}`)
    if (peerNodes.length > 0) {
        console.log(`🌍 Peer nodes: ${peerNodes.join(', ')}`)
    }

    await checkNetworkPulse()
    setInterval(() => {
        void checkNetworkPulse()
    }, DEFAULT_INTERVAL_MS)
}

void startHeartbeat()
