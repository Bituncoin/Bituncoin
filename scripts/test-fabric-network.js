#!/usr/bin/env node

/**
 * BTNG Fabric Network Test
 * Tests fabric network endpoints with JWT authentication
 */

import fetch from 'node-fetch';

const BASE_URL = process.env.BASE_URL || 'http://localhost:3002';
const JWT_SECRET = 'a-string-secret-at-least-256-bits-long';

async function getAuthToken() {
  console.log('🔐 Getting JWT token for fabric operations...');

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
  return loginData.token;
}

async function testFabricNetwork() {
  console.log('🏗️  Testing BTNG Fabric Network Integration...\n');

  try {
    // Test 1: Public fabric network status
    console.log('1️⃣ Testing public fabric network status...');
    const networkResponse = await fetch(`${BASE_URL}/api/btng/fabric/network`);

    if (networkResponse.ok) {
      const networkData = await networkResponse.json();
      console.log('✅ Fabric network status accessible!');
      console.log('📊 Network:', networkData.network.name);
      console.log('🔗 Node ID:', networkData.network.nodeId);
      console.log('👑 Root Member:', networkData.network.rootMember);
    } else {
      console.log('❌ Network status failed:', networkResponse.status);
    }
    console.log();

    // Test 2: Public fabric node status
    console.log('2️⃣ Testing public fabric node status...');
    const nodeResponse = await fetch(`${BASE_URL}/api/btng/fabric/node`);

    if (nodeResponse.ok) {
      const nodeData = await nodeResponse.json();
      console.log('✅ Fabric node status accessible!');
      console.log('🆔 Node ID:', nodeData.nodeId);
      console.log('🏢 Organization:', nodeData.organization);
      console.log('⚡ Status:', nodeData.status);
      console.log('🛡️  Capabilities:', nodeData.capabilities.join(', '));
    } else {
      console.log('❌ Node status failed:', nodeResponse.status);
    }
    console.log();

    // Test 3: Get authentication token
    const token = await getAuthToken();
    console.log('✅ Authentication successful!');
    console.log('🎫 Token acquired for fabric operations\n');

    // Test 4: Deploy chaincode operation
    console.log('3️⃣ Testing fabric chaincode deployment...');
    const deployResponse = await fetch(`${BASE_URL}/api/btng/fabric/network`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        action: 'deploy',
        channel: 'btng-sovereign-channel',
        chaincode: 'btng-gold-token',
        transaction: {
          version: '1.0.0',
          sequence: '1',
          initRequired: true
        }
      })
    });

    if (deployResponse.ok) {
      const deployResult = await deployResponse.json();
      console.log('✅ Chaincode deployment successful!');
      console.log('📋 Transaction ID:', deployResult.operation.result.txId);
      console.log('🔗 Channel:', deployResult.operation.channel);
      console.log('⚙️  Chaincode:', deployResult.operation.chaincode);
    } else {
      console.log('❌ Deployment failed:', deployResponse.status);
      const error = await deployResponse.json();
      console.log('Error:', error.error);
    }
    console.log();

    // Test 5: Test comprehensive chaincode operations
    console.log('4️⃣ Testing fabric chaincode operations...');

    // Test Gold Token operations
    const goldTokenTests = [
      { function: 'Mint', args: ['1000', 'BTNG-SOVEREIGN-001'] },
      { function: 'Transfer', args: ['500', 'BTNG-SOVEREIGN-001', 'BTNG-MEMBER-002'] },
      { function: 'BalanceOf', args: ['BTNG-SOVEREIGN-001'] },
      { function: 'TotalSupply', args: [] }
    ];

    for (const test of goldTokenTests) {
      const goldResponse = await fetch(`${BASE_URL}/api/btng/fabric/chaincode`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          chaincode: 'btng-gold-token',
          function: test.function,
          args: test.args
        })
      });

      if (goldResponse.ok) {
        const goldResult = await goldResponse.json();
        console.log(`✅ Gold Token ${test.function}: ${goldResult.operation.result.message}`);
      } else {
        console.log(`❌ Gold Token ${test.function} failed:`, goldResponse.status);
      }
    }

    // Test Identity operations
    const identityTests = [
      { function: 'RegisterIdentity', args: ['BTNG-IDENTITY-001', 'public-key-123', '{"type":"sovereign","country":"BTNG"}'] },
      { function: 'VerifyIdentity', args: ['BTNG-IDENTITY-001'] },
      { function: 'GetIdentity', args: ['BTNG-IDENTITY-001'] }
    ];

    for (const test of identityTests) {
      const identityResponse = await fetch(`${BASE_URL}/api/btng/fabric/chaincode`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          chaincode: 'btng-sovereign-identity',
          function: test.function,
          args: test.args
        })
      });

      if (identityResponse.ok) {
        const identityResult = await identityResponse.json();
        console.log(`✅ Identity ${test.function}: ${identityResult.operation.result.message}`);
      } else {
        console.log(`❌ Identity ${test.function} failed:`, identityResponse.status);
      }
    }

    console.log('\n🎉 BTNG Fabric Network test complete!');
    console.log('🏗️  Network: btng-fabric-network');
    console.log('👑 Root Member: btng-root-member');
    console.log('🆔 Node ID: nd-6HRNJ6OUIBGP3MV74YAW53NWYQ');

  } catch (error) {
    console.error('❌ Fabric network test failed:', error.message);
    console.log('\n💡 Make sure your Next.js server is running and JWT auth is configured');
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  testFabricNetwork();
}