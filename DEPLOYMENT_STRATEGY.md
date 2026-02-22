# BTNG Deployment Strategy Guide

🇰🇪 **Strategic deployment approaches for the BTNG Sovereign Gold Standard**

## 📊 Current Status: ✅ API Authentication Complete

**Latest Updates (February 22, 2026):**
- ✅ **JWT Authentication**: Fully implemented and tested
- ✅ **API Endpoints**: All routes responding with proper JSON
- ✅ **TypeScript Errors**: Fixed compilation issues
- ✅ **Dependencies**: Installed missing packages (framer-motion, @types/jsonwebtoken)
- ✅ **Server Stability**: Next.js dev server running without errors
- ✅ **Database Fallback**: MongoDB connection with in-memory storage fallback

### 🔐 Authentication Status
- **Login Endpoint**: `POST /api/auth/login` ✅ Working
- **Protected Routes**: `POST /api/btng/gold/price` ✅ JWT-protected
- **Public Routes**: `GET /api/btng/gold/price/status` ✅ Accessible
- **Token Validation**: HS256 with 24-hour expiration ✅ Implemented

### 🧪 API Testing Results
```bash
# Status endpoint (no auth required)
✅ GET /api/btng/gold/price/status → JSON response

# Authentication flow
✅ POST /api/auth/login → JWT token generated
✅ POST /api/btng/gold/price (with valid token) → Success
❌ POST /api/btng/gold/price (no token) → 401 Unauthorized
❌ POST /api/btng/gold/price (invalid token) → 401 Unauthorized
```

## � API Fixes Applied

**Issues Resolved:**
- **TypeScript Compilation**: Fixed JSX syntax errors in `app/(onboarding)/user/page.tsx`
- **Missing Dependencies**: Installed `framer-motion` and `@types/jsonwebtoken`
- **Client Components**: Added `"use client"` directive to interactive components
- **Path Resolution**: Verified `@/lib/*` aliases working correctly
- **Database Fallback**: Implemented MongoDB → in-memory storage graceful degradation

**Server Status:**
- ✅ Next.js dev server running on `http://localhost:3003`
- ✅ Hot reload working without compilation errors
- ✅ API routes responding with proper JSON (not HTML error pages)
- ✅ JWT authentication middleware functioning correctly

## �🚀 Next Steps: Smart Contract Deployment

**Ready for Sepolia Testnet Deployment:**
1. **Environment Check**: `npm run check-env`
2. **Network Test**: `npm run test-connection`
3. **Deploy Oracle**: `npm run deploy:oracle`
4. **Full Verification**: `npm run verify:all`

---

## 📊 Decision Framework

| Scenario | Recommended Approach | Reasoning |
|----------|---------------------|-----------|
| **First deployment** | Single Contract (Oracle) | Test network connectivity, gas costs, verification |
| **Token economics iteration** | Single Contract (Token) | Fast feedback on token logic, supply mechanics |
| **Custody vault testing** | Full Suite | Test inter-contract dependencies |
| **Production release** | Full Suite | Complete system validation |
| **Cost optimization** | Single Contract | Minimize Sepolia ETH usage |
| **CI/CD pipeline** | Full Suite | Comprehensive testing |

## 🚀 Deployment Options

### 1. Single Contract Deployment
**Best for:** Development iteration, cost control, focused testing

```bash
# Deploy Oracle only (no dependencies)
npm run deploy:oracle

# Deploy Token only (requires manual custody setup)
# Edit scripts/deploy-token-only.js with your parameters
```

### 2. Full Suite Deployment
**Best for:** Production releases, complete system testing

```bash
# Deploy all contracts with dependency management
npm run deploy:full
```

### 3. Original Deployment
**Best for:** Quick deployment with existing logic

```bash
# Use existing deployment script
npm run deploy:testnet
```

## 🔗 Contract Dependencies

```
BTNGGoldOracle
    ↓ (no dependencies)

BTNGGoldToken
    ↓ (initially uses ZeroAddress)
    ↓ (updated with custody address after BTNGCustody deployment)

BTNGCustody
    ↓ (depends on: BTNGGoldToken, BTNGGoldOracle)
```

