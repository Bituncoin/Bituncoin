// @ts-nocheck
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
    try {
        const globalRef = (globalThis as any);
        const bridgeUrl = globalRef['process']?.['env']?.['BTNG_API_BASE_URL'] || 'http://localhost:64799';

        const response = await fetch(`${bridgeUrl}/api/wallet/ledger`, {
            cache: 'no-store',
            signal: globalRef.AbortSignal.timeout(15000),
        });

        const data = await response.json().catch(() => ({ error: 'Unable to parse response' }));

        if (!response.ok) {
            return NextResponse.json(data, { status: response.status });
        }

        return NextResponse.json(data);
    } catch (error) {
        console.error('Wallet ledger fetch error:', error);
        return NextResponse.json(
            { error: 'Ledger bridge failure', status: 'failed' },
            { status: 503 }
        );
    }
}
