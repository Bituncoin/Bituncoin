#!/usr/bin/env node

/**
 * BTNG JWT Authentication Test
 * Tests JWT login and protected API endpoints
 */

import fetch from 'node-fetch';

const BASE_URL = process.env.BASE_URL || 'http://localhost:3002';

async function testJWTAuthentication() {
  console.log('🔐 Testing BTNG JWT Authentication...\n');

  try {
    // Step 1: Login to get JWT token
    console.log('1️⃣ Logging in to get JWT token...');
    const loginResponse = await fetch(`${BASE_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: 'admin',
        password: 'password123'
      })
    });

    if (!loginResponse.ok) {
      throw new Error(`Login failed: ${loginResponse.status}`);
    }

    const loginData = await loginResponse.json();
    const token = loginData.token;

    console.log('✅ Login successful!');
    console.log('📋 Token received:', token.substring(0, 50) + '...');
    console.log('👤 User:', loginData.user);
    console.log();

    // Step 2: Test protected gold price endpoint without token
    console.log('2️⃣ Testing protected endpoint WITHOUT token...');
    const unprotectedResponse = await fetch(`${BASE_URL}/api/btng/gold/price`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        base_price_gram: 50.00,
        base_price_ounce: 1600.00,
        base_price_kilo: 50000.00
      })
    });

    console.log('❌ Expected 401, got:', unprotectedResponse.status);
    if (unprotectedResponse.status === 401) {
      console.log('✅ Authentication correctly required!');
    }
    console.log();

    // Step 3: Test protected gold price endpoint WITH token
    console.log('3️⃣ Testing protected endpoint WITH token...');
    const protectedResponse = await fetch(`${BASE_URL}/api/btng/gold/price`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        base_price_gram: 51.25,
        base_price_ounce: 1640.00,
        base_price_kilo: 51250.00,
        currencies: ['EUR', 'GBP', 'JPY'],
        timestamp: Date.now()
      })
    });

    if (protectedResponse.ok) {
      const result = await protectedResponse.json();
      console.log('✅ Gold price update successful!');
      console.log('📊 Updated by:', result.updated_by);
      console.log('💰 New price:', result.stored.base_price_gram, 'USD/gram');
    } else {
      console.log('❌ Update failed:', protectedResponse.status);
      const error = await protectedResponse.json();
      console.log('Error:', error.error);
    }
    console.log();

    // Step 4: Test public status endpoint (should work without auth)
    console.log('4️⃣ Testing public status endpoint...');
    const statusResponse = await fetch(`${BASE_URL}/api/btng/gold/price/status`);

    if (statusResponse.ok) {
      const status = await statusResponse.json();
      console.log('✅ Status endpoint accessible!');
      console.log('📈 Service health:', status.service_health);
      if (status.latest_price) {
        console.log('💰 Latest price:', status.latest_price.base_price_gram, 'USD/gram');
      }
    } else {
      console.log('❌ Status check failed:', statusResponse.status);
    }

    console.log('\n🎉 JWT Authentication test complete!');

  } catch (error) {
    console.error('❌ Test failed:', error.message);
    console.log('\n💡 Make sure your Next.js server is running on the correct port');
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  testJWTAuthentication();
}