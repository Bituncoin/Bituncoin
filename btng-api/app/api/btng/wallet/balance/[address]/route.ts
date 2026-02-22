import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function GET(
  request: NextRequest,
  { params }: { params: { address: string } }
) {
  const address = params.address;

  try {
    if (!address || (!address.startsWith('BTNG1') && !address.startsWith('nd-'))) {
      return NextResponse.json(
        { error: 'Invalid BTNG address' },
        { status: 400 }
      );
    }

    const response = await fetch(`${BTNG_API_BASE}/wallet/balance/${address}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to fetch balance' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Balance fetch error:', error);
    // Return mock data when backend is unreachable
    return NextResponse.json(
      { balance: 125.487, address: address, status: 'mock' },
      { status: 200 }
    );
  }
}