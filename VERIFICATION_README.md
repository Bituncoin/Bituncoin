# BTNG Sovereign Platform - Verification Suite

🇰🇪 **Complete pre-deployment verification for the BTNG Sovereign Gold Standard**

This verification suite ensures your BTNG system is cryptographically secure, sovereign-compliant, and ready for Sepolia testnet deployment.

## 🚀 Quick Start

```bash
# Run complete verification suite
npm run verify:all

# Or run individual checks
npm run verify:sepolia      # Smart contract deployment
npm run verify:gatekeeper   # Security policies
npm run verify:admission    # Access control
npm run verify:crypto       # Cryptographic verification
```

## 📋 Verification Checks

### 1. 🔗 Sepolia Smart Contract Deployment
**File:** `scripts/verify-sepolia-deployment.sh`
- ✅ Environment variables validation
- ✅ Infura API connectivity
- ✅ Wallet balance verification
- ✅ Contract compilation check
- ✅ Gas estimation
- ✅ Deployment simulation

### 2. 🛡 Sovereign Gatekeeper Policies
**File:** `scripts/verify-gatekeeper-policies.sh`
- ✅ Policy file existence
- ✅ Vulnerability + VEX enforcement
- ✅ Signature verification logic
- ✅ Runtime hardening rules
- ✅ Base image allowlist validation

### 3. 🔐 Zero-Trust Admission Control
**File:** `scripts/verify-admission-control.sh`
- ✅ JWT creation and verification
- ✅ Sovereign entity authorization
- ✅ TLS certificate validation
- ✅ Admission request simulation

### 4. 🔒 Cryptographic Verification
**File:** `scripts/verify-cryptography.sh`
- ✅ Ethereum signature verification
- ✅ Hash function consistency
- ✅ HMAC generation
- ✅ Transaction signing simulation
- ✅ Sovereign certificate validation

## 🏗 System Requirements

### Prerequisites
- **Node.js** 18+ with npm
- **OpenSSL** for cryptographic operations
- **jq** for JSON processing (optional)
- **bash** shell (or Windows with WSL/Git Bash)

### Project Structure
```
btng-sovereign-platform/
├── .env                          # Environment variables
├── contracts/                    # Solidity contracts
├── k8s/                          # Gatekeeper policies
├── scripts/
│   ├── verify-*.sh              # Individual verifications
│   ├── btng-verification-suite.sh  # Complete suite
│   └── btng-verification-suite.bat # Windows version
└── package.json                  # npm scripts
```

## 🔧 Manual Verification Steps

If automated scripts fail, verify manually:

### Environment Setup
```bash
# Check .env file
cat .env

# Validate credentials format
npm run check-env
```

### Contract Compilation
```bash
# Compile contracts
npm run compile

# Check artifacts
ls -la artifacts/contracts/
```

### Network Connectivity
```bash
# Test Sepolia connection
npm run test-connection
```

## 📊 Expected Results

### ✅ All Tests Pass
```
🎉 ALL TESTS PASSED! 🇰🇪
🇰🇪 Your BTNG Sovereign Gold Standard is READY for Sepolia deployment!
```

### ❌ Test Failures
- **Environment**: Check `.env` credentials
- **Sepolia**: Verify Infura API key and Sepolia ETH balance
- **Gatekeeper**: Ensure policy files exist in `k8s/`
- **Admission**: Check JWT/TLS configuration
- **Crypto**: Verify private key format

## 🚀 Post-Verification Deployment

After all tests pass:

### 1. Update Credentials
```bash
# Edit .env with real values
code .env
```

### 2. Deploy to Sepolia
```bash
npm run deploy:testnet
```

### 3. Verify on Etherscan
- Check contract addresses
- Verify source code
- Confirm ownership

### 4. Deploy Security Policies
```bash
./scripts/deploy-gatekeeper-policies.sh
```

## 🔍 Troubleshooting

### Common Issues

**"bash: command not found" (Windows)**
```bash
# Use Windows batch file instead
scripts\btng-verification-suite.bat
```

**"INFURA_API_KEY not set"**
```bash
# Update .env file
echo "INFURA_API_KEY=your_key_here" >> .env
```

**"Low balance"**
```bash
# Get Sepolia ETH from faucet
start https://sepoliafaucet.com
```

**"Contract compilation failed"**
```bash
# Check Solidity syntax
npx hardhat compile --verbose
```

## 🛡 Security Validation

The verification suite ensures:

- **Cryptographic Integrity**: All signatures verified
- **Sovereign Compliance**: African nation entities authorized
- **Zero-Trust Security**: Admission control enforced
- **Vulnerability Prevention**: Gatekeeper policies active
- **Production Readiness**: All systems tested

## 🌟 BTNG Sovereign Mission

This verification suite protects the integrity of Africa's first sovereign digital gold standard, ensuring that only cryptographically verified, sovereign-approved transactions can enter the BTNG ecosystem.

**🇰🇪 Sovereign prosperity through digital gold 🇰🇪**

---

**Verification Suite Version:** 1.0.0
**BTNG Protocol:** Sovereign Gold Standard v2.0
**Testnet:** Sepolia
**Mainnet Ready:** ✅