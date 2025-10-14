# Bituncoin Universal Wallet

A comprehensive, secure, multi-currency blockchain wallet for the Bituncoin ecosystem.

## 🌟 Features

### 1. Multi-Currency Support
- Support for multiple cryptocurrencies: BTN (Bituncoin), BTC, ETH, USDT, BNB
- Unified interface to manage all your digital assets
- Real-time balance tracking across all currencies

### 2. Cross-Chain Transactions
- Seamless transactions between different blockchain networks
- Bridge transactions between supported chains
- Automated cross-chain settlement

### 3. Enhanced Security
- **Two-Factor Authentication (2FA)**: Additional layer of security for transactions
- **Biometric Authentication**: Fingerprint/face recognition support
- **Encrypted Key Storage**: Private keys stored with strong encryption
- **Secure Identity Management**: ECDSA key pair generation
- **Transaction Signing**: Cryptographic signing for all transactions

### 4. User-Friendly Interface
- Intuitive React-based web interface
- Easy wallet creation and management
- Clear transaction history
- Multi-tab navigation for different functions

### 5. Bituncoin Blockchain Integration
- Full integration with Bituncoin blockchain
- Consensus mechanism (Validator Echo)
- LevelDB-based persistent storage
- RESTful API for all operations

## 🏗️ Architecture

```
bituncoin-btn/
├── core/                 # Blockchain core functionality
│   └── btnchain.go      # Block and transaction handling
├── api/                  # REST API endpoints
│   └── btnnode.go       # Node API implementation
├── wallet/               # Frontend wallet UI
│   └── Wallet.jsx       # React wallet component
├── identity/             # Identity management
│   └── btnaddress.go    # Address generation and signing
├── storage/              # Data persistence
│   └── leveldb.go       # LevelDB wrapper
├── consensus/            # Consensus mechanism
│   └── validator-echo.go # Validator consensus
├── main.go              # Application entry point
└── go.mod               # Go module definition
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- Node.js 16+ and npm (for React frontend)

### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/Bituncoin/Bituncoin.git
cd Bituncoin/bituncoin-btn
```

2. Build and run the backend:
```bash
go mod download
go build -o bituncoin-node
./bituncoin-node
```

The node will start on `http://localhost:8080`

### Frontend Setup

1. Create a React app (if not already created):
```bash
npx create-react-app bituncoin-wallet-ui
cd bituncoin-wallet-ui
```

2. Copy the Wallet.jsx component:
```bash
cp ../bituncoin-btn/wallet/Wallet.jsx src/
```

3. Update src/App.js to use the Wallet component:
```javascript
import Wallet from './Wallet';

function App() {
  return <Wallet />;
}

export default App;
```

4. Start the frontend:
```bash
npm start
```

The wallet UI will be available at `http://localhost:3000`

## 📡 API Documentation

### Create Wallet
```
POST /api/wallet/create
Content-Type: application/json

{
  "enable_2fa": true,
  "enable_biometric": false,
  "biometric_data": ""
}
```

### Get Balance
```
GET /api/wallet/balance?address=BTN...
```

### Send Transaction
```
POST /api/transaction/send
Content-Type: application/json

{
  "from": "BTN...",
  "to": "BTN...",
  "amount": 10.5,
  "currency": "BTN",
  "cross_chain": false,
  "target_chain": ""
}
```

### Get Transaction History
```
GET /api/transaction/history?address=BTN...
```

### Get Blockchain Info
```
GET /api/blockchain/info
```

### Mine Block
```
POST /api/mine?miner=BTN...
```

### Get Supported Currencies
```
GET /api/currencies
```

### Cross-Chain Bridge
```
POST /api/crosschain/bridge
Content-Type: application/json

{
  "from_chain": "BTN",
  "to_chain": "ETH",
  "amount": 5.0
}
```

## 🔐 Security Features

### Key Management
- ECDSA P256 curve for key generation
- SHA256 hashing for addresses
- Secure private key storage

### Authentication
- Two-Factor Authentication (2FA) support
- Biometric authentication capability
- Transaction signing verification

### Encryption
- Private keys encrypted at rest
- Secure communication protocols
- Protected API endpoints

## 🧪 Testing

### Test Wallet Creation
```bash
curl -X POST http://localhost:8080/api/wallet/create \
  -H "Content-Type: application/json" \
  -d '{"enable_2fa": true, "enable_biometric": false}'
```

### Test Transaction
```bash
curl -X POST http://localhost:8080/api/transaction/send \
  -H "Content-Type: application/json" \
  -d '{
    "from": "BTN...",
    "to": "BTN...",
    "amount": 10,
    "currency": "BTN",
    "cross_chain": false
  }'
```

### Check Balance
```bash
curl http://localhost:8080/api/wallet/balance?address=BTN...
```

## 🌐 Cross-Chain Support

The wallet supports cross-chain transactions between:
- BTN (Bituncoin native)
- BTC (Bitcoin)
- ETH (Ethereum)
- USDT (Tether)
- BNB (Binance Coin)

Cross-chain transactions are handled through the bridge mechanism, ensuring secure and reliable transfers between different blockchain networks.

## 📊 Consensus Mechanism

The Bituncoin network uses a **Validator Echo** consensus mechanism:
- Validators stake BTN tokens to participate
- Minimum stake requirement: 1000 BTN
- Quorum-based validation
- Stake-weighted voting

## 💾 Storage

LevelDB is used for persistent storage:
- Wallet data
- Transaction history
- Blockchain state
- User preferences

Data is stored in `~/.bituncoin/data/`

## 🛠️ Development

### Running Tests
```bash
go test ./...
```

### Code Structure
- **Core**: Blockchain logic and transaction processing
- **API**: HTTP endpoints and request handling
- **Identity**: Key management and authentication
- **Storage**: Data persistence layer
- **Consensus**: Validator and block validation
- **Wallet**: User interface components

## 🤝 Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.

## 🔗 Links

- GitHub: https://github.com/Bituncoin/Bituncoin
- Documentation: See this README
- Issues: https://github.com/Bituncoin/Bituncoin/issues

## 🎯 Roadmap

- [ ] Mobile wallet application (iOS/Android)
- [ ] Hardware wallet integration
- [ ] Advanced multi-sig support
- [ ] DeFi protocol integration
- [ ] NFT support
- [ ] Staking mechanism
- [ ] Governance voting
- [ ] Enhanced cross-chain bridges

## 📞 Support

For support and questions:
- Open an issue on GitHub
- Check the documentation
- Review API examples

---

**Bituncoin Universal Wallet** - Secure, Fast, Reliable Multi-Currency Blockchain Wallet
