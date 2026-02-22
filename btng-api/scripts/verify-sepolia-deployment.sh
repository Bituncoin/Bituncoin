#!/usr/bin/env bash
set -euo pipefail

# BTNG Sepolia Deployment Verification Script
# This script verifies that your BTNG contracts can be deployed to Sepolia testnet

echo "🔍 BTNG Sepolia Deployment Verification"
echo "======================================"

# ── 1️⃣  Load environment
if [ ! -f ".env" ]; then
    echo "❌ .env file not found. Please create it with your credentials."
    exit 1
fi

source .env
echo "✅ Environment loaded"

# Validate required variables
if [ -z "${INFURA_API_KEY:-}" ]; then
    echo "❌ INFURA_API_KEY not set in .env"
    exit 1
fi

if [ -z "${PRIVATE_KEY:-}" ]; then
    echo "❌ PRIVATE_KEY not set in .env"
    exit 1
fi

echo "✅ Required credentials present"

# ── 2️⃣  Check contract artifacts
ARTIFACT_DIR="./artifacts/contracts"
if [ ! -d "$ARTIFACT_DIR" ]; then
    echo "❌ Contract artifacts not found. Run 'npm run compile' first."
    exit 1
fi

echo "✅ Contract artifacts found"

# ── 3️⃣  Test deployment simulation
echo ""
echo "🚀 Testing deployment simulation..."

NODE_URL="https://sepolia.infura.io/v3/${INFURA_API_KEY}"

node - <<EOF
const { ethers } = require('ethers');
const fs = require('fs');

async function testDeployment() {
  try {
    console.log('🔗 Connecting to Sepolia...');
    const provider = new ethers.JsonRpcProvider('${NODE_URL}');
    const network = await provider.getNetwork();
    console.log('✅ Connected to network:', network.name, '(Chain ID:', network.chainId, ')');

    console.log('🔑 Testing wallet...');
    const wallet = new ethers.Wallet('${PRIVATE_KEY}', provider);
    console.log('✅ Wallet address:', wallet.address);

    const balance = await provider.getBalance(wallet.address);
    const balanceEth = ethers.formatEther(balance);
    console.log('💰 Wallet balance:', balanceEth, 'SepoliaETH');

    if (parseFloat(balanceEth) < 0.01) {
      console.log('⚠️  Low balance! Need at least 0.01 SepoliaETH for deployment');
      console.log('   Get test ETH from: https://sepoliafaucet.com');
    } else {
      console.log('✅ Sufficient balance for deployment');
    }

    // Test contract compilation
    console.log('🔨 Checking contract artifacts...');
    const artifactPath = './artifacts/contracts/BTNGGoldToken.sol/BTNGGoldToken.json';
    if (!fs.existsSync(artifactPath)) {
      throw new Error('BTNGGoldToken artifact not found');
    }

    const artifact = JSON.parse(fs.readFileSync(artifactPath, 'utf8'));
    console.log('✅ Contract artifact loaded');

    // Estimate deployment cost
    const factory = new ethers.ContractFactory(artifact.abi, artifact.bytecode, wallet);
    const deployTx = await factory.getDeployTransaction(ethers.ZeroAddress);
    const estimatedGas = await provider.estimateGas(deployTx);
    const feeData = await provider.getFeeData();
    const estimatedCost = estimatedGas * (feeData.gasPrice || 0n);
    const estimatedCostEth = ethers.formatEther(estimatedCost);

    console.log('⛽ Estimated deployment cost:', estimatedCostEth, 'SepoliaETH');

    console.log('🎯 Deployment simulation successful!');
    console.log('Ready to run: npm run deploy:testnet');

  } catch (error) {
    console.error('❌ Deployment test failed:', error.message);
    process.exit(1);
  }
}

testDeployment();
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 BTNG Sepolia deployment verification PASSED!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Ensure you have sufficient SepoliaETH (get from https://sepoliafaucet.com)"
    echo "2. Run: npm run deploy:testnet"
    echo "3. Verify contracts on https://sepolia.etherscan.io/"
else
    echo "❌ BTNG Sepolia deployment verification FAILED!"
    exit 1
fi