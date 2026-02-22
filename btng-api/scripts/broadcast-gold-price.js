#!/usr/bin/env node

import { buildMultiCurrencyGoldPrice } from '../lib/gold-price/service.js';
import { saveGoldPrice, setLatestGoldPrice } from '../lib/gold-price/model.js';
import connectToDatabase from '../lib/mongodb.js';

async function broadcastGoldPrice() {
  try {
    console.log('🔄 Connecting to database...');
    await connectToDatabase();

    console.log('📡 Fetching live gold price from GoldAPI...');
    const goldPriceData = await buildMultiCurrencyGoldPrice();

    console.log('💾 Saving gold price to database...');
    await saveGoldPrice(goldPriceData);

    console.log('⚡ Updating latest price cache...');
    await setLatestGoldPrice(goldPriceData);

    console.log('✅ BTNG Gold Price Broadcast Complete!');
    console.log('📊 Price Data:', {
      timestamp: new Date(goldPriceData.timestamp).toISOString(),
      base_price_gram: `$${goldPriceData.base_price_gram.toFixed(2)}`,
      base_price_ounce: `$${goldPriceData.base_price_ounce.toFixed(2)}`,
      currencies: goldPriceData.currencies.length,
      fx_rates: Object.keys(goldPriceData.fx_rates).length
    });

    process.exit(0);
  } catch (error) {
    console.error('❌ BTNG Gold Price Broadcast Failed:', error.message);
    process.exit(1);
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  broadcastGoldPrice();
}

export { broadcastGoldPrice };