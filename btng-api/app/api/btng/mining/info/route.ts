import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function GET() {
  try {
    const response = await fetch(`${BTNG_API_BASE}/mining/info`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Mining info unavailable' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Mining info error:', error);
    // Return mock data when backend is unreachable
    return NextResponse.json(
      { difficulty: 4, networkHashrate: '24.5 TH/s', status: 'mock' },
      { status: 200 }
    );
  }
}