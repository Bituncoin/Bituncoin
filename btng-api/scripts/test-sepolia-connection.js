#!/usr/bin/env node

/**
 * BTNG Testnet Connection Test
 * Tests if your credentials can connect to Sepolia testnet
 */

import dotenv from 'dotenv';
import { ethers } from 'ethers';

dotenv.config();

async function testConnection() {
  console.log('🔗 Testing BTNG Sepolia Testnet Connection...\n');

  const infuraKey = process.env.INFURA_API_KEY;
  const privateKey = process.env.PRIVATE_KEY;
  const expectedNetwork = process.env.NETWORK;
  const expectedChainId = process.env.CHAIN_ID;

  if (!infuraKey || !privateKey) {
    console.log('❌ Missing credentials in .env file');
    console.log('Please set INFURA_API_KEY and PRIVATE_KEY');
    process.exit(1);
  }

  try {
    // Create provider
    const provider = new ethers.JsonRpcProvider(`https://sepolia.infura.io/v3/${infuraKey}`);
    console.log('✅ Connected to Sepolia via Infura');

    // Create wallet
    const wallet = new ethers.Wallet(privateKey, provider);
    console.log('✅ Wallet created:', wallet.address);

    // Check balance
    const balance = await provider.getBalance(wallet.address);
    const balanceEth = ethers.formatEther(balance);
    console.log('💰 Wallet balance:', balanceEth, 'SepoliaETH');

    if (parseFloat(balanceEth) < 0.01) {
      console.log('⚠️  Low balance! Get test ETH from: https://sepoliafaucet.com');
      console.log('   Need at least 0.01 SepoliaETH for deployment');
    } else {
      console.log('✅ Sufficient balance for deployment');
    }

    // Test network
    const network = await provider.getNetwork();
    console.log('🌐 Network:', network.name, '(Chain ID:', network.chainId, ')');

    // Verify network configuration
    if (expectedNetwork && network.name !== expectedNetwork) {
      console.log(`⚠️  Network mismatch: got ${network.name}, expected ${expectedNetwork}`);
    }
    if (expectedChainId && network.chainId.toString() !== expectedChainId) {
      console.log(`⚠️  Chain ID mismatch: got ${network.chainId}, expected ${expectedChainId}`);
    }

    // Test gas price
    const feeData = await provider.getFeeData();
    console.log('⛽ Gas price:', ethers.formatUnits(feeData.gasPrice || 0, 'gwei'), 'gwei');

    console.log('\n🎯 Connection test successful!');
    console.log('Ready to deploy BTNG to Sepolia testnet 🚀');

  } catch (error) {
    console.log('❌ Connection test failed:', error.message);

    if (error.message.includes('project id') || error.message.includes('api key')) {
      console.log('💡 Check your INFURA_API_KEY');
    } else if (error.message.includes('private key')) {
      console.log('💡 Check your PRIVATE_KEY (should not have 0x prefix)');
    } else {
      console.log('💡 Check your internet connection and try again');
    }

    process.exit(1);
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  testConnection();
}