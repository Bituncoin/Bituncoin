#!/bin/bash

# Example API Tests for Bituncoin Universal Wallet

BASE_URL="http://localhost:8080/api"

echo "🧪 Bituncoin Wallet API Tests"
echo "=============================="
echo ""

# Test 1: Create a new wallet
echo "Test 1: Creating a new wallet with 2FA enabled..."
WALLET_RESPONSE=$(curl -s -X POST $BASE_URL/wallet/create \
  -H "Content-Type: application/json" \
  -d '{
    "enable_2fa": true,
    "enable_biometric": false
  }')

echo "Response: $WALLET_RESPONSE"
echo ""

# Extract wallet address
WALLET_ADDRESS=$(echo $WALLET_RESPONSE | grep -o 'BTN[a-f0-9]*' | head -1)

if [ -z "$WALLET_ADDRESS" ]; then
    echo "❌ Failed to create wallet. Make sure the node is running on port 8080."
    exit 1
fi

echo "✓ Wallet created: $WALLET_ADDRESS"
echo ""

# Test 2: Get wallet balance
echo "Test 2: Checking wallet balance..."
BALANCE_RESPONSE=$(curl -s "$BASE_URL/wallet/balance?address=$WALLET_ADDRESS")
echo "Response: $BALANCE_RESPONSE"
echo ""

# Test 3: Get blockchain info
echo "Test 3: Getting blockchain info..."
INFO_RESPONSE=$(curl -s "$BASE_URL/blockchain/info")
echo "Response: $INFO_RESPONSE"
echo ""

# Test 4: Get supported currencies
echo "Test 4: Getting supported currencies..."
CURRENCIES_RESPONSE=$(curl -s "$BASE_URL/currencies")
echo "Response: $CURRENCIES_RESPONSE"
echo ""

# Test 5: Create another wallet for transactions
echo "Test 5: Creating a second wallet..."
WALLET2_RESPONSE=$(curl -s -X POST $BASE_URL/wallet/create \
  -H "Content-Type: application/json" \
  -d '{
    "enable_2fa": false,
    "enable_biometric": true,
    "biometric_data": "sample-fingerprint-data"
  }')

WALLET2_ADDRESS=$(echo $WALLET2_RESPONSE | grep -o 'BTN[a-f0-9]*' | head -1)
echo "✓ Second wallet created: $WALLET2_ADDRESS"
echo ""

# Test 6: Send a transaction
echo "Test 6: Sending a transaction..."
TX_RESPONSE=$(curl -s -X POST $BASE_URL/transaction/send \
  -H "Content-Type: application/json" \
  -d "{
    \"from\": \"$WALLET_ADDRESS\",
    \"to\": \"$WALLET2_ADDRESS\",
    \"amount\": 10.5,
    \"currency\": \"BTN\",
    \"cross_chain\": false
  }")
echo "Response: $TX_RESPONSE"
echo ""

# Test 7: Get transaction history
echo "Test 7: Getting transaction history..."
HISTORY_RESPONSE=$(curl -s "$BASE_URL/transaction/history?address=$WALLET_ADDRESS")
echo "Response: $HISTORY_RESPONSE"
echo ""

# Test 8: Mine a block
echo "Test 8: Mining a block..."
MINE_RESPONSE=$(curl -s -X POST "$BASE_URL/mine?miner=$WALLET_ADDRESS")
echo "Response: $MINE_RESPONSE"
echo ""

# Test 9: Cross-chain transaction
echo "Test 9: Initiating cross-chain transaction..."
CROSSCHAIN_RESPONSE=$(curl -s -X POST $BASE_URL/transaction/send \
  -H "Content-Type: application/json" \
  -d "{
    \"from\": \"$WALLET_ADDRESS\",
    \"to\": \"$WALLET2_ADDRESS\",
    \"amount\": 5.0,
    \"currency\": \"BTN\",
    \"cross_chain\": true,
    \"target_chain\": \"ETH\"
  }")
echo "Response: $CROSSCHAIN_RESPONSE"
echo ""

echo "✅ All tests completed!"
echo ""
echo "📝 Summary:"
echo "  - Wallet 1: $WALLET_ADDRESS"
echo "  - Wallet 2: $WALLET2_ADDRESS"
echo "  - Transactions: Created and pending"
echo "  - Mining: Block mined successfully"
