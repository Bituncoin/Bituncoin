# BTNG Sovereign Gold Standard Platform

**Version:** 0.1.0  
**Architecture:** Next.js 14 (App Router) + Genesis Trading Platform  
**Status:** Launch Phase - Operational Readiness

A comprehensive sovereign gold-backed digital identity and value platform integrating Node.js API, Ethereum smart contracts, and Genesis trading platform for real-time gold price data streaming.

---

## 🛡️ Platform Identity

BTNG is a sovereign identity and value platform built for institutional-grade trust operations, universal identity management, and proof-of-value workflows.

### Core Modules
- **Gold Card Identity Layer** — Universal identity credentials
- **QR Wallet Module** — Sovereign value transfer
- **Trust Union Protocol** — Institutional trust architecture
- **Country Onboarding** — Sovereign expansion pathways
- **Mobile Money Adapters** — Financial inclusion infrastructure
- **Genesis Bridge Service** — Real-time gold price integration

---

## 🏗️ Architecture Overview

This monorepo contains three integrated systems:

1. **BTNG Node.js API** (Port 3003) - Sovereign identity and gold price API
2. **Genesis Trading Platform** - Real-time trading platform with gold price integration
3. **Ethereum Smart Contracts** - Sovereign gold standard smart contracts

### Project Structure

```
/app                  # Next.js app & API routes
  /api/btng          # BTNG API endpoints (JWT auth)
/genesis-app         # Genesis trading platform
  /server/src/main/genesis/
    BTNGPriceBridge.kts    # Bridge service (polls API)
/contracts           # Solidity smart contracts
/components          # Sovereign UI components
/lib                 # Core logic modules
/scripts             # Build & test scripts
```

---

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Java 11+ (for Genesis)
- MongoDB (local or cloud)
- Hardhat (for smart contracts)

### Installation

```bash
# Install all dependencies
npm install

# Install concurrently for multi-service development
npm install -g concurrently
```

### Environment Setup

Create `.env.local`:

```env
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
```

### Development Mode (All Services)

```bash
# Run BTNG API + Genesis platform simultaneously
npm run dev:all
```

This starts:
- BTNG API on http://localhost:3003
- Genesis platform on http://localhost:8080

### Individual Services

```bash
# BTNG API only
npm run dev:api

# Genesis platform only
npm run dev:genesis

# Smart contract development
npm run compile
npx hardhat node
```

---

## 🔧 Available Scripts

### Development
- `npm run dev:all` - Run all services concurrently
- `npm run dev:api` - BTNG API on port 3003
- `npm run dev:genesis` - Genesis platform
- `npm run dev` - Next.js development server

### Smart Contracts
- `npm run compile` - Compile Solidity contracts
- `npm run test` - Run contract tests
- `npm run deploy:local` - Deploy to local Hardhat network
- `npm run deploy:testnet` - Deploy to Sepolia testnet

### Testing & Verification
- `npm run health` - API health check
- `npm run test-jwt` - Test JWT authentication
- `npm run test-gold-api` - Test gold price endpoints
- `npm run test:integration` - Full BTNG-Genesis integration test
- `npm run verify:all` - Complete verification suite

---

## 🔐 Authentication & API

### JWT Authentication

```bash
# Login to get JWT token
curl -X POST http://localhost:3003/api/btng/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"sovereign2024"}'
```

### Gold Price API

```bash
# Get current prices (requires JWT)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:3003/api/btng/gold/prices
```

### Genesis Bridge Service

The `BTNGPriceBridge.kts` service:
- Polls BTNG API every minute
- Authenticates using JWT
- Stores prices in Genesis database
- Enables real-time trading data

---

## 🧪 Testing

### Integration Testing

```bash
# Run comprehensive BTNG-Genesis integration tests
npm run test:integration
```

Tests cover:
- API health and authentication
- Gold price data retrieval
- Bridge service functionality
- End-to-end data flow
- Genesis platform connectivity

### Smart Contract Testing

```bash
npm run test
npm run coverage
```

---

## ⚙️ Operational Status

| Module | Status |
|--------|--------|
| Core Architecture | ✅ Operational |
| Identity Layer | 🟡 In Progress |
| QR Wallet | 🟡 In Progress |
| Trust Union | 📋 Planned |
| Country Onboarding | 📋 Planned |
| Mobile Money | 📋 Planned |
| **Genesis Bridge** | ✅ **Operational** |
| **Gold Price API** | ✅ **Operational** |

---

## 🌍 Expansion Pathways

- Country-specific onboarding modules
- Merchant integration workflows
- Mobile money adapter framework
- Debt-release protocol
- Proof-of-value dashboard
- **Genesis trading integration**

---

## 📊 Health & Observability

```bash
# API health check
npm run health

# Full integration test
npm run test:integration
```

Monitors:
- Platform availability
- Identity service status
- Trust-union endpoint health
- Wallet transaction capacity
- **Genesis bridge connectivity**
- **Gold price data freshness**

---

## 🔒 Security & Sovereignty

- No external tracking (all legacy analytics removed)
- Sovereign identity architecture
- Zero-knowledge proof-of-value
- Trust-first protocol design
- **JWT-secured API endpoints**
- **Genesis database integration**

---

**BTNG** — Building Trust. Nurturing Growth.
