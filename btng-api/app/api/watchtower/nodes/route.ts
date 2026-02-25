import { createHash, createSign, createVerify } from 'crypto'
import { readFileSync } from 'fs'
import { isAbsolute, resolve } from 'path'
import { NextRequest, NextResponse } from 'next/server'
import { verifyDocumentationIdentity } from '../../../../lib/documentation-identity'

export const runtime = 'nodejs'

type NodeStatus = 'online' | 'offline'

type WatchtowerNode = {
    id: string
    name: string
    endpoint: string
    country: string
    weight: number
}

type RegisteredWatchtowerNode = WatchtowerNode & {
    ip: string
    port: number
    role: string
    lastRegistrationTimestamp: number
}

type NodeRegistrationRequest = {
    nodeId: string
    ip: string
    port: number
    role: string
    timestamp: number
    signature: string
    country?: string
    weight?: number
    name?: string
}

type WatchtowerRegistryState = {
    registry: Map<string, RegisteredWatchtowerNode>
    lastSnapshotHash: string | null
    snapshotSequence: number
}

declare global {
    // eslint-disable-next-line no-var
    var __btngWatchtowerRegistry: WatchtowerRegistryState | undefined
}

type NodePulse = {
    id: string
    name: string
    country: string
    endpoint: string
    status: NodeStatus
    blockHeight: number
    latencyMs: number | null
    heartbeat: 'sync-verified' | 'local-consensus'
    lastSeen: string
}

const DEFAULT_NODES: WatchtowerNode[] = [
    {
        id: 'ghana-parent',
        name: 'Ghana Parent Node',
        endpoint: 'http://74.118.126.72:64799',
        country: 'Ghana',
        weight: 1
    },
    {
        id: 'watchtower-relay',
        name: 'Watchtower Relay',
        endpoint: 'http://localhost:3000',
        country: 'Regional Relay',
        weight: 1
    }
]

const DEFAULT_MAX_TIMESTAMP_SKEW_SECONDS = 300

function getRegistryState(): WatchtowerRegistryState {
    if (!globalThis.__btngWatchtowerRegistry) {
        globalThis.__btngWatchtowerRegistry = {
            registry: new Map<string, RegisteredWatchtowerNode>(),
            lastSnapshotHash: null,
            snapshotSequence: 0
        }
    }

    if (typeof globalThis.__btngWatchtowerRegistry.lastSnapshotHash === 'undefined') {
        globalThis.__btngWatchtowerRegistry.lastSnapshotHash = null
    }

    if (!Number.isFinite(globalThis.__btngWatchtowerRegistry.snapshotSequence)) {
        globalThis.__btngWatchtowerRegistry.snapshotSequence = 0
    }

    return globalThis.__btngWatchtowerRegistry
}

function sortKeysRecursively(value: unknown): unknown {
    if (Array.isArray(value)) {
        return value.map((entry) => sortKeysRecursively(entry))
    }

    if (value && typeof value === 'object') {
        const sortedObject: Record<string, unknown> = {}
        Object.keys(value as Record<string, unknown>)
            .sort()
            .forEach((key) => {
                sortedObject[key] = sortKeysRecursively((value as Record<string, unknown>)[key])
            })
        return sortedObject
    }

    return value
}

function toCanonicalJson(value: unknown): string {
    return JSON.stringify(sortKeysRecursively(value))
}

function hashSnapshot(data: string): string {
    return createHash('sha256').update(data).digest('hex')
}

function loadKeyFromEnvOrFile(inlineVarName: string, pathVarName: string): string | null {
    const inlineValue = process.env[inlineVarName]?.trim()
    if (inlineValue) {
        return inlineValue.replace(/\\n/g, '\n')
    }

    const keyPathValue = process.env[pathVarName]?.trim()
    if (!keyPathValue) {
        return null
    }

    const absolutePath = isAbsolute(keyPathValue) ? keyPathValue : resolve(process.cwd(), keyPathValue)
    try {
        return readFileSync(absolutePath, 'utf8')
    } catch {
        return null
    }
}

function getWatchtowerPublicKey(): string | null {
    return loadKeyFromEnvOrFile('BTNG_ES256_PUBLIC_KEY', 'BTNG_ES256_PUBLIC_KEY_PATH')
}

function getWatchtowerPrivateKey(): string | null {
    return loadKeyFromEnvOrFile('BTNG_ES256_PRIVATE_KEY', 'BTNG_ES256_PRIVATE_KEY_PATH')
}

