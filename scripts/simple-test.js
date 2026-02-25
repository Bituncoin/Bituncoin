#!/usr/bin/env node

/**
 * Simple BTNG Integration Test
 */

const { default: fetch } = require('node-fetch');

async function testLogin() {
  console.log('Testing JWT Login...');
  const response = await fetch('http://localhost:3004/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: 'admin', password: 'btng-admin-2024' })
  });
  const data = await response.json();
  console.log('✅ Login:', response.ok ? 'SUCCESS' : 'FAILED');
  return data.token;
}

async function testGoldPriceUpdate(token) {
  console.log('Testing Gold Price Update with Token Minting...');
  const response = await fetch('http://localhost:3004/api/btng/gold/price', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      base_currency: "USD",
      base_price_gram: 65.50,
      base_price_ounce: 2038.75,
      base_price_kilo: 65500.00,
      currencies: ["EUR", "GBP"],
      fx_rates: { EUR: 0.85, GBP: 0.73 },
      timestamp: Date.now()
    })
  });
  const data = await response.json();
  console.log('✅ Gold Price Update:', response.ok ? 'SUCCESS' : 'FAILED');
  console.log('   Token Minted:', data.token_minted);
  return response.ok;
}

async function testFabricNetwork(token) {
  console.log('Testing Fabric Network...');
  const response = await fetch('http://localhost:3004/api/btng/fabric/network', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  const data = await response.json();
  console.log('✅ Fabric Network:', response.ok ? 'SUCCESS' : 'FAILED');
  if (response.ok) {
    console.log(`   Network ID: ${data.network_id}`);
    console.log(`   Node ID: ${data.node_id}`);
  }
  return response.ok;
}

async function main() {
  console.log('🚀 BTNG Sovereign Gold Token Integration Test\n');

  try {
    const token = await testLogin();
    if (!token) return;

    const goldUpdate = await testGoldPriceUpdate(token);
    const fabricTest = await testFabricNetwork(token);

    console.log('\n📊 Results:');
    console.log(`   Authentication: ✅`);
    console.log(`   Gold Price + Token Minting: ${goldUpdate ? '✅' : '❌'}`);
    console.log(`   Fabric Network: ${fabricTest ? '✅' : '❌'}`);

    if (goldUpdate && fabricTest) {
      console.log('\n🎉 BTNG Sovereign Gold Token System is OPERATIONAL!');
    }
  } catch (error) {
    console.error('❌ Test failed:', error.message);
  }
}

if (require.main === module) {
  main();
}