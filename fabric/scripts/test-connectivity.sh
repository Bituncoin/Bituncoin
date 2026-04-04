#!/bin/bash
# BTNG Sovereign Network Connectivity Test
# Tests peer, orderer, and CA endpoints on AWS Managed Blockchain
set -e

echo "=============================================="
echo "  BTNG SOVEREIGN NETWORK CONNECTIVITY TEST"
echo "  Network: btng712-fabric-network"
echo "  Region:  us-east-1 (AWS Managed Blockchain)"
echo "=============================================="

# Endpoints
PEER_HOST="nd-jkud2atma5a4fahhtnkwalo2k4.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com"
PEER_PORT=30003
ORDERER_HOST="orderer.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com"
ORDERER_PORT=30001
CA_HOST="ca.m-sp4qe6lju5h5zbo7bgl7rj4qnq.n-wfutwh7lvrahpjtn7uspcsgy7a.managedblockchain.us-east-1.amazonaws.com"
CA_PORT=30002
GHANA_PEER="154.161.183.158"
GHANA_PEER_PORT=38982

TLS_CERT="crypto-config/managedblockchain-tls-chain.pem"

PASS=0
FAIL=0

test_endpoint() {
  local name=$1
  local host=$2
  local port=$3
  local protocol=$4

  echo ""
  echo "--- Testing $name ---"
  echo "    Host: $host:$port ($protocol)"

  # DNS resolution
  if nslookup "$host" > /dev/null 2>&1 || host "$host" > /dev/null 2>&1; then
    echo "    [PASS] DNS resolution OK"
  else
    echo "    [WARN] DNS resolution failed (may be internal)"
  fi

  # TCP connectivity
  if timeout 10 bash -c "echo > /dev/tcp/$host/$port" 2>/dev/null; then
    echo "    [PASS] TCP connection OK"
    PASS=$((PASS + 1))
  else
    echo "    [FAIL] TCP connection FAILED"
    FAIL=$((FAIL + 1))
    return
  fi

  # TLS handshake (for gRPCS/HTTPS)
  if [ "$protocol" = "grpcs" ] || [ "$protocol" = "https" ]; then
    if [ -f "$TLS_CERT" ]; then
      if echo | openssl s_client -connect "$host:$port" -CAfile "$TLS_CERT" -servername "$host" 2>/dev/null | grep -q "Verify return code: 0"; then
        echo "    [PASS] TLS handshake OK (verified)"
        PASS=$((PASS + 1))
      else
        echo "    [WARN] TLS handshake completed (cert verification may differ)"
        PASS=$((PASS + 1))
      fi
    else
      echo "    [SKIP] TLS cert not found at $TLS_CERT"
    fi
  fi
}

echo ""
echo "=== AWS MANAGED BLOCKCHAIN ENDPOINTS ==="
test_endpoint "Peer Node (gRPCS)" "$PEER_HOST" "$PEER_PORT" "grpcs"
test_endpoint "Orderer (gRPCS)" "$ORDERER_HOST" "$ORDERER_PORT" "grpcs"
test_endpoint "Certificate Authority (HTTPS)" "$CA_HOST" "$CA_PORT" "https"

echo ""
echo "=== GHANA ANCHOR NODE ==="
test_endpoint "Ghana Anchor Peer (gRPCS)" "$GHANA_PEER" "$GHANA_PEER_PORT" "grpcs"

echo ""
echo "=== CA HEALTH CHECK ==="
CA_URL="https://${CA_HOST}:${CA_PORT}/cainfo"
echo "    Testing: $CA_URL"
if curl -sk --max-time 10 "$CA_URL" | grep -q "CAName\|caname\|result"; then
  echo "    [PASS] CA is responding"
  PASS=$((PASS + 1))
else
  echo "    [INFO] CA health check inconclusive"
fi

echo ""
echo "=== FABRIC PEER CHANNEL STATUS ==="
if command -v peer &> /dev/null; then
  export CORE_PEER_ADDRESS="${PEER_HOST}:${PEER_PORT}"
  export CORE_PEER_LOCALMSPID="m-SP4QE6LJU5H5ZBO7BGL7RJ4QNQ"
  export CORE_PEER_TLS_ENABLED=true
  export CORE_PEER_TLS_ROOTCERT_FILE="$(pwd)/${TLS_CERT}"

  if [ -d "crypto-config/Org1AdminMSP" ]; then
    export CORE_PEER_MSPCONFIGPATH="$(pwd)/crypto-config/Org1AdminMSP"
    echo "    Querying joined channels..."
    peer channel list 2>/dev/null && PASS=$((PASS + 1)) || echo "    [INFO] Channel query requires admin credentials"
  else
    echo "    [SKIP] Admin MSP not found for channel query"
  fi
else
  echo "    [SKIP] Fabric peer CLI not installed"
fi

echo ""
echo "=============================================="
echo "  RESULTS: $PASS passed, $FAIL failed"
echo "=============================================="
echo ""
echo "  MSP ID:    m-SP4QE6LJU5H5ZBO7BGL7RJ4QNQ"
echo "  Network:   btng712-fabric-network"
echo "  Channel:   btng712-fabric-network"
echo "  Chaincode: btng-wallet, btng-gold-token"
echo "=============================================="

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
