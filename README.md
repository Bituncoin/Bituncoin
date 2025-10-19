# BITUNCOIN Blockchain Ecosystem

🌟 **Universal Blockchain Wallet** - A comprehensive cryptocurrency management platform for BITUNCOIN

## Overview

BITUNCOIN is a next-generation blockchain ecosystem featuring a universal wallet with multi-currency support, payment cards, AI-driven insights, and gold-backed cryptocurrency. The platform combines cutting-edge blockchain technology with user-friendly interfaces and advanced security features.

## Key Features

### 💼 Universal Wallet
- **Multi-Currency Support**: BTN, BTC, ETH, USDT, BNB, GLD
- **Real-time Portfolio Tracking**: Live balance updates and USD conversion
- **Transaction Management**: Complete history across all currencies
- **Total Portfolio Value**: Aggregate view of all assets

### 💳 BTN-Pay Payment Cards
- **Virtual Cards**: Instant issuance for online purchases
- **Physical Cards**: Visa/MasterCard delivered to your address
- **Spending Controls**: Customizable daily and monthly limits
- **Real-time Transactions**: Instant payment processing
- **NFC Support**: Tap-to-pay functionality

### 💱 Cryptocurrency Exchange
- **Crypto-to-Crypto**: Trade between all supported currencies
- **Crypto-to-Fiat**: Convert to USD, EUR, GBP
- **Live Rates**: Exchange rates updated every second
- **Low Fees**: 0.1% per transaction
- **Instant Conversion**: No waiting period

### 🏪 Merchant Services
- **QR Code Payments**: Generate and scan for instant payments
- **NFC Payments**: Contactless payment support
- **Mobile Money**: MTN and Airtel integration
- **Payment Requests**: Create and manage merchant invoices
- **Real-time Notifications**: Instant payment confirmations

### 🤖 AI-Driven Insights
- **Market Analysis**: Real-time trend detection and predictions
- **Portfolio Insights**: Diversification and risk assessment
- **Smart Recommendations**: Staking opportunities and rebalancing
- **Performance Tracking**: Historical analysis and projections
- **Automated Alerts**: Price changes and security notifications

### 🔒 Advanced Security
- **Multi-Layer Protection**: 2FA and biometric authentication
- **AES-256 Encryption**: Military-grade wallet encryption
- **Fraud Detection**: Real-time AI-powered fraud prevention
- **Device Management**: Trusted device registration and monitoring
- **Transaction Approval**: Multi-signature and approval workflows

### 🏅 Gold Reserve Backing
- **$2.689 Trillion Reserve**: Every BTN backed by physical gold
- **100% Backing Ratio**: Full reserve maintained at all times
- **Runtime Verification**: Hourly reserve checks
- **Validator Consensus**: Multiple validators verify reserves
- **Transparent Audits**: Public audit trails and records

### 🌟 Gold-Coin (GLD) Cryptocurrency
### 🌟 Gold-Coin (GLD) Cryptocurrency
- **Proof-of-Stake (PoS)** consensus mechanism for energy efficiency
- **100 Million GLD** maximum supply with 8 decimal precision
- **5% Annual Staking Rewards** for validators and stakers
- **Low Transaction Fees** (0.1% per transaction)
- **30-day Lock Period** for staking security

## Supported Cryptocurrencies

| Currency | Symbol | Features |
|----------|--------|----------|
| Bituncoin | BTN | Gold-backed, primary currency |
| Bitcoin | BTC | Original cryptocurrency |
| Ethereum | ETH | Smart contract platform |
| Tether | USDT | Stablecoin (USD-pegged) |
| Binance Coin | BNB | Exchange token |
| Gold-Coin | GLD | PoS cryptocurrency |

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Universal Wallet (React UI)                │
│  Portfolio • Cards • Exchange • Insights • Security     │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   API Gateway                            │
│            RESTful API & WebSocket                      │
└────────────────────┬────────────────────────────────────┘
                     │
       ┌─────────────┼─────────────┬─────────────┐
       │             │              │             │
       ▼             ▼              ▼             ▼
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Wallet   │  │ Cards    │  │ Exchange │  │ Merchant │
│ Manager  │  │ Manager  │  │ System   │  │ Services │
└──────────┘  └──────────┘  └──────────┘  └──────────┘
       │             │              │             │
       └─────────────┼──────────────┴─────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│           Security & AI Services                        │
