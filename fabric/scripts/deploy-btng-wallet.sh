#!/bin/bash
# BTNG Wallet Chaincode Deployment to AWS Managed Blockchain
# Network: btng712-fabric-network
# Channel: btng712-fabric-network
set -e

echo "=============================================="
echo "  DEPLOYING btng-wallet CHAINCODE"
echo "  Network: btng712-fabric-network"
echo "  AWS Managed Blockchain (us-east-1)"
echo "=============================================="

# AWS Managed Blockchain environment
NETWORK_ID="n-wfutwh7lvrahpjtn7uspcsgy7a"
MEMBER_ID="m-SP4QE6LJU5H5ZBO7BGL7RJ4QNQ"
PEER_ENDPOINT="nd-jkud2atma5a4fahhtnkwalo2k4.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30003"
ORDERER_ENDPOINT="orderer.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30001"
CA_ENDPOINT="ca.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com:30002"

CHANNEL_NAME="btng712-fabric-network"
CHAINCODE_NAME="btng-wallet"
CHAINCODE_VERSION="1.0.0"
CHAINCODE_PATH="./chaincode/btng-wallet"
CHAINCODE_LANG="golang"
SEQUENCE=1

TLS_CERT="$(pwd)/../crypto-config/managedblockchain-tls-chain.pem"
MSP_PATH="$(pwd)/../crypto-config/Org1AdminMSP"

# Set Fabric peer environment
export CORE_PEER_ADDRESS="${PEER_ENDPOINT}"
export CORE_PEER_LOCALMSPID="${MEMBER_ID}"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE="${TLS_CERT}"
export CORE_PEER_MSPCONFIGPATH="${MSP_PATH}"

echo ""
echo "Step 1/5: Package chaincode"
peer lifecycle chaincode package ${CHAINCODE_NAME}.tar.gz \
  --path ${CHAINCODE_PATH} \
  --lang ${CHAINCODE_LANG} \
  --label ${CHAINCODE_NAME}_${CHAINCODE_VERSION}

echo "[OK] Packaged: ${CHAINCODE_NAME}.tar.gz"

echo ""
echo "Step 2/5: Install chaincode on peer"
peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz

# Get package ID
PACKAGE_ID=$(peer lifecycle chaincode queryinstalled 2>&1 | grep "${CHAINCODE_NAME}_${CHAINCODE_VERSION}" | awk '{print $3}' | sed 's/,$//')
echo "[OK] Installed. Package ID: ${PACKAGE_ID}"

echo ""
echo "Step 3/5: Approve chaincode for org"
peer lifecycle chaincode approveformyorg \
  --channelID ${CHANNEL_NAME} \
  --name ${CHAINCODE_NAME} \
  --version ${CHAINCODE_VERSION} \
  --package-id ${PACKAGE_ID} \
  --sequence ${SEQUENCE} \
  --tls true \
  --cafile ${TLS_CERT} \
  --orderer ${ORDERER_ENDPOINT}

echo "[OK] Approved for ${MEMBER_ID}"

echo ""
echo "Step 4/5: Check commit readiness"
peer lifecycle chaincode checkcommitreadiness \
  --channelID ${CHANNEL_NAME} \
  --name ${CHAINCODE_NAME} \
  --version ${CHAINCODE_VERSION} \
  --sequence ${SEQUENCE} \
  --tls true \
  --cafile ${TLS_CERT} \
  --orderer ${ORDERER_ENDPOINT}

echo ""
echo "Step 5/5: Commit chaincode"
peer lifecycle chaincode commit \
  --channelID ${CHANNEL_NAME} \
  --name ${CHAINCODE_NAME} \
  --version ${CHAINCODE_VERSION} \
  --sequence ${SEQUENCE} \
  --tls true \
  --cafile ${TLS_CERT} \
  --orderer ${ORDERER_ENDPOINT} \
  --peerAddresses ${PEER_ENDPOINT} \
  --tlsRootCertFiles ${TLS_CERT}

echo ""
echo "=============================================="
echo "  DEPLOYMENT COMPLETE"
echo "  Chaincode: ${CHAINCODE_NAME}"
echo "  Version:   ${CHAINCODE_VERSION}"
echo "  Channel:   ${CHANNEL_NAME}"
echo "  Network:   btng712-fabric-network"
echo "  Functions: Mint, Melt, Transfer, GetBalance,"
echo "             GetTotalSupply, GetTransactionHistory"
echo "=============================================="

# Verify with a test query
echo ""
echo "Verifying deployment with test query..."
peer chaincode query \
  -C ${CHANNEL_NAME} \
  -n ${CHAINCODE_NAME} \
  -c '{"function":"GetTotalSupply","Args":[]}' \
  --tls true \
  --cafile ${TLS_CERT}

echo "[OK] btng-wallet chaincode is LIVE on BTNG Sovereign Network"
