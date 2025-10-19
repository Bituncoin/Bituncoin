# BTNg Wallet - Advanced Features Documentation

## Overview

The BTNg Wallet is a comprehensive blockchain-based wallet system that provides multi-currency support, merchant payment integration, AI-powered assistance, and advanced security features.

## Features

### 1. Multi-Currency Support

The wallet supports the following cryptocurrencies:
- **BTN (Bituncoin)** - Native Bituncoin token
- **GLD (Gold-Coin)** - Proof-of-Stake token
- **BTC (Bitcoin)** - Bitcoin integration
- **ETH (Ethereum)** - Ethereum integration
- **USDT (Tether)** - Stablecoin support
- **BNB (Binance Coin)** - Binance Smart Chain support

All currencies display real-time balances with USD conversion rates.

### 2. Payment Integration

#### BTN-Pay Invoicing System
Create and manage invoices for merchant payments:
```bash
# Create invoice
curl -X POST http://localhost:8080/api/btnpay/invoice \
  -H "Content-Type: application/json" \
  -d '{"merchant":"BTNmerchant123","amount":100.0,"memo":"Order #123"}'

# Get invoice status
curl http://localhost:8080/api/btnpay/invoice/{invoiceId}

# Pay invoice
curl -X POST http://localhost:8080/api/btnpay/pay \
  -H "Content-Type: application/json" \
  -d '{"invoiceId":"btnpay_123","from":"BTNaddr","txId":"tx_456"}'
```

#### Mobile Money Integration
Support for African mobile money services:
- MTN Mobile Money
- AirtelTigo Money
- Vodafone Cash

```bash
# Initiate mobile money payment
curl -X POST http://localhost:8080/api/mobilemoney/pay \
  -H "Content-Type: application/json" \
  -d '{"provider":"MTN","phoneNumber":"+233244123456","amount":100.0,"currency":"GHS"}'

# Check payment status
curl http://localhost:8080/api/mobilemoney/status?id={transactionId}
```

#### QR Code Payments
Generate and parse QR codes for payments:
```bash
# Generate payment QR code
curl -X POST http://localhost:8080/api/qrcode/generate \
  -H "Content-Type: application/json" \
  -d '{"address":"BTNaddr123","amount":50.0,"currency":"BTN","memo":"Payment"}'

# Parse QR code
curl -X POST http://localhost:8080/api/qrcode/parse \
  -H "Content-Type: application/json" \
  -d '{"qrData":"<qr_data_string>"}'
```

#### BTN-Pay Cards (MasterCard/Visa)
Issue virtual and physical payment cards:
```bash
# Issue new card
curl -X POST http://localhost:8080/api/card/issue \
  -H "Content-Type: application/json" \
  -d '{"holderName":"John Doe","linkedAddress":"BTNaddr","cardType":"MasterCard","category":"Virtual"}'

# Load card with funds
curl -X POST http://localhost:8080/api/card/load \
  -H "Content-Type: application/json" \
  -d '{"cardId":"card_123","amount":500.0}'

# Process card payment
curl -X POST http://localhost:8080/api/card/payment \
  -H "Content-Type: application/json" \
  -d '{"cardId":"card_123","amount":50.0,"merchant":"Store ABC"}'
```

### 3. Exchange & Swap Features

Built-in crypto-to-crypto and crypto-to-fiat exchange:
```bash
# Get exchange rate
curl http://localhost:8080/api/exchange/rates?from=BTN&to=BTC

# Estimate swap
curl http://localhost:8080/api/exchange/estimate?from=BTN&to=ETH&amount=100

# Execute swap
curl -X POST http://localhost:8080/api/exchange/swap \
  -H "Content-Type: application/json" \
  -d '{"from":"BTN","to":"ETH","amount":100.0,"userAddress":"BTNaddr"}'
```

Features:
- Real-time exchange rates
- Low fees (0.5% per swap)
- Cross-chain swaps
- Fiat currency support (USD, EUR, GBP, GHS, NGN, KES)

### 4. AI-Powered Assistant

The AI assistant provides intelligent insights and recommendations:

#### Portfolio Analysis
```bash
curl http://localhost:8080/api/ai/insights
```
Returns:
- Portfolio diversification analysis
- Staking recommendations
- Market trend insights
- Risk assessment

#### Security Alerts
```bash
curl http://localhost:8080/api/ai/alerts
```
Returns:
- Fraud detection alerts
- Security status checks
- Backup reminders
- 2FA status

