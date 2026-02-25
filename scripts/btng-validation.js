#!/usr/bin/env node

/**
 * BTNG Real-API Validation Script
 * Tests all endpoints against the MTN backend at 74.118.126.72:64799
 */

const https = require('https');
const http = require('http');

const BTNG_API_BASE = 'http://74.118.126.72:64799';
const NEXTJS_BASE = 'http://localhost:3003';

function makeRequest(url, options = {}) {
  return new Promise((resolve, reject) => {
    const protocol = url.startsWith('https:') ? https : http;
    const req = protocol.get(url, options, (res) => {
      let data = '';
      res.on('data', (chunk) => data += chunk);
      res.on('end', () => {
        try {
          const json = JSON.parse(data);
          resolve({ status: res.statusCode, data: json });
        } catch (e) {
          resolve({ status: res.statusCode, data: data });
        }
      });
    });

    req.on('error', (err) => {
      resolve({ status: 'ERROR', error: err.message });
    });

    req.setTimeout(10000, () => {
      req.destroy();
      resolve({ status: 'TIMEOUT', error: 'Request timeout' });
    });
  });
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

  // 3. Wallet Balance (using a test address)
  total++;
  const testAddress = 'BTNG1TESTADDRESS123456789012345678';
  if (await testEndpoint('Wallet Balance', `${NEXTJS_BASE}/api/btng/wallet/balance/${testAddress}`)) {
    passed++;
  }

  // 4. Wallet Transactions
  total++;
  if (await testEndpoint('Wallet Transactions', `${NEXTJS_BASE}/api/btng/wallet/transactions/${testAddress}`)) {
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
  if (await testEndpoint('Explorer Address', `${NEXTJS_BASE}/api/btng/explorer/address/${testAddress}`)) {
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

  // 12. Fabric Network Status
  total++;
  if (await testEndpoint('Fabric Network Status', `${NEXTJS_BASE}/api/btng/fabric/network`)) {
    passed++;
  }

  // 13. Fabric Node Status
  total++;
  if (await testEndpoint('Fabric Node Status', `${NEXTJS_BASE}/api/btng/fabric/node`)) {
    passed++;
  }

  // 14. Fabric Network Operation (requires auth - will fail without token)
  total++;
  if (await testEndpoint('Fabric Network Operation (No Auth)', `${NEXTJS_BASE}/api/btng/fabric/network`, 401)) {
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
    console.log('   - Port 64799 is forwarded');
    console.log('   - Backend is running and bound to 0.0.0.0:64799');
    console.log('   - Firewall allows inbound connections');
  }

  console.log('\n🔗 Next Steps:');
  console.log('1. If backend tests failed: Check MTN router and backend configuration');
  console.log('2. If Next.js tests failed: Restart the development server');
  console.log('3. If Fabric tests failed: Verify fabric network configuration');
  console.log('4. Visit http://localhost:3001/btng-demo to test the UI');
  console.log('5. Run a real transaction to complete the validation cycle');
  console.log('6. Deploy chaincode to btng-fabric-network for full sovereignty');
}

runValidationChecklist().catch(console.error);