## 📋 Deployment Phases

### Phase 1: Contract Deployment
1. **BTNGGoldOracle** - Independent deployment
2. **BTNGGoldToken** - Deployed with placeholder custody
3. **BTNGCustody** - Deployed with real token/oracle addresses

### Phase 2: Contract Linking
1. Update BTNGGoldToken with BTNGCustody address
2. Transfer ownerships to deployer

### Phase 3: Verification
1. Test contract interactions
2. Verify ownership transfers
3. Run integration tests

## 💰 Gas Cost Estimation

| Deployment Type | Contracts | Est. Gas Cost | Sepolia ETH |
|----------------|-----------|---------------|-------------|
| Oracle Only | 1 | ~100k | ~0.001 ETH |
| Token Only | 1 | ~200k | ~0.002 ETH |
| Full Suite | 3 | ~500k | ~0.005 ETH |

*Costs are estimates. Get Sepolia ETH from: https://sepoliafaucet.com*

## 🧪 Testing Strategy

### Single Contract Testing
```bash
# Deploy Oracle
npm run deploy:oracle

# Test Oracle functionality
npx hardhat test test/BTNGGoldOracle.test.js

# Verify on Etherscan
# https://sepolia.etherscan.io/address/<ORACLE_ADDRESS>
```

### Full Suite Testing
```bash
# Deploy all contracts
npm run deploy:full

# Run complete test suite
npm test

# Run verification suite
npm run verify:all
```

## 🔍 Verification Commands

```bash
# Environment check
npm run check-env

# Network connectivity
npm run test-connection

# Full verification suite
npm run verify:all

# Individual verifications
npm run verify:sepolia
npm run verify:gatekeeper
npm run verify:admission
npm run verify:crypto
```

## 🌟 Recommended Workflow

### Development Phase
1. **Start with Oracle**: `npm run deploy:oracle`
2. **Test thoroughly** with unit tests
3. **Deploy Token**: Create and run token-only script
4. **Test token economics**
5. **Move to Full Suite**: `npm run deploy:full`

### Production Phase
1. **Run verification suite**: `npm run verify:all`
2. **Deploy full suite**: `npm run deploy:full`
3. **Verify on Etherscan**
4. **Run integration tests**
5. **Deploy Gatekeeper policies**

## 🛠 Troubleshooting

### Common Issues

**"Contract deployment failed"**
- Check wallet balance: `npm run test-connection`
- Verify network configuration in `hardhat.config.js`
- Check constructor arguments

**"Dependency resolution failed"**
- Use full suite deployment for interdependent contracts
- Check deployment order in `deploy-config.js`

**"Verification failed"**
- Run `npm run verify:all` to identify specific issues
- Check contract addresses in deployment output
- Verify network connectivity

## 📈 Success Metrics

### Single Contract Success
- ✅ Contract deployed successfully
- ✅ Basic functionality tests pass
- ✅ Gas costs within expected range
- ✅ Contract visible on Etherscan

### Full Suite Success
- ✅ All contracts deployed
- ✅ Inter-contract communication works
- ✅ Ownership transfers complete
- ✅ Integration tests pass
- ✅ Verification suite passes

## 🎯 Next Steps

**For your current development cycle:**

1. **If iterating on Oracle logic**: Use `npm run deploy:oracle`
2. **If testing token economics**: Deploy token separately first
3. **If ready for production**: Use `npm run deploy:full`

**Recommended starting point:** Deploy the Oracle first to establish your deployment pipeline, then move to the full suite.

---

# 🔐 JWT Authentication Guide

**Quick-start guide for JWT operations in the BTNG ecosystem**

## 📋 JWT Operations Overview

| Operation | Browser Tool | Command Line | Use Case |
|-----------|-------------|--------------|----------|
| **Decode** | JWT Decoder | Node/Python scripts | Inspect token contents |
| **Verify** | Signature Validator | CLI verification | Authenticate requests |
| **Create** | JWT Encoder | Programmatic generation | Issue access tokens |

