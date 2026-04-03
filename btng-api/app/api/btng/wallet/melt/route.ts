// @ts-nocheck
import { NextRequest, NextResponse } from 'next/server';

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

        const response = await fetch(`${bridgeUrl}/api/wallet/melt`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ wallet: body.wallet, amount: body.amount }),
            signal: globalRef.AbortSignal.timeout(15000),
        });

        const data = await response.json().catch(() => ({ error: 'Unable to parse response' }));

        if (!response.ok) {
            return NextResponse.json(data, { status: response.status });
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error('Wallet melt error:', error);
        return NextResponse.json(
            { error: 'Melt bridge failure', status: 'failed' },
            { status: 503 }
        );
    }
}
