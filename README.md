BTNG Sovereign Gold Standard Platform
Version: 1.0.0
Architecture: Next.js 14 + BTNG Quantum Core
Status: 🟢 PRODUCTION / FULLY ACTIVATED

A unified sovereign gold-backed digital identity and economic engine. The platform integrates a Quantum Private Banker (AI), Music-Driven Block Mining, and Automated Liquidity Rebalancing (ALR) on a 54-nation mesh network anchored to the Ghana Primary MTN node.

🏛️ Sovereign Activation (2026)
As of March 2026, the BTNG system has transitioned to full production mode under the authority of INTERNATIONAL QUIYTY TRUST and Ekuye Digital Gateway Trust LTD.

Core Activated Modules
Quantum Private Banker — Unified AI for Sovereign Multilingual Guidance.
Sovereign Identity Registry — Full Lifecycle Enforcement (ACTIVE/EXPIRED/REVOKED).
BTNG-GOLD Ecosystem — 100% Gold-backed (1g = 1 BTNG-G).
Music-Driven Block Mining — Listen-to-Earn (1 BTNG-G / 60s) via MusicMiningContract.
Sovereign SMS Alert Service — Automated MTG Ghana Chenosis alerts for on-boarded artists.
NFT Pilot Manager — Master Batch Initialization for Creative Assets.
Zero-Trust Security Gateway — Tiered Biometric + Device Binding Authentication.
Automated Liquidity Rebalancer (ALR) — Dynamic Sovereign Wealth Stewardship.
🛡️ Platform Infrastructure
Primary Anchor: 154.161.183.158:38982 (Ghana Sovereign Node)
Security Key: al-A-cKa4Yh49mpq59IgyPGQd5jO4iMykOtBW2OoSs814R
Network: btng712-fabric-network (AWS Managed Blockchain)
🏗️ Architecture Overview
This monorepo contains three integrated systems:

BTNG Node.js API (Port 3003) - Sovereign identity and gold price API
Genesis Trading Platform - Real-time trading platform with gold price integration
Ethereum Smart Contracts - Sovereign gold standard smart contracts
Developer Master Folder - BTNG-DEVELOPER-MASTER/12_BITUNCOINOS houses the sovereign BituncoinOS modules, CLI, SDK, and docs so the OS shares the Quantum Memory Container.
Project Structure
/app                  # Next.js app & API routes
  /api/btng          # BTNG API endpoints (JWT auth)
/genesis-app         # Genesis trading platform
  /server/src/main/genesis/
    BTNGPriceBridge.kts    # Bridge service (polls API)
/contracts           # Solidity smart contracts
/components          # Sovereign UI components
/lib                 # Core logic modules
/scripts             # Build & test scripts
BTNG-DEVELOPER-MASTER/ # BituncoinOS master container (Quantum Memory + Developer Master Folder)
  12_BITUNCOINOS/     # Sovereign OS modules, CLI, SDK, docs
🚀 Quick Start
Prerequisites
Node.js 18+
Java 11+ (for Genesis)
MongoDB (local or cloud)
Hardhat (for smart contracts)
Installation
# Install all dependencies
npm install

# Install concurrently for multi-service development
npm install -g concurrently
Environment Setup
Create .env.local:

# BTNG API Configuration
BTNG_ADMIN_PASSWORD=sovereign2024
JWT_SECRET=your-super-secure-jwt-secret-here
MONGODB_URI=mongodb://localhost:27017/btng-sovereign

# Genesis Configuration
GENESIS_DB_URL=jdbc:postgresql://localhost:5432/genesis
GENESIS_DB_USER=genesis
GENESIS_DB_PASSWORD=genesis

# Ethereum Configuration
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
PRIVATE_KEY=your-private-key-without-0x-prefix
ETHERSCAN_API_KEY=your-etherscan-api-key
Development Mode (All Services)
# Run BTNG API + Genesis platform simultaneously
npm run dev:all
This starts:

BTNG API on http://localhost:3003
Genesis platform on http://localhost:8080
Individual Services
# BTNG API only
npm run dev:api

# Genesis platform only
npm run dev:genesis

# Smart contract development
npm run compile
npx hardhat node
🔧 Available Scripts
Development
npm run dev:all - Run all services concurrently
npm run dev:api - BTNG API on port 3003
npm run dev:genesis - Genesis platform
npm run dev - Next.js development server
Smart Contracts
npm run compile - Compile Solidity contracts
npm run test - Run contract tests
npm run deploy:local - Deploy to local Hardhat network
npm run deploy:testnet - Deploy to Sepolia testnet
Testing & Verification
npm run health - API health check
npm run test-jwt - Test JWT authentication
npm run test-fabric - Test fabric network integration
npm run test-gold-api - Test gold price endpoints
npm run test:integration - Full BTNG-Genesis integration test
npm run verify:all - Complete verification suite
📚 Operations Docs
BTNG Gold Price System Runbook
BTNG Merchant API Guide
BTNG Gold System Completion Notes
BTNG Fabric Network Configuration Inquiry
PM2 Oracle Service
npm run pm2:oracle:start - Start gold oracle updater service
npm run pm2:oracle:status - View current PM2 service state
npm run pm2:oracle:logs - Tail updater logs
npm run pm2:oracle:save - Persist PM2 process list across reboot
Documentation Identity (SHA-256 + ES256)
npm run docs:identity:refresh - Regenerate, propagate, and verify canonical docs identity
npm run docs:identity:watch - Auto-regenerate identity on canonical docs changes
npm run docs:identity:verify - Verify manifest hash and ES256 signature
Watchtower metadata endpoint: GET /api/watchtower/meta
Watchtower hard-fail gate: snapshots and node registrations are blocked when documentation identity is invalid
Optional override (not recommended for production): set BTNG_ENFORCE_DOCUMENTATION_IDENTITY=false
All-in-One Network Verification
npm run verify:btng-network - Strict consolidated checks (production gate)
npm run verify:btng-network:strict - Explicit strict mode
npm run verify:btng-network:soft - Rollout mode (treats readiness failures as warnings)
npm run verify:btng-network:ci - Strict mode with deterministic CI artifact path
Verifier JSON artifact outputs:

