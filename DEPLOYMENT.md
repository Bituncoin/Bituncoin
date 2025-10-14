# Gold-Coin Deployment Guide

## Prerequisites

- Go 1.18 or higher
- Node.js 16+ and npm
- Git

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Bituncoin/Bituncoin.git
cd Bituncoin
```

### 2. Install Go Dependencies

```bash
go mod init github.com/Bituncoin/Bituncoin
go mod tidy
```

### 3. Install Frontend Dependencies

```bash
cd wallet
npm install react react-dom
npm install --save-dev @types/react @types/react-dom
```

## Configuration

Edit `config.yml` to customize your Gold-Coin deployment:

```yaml
goldcoin:
  name: "Gold-Coin"
  symbol: "GLD"
  max_supply: 100000000
```

## Building

### Backend (Go)

```bash
# Build all modules
go build -o bin/goldcoin ./goldcoin
go build -o bin/api-node ./api
go build -o bin/validator ./consensus
```

### Frontend (React)

```bash
cd wallet
npm run build
```

## Deployment Steps

### 1. Initialize the Blockchain

```bash
./bin/goldcoin init --network mainnet
```

### 2. Start the API Node

```bash
./bin/api-node start --port 8080
```

### 3. Register as a Validator

```bash
./bin/validator register --stake 1000 --address <your-address>
```

### 4. Start the Wallet UI

```bash
cd wallet
npm start
```

## Testing

### Unit Tests

```bash
# Test Gold-Coin core
go test ./goldcoin -v

# Test Proof-of-Stake consensus
go test ./consensus -v

# Test wallet functionality
go test ./wallet -v
```

### Integration Tests

```bash
# Run all tests
go test ./... -v
```

## Network Launch

### Testnet Launch

1. Deploy to testnet environment
2. Distribute test tokens to early adopters
3. Run security audits
4. Gather feedback and fix issues

### Mainnet Launch

1. Complete security audits
2. Freeze smart contracts
3. Initialize genesis block
4. Distribute initial token allocation
5. Enable public validator registration
6. Launch wallet interface

## Monitoring

### Check Node Status

```bash
curl http://localhost:8080/api/info
```

### View Validators

```bash
curl http://localhost:8080/api/goldcoin/validators
```

### Check Health

```bash
curl http://localhost:8080/api/health
```

## Security Best Practices

1. **Private Keys**: Never share or expose private keys
2. **2FA**: Enable two-factor authentication for all accounts
3. **Backups**: Regular encrypted backups of wallet data
4. **Updates**: Keep all software up to date
5. **Audits**: Regular security audits before major releases

## Troubleshooting

### Node Won't Start

- Check if port 8080 is available
- Verify configuration in `config.yml`
- Check logs in `./logs/`

### Wallet Connection Issues

- Ensure API node is running
- Check network connectivity
- Verify RPC endpoint in configuration

### Staking Issues

- Verify minimum stake requirement (1000 GLD for validators)
- Check lock period hasn't expired
- Ensure sufficient balance

## Support

- Documentation: https://docs.goldcoin.bituncoin.io
- Issues: https://github.com/Bituncoin/Bituncoin/issues
- Community: https://discord.gg/goldcoin

## License

GPL-3.0 License - See LICENSE file for details
