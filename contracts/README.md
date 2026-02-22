# BTNG Gold Smart Contracts

## Overview

The BTNG Gold smart contract system implements a sovereign digital gold standard for the BTNG 54 Africa Gold Coin. This system provides:

- **Gold-Backed Token**: 1 BTNG = 1 gram of pure African gold
- **Custody Management**: Secure gold reserve management
- **Oracle Integration**: Real-time gold pricing and reserve verification
- **Sovereign Value**: Economic sovereignty for African nations

## Architecture

### Core Contracts

#### 1. BTNGGoldToken.sol
**ERC-20 token representing gold-backed value**
- Minting controlled by Custody contract
- Burning for physical gold redemption
- Gold backing verification per address
- Emergency pause functionality

#### 2. BTNGCustody.sol
**Gold reserve custody and operations management**
- Gold deposit/withdrawal tracking
- Minting and redemption request processing
- Authorized custodian management
- Oracle integration for verification

#### 3. BTNGGoldOracle.sol
**Real-time gold data and verification oracle**
- Gold price feeds (per gram in USD)
- Country-specific gold reserve tracking
- Sovereign value calculations
- Operation verification system

## Key Features

### Sovereign Gold Standard
```
1 BTNG = 1 Gram of Pure African Gold
```
- Backed by collective reserves of 54 African nations
- Sovereign value calculated per nation: `(Gold Reserves × 1,000,000) ÷ Population`

### Gold Backing Verification
- Every token backed by physical gold in custody
- Real-time verification of gold reserves
- Immutable audit trail on blockchain

### Multi-Role Security
- **Owner**: System administration
- **Custodians**: Gold reserve management
- **Feeders**: Oracle data updates
- **Users**: Token holders and transactors

## Installation & Setup

### Prerequisites
```bash
Node.js >= 18.0.0
npm or yarn
```

### Installation
```bash
cd contracts
npm install
```

### Compilation
```bash
npm run compile
```

### Testing
```bash
npm run test
```

### Deployment
```bash
# Local deployment
npm run deploy:localhost

# Mainnet deployment (configure network first)
npm run deploy
```

## Contract Interfaces

### BTNGGoldToken

```solidity
// Core functions
function mintGoldBacked(address to, uint256 amount, string batchId, uint256 goldGrams) external;
function redeemGold(uint256 amount, string redemptionId) external;
function transfer(address to, uint256 amount) public override returns (bool);
function getGoldBacking(address account) external view returns (uint256);
function verifyGoldBacking() external view returns (bool);
```

### BTNGCustody

```solidity
// Custody operations
function depositGold(string batchId, uint256 grams, string certificateHash) external;
function requestMint(uint256 amount, string batchId) external;
function requestRedemption(uint256 amount, string redemptionId, string deliveryAddress) external;
function executeMint(address to, uint256 amount, string batchId) external;
function executeRedemption(string redemptionId) external;
```

### BTNGGoldOracle

```solidity
// Oracle functions
function updateGoldPrice(uint256 newPrice) external;
function updateCountryReserve(string countryCode, uint256 goldTonnes) external;
function getGoldPrice() external view returns (uint256);
function calculateSovereignValue(string countryCode, uint256 population) external view returns (uint256);
function getTotalAfricanReserves() external view returns (uint256);
```

## Sovereign Value Calculation

The sovereign value represents each nation's economic contribution to the gold standard:

```
Sovereign BTNG Value = (Gold Reserves in Grams × 1,000,000) ÷ Population
```

### Example: Ghana
- Gold Reserves: 8.7 tonnes = 8,700,000 grams
- Population: 31,072,940
- Sovereign Value: (8,700,000 × 1,000,000) ÷ 31,072,940 = **0.00028 BTNG per person**

## Security Features

### Access Control
- **OpenZeppelin Ownable**: Single owner administration
- **Role-based Access**: Custodians, feeders, users
- **Multi-signature**: Critical operations require verification

### Safety Mechanisms
- **Pausable**: Emergency stop functionality
- **ReentrancyGuard**: Protection against reentrancy attacks
- **Input Validation**: Comprehensive parameter checking
- **Circuit Breakers**: Automatic halting on anomalies

### Audit Trail
- **Event Logging**: All operations emit events
- **IPFS Certificates**: Gold certificate storage
- **Immutable Records**: On-chain transaction history

## Testing

### Test Coverage
```bash
npm run coverage
```

