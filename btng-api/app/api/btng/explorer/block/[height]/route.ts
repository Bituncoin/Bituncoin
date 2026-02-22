import { NextRequest, NextResponse } from 'next/server';

const BTNG_API_BASE = 'http://74.118.126.72:54328';

export async function GET(
  request: NextRequest,
  { params }: { params: { height: string } }
) {
  const height = parseInt(params.height);

  try {
    if (isNaN(height) || height < 0) {
      return NextResponse.json(
        { error: 'Invalid block height' },
        { status: 400 }
      );
    }

    const response = await fetch(`${BTNG_API_BASE}/explorer/block/${height}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      signal: AbortSignal.timeout(10000),
    });

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Block not found' },
        { status: response.status }
      );
    }

    const data = await response.json();
    return NextResponse.json(data);
  } catch (error) {
    console.error('Block fetch error:', error);
    // Return mock data when backend is unreachable
    const mockBlock = {
      number: height,
      hash: '0x' + Math.random().toString(16).substr(2, 64),
      txCount: Math.floor(Math.random() * 20) + 1,
      timestamp: Date.now() - Math.floor(Math.random() * 86400000),
      status: 'mock'
    };
    return NextResponse.json(
      { block: mockBlock },
      { status: 200 }
    );
  }
}