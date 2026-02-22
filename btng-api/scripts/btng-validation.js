#!/usr/bin/env node

/**
 * BTNG Real-API Validation Script
 * Tests all endpoints against the MTN backend at 74.118.126.72:54328
 */

const https = require('https');
const http = require('http');

const BTNG_API_BASE = 'http://74.118.126.72:54328';
const NEXTJS_BASE = 'http://localhost:3002';

// Network configuration
const BTNG_FABRIC_NETWORK = 'btng-fabric-network';
const BTNG_ROOT_MEMBER = 'btng-root-member';
const NODE_ID = 'nd-6HRNJ6OUIBGP3MV74YAW53NWYQ';

async function makeRequest(url, options = {}) {
  try {
    const response = await fetch(url, {
      ...options,
      signal: AbortSignal.timeout(20000),
    });
    const data = await response.text();
    let json;
    try {
      json = JSON.parse(data);
    } catch (e) {
      json = data;
    }
    return { status: response.status, data: json };
  } catch (error) {
    return { status: 'ERROR', error: error.message };
  }
}

async function testEndpoint(name, url, expectedStatus = 200) {
  console.log(`\n🔍 Testing ${name}...`);
  console.log(`   URL: ${url}`);

  const result = await makeRequest(url);

  if (result.status === expectedStatus) {
    console.log(`   ✅ Status: ${result.status} - SUCCESS`);
    if (result.data && typeof result.data === 'object') {
      console.log(`   📊 Data: ${JSON.stringify(result.data, null, 2)}`);
    } else {
      console.log(`   📄 Response: ${result.data}`);
    }
    return true;
  } else {
    console.log(`   ❌ Status: ${result.status} - FAILED`);
    if (result.error) {
      console.log(`   🚨 Error: ${result.error}`);
    }
    return false;
  }
}

async function runValidationChecklist() {
  console.log('🚀 BTNG Real-API Validation Checklist');
  console.log('=====================================');
  console.log(`Backend Endpoint: ${BTNG_API_BASE}`);
  console.log(`Next.js Proxy: ${NEXTJS_BASE}`);
  console.log('');

  let passed = 0;
  let total = 0;

  // 1. Backend Health Check
  total++;
  if (await testEndpoint('Backend Health', `${BTNG_API_BASE}/health`)) {
    passed++;
  }

  // 2. Next.js Health Proxy
  total++;
  if (await testEndpoint('Next.js Health Proxy', `${NEXTJS_BASE}/api/btng/health`)) {
    passed++;
  }

  // 3. Wallet Balance (using configured node ID)
  total++;
  if (await testEndpoint('Wallet Balance', `${NEXTJS_BASE}/api/btng/wallet/balance/${NODE_ID}`)) {
    passed++;
  }

  // 4. Wallet Transactions
  total++;
  if (await testEndpoint('Wallet Transactions', `${NEXTJS_BASE}/api/btng/wallet/transactions/${NODE_ID}`)) {
    passed++;
  }

  // 5. Explorer Block
  total++;
  if (await testEndpoint('Explorer Block', `${NEXTJS_BASE}/api/btng/explorer/block/1`)) {
    passed++;
  }

  // 6. Explorer Transaction
  total++;
  const testTxHash = '0x1234567890abcdef1234567890abcdef12345678';
  if (await testEndpoint('Explorer Transaction', `${NEXTJS_BASE}/api/btng/explorer/tx/${testTxHash}`)) {
    passed++;
  }

  // 7. Explorer Address
  total++;
  if (await testEndpoint('Explorer Address', `${NEXTJS_BASE}/api/btng/explorer/address/${NODE_ID}`)) {
    passed++;
  }

  // 8. Mining Info
  total++;
  if (await testEndpoint('Mining Info', `${NEXTJS_BASE}/api/btng/mining/info`)) {
    passed++;
  }

  // 9. Mining Hashrate
  total++;
  if (await testEndpoint('Mining Hashrate', `${NEXTJS_BASE}/api/btng/mining/hashrate`)) {
    passed++;
  }

  // 10. Oracle Price
  total++;
  if (await testEndpoint('Oracle Price', `${NEXTJS_BASE}/api/btng/oracle/price`)) {
    passed++;
  }

  // 11. Oracle Market Cap
  total++;
  if (await testEndpoint('Oracle Market Cap', `${NEXTJS_BASE}/api/btng/oracle/marketcap`)) {
    passed++;
  }

  // Summary
  console.log('\n📊 VALIDATION SUMMARY');
  console.log('====================');
  console.log(`✅ Passed: ${passed}/${total}`);
  console.log(`❌ Failed: ${total - passed}/${total}`);

  if (passed === total) {
    console.log('\n🎉 ALL TESTS PASSED! BTNG Sovereign Backend is LIVE!');
    console.log('🚀 Ready for real transactions and sovereign operations.');
  } else if (passed > 0) {
    console.log('\n⚠️  PARTIAL SUCCESS: Some endpoints working, others need attention.');
    console.log('🔧 Check your backend configuration and network connectivity.');
  } else {
    console.log('\n❌ NO CONNECTION: Backend unreachable or not configured.');
    console.log('🔧 Verify:');
    console.log('   - MTN public IP is correct');
    console.log('   - Port 54328 is forwarded');
    console.log('   - Backend is running and bound to 0.0.0.0:54328');
    console.log('   - Firewall allows inbound connections');
  }

  console.log('\n🔗 Next Steps:');
  console.log('1. If backend tests failed: Check MTN router and backend configuration');
  console.log('2. If Next.js tests failed: Restart the development server');
  console.log('3. Visit http://localhost:3001/btng-demo to test the UI');
  console.log('4. Run a real transaction to complete the validation cycle');
}

runValidationChecklist().catch(console.error);