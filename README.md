# BITUNCOIN Blockchain Ecosystem

🪙 **Gold-Coin (GLD)** - A next-generation cryptocurrency powered by Proof-of-Stake consensus

## Features

### 🌟 BTNg Wallet - Advanced Blockchain Operating System
- **Multi-Currency Support**: BTN, GLD, BTC, ETH, USDT, BNB
- **Merchant Payments**: BTN-Pay invoicing, QR codes, mobile money
- **Payment Cards**: Virtual & Physical MasterCard/Visa integration
- **AI Assistant**: Portfolio insights, security alerts, recommendations
- **Cross-Chain Swaps**: Seamless asset transfers with 0.5% fees
- **Dark/Light Mode**: Customizable UI themes
- **Multi-Language**: Support for English, French, Spanish, Chinese

### 🪙 Gold-Coin Cryptocurrency
- **Proof-of-Stake (PoS)** consensus mechanism for energy efficiency
- **100 Million GLD** maximum supply with 8 decimal precision
- **5% Annual Staking Rewards** for validators and stakers
- **Low Transaction Fees** (0.1% per transaction)
- **30-day Lock Period** for staking security

### 💼 BTNg Universal Wallet
- **Multi-Currency Support**: BTN, GLD, BTC, ETH, USDT, BNB
- **Real-time Balance Tracking**: Live updates for all assets
- **Cross-Chain Transactions**: Seamless asset transfers between blockchains
- **Modern UI**: Dark/Light mode with intuitive interface
- **Payment Integration**: Mobile money, QR codes, payment cards
- **AI-Powered Insights**: Portfolio analysis and recommendations

### 🔒 Security Features
- **AES-256 Encryption**: Military-grade wallet encryption
- **Two-Factor Authentication (2FA)**: Extra layer of account security
- **Biometric Login**: Fingerprint and face recognition support
- **Fraud Monitoring**: Real-time transaction monitoring and alerts
- **Encrypted Backups**: Secure wallet backup and recovery
- **Recovery Phrase**: 12-word mnemonic for wallet restoration

### 💳 Payment & Merchant Integration
- **BTN-Pay Cards**: Virtual & Physical MasterCard/Visa
- **Mobile Money**: MTN, AirtelTigo, Vodafone Cash integration
- **QR Code Payments**: Easy merchant payment solution
- **Invoice System**: Create and manage payment invoices
- **Crypto-to-Fiat**: Built-in exchange for multiple fiat currencies

### 🤖 AI-Powered Features
- **Portfolio Analysis**: Automated diversification insights
- **Security Alerts**: Real-time fraud detection and warnings
- **Staking Recommendations**: Optimize your earning potential
- **Market Insights**: Trend analysis and predictions
- **Smart Suggestions**: Context-aware wallet management tips

### ⛓️ Cross-Chain Bridge
- Support for Bitcoin, Ethereum, Binance Smart Chain, and Gold-Coin
- Seamless token swaps between supported networks
- Competitive 1% cross-chain transaction fees
- Real-time transaction status tracking

## Repository Structure

```
bituncoin-btn/
├── goldcoin/           # Gold-Coin cryptocurrency implementation
│   ├── goldcoin.go     # Core token logic
│   └── staking.go      # Staking pool management
├── consensus/          # Proof-of-Stake consensus
│   └── pos-validator.go # PoS validator logic
├── core/               # Blockchain core
│   └── btnchain.go     # Blockchain implementation
├── api/                # API server
│   └── btnnode.go      # Node API endpoints
├── wallet/             # BTNg Universal Wallet
│   ├── Wallet.jsx      # React wallet UI
│   ├── Wallet.css      # Wallet styling
│   ├── security.go     # Security features
│   ├── crosschain.go   # Cross-chain bridge
│   ├── exchange.go     # Exchange & swap
│   └── ai_assistant.go # AI-powered insights
├── payments/           # Payment integration
│   ├── btnpay.go       # BTN-Pay invoicing
│   ├── mobilemoney.go  # Mobile money integration
│   ├── qrcode.go       # QR code payments
│   └── card.go         # Card payment system
├── identity/           # Address management
│   └── btnaddress.go   # Address generation
├── storage/            # Data persistence
│   └── leveldb.go      # Key-value storage
├── config.yml          # Configuration file
├── DEPLOYMENT.md       # Deployment guide
└── README.md
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
go mod init github.com/Bituncoin/Bituncoin
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

## Tokenomics

| Parameter | Value |
|-----------|-------|
| Name | Gold-Coin |
| Symbol | GLD |
| Max Supply | 100,000,000 GLD |
| Decimals | 8 |
| Consensus | Proof-of-Stake |
| Block Time | 10 seconds |
| Staking Reward | 5% annual |
| Min Stake | 100 GLD |
| Min Validator Stake | 1,000 GLD |
| Transaction Fee | 0.1% |

## Proof-of-Stake Features

- **Energy Efficient**: PoS consumes 99.9% less energy than PoW
- **Scalable**: Higher transaction throughput
- **Secure**: Economic incentives align validator interests
- **Democratic**: Stake-weighted validator selection
- **Rewarding**: Earn 5% annual rewards for staking

## Wallet Features

### Overview Tab
- View balances across all supported cryptocurrencies
- Quick actions: Send, Receive, Swap, Stake
- Real-time USD conversion

### Staking Tab
- View staked amount and rewards
- Stake more GLD
- Claim accumulated rewards
- Unstake after lock period

### Transactions Tab
- Complete transaction history
- Filter by type (Sent, Received, Staked)
- Transaction status tracking

### Security Tab
- Enable/disable 2FA
- Configure biometric authentication
- Create encrypted backups
- Restore from backup
- View encryption status

## API Endpoints

### Core Endpoints
- `GET /api/info` - Node information
- `GET /api/health` - Health check
- `GET /api/goldcoin/balance?address=<addr>` - Get balance
- `POST /api/goldcoin/send` - Send transaction
- `POST /api/goldcoin/stake` - Stake tokens
- `GET /api/goldcoin/validators` - List validators

### Payment Endpoints
- `POST /api/btnpay/invoice` - Create invoice
- `GET /api/btnpay/invoice/{id}` - Get invoice
- `POST /api/btnpay/pay` - Pay invoice
- `POST /api/mobilemoney/pay` - Mobile money payment
- `POST /api/qrcode/generate` - Generate QR code
- `POST /api/card/issue` - Issue payment card
- `POST /api/card/payment` - Process card payment

### Exchange Endpoints
- `GET /api/exchange/rates` - Get exchange rates
- `POST /api/exchange/swap` - Execute swap
- `GET /api/exchange/estimate` - Estimate swap

### AI Assistant Endpoints
- `GET /api/ai/insights` - Get portfolio insights
- `GET /api/ai/alerts` - Get security alerts
- `POST /api/ai/ask` - Ask AI assistant

## Development

### Run Tests

```bash
# Run all tests
go test ./... -v

# Test specific module
go test ./goldcoin -v
go test ./consensus -v
go test ./wallet -v
```

### Build

```bash
# Build all binaries
go build -o bin/goldcoin ./goldcoin
go build -o bin/api-node ./api
go build -o bin/validator ./consensus

# Build wallet
cd wallet
npm run build
```

## Documentation

- [BTNg Wallet Guide](docs/BTNg-WALLET.md) - Complete wallet features and usage
- [BTN-PAY Integration](docs/BTN-PAY.md) - Payment protocol specification
- [Deployment Guide](DEPLOYMENT.md) - Complete deployment instructions
- [Configuration](config.yml) - Configuration options

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
