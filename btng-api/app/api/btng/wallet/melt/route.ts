// @ts-nocheck
import { NextRequest, NextResponse } from 'next/server';
import { fabricMelt } from '@/lib/fabric/wallet-bridge';

export async function POST(request: NextRequest) {
    try {
        const globalRef = (globalThis as any);
        const bridgeUrl = globalRef['process']?.['env']?.['BTNG_API_BASE_URL'] || 'http://localhost:64799';

        const body = await request.json();

        if (!body.wallet || !body.amount || body.amount <= 0) {
            return NextResponse.json(
                { error: 'Wallet address and positive amount are required' },
                { status: 400 }
            );
        }

        // Dual-write: .NET backend + Fabric on-chain
        const [backendResponse, fabricResult] = await Promise.allSettled([
            fetch(`${bridgeUrl}/api/wallet/melt`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ wallet: body.wallet, amount: body.amount }),
                signal: globalRef.AbortSignal.timeout(15000),
            }).then(r => r.json()),
            fabricMelt(body.wallet, body.amount, body.gold_price_usd || 0),
        ]);

        const backendData = backendResponse.status === 'fulfilled' ? backendResponse.value : null;
        const fabricData = fabricResult.status === 'fulfilled' ? fabricResult.value : null;

        return NextResponse.json({
            ...backendData,
            fabric: fabricData || { onChain: false, error: fabricResult.status === 'rejected' ? 'Fabric unavailable' : undefined },
            sovereign_network: 'btng712-fabric-network',
        });
    } catch (error) {
        console.error('Wallet melt error:', error);
        return NextResponse.json(
            { error: 'Melt bridge failure', status: 'failed' },
            { status: 503 }
        );
    }
}
