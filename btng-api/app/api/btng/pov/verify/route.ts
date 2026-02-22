import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { address, signature, payload } = body;

    // Basic validation
    if (!address || !signature || !payload) {
      return NextResponse.json(
        { error: 'Missing required fields' },
        { status: 400 }
      );
    }

    if (!address.startsWith('BTNG1')) {
      return NextResponse.json(
        { error: 'Invalid BTNG address' },
        { status: 400 }
      );
    }

    const response = await fetch(`${BTNG_API_BASE}/pov/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ address, signature, payload }),
      signal: AbortSignal.timeout(15000),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Verification failed' }));
      return NextResponse.json(
        { error: errorData.error || 'Verification failed' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('PoV verification error:', error);
    return NextResponse.json(
      { error: 'Backend unreachable', verified: false },
      { status: 503 }
    );
  }
}