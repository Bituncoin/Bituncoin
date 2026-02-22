import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function GET() {
  try {
    const response = await fetch(`${BTNG_API_BASE}/health`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      // Add timeout
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Backend health check failed' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Health check error:', error);
    // Return mock data when backend is unreachable
    return NextResponse.json(
      { status: 'mock', message: 'Backend unreachable - using mock data', timestamp: new Date().toISOString() },
      { status: 200 }
    );
  }
}