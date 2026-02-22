#!/usr/bin/env node

import axios from 'axios';

const BASE_URL = 'http://localhost:64799/api/btng/gold/price';

async function testGoldPriceAPI() {
  try {
    console.log('🧪 Testing BTNG Gold Price API Endpoints...\n');

    // Test 1: Check latest price (should fail initially)
    console.log('1️⃣ Testing GET /latest (should return 404 initially)...');
    try {
      const latestResponse = await axios.get(`${BASE_URL}/latest`);
      console.log('   ✅ Latest price available:', latestResponse.data.base_price_gram);
    } catch (error) {
      if (error.response?.status === 404) {
        console.log('   ✅ Expected 404 - No price data yet');
      } else {
        console.log('   ❌ Unexpected error:', error.message);
      }
    }

    // Test 2: Post a test gold price
    console.log('\n2️⃣ Testing POST /price (storing test data)...');
    const testData = {
      base_currency: "USD",
      base_price_gram: 75.50,
      base_price_ounce: 2348.90,
      base_price_kilo: 75500.00,
      currencies: [
        { currency: "EUR", price_gram: 69.25, price_ounce: 2153.45, price_kilo: 69250.00 },
        { currency: "GHS", price_gram: 1050.75, price_ounce: 32685.35, price_kilo: 1050750.00 }
      ],
      fx_rates: {
        EUR: 0.917,
        GHS: 13.90
      },
      bid: 2345.50,
      ask: 2352.30,
      spread: 6.80,
      timestamp: Date.now()
    };

    const postResponse = await axios.post(BASE_URL, testData);
    console.log('   ✅ Test data stored successfully');
    console.log('   📊 Stored record ID:', postResponse.data.stored._id);

    // Test 3: Check latest price again (should work now)
    console.log('\n3️⃣ Testing GET /latest (should return test data)...');
    const latestResponse2 = await axios.get(`${BASE_URL}/latest`);
    console.log('   ✅ Latest price retrieved:');
    console.log('   💰 Base price per gram:', `$${latestResponse2.data.base_price_gram}`);
    console.log('   💰 Base price per ounce:', `$${latestResponse2.data.base_price_ounce}`);
    console.log('   🌍 Currencies available:', latestResponse2.data.currencies.length);

    // Test 4: Check status endpoint
    console.log('\n4️⃣ Testing GET /status...');
    const statusResponse = await axios.get(`${BASE_URL}/status`);
    console.log('   ✅ Status retrieved:');
    console.log('   🔄 Broadcaster running:', statusResponse.data.broadcaster.isRunning);
    console.log('   📈 Latest price available:', !!statusResponse.data.latest_price);

    // Test 5: Check history endpoint
    console.log('\n5️⃣ Testing GET /history...');
    const historyResponse = await axios.get(`${BASE_URL}/history?limit=5`);
    console.log('   ✅ History retrieved:', historyResponse.data.count, 'records');

    console.log('\n🎉 All BTNG Gold Price API tests passed!');
    console.log('💡 Next: Run "npm run broadcast-gold" to fetch real GoldAPI data');

  } catch (error) {
    console.error('❌ API Test failed:', error.message);
    if (error.response) {
      console.error('   Response status:', error.response.status);
      console.error('   Response data:', error.response.data);
    }
    process.exit(1);
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  testGoldPriceAPI();
}