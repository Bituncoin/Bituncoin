import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function GET() {
  try {
    const response = await fetch(`${BTNG_API_BASE}/oracle/price`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Price data unavailable' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Price fetch error:', error);
    // Return mock data when backend is unreachable
    return NextResponse.json(
      { price: 4.27, change24h: 8.34, status: 'mock' },
      { status: 200 }
    );
  }
}