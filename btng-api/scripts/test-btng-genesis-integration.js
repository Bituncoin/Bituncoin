#!/usr/bin/env node

/**
 * BTNG Integration Sanity Check
 *
 * Quick verification that BTNG API is operational
 */

import { config } from 'dotenv';
config(); // Load environment variables
import fetch from 'node-fetch';

async function test() {
  console.log('\n🔍 BTNG Integration Sanity Check\n');

  try {
    // Test 1: Status endpoint (no auth)
    console.log('1. Testing status endpoint (no auth):');
    const statusRes = await fetch('http://localhost:3003/api/btng/gold/price/status');
    if (!statusRes.ok) {
      throw new Error(`Status endpoint failed: ${statusRes.status}`);
    }
    const statusData = await statusRes.json();
    console.log(`   Status: ${statusRes.status} - ✅ JSON response received`);
    console.log(`   Data: ${JSON.stringify(statusData).substring(0, 100)}...`);

    // Test 2: Login
    console.log('\n2. Testing login:');
    const loginRes = await fetch('http://localhost:3003/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: 'admin', password: 'password123' })
    });

    if (!loginRes.ok) {
      const errorText = await loginRes.text();
      throw new Error(`Login failed: ${loginRes.status} - ${errorText.substring(0, 100)}`);
    }

    const loginData = await loginRes.json();
    const token = loginData.token;
    if (!token) {
      throw new Error('No token received from login');
    }
    console.log(`   Login: ${loginRes.status} - ✅ Token received`);

    // Test 3: Protected POST endpoint
    console.log('\n3. Testing protected gold price endpoint:');
    const postRes = await fetch('http://localhost:3003/api/btng/gold/price', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        base_price_gram: 75.5,
        base_price_ounce: 2400,
        base_price_kilo: 75000
      })
    });

    if (!postRes.ok) {
      const errorText = await postRes.text();
      throw new Error(`Protected endpoint failed: ${postRes.status} - ${errorText.substring(0, 100)}`);
    }

    const postData = await postRes.json();
    console.log(`   Protected endpoint: ${postRes.status} - ✅`);
    console.log(`   Response: ${postData.message || postData.error || 'Gold price updated'}`);

    console.log('\n🎉 BTNG API is fully operational!');
    console.log('💡 Ready for Genesis integration and production deployment');

  } catch (error) {
    console.error('❌ Test failed:', error.message);
    console.log('\n🔧 Troubleshooting:');
    console.log('   - Make sure BTNG API is running: npm run dev:api');
    console.log('   - Check port 3003 is available');
    console.log('   - Verify .env file exists with JWT_SECRET and other required variables');
    console.log('   - Check that the API can read environment variables');
    process.exit(1);
  }
}

test();