# BITUNCOIN - Universal Blockchain Wallet

![Bituncoin](https://img.shields.io/badge/Bituncoin-Universal%20Wallet-blue)
![License](https://img.shields.io/badge/license-GPL--3.0-green)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)

A comprehensive, secure, multi-currency blockchain wallet with cross-chain capabilities for the Bituncoin ecosystem.

## 🌟 Features

- **Multi-Currency Support**: BTN, BTC, ETH, USDT, BNB
- **Cross-Chain Transactions**: Seamless interoperability between blockchains
- **Enhanced Security**: 2FA, biometric authentication, encrypted storage
- **User-Friendly Interface**: Intuitive React-based wallet UI
- **Full Blockchain Integration**: Complete node with consensus mechanism

## 📁 Project Structure

```
bituncoin-btn/
├── core/              # Blockchain core
│   └── btnchain.go    # Block & transaction handling
├── api/               # REST API
│   └── btnnode.go     # Node endpoints
├── wallet/            # Frontend UI
│   └── Wallet.jsx     # React wallet component
├── identity/          # Identity management
│   └── btnaddress.go  # Address generation & signing
├── storage/           # Data persistence
│   └── leveldb.go     # LevelDB wrapper
├── consensus/         # Consensus mechanism
│   └── validator-echo.go
├── examples/          # Example scripts
│   └── test_api.sh    # API test examples
├── main.go            # Application entry point
├── go.mod             # Go dependencies
├── package.json       # NPM dependencies
└── README.md          # Detailed documentation
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 16+ and npm (optional, for wallet UI)

### 1. Build & Run Backend

```bash
# Clone the repository
git clone https://github.com/Bituncoin/Bituncoin.git
cd Bituncoin

# Use the quick start script
chmod +x start.sh
./start.sh

# Or build manually
cd bituncoin-btn
go build -o bituncoin-node main.go
./bituncoin-node
```

The node starts on `http://localhost:8080`

### 2. Setup Wallet UI (Optional)

```bash
# Create React app
npx create-react-app bituncoin-wallet-ui
cd bituncoin-wallet-ui

# Copy wallet component
cp ../bituncoin-btn/wallet/Wallet.jsx src/

# Update src/App.js to use Wallet component
# Then start the UI
npm start
```

Access wallet at `http://localhost:3000`

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/wallet/create` | Create new wallet |
| GET | `/api/wallet/balance` | Get balance |
| POST | `/api/transaction/send` | Send transaction |
| GET | `/api/transaction/history` | Transaction history |
| GET | `/api/blockchain/info` | Blockchain info |
| POST | `/api/mine` | Mine a block |
| GET | `/api/currencies` | Supported currencies |
| POST | `/api/crosschain/bridge` | Cross-chain bridge |

## 🧪 Testing

Run the included test script:
```bash
cd bituncoin-btn/examples
chmod +x test_api.sh
./test_api.sh
```

## 📖 Documentation

For detailed documentation, see [bituncoin-btn/README.md](bituncoin-btn/README.md)

## 🔐 Security Features

- ECDSA P256 key generation
- SHA256 address hashing
- Two-Factor Authentication (2FA)
- Biometric authentication support
- Encrypted private key storage
- Transaction signing & verification

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request

## 📄 License

GNU General Public License v3.0 - see [LICENSE](LICENSE)

## 🔗 Links

- **Repository**: https://github.com/Bituncoin/Bituncoin
- **Issues**: https://github.com/Bituncoin/Bituncoin/issues

---

**Made with ❤️ by the Bituncoin Team**