### Test Scenarios
- ✅ Gold deposit and minting
- ✅ Token transfers with gold backing
- ✅ Redemption requests and execution
- ✅ Oracle price and reserve updates
- ✅ Sovereign value calculations
- ✅ Access control and security
- ✅ Emergency pause functionality

## Deployment

### Network Configuration
Update `hardhat.config.js` with your network settings:

```javascript
networks: {
  mainnet: {
    url: "https://mainnet.infura.io/v3/YOUR_PROJECT_ID",
    accounts: [process.env.PRIVATE_KEY]
  }
}
```

### Deployment Steps
1. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your private keys and API keys
   ```

2. **Deploy Contracts**
   ```bash
   npm run deploy
   ```

3. **Verify Contracts**
   ```bash
   npx hardhat verify --network mainnet CONTRACT_ADDRESS
   ```

## Integration

### Frontend Integration
```typescript
// Connect to contracts
const goldToken = new ethers.Contract(TOKEN_ADDRESS, BTNGGoldToken.abi, signer);
const custody = new ethers.Contract(CUSTODY_ADDRESS, BTNGCustody.abi, signer);
const oracle = new ethers.Contract(ORACLE_ADDRESS, BTNGGoldOracle.abi, provider);

// Get gold price
const price = await oracle.getGoldPrice();

// Mint gold-backed tokens
await custody.requestMint(amount, batchId);
```

### API Integration
```javascript
// Fetch sovereign values
const response = await fetch('/api/gold-coin');
const data = await response.json();

// Display country data
data.countries.forEach(country => {
  console.log(`${country.name}: ${country.sovereignValue} BTNG`);
});
```

## Economic Model

### Tokenomics
- **Max Supply**: 525,000,000 BTNG (525 tonnes African gold reserves)
- **Backing**: 100% physical gold reserves
- **Redemption**: 1:1 gold gram exchange
- **Fees**: 0% transfer fees, minimal minting fees

### Revenue Streams
- **Minting Fees**: Small percentage on new token creation
- **Custody Fees**: Annual storage fees for gold reserves
- **Oracle Services**: Premium data feeds for institutions

### Governance
- **Nation States**: Equal representation (1 nation = 1 vote)
- **Gold Holders**: Voting power proportional to holdings
- **Custodians**: Technical and operational decisions

## Risk Management

### Gold Reserve Risks
- **Custody Insurance**: Comprehensive coverage for physical gold
- **Geographic Diversification**: Multiple secure storage facilities
- **Regular Audits**: Independent verification of reserves

### Smart Contract Risks
- **Formal Verification**: Mathematical proof of contract correctness
- **Bug Bounties**: Community-driven security testing
- **Upgradeability**: Proxy patterns for contract updates

### Market Risks
- **Price Stabilization**: Algorithmic stabilization mechanisms
- **Liquidity Pools**: DEX integration for trading
- **Hedge Funds**: Institutional participation for stability

## Future Enhancements

### Phase 2: DeFi Integration
- **Yield Farming**: Stake BTNG for gold-backed yields
- **Lending Protocol**: Gold-backed loans and borrowing
- **DEX Integration**: Decentralized gold trading

### Phase 3: Cross-Chain Expansion
- **Bridge Contracts**: Multi-chain gold standard
- **Layer 2 Scaling**: High-throughput transactions
- **Cross-Chain Oracles**: Unified global pricing

### Phase 4: Institutional Adoption
- **Central Bank Integration**: Official reserve asset
- **Commercial Banking**: Gold-backed commercial loans
- **Insurance Products**: Gold-backed insurance policies

## Contributing

### Development Workflow
1. **Fork** the repository
2. **Create** a feature branch
3. **Write** comprehensive tests
4. **Implement** the feature
5. **Test** thoroughly
6. **Submit** a pull request

### Code Standards
- **Solidity Style Guide**: Follow official Solidity guidelines
- **Comprehensive Testing**: Minimum 95% test coverage
- **Security Audits**: All changes require security review

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For technical support or questions:
- **Documentation**: [BTNG Gold Docs](https://docs.btng.africa)
- **Discord**: [BTNG Community](https://discord.gg/btng)
- **Email**: support@btng.africa

## Disclaimer

This system represents a sovereign digital gold standard backed by African gold reserves. Users should understand the risks associated with cryptocurrency and gold investments. Always conduct your own research and consult with financial advisors before participating.

---

**BTNG Gold - Sovereign Digital Gold for Africa** 🌍✨

*Building economic sovereignty through blockchain technology*