#!/bin/bash

# BTNG Fabric Channel Creation Script
# Creates the btng-sovereign-channel for the BTNG network

set -e

FABRIC_NETWORK="btng-fabric-network"
CHANNEL_NAME="btng-sovereign-channel"
ROOT_MEMBER="btng-root-member"

echo "🏗️  Creating ${CHANNEL_NAME} for ${FABRIC_NETWORK}"

# Set environment variables
export FABRIC_CFG_PATH=${PWD}/fabric
export CORE_PEER_LOCALMSPID="${ROOT_MEMBER}MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/crypto-config/peerOrganizations/${ROOT_MEMBER}.${FABRIC_NETWORK}.com/users/Admin@${ROOT_MEMBER}.${FABRIC_NETWORK}.com/msp
export CORE_PEER_ADDRESS=peer0.${ROOT_MEMBER}.${FABRIC_NETWORK}.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/crypto-config/peerOrganizations/${ROOT_MEMBER}.${FABRIC_NETWORK}.com/peers/peer0.${ROOT_MEMBER}.${FABRIC_NETWORK}.com/tls/ca.crt

# Create channel transaction
echo "📝 Creating channel transaction..."
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID ${CHANNEL_NAME}

# Create anchor peer transaction
echo "⚓ Creating anchor peer transaction..."
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/${ROOT_MEMBER}MSPanchors.tx -channelID ${CHANNEL_NAME} -asOrg ${ROOT_MEMBER}

# Join peer to channel
echo "🤝 Joining peer to channel..."
peer channel create -o orderer.${FABRIC_NETWORK}.com:7050 -c ${CHANNEL_NAME} -f ./channel-artifacts/${CHANNEL_NAME}.tx --tls true --cafile ${PWD}/crypto-config/ordererOrganizations/${FABRIC_NETWORK}.com/orderers/orderer.${FABRIC_NETWORK}.com/tls/ca.crt

peer channel join -b ${CHANNEL_NAME}.block

# Update anchor peer
echo "🔗 Updating anchor peer..."
peer channel update -o orderer.${FABRIC_NETWORK}.com:7050 -c ${CHANNEL_NAME} -f ./channel-artifacts/${ROOT_MEMBER}MSPanchors.tx --tls true --cafile ${PWD}/crypto-config/ordererOrganizations/${FABRIC_NETWORK}.com/orderers/orderer.${FABRIC_NETWORK}.com/tls/ca.crt

echo "🎉 Channel ${CHANNEL_NAME} created successfully!"
echo "🌐 Network: ${FABRIC_NETWORK}"
echo "📺 Channel: ${CHANNEL_NAME}"
echo "👑 Root Member: ${ROOT_MEMBER}"