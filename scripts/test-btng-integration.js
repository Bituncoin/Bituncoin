#!/usr/bin/env node

/**
 * BTNG Sovereign Gold Token Integration Test
 * Tests the complete flow: JWT auth -> Gold price update -> Token minting -> Token operations
 */

const { default: fetch } = require('node-fetch');

const API_BASE = 'http://localhost:3003';
let authToken = '';

async function logStep(step, message) {
  console.log(`\n[${step}] ${message}`);
}

async function login() {
  logStep('AUTH', 'Logging in as admin...');

  try {
    const response = await fetch(`${API_BASE}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: 'admin', password: 'btng-admin-2024' })
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(`Login failed: ${data.error}`);
    }

    authToken = data.token;
    console.log('✅ Login successful, token received');
    return true;
  } catch (error) {
    console.error('❌ Login failed:', error.message);
    return false;
  }
}

async function updateGoldPrice() {
  logStep('GOLD PRICE', 'Updating gold price with token minting...');

  const goldPriceData = {
    base_currency: "USD",
    base_price_gram: 65.50,
    base_price_ounce: 2038.75,
    base_price_kilo: 65500.00,
    currencies: ["EUR", "GBP", "JPY"],
    fx_rates: { EUR: 0.85, GBP: 0.73, JPY: 110.0 },
    bid: 65.45,
    ask: 65.55,
    spread: 0.10,
    timestamp: Date.now()
  };

  try {
    const response = await fetch(`${API_BASE}/api/btng/gold/price`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      },
      body: JSON.stringify(goldPriceData)
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(`Gold price update failed: ${data.error}`);
    }

    console.log('✅ Gold price updated successfully');
    console.log(`   Token minted: ${data.token_minted}`);
    console.log(`   Transaction ID: ${data.transaction_id || 'N/A'}`);
    return data.token_minted;
  } catch (error) {
    console.error('❌ Gold price update failed:', error.message);
    return false;
  }
}

async function checkTokenBalance() {
  logStep('TOKEN BALANCE', 'Checking admin token balance...');

  try {
    const response = await fetch(`${API_BASE}/api/btng/tokens/balance`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(`Balance check failed: ${data.error}`);
    }

    console.log('✅ Token balance retrieved');
    console.log(`   Balance: ${data.balance} BTNG tokens`);
    return data.balance;
  } catch (error) {
    console.error('❌ Balance check failed:', error.message);
    return null;
  }
}

async function transferTokens(amount) {
  logStep('TOKEN TRANSFER', `Transferring ${amount} tokens to test-user...`);

  try {
    const response = await fetch(`${API_BASE}/api/btng/tokens/transfer`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      },
      body: JSON.stringify({
        to_user_id: 'test-user',
        amount: amount
      })
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(`Token transfer failed: ${data.error}`);
    }

    console.log('✅ Token transfer successful');
    console.log(`   Transaction ID: ${data.transaction_id}`);
    return true;
  } catch (error) {
    console.error('❌ Token transfer failed:', error.message);
    return false;
  }
}

async function testFabricNetwork() {
  logStep('FABRIC NETWORK', 'Testing Fabric network connectivity...');

  try {
    const response = await fetch(`${API_BASE}/api/btng/fabric/network`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(`Fabric network test failed: ${data.error}`);
    }

    console.log('✅ Fabric network operational');
    console.log(`   Network ID: ${data.network_id}`);
    console.log(`   Node ID: ${data.node_id}`);
    return true;
  } catch (error) {
    console.error('❌ Fabric network test failed:', error.message);
    return false;
  }
}

async function testChaincodeOperations() {
  logStep('CHAINCODE', 'Testing chaincode operations...');

  try {
    const response = await fetch(`${API_BASE}/api/btng/fabric/chaincode`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`
      }
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(`Chaincode test failed: ${data.error}`);
    }

    console.log('✅ Chaincode operations available');
    console.log(`   Gold Token Chaincode: ${data.chaincodes.includes('btng-gold-token')}`);
    console.log(`   Identity Chaincode: ${data.chaincodes.includes('btng-sovereign-identity')}`);
    return true;
  } catch (error) {
    console.error('❌ Chaincode test failed:', error.message);
    return false;
  }
}

async function runIntegrationTest() {
  console.log('🚀 Starting BTNG Sovereign Gold Token Integration Test');
  console.log('=' .repeat(60));

  const results = {
    login: false,
    fabricNetwork: false,
    chaincode: false,
    goldPriceUpdate: false,
    tokenMinting: false,
    balanceCheck: false,
    tokenTransfer: false
  };

  // Test authentication
  results.login = await login();
  if (!results.login) return results;

  // Test Fabric infrastructure
  results.fabricNetwork = await testFabricNetwork();
  results.chaincode = await testChaincodeOperations();

  // Test gold price update with token minting
  results.goldPriceUpdate = await updateGoldPrice();
  results.tokenMinting = results.goldPriceUpdate; // Assuming if update succeeds, minting was attempted

  // Test token operations
  const balance = await checkTokenBalance();
  results.balanceCheck = balance !== null;

  if (balance > 0) {
    results.tokenTransfer = await transferTokens(100);
  } else {
    console.log('⚠️  Skipping transfer test - no tokens available');
  }

  // Summary
  console.log('\n' + '=' .repeat(60));
  console.log('📊 TEST RESULTS SUMMARY');
  console.log('=' .repeat(60));

  const passed = Object.values(results).filter(Boolean).length;
  const total = Object.keys(results).length;

  Object.entries(results).forEach(([test, passed]) => {
    const status = passed ? '✅' : '❌';
    console.log(`${status} ${test.replace(/([A-Z])/g, ' $1').toLowerCase()}`);
  });

  console.log(`\n🎯 Overall: ${passed}/${total} tests passed`);

  if (passed === total) {
    console.log('🎉 All tests passed! BTNG Sovereign Gold Token system is fully operational.');
  } else {
    console.log('⚠️  Some tests failed. Check the output above for details.');
  }

  return results;
}

// Run the test
if (require.main === module) {
  runIntegrationTest().catch(console.error);
}

module.exports = { runIntegrationTest };