function toUnixSeconds(timestampValue: number): number {
    return timestampValue > 1_000_000_000_000 ? Math.floor(timestampValue / 1000) : Math.floor(timestampValue)
}

function normalizeRegistrationPayload(payload: NodeRegistrationRequest): NodeRegistrationRequest {
    return {
        ...payload,
        nodeId: payload.nodeId.trim(),
        ip: payload.ip.trim(),
        role: payload.role.trim(),
        timestamp: toUnixSeconds(payload.timestamp)
    }
}

function serializeRegistrationPayload(payload: {
    nodeId: string
    ip: string
    port: number
    role: string
    timestamp: number
}): string {
    return JSON.stringify({
        nodeId: payload.nodeId,
        ip: payload.ip,
        port: payload.port,
        role: payload.role,
        timestamp: payload.timestamp
    })
}

function verifyNodeSignature(payloadData: string, signature: string, publicKey: string): boolean {
    const verifier = createVerify('SHA256')
    verifier.update(payloadData)
    verifier.end()
    return verifier.verify(publicKey, signature, 'base64')
}

function signSnapshot(data: string, privateKey: string): string {
    const signer = createSign('SHA256')
    signer.update(data)
    signer.end()
    return signer.sign(privateKey, 'base64')
}

function isDocumentationIdentityEnforcementEnabled(): boolean {
    return String(process.env.BTNG_ENFORCE_DOCUMENTATION_IDENTITY || 'true').toLowerCase() !== 'false'
}

function getDocumentationIdentityErrorMessage(error?: string): string {
    return error || 'Documentation identity verification failed.'
}

function getRegisteredNodes(): WatchtowerNode[] {
    const { registry } = getRegistryState()
    return Array.from(registry.values()).map((node) => ({
        id: node.id,
        name: node.name,
        endpoint: node.endpoint,
        country: node.country,
        weight: node.weight
    }))
}

function mergeNodes(staticNodes: WatchtowerNode[], dynamicNodes: WatchtowerNode[]): WatchtowerNode[] {
    const merged = new Map<string, WatchtowerNode>()
    staticNodes.forEach((node) => merged.set(node.id, node))
    dynamicNodes.forEach((node) => merged.set(node.id, node))
    return Array.from(merged.values())
}

function normalizeNodesFromEnv(): WatchtowerNode[] {
    const envValue = process.env.BTNG_WATCHTOWER_NODES
    if (!envValue) {
        return DEFAULT_NODES
    }

    const parsed = envValue
        .split(',')
        .map((entry) => entry.trim())
        .filter(Boolean)
        .map((entry, index) => {
            const [name, endpoint, country, weightValue] = entry.split('|').map((part) => part.trim())
            if (!name || !endpoint) {
                return null
            }

            const safeWeight = Number.isFinite(Number(weightValue)) ? Math.max(0, Number(weightValue)) : 1

            return {
                id: `node-${index + 1}`,
                name,
                endpoint,
                country: country || 'Unknown',
                weight: safeWeight
            } satisfies WatchtowerNode
        })
        .filter((node): node is WatchtowerNode => node !== null)

    return parsed.length > 0 ? parsed : DEFAULT_NODES
}

