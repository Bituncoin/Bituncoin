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

    const response = await fetch(`${BTNG_API_BASE}/explorer/address/${address}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Address not found' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Address fetch error:', error);
    // Return mock data when backend is unreachable
    const mockAddress = {
      address: address,
      balance: Math.random() * 1000,
      txCount: Math.floor(Math.random() * 50),
      status: 'mock'
    };
    return NextResponse.json(
      { address: mockAddress },
      { status: 200 }
    );
  }
}