#!/usr/bin/env node

/**
 * BTNG Testnet Environment Checker
 * Verifies that all required environment variables are set for testnet deployment
 */

import dotenv from 'dotenv';

dotenv.config();

const requiredVars = [
  'INFURA_API_KEY',
  'PRIVATE_KEY'
];

const optionalVars = [
  'ETHERSCAN_API_KEY'
];

const configVars = [
  'NETWORK',
  'CHAIN_ID'
];

console.log('🔍 BTNG Testnet Environment Check\n');

let allRequired = true;

console.log('📋 Required Variables:');
requiredVars.forEach(varName => {
  const value = process.env[varName];
  if (value && value !== 'your_' + varName.toLowerCase() + '_here' && value !== 'demo_' + varName.toLowerCase().replace(/_/g, '_')) {
    console.log(`✅ ${varName}: Set`);
  } else if (value && value.startsWith('demo_')) {
    console.log(`⚠️  ${varName}: Demo value (replace for real deployment)`);
  } else {
    console.log(`❌ ${varName}: Not set or using placeholder`);
    allRequired = false;
  }
});

console.log('\n📋 Optional Variables:');
optionalVars.forEach(varName => {
  const value = process.env[varName];
  if (value && value !== 'your_' + varName.toLowerCase() + '_here' && !value.startsWith('demo_')) {
    console.log(`✅ ${varName}: Set`);
  } else {
    console.log(`⚠️  ${varName}: Not set or demo value (optional)`);
  }
});

console.log('\n📋 Configuration Variables:');
configVars.forEach(varName => {
  const value = process.env[varName];
  const expectedValue = varName === 'NETWORK' ? 'sepolia' : varName === 'CHAIN_ID' ? '11155111' : '';
  if (value === expectedValue) {
    console.log(`✅ ${varName}: ${value}`);
  } else {
    console.log(`⚠️  ${varName}: ${value || 'Not set'} (expected: ${expectedValue})`);
  }
});

console.log('\n🎯 Deployment Status:');
if (allRequired) {
  console.log('✅ Ready for testnet deployment!');
  console.log('⚠️  Note: Using demo credentials - replace with real values for production');
  console.log('Run: npm run deploy:testnet');
} else {
  console.log('❌ Missing required environment variables');
  console.log('Please set up your .env file with the required credentials');
  console.log('See: TESTNET_DEPLOYMENT_GUIDE.md');
}

process.exit(allRequired ? 0 : 1);