#### Ask AI Assistant
```bash
curl -X POST http://localhost:8080/api/ai/ask \
  -H "Content-Type: application/json" \
  -d '{"query":"How can I improve my portfolio?"}'
```

Sample queries:
- "How can I improve my portfolio?"
- "What are the best staking strategies?"
- "How do I secure my wallet?"
- "Should I swap my BTN to BTC?"

### 5. Enhanced Security Features

#### Two-Factor Authentication (2FA)
- Time-based OTP support
- SMS verification
- Email verification

#### Biometric Authentication
- Fingerprint recognition
- Face recognition
- Voice authentication (future)

#### Fraud Monitoring
- Real-time transaction monitoring
- Suspicious activity alerts
- Automatic transaction blocking
- Daily/monthly spending limits

#### Encryption
- AES-256 encryption for wallet data
- ECDSA for key management
- SHA-256 for transaction hashing
- Secure key storage

### 6. Staking Features

Stake BTN or GLD tokens to earn rewards:
- 5% Annual Percentage Yield (APY)
- Minimum stake: 100 tokens
- Lock period: 30 days
- Daily reward calculation
- Automatic compounding option

### 7. User Interface Features

#### Dark/Light Mode
Toggle between dark and light themes for comfortable viewing in any environment.

#### Multi-Language Support
Currently supported languages:
- English
- Français (French)
- Español (Spanish)
- 中文 (Chinese)

#### Responsive Design
Optimized for:
- Desktop browsers
- Mobile devices
- Tablet devices

#### Navigation
Six main sections:
1. **Overview** - View all balances and quick actions
2. **Pay** - Merchant payments, invoices, QR codes, mobile money
3. **Staking** - Stake tokens and earn rewards
4. **Transactions** - Transaction history
5. **Security** - Security settings and encryption status
6. **AI Assistant** - AI-powered insights and recommendations

## API Endpoints Reference

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
- `GET /api/mobilemoney/status?id=<id>` - Payment status
- `POST /api/qrcode/generate` - Generate QR code
- `POST /api/qrcode/parse` - Parse QR code

### Card Endpoints
- `POST /api/card/issue` - Issue new card
- `POST /api/card/load` - Load card
- `POST /api/card/payment` - Process payment
- `GET /api/card/transactions?cardId=<id>` - Get transactions

### Exchange Endpoints
- `GET /api/exchange/rates?from=<cur>&to=<cur>` - Get exchange rate
- `GET /api/exchange/estimate?from=<cur>&to=<cur>&amount=<amt>` - Estimate swap
- `POST /api/exchange/swap` - Execute swap

### AI Assistant Endpoints
- `GET /api/ai/insights` - Get portfolio insights
- `GET /api/ai/alerts` - Get security alerts
- `POST /api/ai/ask` - Ask AI assistant

## Security Best Practices

1. **Enable 2FA** - Always enable two-factor authentication
2. **Regular Backups** - Create encrypted backups weekly
3. **Strong Passwords** - Use unique, strong passwords
4. **Verify Addresses** - Always double-check recipient addresses
5. **Monitor Transactions** - Enable fraud monitoring
6. **Secure Storage** - Store recovery phrases offline
7. **Update Software** - Keep wallet software up to date
8. **Use Hardware Wallets** - For large holdings

## Performance Characteristics

- **Transaction Speed**: 2-30 seconds (varies by blockchain)
- **Swap Processing**: 2-30 minutes (depends on currencies)
- **Card Transactions**: Instant
- **Mobile Money**: 2-5 minutes
- **API Response Time**: <100ms
- **Daily Limits**: $5,000 (cards), configurable

## Scalability

The wallet is designed for scalability:
- Modular architecture
- Microservices-ready
- Database-agnostic storage layer
- Horizontal scaling support
- Load balancer compatible
- Caching layer support

## Future Enhancements

- [ ] Hardware wallet integration (Ledger, Trezor)
- [ ] DEX integration
- [ ] NFT support
- [ ] DeFi protocols integration
- [ ] Mobile apps (iOS, Android)
- [ ] Desktop apps (Windows, macOS, Linux)
- [ ] Voice commands
- [ ] Advanced charting
- [ ] Tax reporting
- [ ] Governance voting

## Support

For issues and questions:
- GitHub Issues: https://github.com/Bituncoin/Bituncoin/issues
- Email: support@bituncoin.com
- Documentation: https://docs.bituncoin.com

## License

GPL-3.0 - See LICENSE file for details
