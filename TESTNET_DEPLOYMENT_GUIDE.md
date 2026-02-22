# BTNG Testnet Deployment Guide

## 🚀 Quick Setup (5 minutes)

### 1. Get Infura Project ID
- Visit: https://infura.io/dashboard
- Sign up/Login (free)
- Create new project → Select "Ethereum"
- Copy the Project ID

### 2. Get Sepolia Test ETH
- Visit: https://sepoliafaucet.com
- Connect your wallet (MetaMask recommended)
- Request 0.5 ETH (free, instant)

### 3. Get Your Private Key
- Open MetaMask
- Account Details → Export Private Key
- Copy the key (keep secure!)

### 4. Configure Environment
- Edit the `.env` file in project root
- Replace placeholders with your actual values

### 5. Verify & Deploy
```bash
npm run check-env
npm run deploy -- --network sepolia
```

## 🚀 Deploying BTNG Gold System to Sepolia Testnet

### Prerequisites
1. **Infura Account**: Free account at https://infura.io
2. **Sepolia ETH**: Get test ETH from https://sepoliafaucet.com
3. **Private Key**: From a wallet with Sepolia ETH

### Step 1: Get Infura Project ID
1. Go to https://infura.io/dashboard
2. Create a new project
3. Select "Ethereum" network
4. Copy the Project ID

### Step 2: Get Sepolia Test ETH
1. Go to https://sepoliafaucet.com
2. Connect your wallet
3. Request test ETH (0.5 ETH should be enough)

### Step 3: Set Environment Variables
Create a `.env` file in the project root:

```bash
# Copy from .env.testnet and fill in your values
INFURA_PROJECT_ID=your_infura_project_id
PRIVATE_KEY=your_private_key_without_0x
ETHERSCAN_API_KEY=your_etherscan_key_optional
```

### Step 4: Verify Environment Setup
```bash
npm run check-env
```

This will confirm all required variables are set correctly.

### Step 5: Deploy to Testnet
```bash
npm run deploy -- --network sepolia
```

### Expected Output
When successful, you'll see:
```
🚀 Deploying BTNG Gold System...
✅ Oracle deployed to: 0x...
✅ Gold Token deployed to: 0x...
✅ Custody deployed to: 0x...
🎉 BTNG Gold System deployed successfully!
```

### Step 5: Verify Contracts (Optional)
```bash
npx hardhat verify --network sepolia CONTRACT_ADDRESS
```

### Troubleshooting
- **"insufficient funds"**: Get more Sepolia ETH
- **"invalid project id"**: Check INFURA_PROJECT_ID
- **"network error"**: Check internet connection

### Security Notes
- Never commit `.env` file to version control
- Use a dedicated wallet for testnet deployments
- Keep private keys secure

---
*Ready for sovereign gold standard activation!*</content>
<parameter name="filePath">c:\BTNGAI_files\TESTNET_DEPLOYMENT_GUIDE.md