# BITUNCOIN Universal Wallet API Documentation

## Base URL
```
Production: https://api.bituncoin.io/v1
Testnet: https://testnet-api.bituncoin.io/v1
```

## Authentication

All API requests must include an API key in the Authorization header:

```http
Authorization: Bearer YOUR_API_KEY
```

### Obtaining API Keys

1. Log in to your BITUNCOIN account
2. Navigate to Settings > API Keys
3. Click "Generate New API Key"
4. Save the key securely (shown only once)

## Rate Limiting

- **Free Tier**: 100 requests per minute
- **Pro Tier**: 1,000 requests per minute
- **Enterprise**: Custom limits

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1634567890
```

## Wallet Management

### Get Wallet Balance

Retrieve the balance for a specific currency or all currencies.

**Endpoint**: `GET /wallet/balance`

**Parameters**:
- `currency` (optional): Specific currency (BTN, BTC, ETH, USDT, BNB, GLD)

**Example Request**:
```bash
curl -X GET "https://api.bituncoin.io/v1/wallet/balance?currency=BTN" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response**:
```json
{
  "currency": "BTN",
  "amount": 5000.0,
  "usdValue": 50000.0,
  "lastSync": 1634567890
}
```

### Get Portfolio Summary

Retrieve complete portfolio information.

**Endpoint**: `GET /wallet/portfolio`

**Example Response**:
```json
{
  "walletId": "wallet_123",
  "owner": "user@example.com",
  "balances": {
    "BTN": {
      "amount": 5000.0,
      "usdValue": 50000.0,
      "lastSync": 1634567890
    },
    "BTC": {
      "amount": 0.5,
      "usdValue": 22500.0,
      "lastSync": 1634567890
    }
  },
  "totalUSD": 72500.0,
  "lastUpdate": 1634567890,
  "txCount": 150
}
```

### Create Transaction

Send cryptocurrency to another address.

**Endpoint**: `POST /wallet/transaction`

**Request Body**:
```json
{
  "from": "wallet_address_123",
  "to": "recipient_address_456",
  "amount": 100.0,
  "currency": "BTN",
  "memo": "Payment for services"
}
```

**Response**:
```json
{
  "id": "tx_789",
  "type": "send",
  "currency": "BTN",
  "amount": 100.0,
  "fee": 0.1,
  "from": "wallet_address_123",
  "to": "recipient_address_456",
  "status": "pending",
  "timestamp": 1634567890,
  "txHash": "0x1234567890abcdef"
}
```

### Get Transaction History

Retrieve transaction history with optional filters.

**Endpoint**: `GET /wallet/transactions`

**Parameters**:
- `currency` (optional): Filter by currency
- `type` (optional): Filter by type (send, receive, stake, swap, exchange, card)
- `limit` (optional): Number of results (default: 50, max: 100)
- `offset` (optional): Pagination offset

