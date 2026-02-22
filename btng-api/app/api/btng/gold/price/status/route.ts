import { NextRequest, NextResponse } from 'next/server';
import { goldPriceBroadcaster } from '@/lib/gold-price/broadcaster';
import { getLatestGoldPrice } from '@/lib/gold-price/model';
import connectToDatabase from '@/lib/mongodb';

export async function GET(request: NextRequest) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const broadcasterStatus = goldPriceBroadcaster.getStatus();
    const latestPrice = getLatestGoldPrice();

    const status = {
      broadcaster: broadcasterStatus,
      latest_price: latestPrice ? {
        timestamp: new Date(latestPrice.timestamp).toISOString(),
        base_price_gram: latestPrice.base_price_gram,
        base_price_ounce: latestPrice.base_price_ounce,
        currencies_available: latestPrice.currencies.length,
        last_update_age_seconds: Math.floor((Date.now() - latestPrice.timestamp) / 1000)
      } : null,
      service_health: {
        database: 'connected',
        goldapi: latestPrice ? 'operational' : 'pending_first_update',
        broadcaster: broadcasterStatus.isRunning ? 'running' : 'stopped'
      }
    };

    return NextResponse.json(status, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Operational'
      }
    });

  } catch (error: any) {
    console.error("BTNG gold price status endpoint error:", error.message);
    return NextResponse.json(
      {
        error: "Internal server error.",
        details: error.message,
        service_health: {
          database: 'error',
          goldapi: 'unknown',
          broadcaster: 'error'
        }
      },
      { status: 500 }
    );
  }
}