strict: cache/verify-btng-network-strict.json
soft: cache/verify-btng-network-soft.json
custom: --output <path> or BTNG_NETWORK_VERIFY_OUTPUT
🏗️ Hyperledger Fabric Network
Network ID: btng-fabric-network
Root Member: btng-root-member
Node ID: nd-6HRNJ6OUIBGP3MV74YAW53NWYQ

Fabric Integration
The BTNG platform integrates with Hyperledger Fabric for distributed ledger operations:

Network Configuration
Channel: btng-sovereign-channel
Chaincodes: btng-gold-token, btng-sovereign-identity
Organizations: Sovereign member nodes
Consensus: Raft ordering service
API Endpoints
# Network status (public)
curl http://localhost:3003/api/btng/fabric/network

# Node status (public)
curl http://localhost:3003/api/btng/fabric/node

# Chaincode operations (requires JWT)
curl -X POST http://localhost:3003/api/btng/fabric/network \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"invoke","chaincode":"btng-gold-token","function":"mint"}'
Testing Fabric Integration
# Test fabric network endpoints
npm run test-fabric
Deploying Chaincode
# Generate fabric deployment files
npm run deploy-fabric

# Start fabric network (requires Docker)
docker-compose -f fabric/docker-compose.yml up -d

# Create channel
./fabric/scripts/create-channel.sh

# Deploy chaincodes
./fabric/scripts/deploy-btng-gold-token.sh
./fabric/scripts/deploy-btng-sovereign-identity.sh
Sovereign Operations
Gold Token Minting - Fabric-based gold token issuance
Identity Verification - Distributed identity validation
Transaction Ledger - Immutable transaction records
Consensus Validation - Multi-party transaction validation
Chaincode Operations via API
# Mint gold tokens
curl -X POST http://localhost:3003/api/btng/fabric/chaincode \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chaincode": "btng-gold-token",
    "function": "Mint",
    "args": ["1000", "BTNG-SOVEREIGN-001"]
  }'

# Register sovereign identity
curl -X POST http://localhost:3003/api/btng/fabric/chaincode \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chaincode": "btng-sovereign-identity",
    "function": "RegisterIdentity",
    "args": ["BTNG-IDENTITY-001", "public-key", "{\"type\":\"sovereign\"}"]
  }'
🔐 Authentication & API
JWT Authentication
# Login to get JWT token
curl -X POST http://localhost:3003/api/btng/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"sovereign2024"}'
Gold Price API
# Get current prices (requires JWT)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3003/api/btng/gold/prices
Genesis Bridge Service
The BTNGPriceBridge.kts service:

Polls BTNG API every minute
Authenticates using JWT
Stores prices in Genesis database
Enables real-time trading data
🧪 Testing
Integration Testing
# Run comprehensive BTNG-Genesis integration tests
npm run test:integration
Tests cover:

API health and authentication
Gold price data retrieval
Bridge service functionality
End-to-end data flow
Genesis platform connectivity
Smart Contract Testing
npm run test
npm run coverage
⚙️ Operational Status
Module	Status
Core Architecture	✅ Operational
Identity Layer	🟡 In Progress
QR Wallet	🟡 In Progress
Trust Union	📋 Planned
Country Onboarding	📋 Planned
Mobile Money	📋 Planned
Genesis Bridge	✅ Operational
Gold Price API	✅ Operational
🌍 Expansion Pathways
Country-specific onboarding modules
Merchant integration workflows
Mobile money adapter framework
Debt-release protocol
Proof-of-value dashboard
Genesis trading integration
📊 Health & Observability
# API health check
npm run health

# Full integration test
npm run test:integration
Monitors:

Platform availability
Identity service status
Trust-union endpoint health
Wallet transaction capacity
Genesis bridge connectivity
Gold price data freshness
🔒 Security & Sovereignty
No external tracking (all legacy analytics removed)
Sovereign identity architecture
Zero-knowledge proof-of-value
Trust-first protocol design
JWT-secured API endpoints
Genesis database integration
🌐 Network Anchors
Primary Backend Endpoint: http://154.161.183.158:38982
Genesis Transaction Hash: 0x1111111111111111111111111111111111111111111111111111111111111111
Genesis Explorer URL: http://154.161.183.158:38982/explorer/tx/0x1111111111111111111111111111111111111111111111111111111111111111
Genesis Block Height: 12458
Genesis Timestamp: 1771457774
📧 Contact
admin@bituncoin.africa
gold@btng.africa
BTNG — Building Trust. Nurturing Growth.