**Example Response**:
```json
{
  "transactions": [
    {
      "id": "tx_789",
      "type": "send",
      "currency": "BTN",
      "amount": 100.0,
      "fee": 0.1,
      "status": "completed",
      "timestamp": 1634567890
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

## Payment Cards

### Create Card

Create a new virtual or physical payment card.

**Endpoint**: `POST /cards/create`

**Request Body**:
```json
{
  "walletId": "wallet_123",
  "cardholderName": "John Doe",
  "cardType": "virtual_visa",
  "spendingLimit": 5000.0,
  "dailyLimit": 1000.0,
  "monthlyLimit": 10000.0
}
```

**Response**:
```json
{
  "id": "card_456",
  "walletId": "wallet_123",
  "cardNumber": "5399************4532",
  "cardType": "virtual_visa",
  "expiryMonth": 12,
  "expiryYear": 2028,
  "cvv": "***",
  "status": "pending",
  "createdAt": 1634567890
}
```

### Activate Card

Activate a pending card.

**Endpoint**: `POST /cards/{cardId}/activate`

**Response**:
```json
{
  "id": "card_456",
  "status": "active",
  "activatedAt": 1634567890
}
```

### Get Card Details

**Endpoint**: `GET /cards/{cardId}`

**Response**:
```json
{
  "id": "card_456",
  "cardType": "virtual_visa",
  "status": "active",
  "spendingLimit": 5000.0,
  "spentToday": 250.0,
  "spentThisMonth": 1500.0
}
```

### Get Card Transactions

**Endpoint**: `GET /cards/{cardId}/transactions`

**Parameters**:
- `limit` (optional): Number of results (default: 50)

**Response**:
```json
{
  "transactions": [
    {
      "id": "ctx_123",
      "amount": 125.50,
      "merchant": "Amazon",
      "status": "completed",
      "timestamp": 1634567890
    }
  ]
}
```

## Exchange

### Get Exchange Rate

Get the current exchange rate between two currencies.

**Endpoint**: `GET /exchange/rate`

**Parameters**:
- `from`: Source currency
- `to`: Target currency

**Example Request**:
```bash
curl "https://api.bituncoin.io/v1/exchange/rate?from=BTN&to=USD"
```

**Response**:
```json
{
  "fromCurrency": "BTN",
  "toCurrency": "USD",
  "rate": 10.0,
  "lastUpdate": 1634567890
}
```

### Calculate Exchange

Calculate exchange amount without executing.

**Endpoint**: `POST /exchange/calculate`

**Request Body**:
```json
{
  "fromCurrency": "BTN",
  "toCurrency": "USD",
  "fromAmount": 1000.0
}
```

**Response**:
```json
{
  "fromCurrency": "BTN",
  "toCurrency": "USD",
  "fromAmount": 1000.0,
  "toAmount": 10000.0,
  "fee": 1.0,
  "rate": 10.0
}
```

### Create Exchange Order

Execute a currency exchange.

**Endpoint**: `POST /exchange/order`

**Request Body**:
```json
{
  "walletId": "wallet_123",
  "fromCurrency": "BTN",
  "toCurrency": "USD",
  "fromAmount": 1000.0
}
```

**Response**:
```json
{
  "id": "exchange_789",
  "status": "completed",
  "fromCurrency": "BTN",
  "toCurrency": "USD",
  "fromAmount": 1000.0,
  "toAmount": 10000.0,
  "fee": 1.0,
  "rate": 10.0,
  "createdAt": 1634567890,
  "completedAt": 1634567891
}
```

## Merchant Services

### Register Merchant

Register as a merchant to accept payments.

**Endpoint**: `POST /merchant/register`

**Request Body**:
```json
{
  "name": "My Store",
  "email": "store@example.com",
  "phone": "+1234567890",
  "businessType": "retail",
  "walletAddress": "wallet_merchant_123",
  "paymentMethods": ["qr_code", "nfc", "mtn_mobile_money"]
}
```

**Response**:
```json
{
  "id": "merchant_456",
  "name": "My Store",
  "status": "active",
  "apiKey": "sk_merchant_key_789",
  "registeredAt": 1634567890
}
```

### Create Payment Request

Create a payment request for a customer.

**Endpoint**: `POST /merchant/payment-request`

**Request Body**:
```json
{
  "merchantId": "merchant_456",
  "amount": 50.0,
  "currency": "BTN",
  "description": "Order #12345",
  "paymentMethods": ["qr_code", "nfc"],
  "expiryMinutes": 15
}
```

**Response**:
```json
{
  "id": "pr_789",
  "merchantId": "merchant_456",
  "amount": 50.0,
  "currency": "BTN",
  "qrCode": "base64_encoded_qr_code",
  "nfcData": "base64_encoded_nfc_data",
  "status": "pending",
  "createdAt": 1634567890,
  "expiresAt": 1634568790
}
```

### Process Payment

Mark a payment request as paid.

**Endpoint**: `POST /merchant/process-payment`

**Request Body**:
```json
{
  "requestId": "pr_789",
  "customerId": "customer_123",
  "transactionHash": "0xabcdef123456"
}
```

**Response**:
```json
{
  "requestId": "pr_789",
  "status": "completed",
  "paidAt": 1634567890
}
```

## AI Insights

### Get Portfolio Insights

Get AI-generated insights for a portfolio.

**Endpoint**: `GET /insights/portfolio`

**Parameters**:
- `walletId`: Wallet ID

**Response**:
```json
{
  "insights": [
    {
      "id": "ins_123",
      "category": "diversification",
      "title": "Low Portfolio Diversification",
      "description": "Your portfolio is concentrated in only 2 assets...",
      "confidence": 85.0,
      "impact": "high",
      "timestamp": 1634567890
    }
  ]
}
```

### Get Recommendations

Get AI-driven recommendations.

**Endpoint**: `GET /insights/recommendations`

**Parameters**:
- `walletId`: Wallet ID

**Response**:
```json
{
  "recommendations": [
    {
      "id": "rec_456",
      "type": "staking",
      "title": "Stake BTN for 5% APY",
      "description": "You have 5000 BTN available...",
      "priority": 8,
      "confidence": 95.0,
      "actions": ["stake_btn", "learn_more"]
    }
  ]
}
```

### Get Market Trends

Get market trend analysis for a currency.

**Endpoint**: `GET /insights/market-trend`

**Parameters**:
- `currency`: Currency code

**Response**:
```json
{
  "currency": "BTN",
  "trend": "up",
  "changePercent": 5.2,
  "volume24h": 1500000.0,
  "prediction": "continued_growth",
  "confidence": 70.0,
  "timestamp": 1634567890
}
```

## Security

### Enable 2FA

Enable two-factor authentication.

**Endpoint**: `POST /security/2fa/enable`

**Request Body**:
```json
{
  "walletId": "wallet_123",
  "method": "authenticator",
  "secret": "base32_secret"
}
```

**Response**:
```json
{
  "enabled": true,
  "backupCodes": ["12345678", "87654321", "..."]
}
```

### Register Device

Register a trusted device.

**Endpoint**: `POST /security/device/register`

**Request Body**:
```json
{
  "walletId": "wallet_123",
  "deviceType": "mobile",
  "os": "iOS",
  "browser": "Safari",
  "ipAddress": "192.168.1.1"
}
```

**Response**:
```json
{
  "deviceId": "dev_456",
  "isTrusted": false,
  "firstSeen": 1634567890
}
```

### Check Fraud

Check a transaction for fraud indicators.

**Endpoint**: `POST /security/fraud/check`

**Request Body**:
```json
{
  "walletId": "wallet_123",
  "txAmount": 10000.0,
  "currency": "BTN",
  "deviceId": "dev_456",
  "ipAddress": "192.168.1.1"
}
```

**Response**:
```json
{
  "isFraud": false,
  "indicators": [],
  "riskScore": 0.15
}
```

## Gold Reserve

### Get Reserve Info

Get information about gold reserves.

**Endpoint**: `GET /gold-reserve/info`

**Response**:
```json
{
  "totalReserveUsd": 2689000000000.0,
  "reservePerBtn": 537800.0,
  "totalBtnBacked": 5000000.0,
  "backingRatio": 100.0,
  "lastVerified": 1634567890,
  "status": "verified"
}
```

### Get Backing Proof

Get proof of backing for specific BTN amount.

**Endpoint**: `GET /gold-reserve/backing-proof`

**Parameters**:
- `amount`: BTN amount

**Response**:
```json
{
  "btnAmount": 100.0,
  "requiredReserveUsd": 53780000.0,
  "reservePerBtn": 537800.0,
  "isFullyBacked": true,
  "lastVerified": 1634567890
}
```

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid parameters |
| 401 | Unauthorized - Invalid API key |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |
| 503 | Service Unavailable |

## Error Response Format

```json
{
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "Invalid currency code",
    "details": {
      "parameter": "currency",
      "value": "INVALID"
    }
  }
}
```

## WebSocket API

### Real-Time Updates

Connect to WebSocket for real-time updates:

```javascript
const ws = new WebSocket('wss://api.bituncoin.io/v1/ws');

ws.onopen = () => {
  ws.send(JSON.stringify({
    action: 'subscribe',
    channel: 'balance',
    apiKey: 'YOUR_API_KEY'
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Balance update:', data);
};
```

### Channels

- `balance` - Balance updates
- `transactions` - Transaction notifications
- `rates` - Exchange rate updates
- `alerts` - Security and market alerts

## SDK Libraries

### JavaScript/Node.js

```bash
npm install @bituncoin/wallet-sdk
```

```javascript
const Bituncoin = require('@bituncoin/wallet-sdk');

const client = new Bituncoin('YOUR_API_KEY');

const balance = await client.wallet.getBalance('BTN');
```

### Python

```bash
pip install bituncoin-wallet
```

```python
from bituncoin import WalletClient

client = WalletClient('YOUR_API_KEY')
balance = client.wallet.get_balance('BTN')
```

---

**API Version**: v1  
**Last Updated**: October 19, 2025  
**Support**: api-support@bituncoin.io
