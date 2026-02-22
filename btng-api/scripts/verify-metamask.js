#!/usr/bin/env node

/**
 * BTNG MetaMask Sepolia Network Verifier
 * Helps verify MetaMask is connected to Sepolia without handling private keys
 */

import readline from 'readline';

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

function ask(question) {
  return new Promise(resolve => rl.question(question, resolve));
}

async function verifyMetaMaskSetup() {
  console.log('🦊 BTNG MetaMask Sepolia Network Verifier');
  console.log('==========================================\n');

  console.log('This tool helps you verify your MetaMask setup for BTNG deployment.\n');

  console.log('📋 MetaMask Sepolia Network Configuration:');
  console.log('------------------------------------------');
  console.log('Network Name: Sepolia');
  console.log('RPC URL: https://sepolia.infura.io/v3/YOUR_INFURA_API_KEY');
  console.log('Chain ID: 11155111');
  console.log('Currency Symbol: SepoliaETH');
  console.log('Block Explorer: https://sepolia.etherscan.io/\n');

  console.log('✅ To verify your MetaMask setup:');
  console.log('1. Open MetaMask extension');
  console.log('2. Click the network dropdown (should show "Sepolia")');
  console.log('3. Verify the network details match above');
  console.log('4. Check your account balance (need ~0.01 SepoliaETH)');
  console.log('5. Copy your account address for reference\n');

  const hasSepolia = await ask('Do you see "Sepolia" in your MetaMask network list? (y/n): ');
  if (hasSepolia.toLowerCase() !== 'y') {
    console.log('\n❌ Sepolia network not found in MetaMask.');
    console.log('Please add Sepolia network:');
    console.log('1. Click network dropdown → "Add Network"');
    console.log('2. Use the configuration shown above');
    console.log('3. Replace YOUR_INFURA_API_KEY with your actual Infura key');
    return;
  }

  console.log('\n✅ Sepolia network detected!');

  const isConnected = await ask('Is MetaMask connected to Sepolia? (y/n): ');
  if (isConnected.toLowerCase() !== 'y') {
    console.log('\n⚠️ Please switch to Sepolia network in MetaMask.');
    return;
  }

  console.log('\n✅ MetaMask connected to Sepolia!');

  const accountAddress = await ask('What is your MetaMask account address? (0x...): ');
  if (!accountAddress.startsWith('0x') || accountAddress.length !== 42) {
    console.log('\n❌ Invalid Ethereum address format.');
    return;
  }

  console.log(`\n✅ Valid Ethereum address: ${accountAddress}`);

  const hasBalance = await ask('Do you have SepoliaETH in this account? (y/n): ');
  if (hasBalance.toLowerCase() !== 'y') {
    console.log('\n💰 You need SepoliaETH for deployment.');
    console.log('Get free SepoliaETH from: https://sepoliafaucet.com');
    console.log('Recommended: 0.01 SepoliaETH minimum');
    return;
  }

  const balance = await ask('Approximately how much SepoliaETH do you have? (e.g., 0.5): ');
  const balanceNum = parseFloat(balance);

  if (balanceNum < 0.01) {
    console.log('\n⚠️ Low balance detected.');
    console.log('You need at least 0.01 SepoliaETH for deployment.');
    console.log('Get more from: https://sepoliafaucet.com');
  } else {
    console.log(`\n✅ Sufficient balance: ${balanceNum} SepoliaETH`);
  }

  console.log('\n🎉 MetaMask Sepolia setup verified!');
  console.log('\n📋 Next Steps:');
  console.log('1. Update your .env file with real credentials:');
  console.log('   node scripts/update-env-cli.js [INFURA_KEY] [PRIVATE_KEY]');
  console.log('2. Run environment check: npm run check-env');
  console.log('3. Test connection: npm run test-connection');
  console.log('4. Deploy Oracle: npm run deploy:oracle');

  rl.close();
}

if (import.meta.url === `file://${process.argv[1]}`) {
  verifyMetaMaskSetup().catch(console.error);
}