import { createHash, createVerify } from 'crypto'
import { NextRequest, NextResponse } from 'next/server'
import { verifyDocumentationIdentity } from '../../../../lib/documentation-identity'

export const runtime = 'nodejs'

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

function verifySignature(snapshotHash: string, signature: string, publicKey: string): boolean {
    const verifier = createVerify('SHA256')
    verifier.update(snapshotHash)
    verifier.end()
    return verifier.verify(publicKey, signature, 'base64')
}

export async function GET(request: NextRequest) {
    const { searchParams } = new URL(request.url)
    const watchtowerUrl = searchParams.get('url') || process.env.BTNG_WATCHTOWER_URL || 'http://localhost:3001'

    try {
        console.log(`Verifying snapshot from: ${watchtowerUrl}`)

        // Fetch snapshot from watchtower
        const response = await fetch(`${watchtowerUrl}/api/watchtower/nodes`, {
            headers: { 'Accept': 'application/json' },
            cache: 'no-store'
        })

        if (!response.ok) {
            return NextResponse.json(
                {
                    ok: false,
                    error: `Watchtower returned HTTP ${response.status}`,
                    watchtowerUrl,
                    statusCode: response.status
                },
                { status: 502 }
            )
        }

        const data = await response.json()

        if (!data.ok) {
            return NextResponse.json(
                {
                    ok: false,
                    error: `Watchtower API error: ${data.error || 'Unknown error'}`,
                    watchtowerUrl
                },
                { status: 502 }
            )
        }

        if (!data.snapshot) {
            return NextResponse.json(
                {
                    ok: false,
                    error: 'No snapshot data in watchtower response',
                    watchtowerUrl
                },
                { status: 502 }
            )
        }

        const snapshotCore = {
            timestamp: data.snapshot.timestamp,
            metrics: data.snapshot.metrics,
            nodes: data.snapshot.nodes,
            previousSnapshotHash: data.snapshot.previousSnapshotHash
        }

        const canonicalJson = toCanonicalJson(snapshotCore)
        const computedHash = hashSnapshot(canonicalJson)
        const hashValid = computedHash === data.snapshotHash

        const signature = data.snapshotSignature || data.signature
        const publicKey = data.publicKey
        const signatureValid = Boolean(
            data.snapshotSigned &&
            signature &&
            publicKey &&
            hashValid &&
            verifySignature(data.snapshotHash, signature, publicKey)
        )

        const documentationIdentity = verifyDocumentationIdentity()

        return NextResponse.json({
            ok: true,
            verified: hashValid && signatureValid,
            watchtowerUrl,
            snapshotSequence: data.snapshot.snapshotSequence,
            snapshotHash: data.snapshotHash,
            hashValid,
            signatureValid,
            anchoringValid: Boolean(data.snapshot.previousSnapshotHash),
            networkStatus: `${data.metrics.onlineNodes}/${data.metrics.totalNodes} nodes online`,
            timestamp: data.timestamp,
            verificationTime: new Date().toISOString(),
            documentation: {
                available: documentationIdentity.available,
                verified: documentationIdentity.verified,
                hashValid: documentationIdentity.hashValid,
                signatureValid: documentationIdentity.signatureValid,
                versionHash: documentationIdentity.documentation?.versionHash || null
            }
        })

    } catch (error: any) {
        console.error('Snapshot verification error:', error.message)
        return NextResponse.json(
            {
                ok: false,
                error: `Verification failed: ${error.message}`,
                watchtowerUrl
            },
            { status: 500 }
        )
    }
}