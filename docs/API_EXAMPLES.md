# Bituncoin Wallet - API Integration Examples

This document provides practical examples of integrating with the Bituncoin Wallet API.

## Table of Contents
1. [Portfolio Management](#portfolio-management)
2. [Exchange Operations](#exchange-operations)
3. [Payment Cards](#payment-cards)
4. [Merchant Integration](#merchant-integration)
5. [Security & Authentication](#security--authentication)

## Portfolio Management

### Get User Portfolio

**Request:**
```http
GET /api/wallet/portfolio/BTN123abc456def
Authorization: Bearer <token>
```

**Response:**
```json
{
  "address": "BTN123abc456def",
  "totalValueUSD": 82500.00,
  "assets": [
    {
      "symbol": "BTN",
      "name": "Bituncoin",
      "balance": 1000.0,
      "priceUSD": 1.0,
      "usdValue": 1000.0,
      "change24h": 2.5,
      "lastUpdated": "2025-10-19T13:00:00Z"
    },
    {
      "symbol": "BTC",
      "name": "Bitcoin",
      "balance": 0.5,
      "priceUSD": 50000.0,
      "usdValue": 25000.0,
      "change24h": -1.2,
      "lastUpdated": "2025-10-19T13:00:00Z"
    }
  ],
  "performance": {
    "avgChange24h": 0.65,
    "assetCount": 2
  }
}
```

### Add Asset to Portfolio

**Request:**
```http
POST /api/wallet/portfolio/add
Authorization: Bearer <token>
Content-Type: application/json

{
  "address": "BTN123abc456def",
  "symbol": "ETH",
  "name": "Ethereum",
  "balance": 2.0,
  "priceUSD": 3000.0
}
```

**Response:**
```json
{
  "success": true,
  "message": "Asset added successfully",
  "asset": {
    "symbol": "ETH",
    "balance": 2.0,
    "usdValue": 6000.0
  }
}
```

## Exchange Operations

### Get Exchange Quote

**Request:**
```http
POST /api/exchange/quote
Authorization: Bearer <token>
Content-Type: application/json

{
  "fromAsset": "BTC",
  "toAsset": "ETH",
  "amount": 0.5
}
```

**Response:**
```json
{
  "fromAsset": "BTC",
  "toAsset": "ETH",
  "fromAmount": 0.5,
  "toAmount": 7.4625,
  "rate": 15.0,
  "fee": 0.0375,
  "feePercentage": 0.5,
  "expiresAt": "2025-10-19T13:15:00Z"
}
```

### Create Exchange Order

**Request:**
```http
POST /api/exchange/order
Authorization: Bearer <token>
Content-Type: application/json

{
  "userAddress": "BTN123abc456def",
  "fromAsset": "BTC",
  "toAsset": "ETH",
  "amount": 0.5
}
```

**Response:**
```json
{
  "orderId": "EX-1760879661100885477",
  "type": "crypto-to-crypto",
  "fromAsset": "BTC",
  "toAsset": "ETH",
  "fromAmount": 0.5,
  "toAmount": 7.4625,
  "rate": 15.0,
  "fee": 0.0375,
  "status": "pending",
  "timestamp": "2025-10-19T13:00:00Z"
}
```

### Get Exchange Orders

**Request:**
```http
GET /api/exchange/orders/BTN123abc456def
Authorization: Bearer <token>
```

**Response:**
```json
{
  "orders": [
    {
      "orderId": "EX-1760879661100885477",
      "type": "crypto-to-crypto",
      "fromAsset": "BTC",
      "toAsset": "ETH",
      "fromAmount": 0.5,
      "toAmount": 7.4625,
      "status": "completed",
      "timestamp": "2025-10-19T12:00:00Z"
    }
  ],
  "total": 1
}
```

## Payment Cards

### Create Payment Card

**Request:**
```http
POST /api/cards/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "userAddress": "BTN123abc456def",
  "cardType": "virtual",
  "provider": "visa",
  "dailyLimit": 1000.0
}
```

**Response:**
```json
{
  "cardId": "CARD-1760879661100970305",
  "cardNumber": "4879661100915182",
  "cardType": "virtual",
  "provider": "visa",
  "status": "active",
  "balance": 0.0,
  "dailyLimit": 1000.0,
  "expiryDate": "10/28",
  "cvv": "123",
  "createdAt": "2025-10-19T13:00:00Z"
}
```

### Top-up Card

**Request:**
```http
POST /api/cards/CARD-1760879661100970305/topup
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": 500.0,
  "source": "wallet"
}
```

**Response:**
```json
{
  "success": true,
  "cardId": "CARD-1760879661100970305",
  "newBalance": 500.0,
  "transactionId": "TX-topup-123456"
}
```

### Get Card Transactions

**Request:**
```http
GET /api/cards/CARD-1760879661100970305/transactions
Authorization: Bearer <token>
```

**Response:**
```json
{
  "cardId": "CARD-1760879661100970305",
  "transactions": [
    {
      "id": "CTX-1760879661101012345",
      "merchant": "Amazon Store",
      "amount": 99.99,
      "currency": "USD",
      "status": "completed",
      "type": "purchase",
      "timestamp": "2025-10-19T12:30:00Z"
    }
  ],
  "total": 1
}
```

## Merchant Integration

### Register as Merchant

**Request:**
```http
POST /api/merchant/register
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Coffee Shop",
  "walletAddress": "GLD789xyz012abc",
  "email": "shop@example.com",
  "businessType": "retail"
}
```

**Response:**
```json
{
  "merchantId": "MERCH-1760879661101055944",
  "name": "Coffee Shop",
  "walletAddress": "GLD789xyz012abc",
  "status": "active",
  "acceptedMethods": ["qr_code", "nfc", "wallet_transfer"],
  "acceptedAssets": ["BTN", "GLD", "BTC", "ETH", "USDT"],
  "registrationDate": "2025-10-19T13:00:00Z"
}
```

### Create Payment Request

**Request:**
```http
POST /api/merchant/payment/request
Authorization: Bearer <token>
Content-Type: application/json

{
  "merchantId": "MERCH-1760879661101055944",
  "amount": 25.50,
  "asset": "GLD",
  "paymentMethod": "qr_code",
  "description": "Coffee and pastry"
}
```

**Response:**
```json
{
  "requestId": "PAY-1760879661101070571",
  "merchantId": "MERCH-1760879661101055944",
  "amount": 25.50,
  "asset": "GLD",
  "paymentMethod": "qr_code",
  "qrCode": "YnRuOkdMRDc4OS4uLj9hbW91bnQ9MjUuNTAmYXNzZXQ9R0xEJnJlcXVlc3Q9UEFZLTE3NjA4Nzk2NjExMDEwNzA1NzE=",
  "status": "pending",
  "createdAt": "2025-10-19T13:00:00Z",
  "expiresAt": "2025-10-19T13:15:00Z"
}
```

### Complete Payment

**Request:**
```http
POST /api/merchant/payment/complete
Authorization: Bearer <token>
Content-Type: application/json

{
  "requestId": "PAY-1760879661101070571",
  "customerAddress": "BTN123abc456def",
  "txHash": "0x123456789abcdef"
}
```

**Response:**
```json
{
  "success": true,
  "requestId": "PAY-1760879661101070571",
  "status": "completed",
  "completedAt": "2025-10-19T13:05:00Z"
}
```

### Process Mobile Money Payment

**Request:**
```http
POST /api/merchant/payment/mobile-money
Authorization: Bearer <token>
Content-Type: application/json

{
  "merchantId": "MERCH-1760879661101055944",
  "provider": "mtn_mobile_money",
  "phoneNumber": "+233123456789",
  "amount": 50.0,
  "currency": "GHS"
}
```

**Response:**
```json
{
  "paymentId": "MM-1760879661101097071",
  "provider": "mtn_mobile_money",
  "status": "pending",
  "transactionRef": "TXN-abc123def456",
  "timestamp": "2025-10-19T13:00:00Z",
  "instructions": "Please authorize payment on your phone"
}
```

## Security & Authentication

### Enable Two-Factor Authentication

**Request:**
```http
POST /api/security/2fa/enable
Authorization: Bearer <token>
Content-Type: application/json

{
  "userAddress": "BTN123abc456def"
}
```

**Response:**
```json
{
  "success": true,
  "secret": "JBSWY3DPEHPK3PXP",
  "qrCode": "data:image/png;base64,iVBORw0KG...",
  "backupCodes": [
    "12345678",
    "87654321",
    "11223344"
  ]
}
```

### Verify 2FA Code

**Request:**
```http
POST /api/security/2fa/verify
Authorization: Bearer <token>
Content-Type: application/json

{
  "userAddress": "BTN123abc456def",
  "code": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "2FA verified successfully"
}
```

### Get Security Alerts

**Request:**
```http
GET /api/security/alerts/BTN123abc456def
Authorization: Bearer <token>
```

**Response:**
```json
{
  "alerts": [
    {
      "id": "ALERT-20251019130000",
      "type": "security",
      "severity": "high",
      "message": "Unusual login attempt detected",
      "timestamp": "2025-10-19T12:00:00Z",
      "resolved": false
    }
  ],
  "total": 1
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "success": false,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "Amount must be greater than zero",
    "details": {}
  }
}
```

### Common Error Codes

- `UNAUTHORIZED` - Invalid or missing authentication token
- `INVALID_REQUEST` - Request validation failed
- `INSUFFICIENT_BALANCE` - Not enough balance for operation
- `ASSET_NOT_FOUND` - Requested asset doesn't exist
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `INTERNAL_ERROR` - Server error

## Rate Limiting

API requests are limited to:
- 100 requests per minute for standard endpoints
- 10 requests per minute for exchange operations
- 20 requests per minute for payment operations

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1697727600
```

## Webhooks

Configure webhooks to receive real-time notifications:

### Webhook Events

- `transaction.completed` - Transaction completed
- `exchange.completed` - Exchange order completed
- `payment.received` - Payment received by merchant
- `card.transaction` - Card transaction processed
- `security.alert` - Security alert triggered

### Webhook Payload Example

```json
{
  "event": "payment.received",
  "timestamp": "2025-10-19T13:00:00Z",
  "data": {
    "paymentId": "PAY-123456",
    "merchantId": "MERCH-789012",
    "amount": 25.50,
    "asset": "GLD",
    "customerAddress": "BTN123abc456def"
  }
}
```

## Testing

Use the test environment for development:
- Base URL: `https://test-api.bituncoin.io`
- Test tokens available in dashboard
- All transactions use test assets

---

For more information, visit:
- Developer Portal: https://developers.bituncoin.io
- API Documentation: https://api.bituncoin.io/docs
- Support: support@bituncoin.io
