#!/bin/bash

# BTNG Fabric Chaincode Deployment Script
# Deploys btng-gold-token to btng-fabric-network

set -e

echo "🚀 Deploying btng-gold-token to btng-fabric-network"

# Set environment variables
export FABRIC_CFG_PATH=${PWD}/config
export CORE_PEER_LOCALMSPID="btng-root-memberMSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/btng-root-member.btng-fabric-network.com/users/Admin@btng-root-member.btng-fabric-network.com/msp
export CORE_PEER_ADDRESS=peer0.btng-root-member.btng-fabric-network.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/crypto-config/peerOrganizations/btng-root-member.btng-fabric-network.com/peers/peer0.btng-root-member.btng-fabric-network.com/tls/ca.crt

# Package chaincode
echo "📦 Packaging btng-gold-token..."
peer lifecycle chaincode package btng-gold-token.tar.gz --path ./chaincode/btng-gold-token --lang golang --label btng-gold-token_1.0.0

# Install chaincode
echo "📥 Installing btng-gold-token..."
peer lifecycle chaincode install btng-gold-token.tar.gz

# Get package ID
PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep btng-gold-token | awk '{print $3}' | sed 's/,$//')

echo "📋 Package ID: ${PACKAGE_ID}"

# Approve chaincode
echo "✅ Approving btng-gold-token..."
peer lifecycle chaincode approveformyorg \
  --channelID btng-sovereign-channel \
  --name btng-gold-token \
  --version 1.0.0 \
  --package-id ${PACKAGE_ID} \
  --sequence 1 \
  --tls true \
  --cafile ${PWD}/crypto-config/ordererOrganizations/btng-fabric-network.com/orderers/orderer.btng-fabric-network.com/tls/ca.crt

# Check commit readiness
echo "🔍 Checking commit readiness..."
peer lifecycle chaincode checkcommitreadiness \
  --channelID btng-sovereign-channel \
  --name btng-gold-token \
  --version 1.0.0 \
  --sequence 1 \
  --tls true \
  --cafile ${PWD}/crypto-config/ordererOrganizations/btng-fabric-network.com/orderers/orderer.btng-fabric-network.com/tls/ca.crt

# Commit chaincode
echo "📝 Committing btng-gold-token..."
peer lifecycle chaincode commit \
  --channelID btng-sovereign-channel \
  --name btng-gold-token \
  --version 1.0.0 \
  --sequence 1 \
  --tls true \
  --cafile ${PWD}/crypto-config/ordererOrganizations/btng-fabric-network.com/orderers/orderer.btng-fabric-network.com/tls/ca.crt

echo "🎉 btng-gold-token deployment completed successfully!"
echo "🔗 Chaincode: btng-gold-token"
echo "📊 Version: 1.0.0"
echo "🔢 Sequence: 1"
echo "🌐 Network: btng-fabric-network"
echo "📺 Channel: btng-sovereign-channel"