│  Fraud Detection • Insights • Gold Reserve Verification│
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│            Blockchain & Storage Layer                   │
│     PoS Consensus • LevelDB • Transaction Pool         │
└─────────────────────────────────────────────────────────┘
```

## Repository Structure

```
Bituncoin/
├── wallet/                  # Universal Wallet (Backend & Frontend)
│   ├── universal.go         # Multi-currency wallet manager
│   ├── cards.go            # Payment card management
│   ├── exchange.go         # Cryptocurrency exchange
│   ├── merchant.go         # Merchant payment services
│   ├── ai_insights.go      # AI-driven insights engine
│   ├── advanced_security.go # Fraud detection & security
│   ├── gold_reserve.go     # Gold reserve integration
│   ├── security.go         # Core security features
│   ├── crosschain.go       # Cross-chain bridge
│   ├── Wallet.jsx          # React wallet UI
│   ├── Wallet.css          # Wallet styling
│   └── package.json        # NPM dependencies
├── goldcoin/               # Gold-Coin cryptocurrency
│   ├── goldcoin.go         # Core token logic
│   └── staking.go          # Staking pool management
├── consensus/              # Proof-of-Stake consensus
│   └── pos-validator.go    # PoS validator logic
├── core/                   # Blockchain core
│   └── btnchain.go         # Blockchain implementation
├── api/                    # API server
│   └── btnnode.go          # Node API endpoints
├── identity/               # Address management
│   └── btnaddress.go       # Address generation
├── storage/                # Data persistence
│   └── leveldb.go          # Key-value storage
├── payments/               # Payment processing
│   └── btnpay.go           # BTN-PAY service
├── docs/                   # Documentation
│   ├── UNIVERSAL_WALLET_GUIDE.md  # User guide
│   ├── API_DOCUMENTATION.md       # API reference
│   ├── BTN-PAY.md                 # Payment protocol
│   └── DEPLOYMENT.md              # Deployment guide
├── examples/               # Example applications
│   └── demo.go             # Demo program
├── config.yml              # Configuration file
└── README.md               # This file
```

## Quick Start

### Prerequisites
- Go 1.18+
- Node.js 16+
- npm or yarn

### Installation

```bash
# Clone the repository
git clone https://github.com/Bituncoin/Bituncoin.git
cd Bituncoin

# Install Go dependencies
go mod tidy

# Install wallet dependencies
cd wallet
npm install
```

### Running the Node

```bash
# Start the API node
go run api/btnnode.go
```

### Running the Wallet

```bash
# Start the wallet UI
cd wallet
npm start
```

The wallet will be available at `http://localhost:3000`

## Wallet Features

### Overview Tab
- View balances across all 6 supported cryptocurrencies
- Total portfolio value in USD
- Quick actions: Send, Receive, Exchange, Stake, Cards, Pay Merchant

### Cards Tab
- Create virtual and physical payment cards
- Manage card status (activate, suspend, cancel)
- View card transactions and spending limits
- Set daily and monthly spending caps

### Exchange Tab
- Real-time exchange rates
- Crypto-to-crypto and crypto-to-fiat conversion
- Calculate exchange amounts before executing
- Low fee structure (0.1%)

### Staking Tab
- View staked amount and rewards
- Stake more tokens
- Claim accumulated rewards
- Unstake after lock period
- 5% APY on staked BTN/GLD

### Transactions Tab
- Complete transaction history
- Filter by type (Send, Receive, Stake, Exchange, Card)
- Transaction status tracking
- Export transaction history

### AI Insights Tab
- Market trend alerts and predictions
- Portfolio diversification analysis
- Risk assessment and recommendations
- Staking opportunities
- Gold reserve verification status

### Security Tab
- Enable/disable 2FA
- Configure biometric authentication
- Manage trusted devices
- View security events and alerts
- Real-time fraud detection status
- Create encrypted backups
- Restore from backup

## API Endpoints

### Wallet Management
- `GET /api/v1/wallet/balance` - Get wallet balance
- `GET /api/v1/wallet/portfolio` - Get portfolio summary
- `POST /api/v1/wallet/transaction` - Create transaction
- `GET /api/v1/wallet/transactions` - Get transaction history

### Payment Cards
- `POST /api/v1/cards/create` - Create payment card
- `POST /api/v1/cards/{id}/activate` - Activate card
- `GET /api/v1/cards/{id}` - Get card details
- `GET /api/v1/cards/{id}/transactions` - Get card transactions

### Exchange
- `GET /api/v1/exchange/rate` - Get exchange rate
- `POST /api/v1/exchange/calculate` - Calculate exchange
- `POST /api/v1/exchange/order` - Execute exchange