async function readNodePulse(node: WatchtowerNode, referenceHeight: number): Promise<NodePulse> {
    const startedAt = Date.now()
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 4000)

    try {
        const response = await fetch(`${node.endpoint}/health`, {
            method: 'GET',
            signal: controller.signal,
            headers: {
                Accept: 'application/json'
            },
            cache: 'no-store'
        })

        clearTimeout(timeout)

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`)
        }

        const payload = (await response.json()) as Record<string, unknown>
        const payloadHeight = payload.blockHeight
        const blockHeight = typeof payloadHeight === 'number' && Number.isFinite(payloadHeight)
            ? payloadHeight
            : referenceHeight

        return {
            id: node.id,
            name: node.name,
            country: node.country,
            endpoint: node.endpoint,
            status: 'online',
            blockHeight,
            latencyMs: Date.now() - startedAt,
            heartbeat: blockHeight >= referenceHeight ? 'sync-verified' : 'local-consensus',
            lastSeen: new Date().toISOString()
        }
    } catch {
        clearTimeout(timeout)
        return {
            id: node.id,
            name: node.name,
            country: node.country,
            endpoint: node.endpoint,
            status: 'offline',
            blockHeight: referenceHeight,
            latencyMs: null,
            heartbeat: 'local-consensus',
            lastSeen: new Date().toISOString()
        }
    }
}

export async function GET() {
    const documentationIdentity = verifyDocumentationIdentity()
    if (isDocumentationIdentityEnforcementEnabled() && !documentationIdentity.verified) {
        return NextResponse.json(
            {
                ok: false,
                error: 'Documentation identity gate blocked snapshot publication.',
                documentation: {
                    available: documentationIdentity.available,
                    verified: documentationIdentity.verified,
                    hashValid: documentationIdentity.hashValid,
                    signatureValid: documentationIdentity.signatureValid,
                    versionHash: documentationIdentity.documentation?.versionHash || null,
                    reason: getDocumentationIdentityErrorMessage(documentationIdentity.error)
                }
            },
            {
                status: 503,
                headers: {
                    'Cache-Control': 'no-store'
                }
            }
        )
    }

    const registryState = getRegistryState()
    const { registry } = registryState
    const staticNodes = normalizeNodesFromEnv()
    const dynamicNodes = getRegisteredNodes()
    const nodes = mergeNodes(staticNodes, dynamicNodes)
    const baseHeight = Number.isFinite(Number(process.env.BTNG_GENESIS_LEDGER_HEIGHT))
        ? Number(process.env.BTNG_GENESIS_LEDGER_HEIGHT)
        : 12459

    const pulses = await Promise.all(nodes.map((node) => readNodePulse(node, baseHeight)))
    const onlineCount = pulses.filter((pulse) => pulse.status === 'online').length
    const consensusDepth = onlineCount
    const ledgerHeight = pulses.reduce(
        (max, pulse) => (pulse.blockHeight > max ? pulse.blockHeight : max),
        baseHeight
    )
    const networkWeight = nodes
        .filter((node) => pulses.find((pulse) => pulse.id === node.id)?.status === 'online')
        .reduce((sum, node) => sum + node.weight, 0)

    const timestampIso = new Date().toISOString()
    const timestamp = Math.floor(Date.now() / 1000)

    const snapshotNodes = pulses.map((pulse) => {
        const registeredNode = registry.get(pulse.id)
        const endpointUrl = new URL(pulse.endpoint)

        return {
            nodeId: pulse.id,
            ip: registeredNode?.ip || endpointUrl.hostname,
            port: registeredNode?.port || Number(endpointUrl.port || 80),
            role: registeredNode?.role || 'watchtower-node',
            sovereign: true,
            status: pulse.status,
            country: pulse.country,
            endpoint: pulse.endpoint,
            blockHeight: pulse.blockHeight,
            latencyMs: pulse.latencyMs,
            heartbeat: pulse.heartbeat,
            lastSeen: pulse.lastSeen
        }
    })

    const previousSnapshotHash = registryState.lastSnapshotHash
    const snapshotCore = {
        timestamp,
        metrics: {
            consensusDepth,
            ledgerHeight,
            networkWeight,
            totalNodes: nodes.length,
            onlineNodes: onlineCount
        },
        nodes: snapshotNodes,
        previousSnapshotHash
    }

    const snapshotCanonical = toCanonicalJson(snapshotCore)
    const snapshotHash = hashSnapshot(snapshotCanonical)
    const snapshotSequence = registryState.snapshotSequence + 1

    const anchoredSnapshot = {
        ...snapshotCore,
        snapshotSequence,
        snapshotHash,
        producedAt: timestampIso
    }

    const privateKey = getWatchtowerPrivateKey()
    const publicKey = getWatchtowerPublicKey()
    const snapshotSignature = privateKey ? signSnapshot(snapshotHash, privateKey) : null

    registryState.snapshotSequence = snapshotSequence
    registryState.lastSnapshotHash = snapshotHash

    return NextResponse.json(
        {
            ok: true,
            timestamp: timestampIso,
            metrics: {
                consensusDepth,
                ledgerHeight,
                networkWeight,
                totalNodes: nodes.length,
                onlineNodes: onlineCount
            },
            nodes: pulses,
            snapshot: anchoredSnapshot,
            signature: snapshotSignature,
            publicKey,
            snapshotHash,
            snapshotSignature,
            snapshotAlgorithm: 'ES256/SHA256',
            snapshotSigned: Boolean(snapshotSignature),
            documentation: {
                available: documentationIdentity.available,
                verified: documentationIdentity.verified,
                hashValid: documentationIdentity.hashValid,
                signatureValid: documentationIdentity.signatureValid,
                versionHash: documentationIdentity.documentation?.versionHash || null,
                timestamp: documentationIdentity.documentation?.timestamp || null,
                source: documentationIdentity.documentation?.source || null,
                signatureAlgorithm: documentationIdentity.documentation?.signatureAlgorithm || 'ES256/SHA256',
                fileCount: documentationIdentity.documentation?.files?.length || 0,
                error: documentationIdentity.error
            }
        },
        {
            headers: {
                'Cache-Control': 'no-store'
            }
        }
    )
}

export async function POST(request: NextRequest) {
    const documentationIdentity = verifyDocumentationIdentity()
    if (isDocumentationIdentityEnforcementEnabled() && !documentationIdentity.verified) {
        return NextResponse.json(
            {
                ok: false,
                error: 'Documentation identity gate blocked node registration.',
                documentation: {
                    available: documentationIdentity.available,
                    verified: documentationIdentity.verified,
                    hashValid: documentationIdentity.hashValid,
                    signatureValid: documentationIdentity.signatureValid,
                    versionHash: documentationIdentity.documentation?.versionHash || null,
                    reason: getDocumentationIdentityErrorMessage(documentationIdentity.error)
                }
            },
            { status: 503 }
        )
    }

    const publicKey = getWatchtowerPublicKey()
    if (!publicKey) {
        return NextResponse.json(
            {
                ok: false,
                error: 'BTNG ES256 public key is not configured for node registration verification.'
            },
            { status: 503 }
        )
    }

    let payload: NodeRegistrationRequest
    try {
        payload = (await request.json()) as NodeRegistrationRequest
    } catch {
        return NextResponse.json(
            {
                ok: false,
                error: 'Invalid JSON payload.'
            },
            { status: 400 }
        )
    }

    const normalizedPayload = normalizeRegistrationPayload(payload)
    const {
        nodeId,
        ip,
        port,
        role,
        timestamp,
        signature,
        country = 'Unknown',
        weight,
        name
    } = normalizedPayload

    if (!nodeId || !ip || !role || !signature || !Number.isFinite(Number(port)) || !Number.isFinite(Number(timestamp))) {
        return NextResponse.json(
            {
                ok: false,
                error: 'Missing or invalid required fields: nodeId, ip, port, role, timestamp, signature.'
            },
            { status: 400 }
        )
    }

    const maxSkewSeconds = Number.isFinite(Number(process.env.BTNG_WATCHTOWER_MAX_TIMESTAMP_SKEW_SECONDS))
        ? Number(process.env.BTNG_WATCHTOWER_MAX_TIMESTAMP_SKEW_SECONDS)
        : DEFAULT_MAX_TIMESTAMP_SKEW_SECONDS

    const nowSeconds = Math.floor(Date.now() / 1000)
    const timestampDelta = Math.abs(nowSeconds - timestamp)
    if (timestampDelta > maxSkewSeconds) {
        return NextResponse.json(
            {
                ok: false,
                error: `Timestamp outside allowed skew window (${maxSkewSeconds}s).`
            },
            { status: 401 }
        )
    }

    const { registry } = getRegistryState()
    const previous = registry.get(nodeId)
    if (previous && timestamp <= previous.lastRegistrationTimestamp) {
        return NextResponse.json(
            {
                ok: false,
                error: 'Replay protection triggered: registration timestamp is not newer than current record.'
            },
            { status: 409 }
        )
    }

    const serializedPayload = serializeRegistrationPayload({
        nodeId,
        ip,
        port: Number(port),
        role,
        timestamp
    })

    const signatureValid = verifyNodeSignature(serializedPayload, signature, publicKey)
    if (!signatureValid) {
        return NextResponse.json(
            {
                ok: false,
                error: 'Signature verification failed. Node registration rejected.'
            },
            { status: 401 }
        )
    }

    const safeWeight = Number.isFinite(Number(weight)) ? Math.max(0, Number(weight)) : 1
    const registeredNode: RegisteredWatchtowerNode = {
        id: nodeId,
        name: name?.trim() || `${role}-${nodeId}`,
        endpoint: `http://${ip}:${Number(port)}`,
        country: country.trim() || 'Unknown',
        weight: safeWeight,
        ip,
        port: Number(port),
        role,
        lastRegistrationTimestamp: timestamp
    }

    registry.set(nodeId, registeredNode)

    return NextResponse.json(
        {
            ok: true,
            message: 'Node registration accepted and signature verified.',
            node: {
                id: registeredNode.id,
                endpoint: registeredNode.endpoint,
                role: registeredNode.role,
                country: registeredNode.country,
                registeredAt: new Date().toISOString()
            }
        },
        {
            status: 201,
            headers: {
                'Cache-Control': 'no-store'
            }
        }
    )
}
