#!/usr/bin/env node

/**
 * Test Infura Connection using ethers
 */

import { ethers } from 'ethers';
import dotenv from 'dotenv';

dotenv.config();

async function testInfura() {
  const apiKey = process.env.INFURA_API_KEY;
  const privateKey = process.env.PRIVATE_KEY;

  console.log('🔗 Testing Infura Connection...');
  console.log('API Key:', apiKey ? `${apiKey.substring(0, 8)}...` : 'Not set');
  console.log('Private Key:', privateKey ? `${privateKey.substring(0, 8)}...` : 'Not set');

  if (!apiKey || !privateKey) {
    console.log('❌ Missing environment variables');
    return;
  }

  try {
    // Create provider - using public Sepolia RPC
    const provider = new ethers.JsonRpcProvider("https://ethereum-sepolia.publicnode.com");
    console.log('✅ Provider created (using public Sepolia RPC)');

    // Test connection
    const network = await provider.getNetwork();
    console.log('✅ Network:', network.name, 'Chain ID:', network.chainId);

    // Test wallet
    const wallet = new ethers.Wallet(privateKey, provider);
    console.log('✅ Wallet address:', wallet.address);

    // Check balance
    console.log('🔍 Checking balance...');
    const balance = await provider.getBalance(wallet.address);
    console.log('✅ Balance:', ethers.formatEther(balance), 'ETH');

    if (balance < ethers.parseEther('0.01')) {
      console.log('⚠️ Low balance - need at least 0.01 ETH for deployment');
      console.log('💰 Get Sepolia ETH from: https://sepoliafaucet.com');
    } else {
      console.log('✅ Sufficient balance for deployment');
    }

  } catch (error) {
    console.log('❌ Connection failed:', error.message);
    if (error.message.includes('project id')) {
      console.log('💡 This usually means:');
      console.log('   1. Infura API key is invalid');
      console.log('   2. Infura project not configured for Sepolia');
      console.log('   3. Infura account has issues');
    }
  }
}

testInfura().catch(console.error);