import { NextRequest, NextResponse } from 'next/server';
import { saveGoldPrice, setLatestGoldPrice, getLatestGoldPrice, getGoldPriceHistory } from '@/lib/gold-price/model';
import connectToDatabase from '@/lib/mongodb';

export async function GET(request: NextRequest) {
  try {
    // Ensure database connection
    await connectToDatabase();

    // Create test data
    const testData = {
      base_currency: "USD",
      base_price_gram: 76.25,
      base_price_ounce: 2371.50,
      base_price_kilo: 76250.00,
      currencies: [
        { currency: "GHS", price_gram: 1065.00, price_ounce: 33100.00, price_kilo: 1065000.00 },
        { currency: "NGN", price_gram: 118000.00, price_ounce: 3665000.00, price_kilo: 118000000.00 }
      ],
      fx_rates: { GHS: 13.97, NGN: 1548.50 },
      bid: 2368.00,
      ask: 2375.00,
      spread: 7.00,
      timestamp: Date.now()
    };

    // Save to database
    const saved = await saveGoldPrice(testData);

    // Update cache
    await setLatestGoldPrice(testData);

    // Retrieve from cache
    const latest = getLatestGoldPrice();

    // Retrieve from history
    const history = await getGoldPriceHistory(5);

    return NextResponse.json({
      test: "gold_price_system",
      saved: saved ? "success" : "failed",
      cache_updated: latest ? "success" : "failed",
      history_count: history.length,
      latest_price: latest ? {
        base_price_gram: latest.base_price_gram,
        currencies_count: latest.currencies.length
      } : null,
      timestamp: new Date().toISOString()
    }, {
      status: 200,
      headers: {
        'X-BTNG-Platform': 'Sovereign',
        'X-Service-Status': 'Test'
      }
    });

  } catch (error: any) {
    console.error("BTNG gold price test endpoint error:", error.message);
    return NextResponse.json(
      { error: "Test failed", details: error.message },
      { status: 500 }
    );
  }
}