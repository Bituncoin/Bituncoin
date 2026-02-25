'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import styles from './page.module.css'

type Heartbeat = 'sync-verified' | 'local-consensus'
type NodeStatus = 'online' | 'offline'

type NodePulse = {
    id: string
    name: string
    country: string
    endpoint: string
    status: NodeStatus
    blockHeight: number
    latencyMs: number | null
    heartbeat: Heartbeat
    lastSeen: string
}

type WatchtowerPayload = {
    ok: boolean
    timestamp: string
    metrics: {
        consensusDepth: number
        ledgerHeight: number
        networkWeight: number
        totalNodes: number
        onlineNodes: number
    }
    nodes: NodePulse[]
}

const REFRESH_MS = 10000

export default function WatchtowerPage() {
    const [data, setData] = useState<WatchtowerPayload | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const fetchWatchtower = useCallback(async () => {
        try {
            const response = await fetch('/api/watchtower/nodes', { cache: 'no-store' })
            if (!response.ok) {
                throw new Error(`Watchtower request failed: ${response.status}`)
            }

            const payload = (await response.json()) as WatchtowerPayload
            setData(payload)
            setError(null)
        } catch (fetchError) {
            const message = fetchError instanceof Error ? fetchError.message : 'Unable to load watchtower data'
            setError(message)
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        void fetchWatchtower()
        const timerId = setInterval(() => {
            void fetchWatchtower()
        }, REFRESH_MS)

        return () => clearInterval(timerId)
    }, [fetchWatchtower])

    const continentNodes = useMemo(() => Array.from({ length: 54 }), [])

    return (
        <div className={styles.container}>
            <section className={styles.hero}>
                <div className="container">
                    <h1 className={styles.title}>Sovereign Watchtower Map</h1>
                    <p className={styles.subtitle}>
                        Continental node pulse, consensus depth, and heartbeat verification for the independent BTNG mesh.
                    </p>
                </div>
            </section>

            <section className={styles.content}>
                <div className="container">
                    <div className={styles.metricsGrid}>
                        <article className={styles.metricCard}>
                            <h2>Consensus Depth</h2>
                            <p>{data?.metrics.consensusDepth ?? 0}</p>
                        </article>
                        <article className={styles.metricCard}>
                            <h2>Ledger Height</h2>
                            <p>{data?.metrics.ledgerHeight ?? 12459}</p>
                        </article>
                        <article className={styles.metricCard}>
                            <h2>Network Weight</h2>
                            <p>{data?.metrics.networkWeight ?? 0}</p>
                        </article>
                    </div>

                    <div className={styles.mapCard}>
                        <h3>BTNG Continental Mesh</h3>
                        <div className={styles.mesh}>
                            {continentNodes.map((_, index) => (
                                <span key={`mesh-node-${index}`} className={styles.meshNode} />
                            ))}
                        </div>
                        <p className={styles.mapLabel}>BTNG - CONTINENTAL RESERVE</p>
                    </div>

                    <div className={styles.statusHeader}>
                        <h3>Heartbeat Protocol</h3>
                        <span className={styles.refreshMeta}>
                            {loading ? 'Loading...' : `Updated: ${data ? new Date(data.timestamp).toLocaleString() : 'n/a'}`}
                        </span>
                    </div>

                    {error && <p className={styles.error}>{error}</p>}

                    <div className={styles.nodeList}>
                        {(data?.nodes ?? []).map((node) => (
                            <article key={node.id} className={styles.nodeCard}>
                                <div className={styles.nodeTop}>
                                    <h4>{node.name}</h4>
                                    <span className={node.status === 'online' ? styles.online : styles.offline}>
                                        {node.status.toUpperCase()}
                                    </span>
                                </div>
                                <p>{node.country}</p>
                                <p className={styles.endpoint}>{node.endpoint}</p>
                                <div className={styles.nodeMeta}>
                                    <span>Height: {node.blockHeight}</span>
                                    <span>Latency: {node.latencyMs === null ? 'n/a' : `${node.latencyMs} ms`}</span>
                                </div>
                                <p className={styles.heartbeat}>
                                    {node.heartbeat === 'sync-verified'
                                        ? '✅ Sync Verified: Ancient record aligned'
                                        : '⚠️ Local Consensus Mode: waiting lighthouse sync'}
                                </p>
                            </article>
                        ))}
                    </div>
                </div>
            </section>
        </div>
    )
}
