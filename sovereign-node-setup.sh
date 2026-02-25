#!/bin/bash
# BTNG Sovereign Node Deployment - Version: Independent
echo "🌍 Initializing Sovereign Node for the 54-Nation Reserve..."

# 1. Generate Sovereign Crypto (No Amazon CA)
./bin/cryptogen generate --config=./crypto-config.yaml

# 2. Set the Sovereign Orderer (The "Ancient" Authority)
# This ensures YOU control the transaction sequence
export ORDERER_CONTAINER=btng-orderer.sovereign.africa

# 3. Boot the Node with P2P Discovery
# We point this to your independent IP, bypassing corporate DNS
docker-compose -f docker-compose-sovereign.yaml up -d

echo "✅ Node is LIVE. Awaiting P2P handshake from other African peers."