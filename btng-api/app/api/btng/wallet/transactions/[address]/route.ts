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

    const response = await fetch(`${BTNG_API_BASE}/wallet/transactions/${address}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(15000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to fetch transactions' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Transactions fetch error:', error);
    // Return mock data when backend is unreachable
    const mockTransactions = [
      {
        hash: '0x1234567890abcdef',
        from: 'BTNG1MOCKFROM12345678901234567890',
        to: address,
        amount: 1.5,
        status: 'confirmed',
        timestamp: Date.now() - 3600000
      }
    ];
    return NextResponse.json(
      { transactions: mockTransactions, address: address, status: 'mock' },
      { status: 200 }
    );
  }
}