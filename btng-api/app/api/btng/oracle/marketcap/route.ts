import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function GET() {
  try {
    const response = await fetch(`${BTNG_API_BASE}/oracle/marketcap`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Market cap data unavailable' },
        { status: response.status }
      );
    }

    return NextResponse.json(data);
  } catch (error) {
    console.error('Market cap fetch error:', error);
    // Return mock data when backend is unreachable
    return NextResponse.json(
      { marketcap: 5108400, supply: 1240000, status: 'mock' },
      { status: 200 }
    );
  }
}