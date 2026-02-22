import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { from, to, amount, fee, signature } = body;

    // Basic validation
    if (!from || !to || !amount || !signature) {
      return NextResponse.json(
        { error: 'Missing required fields' },
        { status: 400 }
      );
    }

    if (!from.startsWith('BTNG1') || !to.startsWith('BTNG1')) {
      return NextResponse.json(
        { error: 'Invalid BTNG addresses' },
        { status: 400 }
      );
    }

    if (amount <= 0 || fee < 0) {
      return NextResponse.json(
        { error: 'Invalid amount or fee' },
        { status: 400 }
      );
    }

    const response = await fetch(`${BTNG_API_BASE}/wallet/send`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ from, to, amount, fee, signature }),
      signal: AbortSignal.timeout(30000), // Longer timeout for transactions
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Transaction failed' }));
      return NextResponse.json(
        { error: errorData.error || 'Transaction failed' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Send transaction error:', error);
    return NextResponse.json(
      { error: 'Backend unreachable', status: 'failed' },
      { status: 503 }
    );
  }
}