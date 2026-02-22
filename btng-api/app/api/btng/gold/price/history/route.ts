import { NextRequest, NextResponse } from 'next/server';
import { getGoldPriceHistory, getGoldPriceByTimestampRange } from '@/lib/gold-price/model';
import connectToDatabase from '@/lib/mongodb';

export async function GET(request: NextRequest) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const { searchParams } = new URL(request.url);
    const limit = parseInt(searchParams.get('limit') || '100');
    const startTime = searchParams.get('startTime');
    const endTime = searchParams.get('endTime');

    let history;

    if (startTime && endTime) {
      // Fetch by timestamp range
      const start = parseInt(startTime);
      const end = parseInt(endTime);

      if (isNaN(start) || isNaN(end)) {
        return NextResponse.json(
          { error: "Invalid timestamp parameters. Use Unix timestamps in milliseconds." },
          { status: 400 }
        );
      }

      history = await getGoldPriceByTimestampRange(start, end);
    } else {
      // Fetch recent history with limit
      const safeLimit = Math.min(Math.max(limit, 1), 1000); // Limit between 1 and 1000
      history = await getGoldPriceHistory(safeLimit);
    }

    return NextResponse.json({
      data: history,
      count: history.length,
      message: `Retrieved ${history.length} gold price records`
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Operational',
        'Cache-Control': 'public, max-age=60' // Cache for 1 minute
      }
    });

  } catch (error: any) {
    console.error("BTNG gold price history endpoint error:", error.message);
    return NextResponse.json(
      { error: "Internal server error.", details: error.message },
      { status: 500 }
    );
  }
}