## 🌐 Browser-Based JWT Tools

### 1️⃣ Decoding a JWT
- **Tool**: [jwt.io](https://jwt.io) or similar JWT decoder
- **Input**: Paste the JWT token
- **Output**: Automatically displays:
  - **Header**: Algorithm (`alg`), type (`typ`), etc.
  - **Payload**: Claims (`sub`, `name`, `iat`, `exp`, etc.)
  - **Signature**: Base64URL-encoded signature

### 2️⃣ Verifying a JWT Signature
**Step-by-step process:**
1. **Secret/Key**: Paste HMAC secret or RSA/ECDSA public key
2. **Validate**: Click "Verify Signature"
3. **Result**: "Signature Verified" ✅ or "Invalid Signature" ❌

⚠️ **Security Note**: Never paste production secrets into public tools!

### 3️⃣ Creating a JWT
**Browser UI workflow:**
1. **Payload**: Enter JSON claims (sub, name, admin, etc.)
2. **Algorithm**: Choose HS256, RS256, ES256, etc.
3. **Secret/Key**: Paste signing key
4. **Generate**: Click to create token
5. **Copy**: Use the generated JWT string

## 💻 Command Line Implementations

### 4.1 Node.js (jsonwebtoken library)

**Install:**
```bash
npm install -g jsonwebtoken
```

**Create JWT:**
```bash
node -e "
  const jwt = require('jsonwebtoken');
  const payload = { sub: '123', name: 'Alice', admin: true, iat: Math.floor(Date.now()/1000) };
  const secret = 'a-string-secret-at-least-256-bits-long';
  const token = jwt.sign(payload, secret, { algorithm: 'HS256', expiresIn: '1h' });
  console.log('JWT:', token);
"
```

**Verify JWT:**
```bash
node -e "
  const jwt = require('jsonwebtoken');
  const token = 'eyJh...';
  const secret = 'a-string-secret-at-least-256-bits-long';
  try {
    const payload = jwt.verify(token, secret, { algorithms:['HS256']});
    console.log('✅ Signature OK', payload);
  } catch(e) {
    console.error('❌ Signature Invalid:', e.message);
  }
"
```

### 4.2 OpenSSL + HMAC (No External Libraries)

**Create JWT:**
```bash
# Header + Payload as JSON strings
header='{"alg":"HS256","typ":"JWT"}'
payload='{"sub":"123","name":"Alice","admin":true}'

# Base64URL encode function
b64() { echo -n "$1" | openssl base64 -e -A | tr '+/' '-_' | tr -d '='; }

header64=$(b64 "$header")
payload64=$(b64 "$payload")

# Sign (HMAC-SHA256)
secret='a-string-secret-at-least-256-bits-long'
sig=$(printf "%s.%s" "$header64" "$payload64" | openssl dgst -sha256 -hmac "$secret" -binary | base64 | tr '+/' '-_' | tr -d '=')

# Final JWT
jwt=$(printf "%s.%s.%s" "$header64" "$payload64" "$sig")
echo "$jwt"
```

### 4.3 Python (PyJWT)

**Verify JWT:**
```python
python - <<'PY'
import jwt, sys

token = sys.argv[1]
secret = 'a-string-secret-at-least-256-bits-long'

try:
    payload = jwt.decode(token, secret, algorithms=['HS256'])
    print('✅ OK', payload)
except jwt.exceptions.InvalidSignatureError:
    print('❌ Signatures do NOT match')
PY eyJhMcD...
```

## 🛡️ Security Best Practices

### ✅ Essential Guidelines
- **🔐 Secret Management**: Store secrets in environment variables, never hardcode
- **🔍 Minimal Payloads**: Keep JWT size small, include only necessary claims
- **⏰ Expiration**: Always set `exp` claim, don't rely on `iat` alone
- **🔝 Strong Algorithms**: Use RS256/ES256 for asymmetric keys when possible
- **⚠️ Validation Required**: Never trust incoming JWTs without signature verification
- **📦 Trusted Libraries**: Use well-tested libraries (jsonwebtoken, PyJWT, jose), pin versions

### 🚨 Security Checklist
- [ ] Secrets stored in environment variables
- [ ] JWTs include expiration (`exp`) claims
- [ ] Strong signing algorithms used
- [ ] All incoming tokens verified before use
- [ ] No sensitive data in JWT payloads
- [ ] Regular key rotation implemented

## ⚡ Quick Validation One-Liner

**Node.js validation:**
```bash
node -e "
  const jwt=require('jsonwebtoken');
  const token='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30';
  const secret='a-string-secret-at-least-256-bits-long';
  try{
    const payload=jwt.verify(token, secret, {algorithms:['HS256']});
    console.log('✓ Valid!  Payload:', payload);
  }catch(e){
    console.error('✗ Invalid:', e.message);
  }
"
```

**Expected Output:**
```
✓ Valid!  Payload: { sub: '1234567890', name: 'John Doe', admin: true, iat: 1516239022 }
```

##  Integration with BTNG APIs

**For BTNG API endpoints:**
1. **Decode incoming JWTs** to extract user claims
2. **Verify signatures** before processing requests
3. **Generate tokens** for authenticated users
4. **Validate permissions** based on JWT payload claims

**Example API middleware:**
```javascript
// Verify JWT in API routes
const jwt = require('jsonwebtoken');

function authenticateToken(req, res, next) {
  const token = req.headers['authorization']?.split(' ')[1];
  if (!token) return res.status(401).json({ error: 'Access token required' });

  jwt.verify(token, process.env.JWT_SECRET, (err, user) => {
    if (err) return res.status(403).json({ error: 'Invalid token' });
    req.user = user;
    next();
  });
}
```

### 🛠️ Implementation in BTNG

**JWT Authentication has been implemented for your BTNG APIs:**

#### 🔐 Authentication Endpoints
- **POST** `/api/auth/login` - Generate JWT tokens
- **Protected routes** require `Authorization: Bearer <token>` header

#### 🛡️ Protected Endpoints
- **POST** `/api/btng/gold/price` - Requires admin authentication
- **Public endpoints** like status checks remain accessible

#### 🧪 Testing JWT Auth
```bash
# Install dependencies
npm install

# Test JWT authentication
npm run test-jwt
```

**Example API Usage:**
```bash
# 1. Login to get token
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'

# 2. Use token for protected endpoints
curl -X POST http://localhost:3000/api/btng/gold/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"base_price_gram":50.00,"base_price_ounce":1600.00,"base_price_kilo":50000.00}'
```

---

**Ready to implement JWT authentication?** Use the browser tools for testing, then integrate the command-line approaches into your BTNG API endpoints! 🔐

---

# 🔧 API Endpoint Troubleshooting Guide

**Debugging Next.js API routes and request issues**

## 🚨 Common Symptoms & Root Causes

| Symptom | What's Actually Happening |
|---------|---------------------------|
| `curl -s http://localhost:3000/api/btng/gold/price/status` prints "Supply values for the following parameters: Uri:" | Your **Node/Next.js** app is running an **interactive prompt** instead of the API route |
| `Invoke-WebRequest ... -Body {"username":...}` returns nothing | PowerShell command malformed - JSON body parsed as separate arguments |
| `netstat -ano \| findstr :3000` shows *LISTENING* on 0.0.0.0:3000 | TCP socket exists, process is listening |
| `Get-Process -Name node \| Where-Object ...` prints nothing | PowerShell only sees *node* executables; your process is `node.exe` |

**Result**: API endpoints are **not being hit** because URLs point to CLI prompts or requests are malformed.

## 🔧 Quick Fixes

### 1️⃣ Verify Route Implementation

**File Structure:**
```
app/api/btng/gold/price/status/route.ts  ✅ (App Router)
pages/api/btng/gold/price/status.ts      ✅ (Pages Router)
```

**Route Handler:**
```typescript
// app/api/btng/gold/price/status/route.ts
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  // Your logic here
  return NextResponse.json({
    status: 'ok',
    price: 65.38,
    timestamp: Date.now()
  });
}
```

**Common Issues:**
- File must export handler function
- Check `next.config.js` for route rewrites
- Watch dev server logs for build errors

### 2️⃣ Correct API Calls

#### PowerShell (Windows)
```powershell
# POST request
Invoke-WebRequest `
    -Uri 'http://localhost:3000/api/auth/login' `
    -Method POST `
    -Body '{"username":"admin","password":"password123"}' `
    -ContentType 'application/json' | Select-Object -ExpandProperty Content
```

#### Bash/CMD (Cross-platform)
```bash
# GET request
curl -s http://localhost:3000/api/btng/gold/price/status

# POST request
curl -s -X POST http://localhost:3000/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"password123"}'
```

### 3️⃣ Verify Node Process
```powershell
# Correct process name matching
Get-Process node, node.exe | Where-Object {$_.CommandLine -like "*next*"}
```

**Expected Output:**
```
Handles  NPM(K)    PM(K)      WS(K)     CPU(s)  Id ProcessName
-------  ------  -------  --------  -------  -- -----------
   576    142    74248     50786       70.3  82516 node
```

### 4️⃣ Check Server Logs

**Watch for these in dev server console:**
- `SyntaxError: Unexpected token` → Bug in API route code
- `Error: Host not found` → Wrong `process.env.API_BASE_URL`
- `Request to /api/btng/gold/price/status` → Request reaching server
- Add `console.log('REACHED')` in handlers to confirm execution

### 5️⃣ Remove Interactive Prompts

**Problem Code:**
```javascript
const rl = require('readline').createInterface(process.stdin, process.stdout);
rl.question('Supply values for the following parameters: \n', (ans) => { ... });
```

**Solution - Guard behind dev flag:**
```javascript
if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
  // Interactive prompts only in development
  const rl = require('readline').createInterface(process.stdin, process.stdout);
  // ... prompt logic
}
```

## 📦 Run-to-Success Checklist

```bash
# 0️⃣ Navigate to project root
cd C:\BTNGAI_files

# 1️⃣ Start dev server (if not running)
npm run dev

# 2️⃣ In separate terminal, test endpoints
curl -s http://localhost:3000/api/btng/gold/price/status
# OR
Invoke-WebRequest 'http://localhost:3000/api/auth/login' -Method POST -Body '{"username":"admin","password":"password123"}' -ContentType 'application/json'
```

**Success Indicators:**
- JSON response (not CLI prompt)
- No "Supply values" text
- Proper HTTP status codes

## 🎯 Advanced Troubleshooting

| Issue | Likely Fix |
|-------|------------|
| **"Supply values for the following parameters"** | Remove `readline` prompts from API handlers |
| **Empty curl response** | Verify route path (case-sensitive), ensure `res.json()` |
| **PowerShell errors** | Check quoting: single quotes around URLs, proper JSON |
| **"Cannot find module"** | Run `npm install`, verify `node_modules` |
| **Process not found** | Use `Get-Process node, node.exe` |
| **"Cannot read property 'json' of undefined"** | Handler not receiving proper request object |

## 🧪 JWT-Specific Testing

**Test JWT Authentication Flow:**
```bash
# 1. Login and get token
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}' | jq -r '.token')

# 2. Use token for protected endpoint
curl -s -X POST http://localhost:3000/api/btng/gold/price \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"base_price_gram":50.00,"base_price_ounce":1600.00,"base_price_kilo":50000.00}'
```

**Verify Token:**
```bash
node -e "
  const jwt = require('jsonwebtoken');
  const token = '$TOKEN';
  const secret = 'a-string-secret-at-least-256-bits-long';
  console.log(jwt.verify(token, secret, {algorithms:['HS256']}));
"
```

## 🚀 Production Deployment Notes

**Environment Variables:**
```bash
# .env.local
JWT_SECRET=your-production-secret-here
NODE_ENV=production
API_BASE_URL=https://your-domain.com
```

**Health Check:**
```bash
# Test production endpoints
curl -s https://your-domain.com/api/btng/gold/price/status | jq .
```

---

**API debugging complete?** Your endpoints should now return proper JSON instead of CLI prompts! 🎉