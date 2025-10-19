# BITUNCOIN Blockchain Ecosystem

🪙 **Gold-Coin (GLD)** - A next-generation cryptocurrency powered by Proof-of-Stake consensus

## Features

### 🌟 Gold-Coin Cryptocurrency
- **Proof-of-Stake (PoS)** consensus mechanism for energy efficiency
- **100 Million GLD** maximum supply with 8 decimal precision
- **5% Annual Staking Rewards** for validators and stakers
- **Low Transaction Fees** (0.1% per transaction)
- **30-day Lock Period** for staking security

### 👥 User and Admin Accounts
- **Role-Based Access Control (RBAC)**: User, Admin, Merchant, and Validator roles
- **Secure Authentication**: Password hashing, session management, and 2FA support
- **Permission System**: Granular permissions for different user roles
- **User Management**: Create, update, deactivate user accounts
- **Session Management**: Secure 24-hour sessions with automatic expiry

### 💼 Comprehensive Universal Wallet
- **Multi-Currency Support**: BTN, GLD, BTC, ETH, USDT, BNB
- **Real-time Portfolio Tracking**: Live balance updates and performance metrics
- **Cross-Chain Transactions**: Seamless asset transfers between blockchains
- **Transaction History**: Complete tracking with filtering and search
- **Modern UI**: Intuitive and user-friendly interface across all platforms

### 🔄 Built-in Cryptocurrency Exchange
- **Crypto-to-Crypto Exchange**: Swap between supported cryptocurrencies
- **Crypto-to-Fiat Exchange**: Convert to USD, EUR, GBP, and more
- **Competitive Fees**: 0.3-1% depending on pair
- **Real-time Rates**: Live exchange rates from multiple sources
- **Exchange History**: Track all your trades

### 💳 BTN-Pay Payment Cards
- **Virtual Cards**: Instant creation for online payments
- **Physical Cards**: MasterCard and Visa support
- **Real-time Transaction Tracking**: Monitor all card transactions
- **Daily Spending Limits**: Customizable for security
- **Instant Top-up**: Fund cards directly from wallet

### 🏪 Merchant Services
- **QR Code Payments**: Simple scan-to-pay functionality
- **NFC Payments**: Contactless payment support
- **Direct Wallet Transfers**: Peer-to-peer transactions
- **Mobile Money Integration**: MTN, Vodafone, Airtel, Tigo
- **Invoice System**: Create and manage payment requests
- **Merchant Dashboard**: Track payments and transactions

### 🤖 AI-Driven Wallet Management
- **Spending Insights**: Analyze your spending patterns
- **Market Alerts**: Real-time notifications on price changes
- **Trading Recommendations**: AI-powered buy/sell suggestions
- **Staking Optimization**: Maximize your staking rewards
- **Portfolio Optimization**: Suggestions to improve diversification

### 🔒 Advanced Security Features
- **AES-256 Encryption**: Military-grade wallet encryption
- **Two-Factor Authentication (2FA)**: Extra layer of account security
- **Biometric Login**: Fingerprint and face recognition support
- **Fraud Detection**: Real-time monitoring for suspicious activity
- **Security Alerts**: Instant notifications for security events
- **Encrypted Backups**: Secure wallet backup and recovery
- **Recovery Phrase**: 12-word mnemonic for wallet restoration

### 📊 Unified Operations Dashboard
- **System Monitoring**: Real-time health checks
- **Network Status**: Blockchain network connectivity
- **Performance Metrics**: Transaction volume, active users, uptime
- **Alert Management**: Centralized system alerts
- **Update Management**: Schedule and deploy updates
- **Admin Dashboard**: User management, token configuration, system settings

### 🔌 Add-On Module System
- **Plug-and-Play Architecture**: Easily extend wallet functionality
- **Module Categories**: DeFi, Staking, Lending, Trading, Payment, Analytics, Security, Utility
- **Built-in Modules**:
  - Advanced Staking: Multiple pools, auto-compounding, flexible lock periods
  - DeFi Lending: Collateral-based lending and borrowing
- **Module Management**: Enable/disable modules via API
- **Developer-Friendly**: Simple interface for creating custom modules

### 🌍 Multi-Platform Support
- **iOS**: Native iPhone and iPad app (React Native)
- **Android**: Native Android app (React Native)
- **Windows**: Desktop application (Electron)
- **macOS**: Desktop application (Electron)
- **Linux**: Desktop application (Electron) - AppImage, DEB, RPM
- **Web**: Responsive web interface (React)

