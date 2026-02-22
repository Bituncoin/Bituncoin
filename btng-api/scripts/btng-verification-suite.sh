#!/usr/bin/env bash
set -euo pipefail

# BTNG Sovereign Platform - Complete Verification Suite
# Runs all verification checks for Sepolia deployment readiness

echo "🇰🇪 BTNG SOVEREIGN PLATFORM - VERIFICATION SUITE"
echo "==============================================="
echo "Testing complete BTNG sovereign gold standard system"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track results
TOTAL_TESTS=0
PASSED_TESTS=0

function run_test() {
    local test_name="$1"
    local test_script="$2"

    echo ""
    echo -e "${YELLOW}▶ Running: ${test_name}${NC}"
    echo "----------------------------------------"

    ((TOTAL_TESTS++))

    if [ -f "$test_script" ]; then
        if bash "$test_script"; then
            echo -e "${GREEN}✅ PASSED: ${test_name}${NC}"
            ((PASSED_TESTS++))
        else
            echo -e "${RED}❌ FAILED: ${test_name}${NC}"
        fi
    else
        echo -e "${RED}❌ MISSING: ${test_script}${NC}"
    fi
}

# Pre-flight checks
echo "🔍 Pre-flight checks..."
echo "----------------------"

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo -e "${RED}❌ Not in BTNG project directory (missing package.json)${NC}"
    exit 1
fi

if [ ! -d "contracts" ]; then
    echo -e "${RED}❌ Contracts directory not found${NC}"
    exit 1
fi

if [ ! -d "k8s" ]; then
    echo -e "${RED}❌ Kubernetes policies directory not found${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Project structure verified${NC}"

# Check Node.js and npm
if command -v node &> /dev/null && command -v npm &> /dev/null; then
    NODE_VERSION=$(node --version)
    NPM_VERSION=$(npm --version)
    echo -e "${GREEN}✅ Node.js ${NODE_VERSION}, npm ${NPM_VERSION}${NC}"
else
    echo -e "${RED}❌ Node.js or npm not found${NC}"
    exit 1
fi

# Check if dependencies are installed
if [ -d "node_modules" ]; then
    echo -e "${GREEN}✅ Dependencies installed${NC}"
else
    echo -e "${YELLOW}⚠️ Installing dependencies...${NC}"
    npm install
fi

echo ""
echo "🚀 Starting BTNG Verification Suite..."
echo "====================================="

# 1. Environment and Credentials Check
run_test "Environment Configuration" "scripts/check-testnet-env.js"

# 2. Sepolia Deployment Verification
run_test "Sepolia Smart Contract Deployment" "scripts/verify-sepolia-deployment.sh"

# 3. Gatekeeper Policy Authentication
run_test "Sovereign Gatekeeper Policies" "scripts/verify-gatekeeper-policies.sh"

# 4. Zero-Trust Admission Control
run_test "Zero-Trust Admission Control (JWT/TLS)" "scripts/verify-admission-control.sh"

# 5. Cryptographic Verification
run_test "Cryptographic Signature Verification" "scripts/verify-cryptography.sh"

# Summary
echo ""
echo "📊 VERIFICATION SUITE RESULTS"
echo "============================"
echo "Total Tests: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $((TOTAL_TESTS - PASSED_TESTS))"

if [ "$PASSED_TESTS" -eq "$TOTAL_TESTS" ]; then
    echo ""
    echo -e "${GREEN}🎉 ALL TESTS PASSED! 🇰🇪${NC}"
    echo ""
    echo "🇰🇪 Your BTNG Sovereign Gold Standard is READY for Sepolia deployment!"
    echo ""
    echo "📋 Final Deployment Steps:"
    echo "1. Update .env with your real credentials"
    echo "2. Run: npm run deploy:testnet"
    echo "3. Verify contracts on Etherscan"
    echo "4. Deploy Gatekeeper policies: ./scripts/deploy-gatekeeper-policies.sh"
    echo ""
    echo "🌟 BTNG will bring sovereign prosperity to African nations! 🌟"
    exit 0
else
    echo ""
    echo -e "${RED}❌ SOME TESTS FAILED${NC}"
    echo ""
    echo "🔧 Please fix the failed tests before deploying to production."
    echo "   Check the error messages above for guidance."
    exit 1
fi