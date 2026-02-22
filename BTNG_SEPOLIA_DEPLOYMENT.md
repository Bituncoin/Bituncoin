# 🚀 BTNG Sepolia Testnet Deployment - Step-by-Step Setup

## 📋 Prerequisites Checklist

### ✅ Step 1: Install MetaMask Wallet
- [ ] Download from: https://metamask.io/download/
- [ ] Install browser extension
- [ ] Create new wallet or import existing
- [ ] Set strong password

### ✅ Step 2: Add Sepolia Testnet to MetaMask
1. Open MetaMask
2. Click network dropdown (shows "Ethereum Mainnet")
3. Click "Add Network"
4. Click "Add a network manually"
5. Enter these details:
   ```
   Network Name: Sepolia Testnet
   New RPC URL: https://rpc.sepolia.org
   Chain ID: 11155111
   Currency Symbol: SepoliaETH
   Block Explorer URL: https://sepolia.etherscan.io
   ```

### ✅ Step 3: Get Sepolia Test ETH
- [ ] Visit: https://sepoliafaucet.com
- [ ] Connect your MetaMask wallet
- [ ] Request 0.5 SepoliaETH (free, instant)
- [ ] Verify balance in MetaMask

### ✅ Step 4: Get Your Private Key
1. Open MetaMask
2. Click account icon (top right)
3. Click "Account Details"
4. Click "Show Private Key"
5. Enter your MetaMask password
6. Copy the private key (starts with 0x)

### ✅ Step 5: Get Infura Project ID
1. Visit: https://infura.io/dashboard
2. Sign up/Login (free)
3. Click "Create New Project"
4. Select "Ethereum" network
5. Copy the Project ID

### ✅ Step 6: Configure Environment
Update your `.env` file with real values:

```bash
# Replace demo values with your real credentials
INFURA_PROJECT_ID=your_actual_infura_project_id
PRIVATE_KEY=your_actual_private_key_without_0x_prefix
ETHERSCAN_API_KEY=your_etherscan_api_key_optional
```

### ✅ Step 7: Verify Setup
```bash
npm run check-env
```

### ✅ Step 8: Deploy to Testnet
```bash
npm run deploy:testnet
```

## 🔧 Troubleshooting

### "insufficient funds for gas"
**Solution:** Get more Sepolia ETH from https://sepoliafaucet.com

### "could not detect network"
**Solution:** Check INFURA_PROJECT_ID is correct

### "invalid private key"
**Solution:** Make sure private key doesn't have 0x prefix in .env

### "nonce too low"
**Solution:** Wait a few minutes, or reset account in MetaMask

## 📊 Expected Deployment Output

```
🚀 Deploying BTNG Gold System...
Deploying contracts with account: 0x...
Account balance: 500000000000000000
📡 Deploying BTNG Gold Oracle...
✅ Oracle deployed to: 0x...
🪙 Deploying BTNG Gold Token...
✅ Gold Token deployed to: 0x...
🏦 Deploying BTNG Custody...
✅ Custody deployed to: 0x...
🎉 BTNG Gold System deployed successfully!
```

## 🔗 Contract Addresses (After Deployment)

Save these addresses for your frontend integration:

- **BTNG Gold Token**: `0x...`
- **BTNG Custody**: `0x...`
- **BTNG Oracle**: `0x...`

## 🎯 Next Steps After Deployment

1. **Update Frontend**: Connect your demo to testnet contracts
2. **Test Transactions**: Mint/burn gold tokens on testnet
3. **Verify Contracts**: Use Etherscan to verify source code
4. **Mainnet Ready**: Same process works for Ethereum mainnet

---
*Ready to deploy BTNG to Sepolia testnet!* 🚀</content>
<parameter name="filePath">c:\BTNGAI_files\BTNG_SEPOLIA_DEPLOYMENT.md