### Merchant Services
- `POST /api/v1/merchant/register` - Register merchant
- `POST /api/v1/merchant/payment-request` - Create payment request
- `POST /api/v1/merchant/process-payment` - Process payment

### AI Insights
- `GET /api/v1/insights/portfolio` - Get portfolio insights
- `GET /api/v1/insights/recommendations` - Get recommendations
- `GET /api/v1/insights/market-trend` - Get market trends

### Gold Reserve
- `GET /api/v1/gold-reserve/info` - Get reserve information
- `GET /api/v1/gold-reserve/backing-proof` - Get backing proof

### Legacy Endpoints (Gold-Coin)
- `GET /api/info` - Node information
- `GET /api/health` - Health check
- `GET /api/goldcoin/balance?address=<addr>` - Get GLD balance
- `POST /api/goldcoin/send` - Send GLD transaction
- `POST /api/goldcoin/stake` - Stake tokens
- `GET /api/goldcoin/validators` - List validators

For complete API documentation, see [API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md)

## Development

### Run Tests

```bash
# Run all tests
go test ./... -v

# Test specific module
go test ./wallet -v
go test ./goldcoin -v
go test ./consensus -v

# Run with coverage
go test ./wallet -cover
```

### Build

```bash
# Build backend binaries
go build -o bin/api-node ./api
go build -o bin/goldcoin ./goldcoin

# Build wallet frontend
cd wallet
npm run build
```

### Code Structure

The codebase follows Go best practices with clear separation of concerns:

- **wallet/**: Universal wallet implementation with all core services
- **goldcoin/**: Gold-Coin cryptocurrency implementation
- **consensus/**: Proof-of-Stake consensus mechanism
- **api/**: RESTful API server
- **docs/**: Comprehensive documentation

## Testing

### Test Coverage

- **Wallet Tests**: 24 tests covering all wallet functionality
- **Card Tests**: Complete test suite for payment cards
- **Gold-Coin Tests**: 17 tests for token operations
- **Consensus Tests**: 11 tests for PoS validator logic
- **Total**: 50+ tests with comprehensive coverage

### Running Specific Tests

```bash
# Universal wallet tests
go test ./wallet -v -run TestUniversal

# Payment card tests
go test ./wallet -v -run TestCard

# Exchange tests
go test ./wallet -v -run TestExchange
```

## Documentation

- [Universal Wallet User Guide](docs/UNIVERSAL_WALLET_GUIDE.md) - Complete user documentation
- [API Documentation](docs/API_DOCUMENTATION.md) - API reference and examples
- [Deployment Guide](DEPLOYMENT.md) - Deployment instructions
- [BTN-PAY Protocol](docs/BTN-PAY.md) - Payment protocol specification
- [Configuration](config.yml) - System configuration options

## Key Technologies

- **Backend**: Go 1.18+
- **Frontend**: React 18, JavaScript/JSX
- **Blockchain**: Custom implementation with PoS consensus
- **Storage**: LevelDB for persistence
- **Encryption**: AES-256 for wallet security
- **API**: RESTful HTTP + WebSocket
- **Testing**: Go testing framework

## Performance

- **Transaction Speed**: ~10 seconds per block
- **Throughput**: Scalable to 10,000+ concurrent users
- **Energy Efficiency**: 99.9% less than Proof-of-Work
- **Consensus**: Stake-weighted validator selection
- **Availability**: 99.9% uptime target

## Security

### Multi-Layer Security
1. **Wallet Encryption**: AES-256 encryption for all wallet data
2. **Two-Factor Authentication**: SMS, email, or authenticator app
3. **Biometric Login**: Fingerprint and face recognition
4. **Fraud Detection**: Real-time AI-powered fraud prevention
5. **Device Fingerprinting**: Track and manage trusted devices
6. **Transaction Approval**: Multi-signature support for high-value transactions

### Compliance
- KYC/AML procedures implemented
- GDPR compliant data handling
- Regular security audits
- Encrypted data transmission (TLS 1.3)

## Gold Reserve

The BITUNCOIN ecosystem is backed by **$2.689 trillion in physical gold reserves**:

- **100% Backing**: Every BTN is fully backed by gold
- **Runtime Verification**: Hourly reserve checks
- **Validator Consensus**: Multiple independent validators
- **Transparent Audits**: Public audit trails
- **Reserve Per BTN**: Dynamically calculated based on circulation

This makes BTN one of the most stable and valuable cryptocurrencies in the market.

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: https://github.com/Bituncoin/Bituncoin/issues
- **Documentation**: Coming soon
- **Community**: Join our Discord

---

Built with ❤️ by the Bituncoin Team
