import { NextRequest, NextResponse } from 'next/server';
import { getLatestGoldPrice } from '@/lib/gold-price/model';
import connectToDatabase from '@/lib/mongodb';

export async function GET(request: NextRequest) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const latest = getLatestGoldPrice();

    if (!latest) {
      return NextResponse.json(
        {
          error: "No gold price available yet.",
          message: "The gold price feed may still be initializing. Please try again in a few moments."
        },
        {
          status: 404,
          headers: {
            'X-BTNG-Platform': 'Sovereign',
            'X-Service-Status': 'Initializing'
          }
        }
      );
    }

    return NextResponse.json(latest, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Operational',
        'Cache-Control': 'public, max-age=30' // Cache for 30 seconds
      }
    });

  } catch (error: any) {
    console.error("BTNG gold price latest endpoint error:", error.message);
    return NextResponse.json(
      { error: "Internal server error.", details: error.message },
      { status: 500 }
    );
  }
}