### 🚀 CI/CD Automation
- **GitHub Actions**: Automated testing on every commit
- **Multi-Platform Builds**: Automatic builds for all platforms
- **Automated Releases**: Tag-based releases with artifacts
- **Code Quality**: Automated linting and testing
- **Docker Support**: Containerized deployment options

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
├── auth/               # Authentication & authorization
│   ├── accounts.go     # User account management
│   └── accounts_test.go # Auth tests
├── addons/             # Add-on module system
│   ├── registry.go     # Module registry
│   ├── staking_module.go # Advanced staking module
│   ├── lending_module.go # DeFi lending module
│   └── addons_test.go  # Module tests
├── wallet/             # Comprehensive wallet system
│   ├── Wallet.jsx      # React wallet UI
│   ├── Wallet.css      # Wallet styling
│   ├── portfolio.go    # Portfolio management
│   ├── transactions.go # Transaction history
│   ├── exchange.go     # Cryptocurrency exchange
│   ├── cards.go        # Payment card system
│   ├── merchant.go     # Merchant services
│   ├── platform.go     # Platform detection
│   ├── ai_manager.go   # AI-driven insights
│   ├── dashboard.go    # Operations dashboard
│   ├── security.go     # Security & fraud detection
│   └── crosschain.go   # Cross-chain bridge
├── payments/           # Payment protocols
│   └── btnpay.go       # BTN-Pay implementation
├── identity/           # Address management
│   └── btnaddress.go   # Address generation
├── storage/            # Data persistence
│   └── leveldb.go      # Key-value storage
├── .github/workflows/  # CI/CD pipelines
│   ├── test.yml        # Automated testing
│   └── build.yml       # Multi-platform builds
├── docs/               # Documentation
│   ├── BTN-PAY.md      # Payment protocol docs
│   ├── DEVELOPER_GUIDE.md # Developer documentation
│   ├── USER_GUIDE.md   # End-user guide
│   ├── ADMIN_GUIDE.md  # Administrator guide
│   ├── MODULE_DEVELOPER_GUIDE.md # Module development
│   ├── PLATFORM_DEPLOYMENT.md # Platform deployment
│   └── LAUNCH_STRATEGY.md # Launch plan
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
- Real-time USD conversion
- Portfolio performance metrics
- Quick actions: Send, Receive, Swap, Stake

### Exchange Tab
- Crypto-to-crypto swaps
- Crypto-to-fiat conversion
- Live exchange rates
- Transaction history

### Cards Tab
- Create virtual/physical cards
- View card transactions
- Top-up cards
- Manage spending limits

### Merchant Tab
- Register as merchant
- Accept payments via QR/NFC
- Mobile money integration
- Payment tracking

### Staking Tab
- View staked amount and rewards
- Stake more GLD
- Claim accumulated rewards
- Unstake after lock period

### Transactions Tab
- Complete transaction history
- Filter by type (Sent, Received, Staked, Swapped)
- Transaction status tracking
- Export capabilities

### Security Tab
- Enable/disable 2FA
- Configure biometric authentication
- Create encrypted backups
- Restore from backup
- View encryption status
- Security alerts

### Insights Tab
- AI-powered spending analysis
- Market trend alerts
- Trading recommendations
- Staking optimization
- Portfolio suggestions

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `GET /api/auth/validate` - Validate session

### User Management (Admin)
- `GET /api/users/list` - List all users
- `POST /api/users/update-role` - Update user role
- `POST /api/users/deactivate` - Deactivate user

### Wallet & Portfolio
- `GET /api/wallet/portfolio/:address` - Get portfolio
- `POST /api/wallet/portfolio/add` - Add asset
- `GET /api/wallet/portfolio/:address/performance` - Get performance

### Exchange
- `POST /api/exchange/quote` - Get exchange quote
- `POST /api/exchange/order` - Create exchange order
- `GET /api/exchange/orders/:address` - Get user orders

### Cards
- `POST /api/cards/create` - Create payment card
- `POST /api/cards/:cardId/topup` - Top-up card
- `GET /api/cards/:address` - Get user cards
- `GET /api/cards/:cardId/transactions` - Get card transactions

### Merchant
- `POST /api/merchant/register` - Register merchant
- `POST /api/merchant/payment/request` - Create payment request
- `POST /api/merchant/payment/complete` - Complete payment
- `GET /api/merchant/:merchantId/payments` - Get payments

### Gold-Coin
- `GET /api/goldcoin/balance?address=<addr>` - Get balance
- `POST /api/goldcoin/send` - Send transaction
- `POST /api/goldcoin/stake` - Stake tokens
- `GET /api/goldcoin/validators` - List validators

### System
- `GET /api/info` - Node information
- `GET /api/health` - Health check
- `GET /api/dashboard/status` - System status
- `GET /api/dashboard/metrics` - System metrics

### Add-On Modules
- `GET /api/addons/list` - List available modules
- `POST /api/addons/enable` - Enable a module
- `POST /api/addons/disable` - Disable a module
- `POST /api/addons/execute` - Execute module action

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

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o bin/api-node-linux ./api
GOOS=windows GOARCH=amd64 go build -o bin/api-node-windows.exe ./api
GOOS=darwin GOARCH=amd64 go build -o bin/api-node-macos ./api
```

### Multi-Platform Builds

See [Platform Deployment Guide](docs/PLATFORM_DEPLOYMENT.md) for detailed instructions on building for:
- Web (React build)
- Desktop (Electron for Windows, macOS, Linux)
- Mobile (React Native for iOS and Android)
- Docker containers


## Documentation

- [Platform Deployment Guide](docs/PLATFORM_DEPLOYMENT.md) - Multi-platform deployment instructions
- [Admin Guide](docs/ADMIN_GUIDE.md) - Administrator documentation
- [Module Developer Guide](docs/MODULE_DEVELOPER_GUIDE.md) - Creating add-on modules
- [Deployment Guide](DEPLOYMENT.md) - Complete deployment instructions
- [Developer Guide](docs/DEVELOPER_GUIDE.md) - Developer documentation
- [User Guide](docs/USER_GUIDE.md) - End-user documentation
- [Launch Strategy](docs/LAUNCH_STRATEGY.md) - Launch and marketing plan
- [BTN-Pay Protocol](docs/BTN-PAY.md) - Payment protocol specification
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
