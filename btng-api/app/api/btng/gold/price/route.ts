import { NextRequest, NextResponse } from 'next/server';
import { saveGoldPrice, setLatestGoldPrice } from '@/lib/gold-price/model';
import { requireAdmin } from '@/lib/auth/jwt';
import connectToDatabase from '@/lib/mongodb';

import { JWTPayload } from '@/lib/auth/jwt';

async function handleGoldPriceUpdate(request: NextRequest, user: JWTPayload) {
  try {
    // Ensure database connection
    await connectToDatabase();

    const body = await request.json();

    const {
      base_currency,
      base_price_gram,
      base_price_ounce,
      base_price_kilo,
      currencies,
      fx_rates,
      bid,
      ask,
      spread,
      timestamp
    } = body;

    // Validate required base price fields
    if (!base_price_gram || !base_price_ounce || !base_price_kilo) {
      return NextResponse.json(
        { error: "Missing required base gold price fields (base_price_gram, base_price_ounce, base_price_kilo)." },
        { status: 400 }
      );
    }

    // Validate base prices are positive numbers
    if (base_price_gram <= 0 || base_price_ounce <= 0 || base_price_kilo <= 0) {
      return NextResponse.json(
        { error: "Base gold prices must be positive numbers." },
        { status: 400 }
      );
    }

    const record = {
      base_currency: base_currency || "USD",
      base_price_gram,
      base_price_ounce,
      base_price_kilo,
      currencies: currencies || [],
      fx_rates: fx_rates || {},
      bid: bid || null,
      ask: ask || null,
      spread: spread || null,
      timestamp: timestamp || Date.now()
    };

    // Save to database
    const savedRecord = await saveGoldPrice(record);

    // Update in-memory cache
    await setLatestGoldPrice(record);

    console.log('BTNG Gold Price Updated by:', user.sub, {
      timestamp: record.timestamp,
      base_price_gram: record.base_price_gram,
      currencies_count: record.currencies.length
    });

    return NextResponse.json({
      status: "ok",
      stored: savedRecord,
      updated_by: user.sub,
      message: "Gold price update stored successfully"
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Operational',
        'X-Auth-User': user.sub
      }
    });

  } catch (error: any) {
    console.error("BTNG gold price endpoint error:", error.message);
    return NextResponse.json(
      { error: "Internal server error.", details: error.message },
      { status: 500 }
    );
  }
}

export const POST = requireAdmin(handleGoldPriceUpdate);