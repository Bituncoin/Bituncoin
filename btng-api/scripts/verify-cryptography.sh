#!/usr/bin/env bash
set -euo pipefail

# BTNG Cryptographic Verification Script
# Tests signature verification for sovereign transactions

echo "🔐 BTNG Cryptographic Verification"
echo "=================================="

# Check for required tools
if ! command -v openssl &> /dev/null; then
    echo "❌ openssl not found. Please install OpenSSL."
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "❌ node not found. Please install Node.js."
    exit 1
fi

echo "✅ Required tools available"

# Load environment
if [ ! -f ".env" ]; then
    echo "❌ .env file not found"
    exit 1
fi

source .env

if [ -z "${PRIVATE_KEY:-}" ]; then
    echo "❌ PRIVATE_KEY not set in .env"
    exit 1
fi

echo "✅ Environment loaded"

# Create cryptographic verification tests
echo ""
echo "🔧 Testing cryptographic verification..."

node - <<EOF
const crypto = require('crypto');
const { ethers } = require('ethers');

function testCryptographicVerification() {
  console.log('🧪 Testing BTNG cryptographic verification...');

  // Test 1: Ethereum signature verification
  console.log('\\n1️⃣ Testing Ethereum signature verification...');

  const message = 'BTNG sovereign transaction: mint 1000 tokens';
  const privateKey = '${PRIVATE_KEY}';

  try {
    // Create wallet and sign message
    const wallet = new ethers.Wallet(privateKey);
    const signature = await wallet.signMessage(message);

    console.log('✅ Message signed with private key');

    // Verify signature
    const recoveredAddress = ethers.verifyMessage(message, signature);
    const expectedAddress = wallet.address;

    if (recoveredAddress.toLowerCase() === expectedAddress.toLowerCase()) {
      console.log('✅ Ethereum signature verified');
      console.log('Signer address:', recoveredAddress);
    } else {
      console.log('❌ Ethereum signature verification failed');
      return false;
    }
  } catch (error) {
    console.log('❌ Ethereum signing failed:', error.message);
    return false;
  }

  // Test 2: Hash verification
  console.log('\\n2️⃣ Testing cryptographic hash functions...');

  const testData = 'BTNG sovereign gold standard';
  const hash256 = crypto.createHash('sha256').update(testData).digest('hex');
  const hash512 = crypto.createHash('sha512').update(testData).digest('hex');

  console.log('✅ SHA-256 hash:', hash256.substring(0, 16) + '...');
  console.log('✅ SHA-512 hash:', hash512.substring(0, 16) + '...');

  // Verify hash consistency
  const hash256Again = crypto.createHash('sha256').update(testData).digest('hex');
  if (hash256 === hash256Again) {
    console.log('✅ Hash function consistent');
  } else {
    console.log('❌ Hash function inconsistent');
    return false;
  }

  // Test 3: HMAC verification
  console.log('\\n3️⃣ Testing HMAC for sovereign secrets...');

  const secret = 'btng-sovereign-master-key';
  const hmac = crypto.createHmac('sha256', secret).update(testData).digest('hex');

  console.log('✅ HMAC generated:', hmac.substring(0, 16) + '...');

  // Verify HMAC consistency
  const hmacAgain = crypto.createHmac('sha256', secret).update(testData).digest('hex');
  if (hmac === hmacAgain) {
    console.log('✅ HMAC function consistent');
  } else {
    console.log('❌ HMAC function inconsistent');
    return false;
  }

  // Test 4: Transaction signature simulation
  console.log('\\n4️⃣ Testing BTNG transaction signature simulation...');

  const transactionData = {
    type: 'mint',
    amount: '1000',
    currency: 'BTNG',
    recipient: '0x742d35Cc6634C0532925a3b844Bc454e4438f44e',
    timestamp: Date.now(),
    nonce: Math.floor(Math.random() * 1000000)
  };

  const txString = JSON.stringify(transactionData, Object.keys(transactionData).sort());
  const txHash = crypto.createHash('sha256').update(txString).digest('hex');

  console.log('✅ Transaction data hashed');
  console.log('TX hash:', txHash);

  // Sign transaction hash
  try {
    const wallet = new ethers.Wallet(privateKey);
    const txSignature = await wallet.signMessage(txHash);

    console.log('✅ Transaction hash signed');

    // Verify transaction signature
    const recoveredSigner = ethers.verifyMessage(txHash, txSignature);

    if (recoveredSigner.toLowerCase() === wallet.address.toLowerCase()) {
      console.log('✅ Transaction signature verified');
      console.log('Signer:', recoveredSigner);
    } else {
      console.log('❌ Transaction signature verification failed');
      return false;
    }

    // Test 5: Sovereign certificate simulation
    console.log('\\n5️⃣ Testing sovereign certificate validation...');

    const certificateData = {
      issuer: 'BTNG Sovereign Authority',
      subject: 'Country Entity - Kenya',
      publicKey: wallet.publicKey,
      validFrom: Date.now(),
      validTo: Date.now() + (365 * 24 * 60 * 60 * 1000), // 1 year
      permissions: ['mint-tokens', 'update-oracle', 'transfer-reserves']
    };

    const certString = JSON.stringify(certificateData, Object.keys(certificateData).sort());
    const certHash = crypto.createHash('sha256').update(certString).digest('hex');
    const certSignature = await wallet.signMessage(certHash);

    console.log('✅ Sovereign certificate created and signed');
    console.log('Certificate hash:', certHash.substring(0, 16) + '...');

    // Verify certificate signature
    const certSigner = ethers.verifyMessage(certHash, certSignature);

    if (certSigner.toLowerCase() === wallet.address.toLowerCase()) {
      console.log('✅ Sovereign certificate signature verified');
      console.log('Authority:', certSigner);
      console.log('Permissions:', certificateData.permissions.join(', '));
    } else {
      console.log('❌ Sovereign certificate verification failed');
      return false;
    }

  } catch (error) {
    console.log('❌ Transaction signing failed:', error.message);
    return false;
  }

  console.log('\\n🎉 All cryptographic verification tests PASSED!');
  return true;
}

async function runTests() {
  const success = await testCryptographicVerification();
  process.exit(success ? 0 : 1);
}

runTests().catch(error => {
  console.error('❌ Test execution failed:', error);
  process.exit(1);
});
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 BTNG cryptographic verification PASSED!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Implement production cryptographic services"
    echo "2. Set up hardware security modules (HSM) for key management"
    echo "3. Deploy signature verification services"
else
    echo "❌ BTNG cryptographic verification FAILED!"
    exit 1
fi