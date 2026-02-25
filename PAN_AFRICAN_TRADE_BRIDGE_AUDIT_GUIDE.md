# 🏛️ BTNG Pan-African Trade Bridge: Settlement Flow Documentation

## Overview

The Pan-African Trade Bridge enables instant gold-backed settlements between merchants across all 54 BTNG sovereign nations, bypassing traditional banking rails and fiat currency volatility.

## Architecture

### Core Components

- **Settlement API**: `/api/btng/bridge/settle` (POST)
- **Gold Token Transfers**: Hyperledger Fabric chaincode operations
- **Cross-Node Communication**: RESTful API calls between sovereign nodes
- **Audit Trail**: MongoDB settlement records with cryptographic verification

### Settlement Flow Diagram

```
Merchant A (Nation X) → Bridge API → Gold Token Transfer → Merchant B (Nation Y)
       ↓                        ↓              ↓                      ↓
   Invoice Generation    Fiat→Gold Conversion  Fabric TX         Credit Received
   (Local Currency)      (Real-time Pricing)   (Immutable)        (Instant)
```

## API Specification

### Endpoint: `POST /api/btng/bridge/settle`

**Authentication**: Required (JWT Bearer Token)

**Request Body**:

```json
{
  "fromMerchant": "merchant-uuid-123",
  "toMerchant": "merchant-uuid-456",
  "fromNation": "GH",
  "toNation": "NG",
  "amount": 1000.50,
  "currency": "GHS",
  "invoiceId": "INV-2026-00123"
}
```

**Response (Success)**:

```json
{
  "status": "ok",
  "settlement": {
    "id": "bridge-1771457774000-abc123def",
    "fromMerchant": "merchant-uuid-123",
    "toMerchant": "merchant-uuid-456",
    "fromNation": "GH",
    "toNation": "NG",
    "amount": 1000.50,
    "currency": "GHS",
    "goldEquivalent": 0.547,
    "status": "settled",
    "timestamp": "2026-02-24T09:45:00.000Z",
    "transactionId": "btng-tx-1771457774000"
  },
  "message": "Pan-African settlement completed: 1000.5 GHS transferred in gold",
  "goldTransferred": 0.547
}
```

**Response (Error)**:

```json
{
  "error": "Settlement failed",
  "settlement": {
    "id": "bridge-1771457774000-abc123def",
    "status": "failed",
    "timestamp": "2026-02-24T09:45:00.000Z"
  },
  "details": "Insufficient gold balance"
}
```

## Settlement Process Steps

### 1. Authentication & Validation

- JWT token verification
- Merchant identity validation
- Nation code verification
- Amount and currency validation

### 2. Gold Price Conversion

- Fetch latest gold price from oracle
- Convert fiat amount to gold grams
- Apply current spot price: $5,185.76/oz ($1,832.31/g)

### 3. Balance Verification

- Check sender's gold token balance
- Verify sufficient funds for settlement
- Reserve gold tokens during processing

### 4. Cross-Node Transfer

- Execute Fabric chaincode `Transfer` function
- Generate unique transaction ID
- Update both merchant balances atomically

### 5. Settlement Record

- Store complete audit trail in MongoDB
- Include all conversion rates and timestamps
- Link to original invoice for reconciliation

## Audit Trail Structure

Each settlement creates an immutable record:

```javascript
{
  _id: ObjectId("507f1f77bcf86cd799439011"),
  settlementId: "bridge-1771457774000-abc123def",
  fromMerchant: "merchant-uuid-123",
  toMerchant: "merchant-uuid-456",
  fromNation: "GH",
  toNation: "NG",
  originalAmount: 1000.50,
  originalCurrency: "GHS",
  goldEquivalent: 0.547,
  goldPriceAtSettlement: {
    usdPerOunce: 5185.76,
    usdPerGram: 1832.31,
    timestamp: "2026-02-24T09:45:00.000Z"
  },
  fabricTransaction: {
    txId: "btng-tx-1771457774000",
    blockNumber: 12488,
    status: "committed"
  },
  status: "settled",
  createdAt: ISODate("2026-02-24T09:45:00.000Z"),
  completedAt: ISODate("2026-02-24T09:45:01.234Z")
}
```

## Security Measures

### Cryptographic Verification

- All settlements signed with node private keys
- Transaction hashes verified on Fabric network
- End-to-end encryption for cross-node communication

### Fraud Prevention

- Real-time balance checks
- Rate limiting per merchant (100 settlements/minute)
- Duplicate transaction detection
- Geographic risk scoring

### Compliance

- Full audit trail retention (7 years)
- Regulatory reporting capabilities
- AML/KYC integration points
- Cross-border transaction logging

## Monitoring & Alerting

### Key Metrics

- Settlement success rate (>99.9%)
- Average settlement time (<2 seconds)
- Cross-node transfer latency
- Gold token balance accuracy

### Alert Triggers

- Settlement failure rate >0.1%
- Average latency >5 seconds
- Gold balance discrepancies
- Network partition events

## Integration Examples

### cURL Example

```bash
curl -X POST "http://74.118.126.72:64799/api/btng/bridge/settle" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "fromMerchant": "merchant-gh-001",
    "toMerchant": "merchant-ng-001",
    "fromNation": "GH",
    "toNation": "NG",
    "amount": 1000.50,
    "currency": "GHS",
    "invoiceId": "INV-2026-00123"
  }'
```

### Node.js Integration

```javascript
const BTNG = require('btng-merchant-sdk');

const bridge = new BTNG.PanAfricanBridge({
  apiKey: 'your-jwt-token',
  baseUrl: 'http://74.118.126.72:64799'
});

const settlement = await bridge.settle({
  fromMerchant: 'merchant-gh-001',
  toMerchant: 'merchant-ng-001',
  fromNation: 'GH',
  toNation: 'NG',
  amount: 1000.50,
  currency: 'GHS',
  invoiceId: 'INV-2026-00123'
});

console.log('Settlement completed:', settlement.id);
```

## Reconciliation & Reporting

### Daily Reconciliation

- Compare settlement records with merchant statements
- Verify gold token balances across all nodes
- Reconcile with Fabric blockchain state

### Monthly Reporting

- Cross-border trade volume by nation pairs
- Settlement success rates and failure analysis
- Gold reserve utilization trends

### Audit Procedures

- Random sampling of settlement records
- Cryptographic verification of transaction signatures
- Balance proof validation with Merkle trees

## Future Enhancements

### Planned Features

- **Multi-hop routing**: Settlements through intermediate nodes
- **Atomic swaps**: Cross-chain settlement capabilities
- **Smart contracts**: Automated escrow and conditional payments
- **Real-time FX**: Dynamic currency conversion during settlement

### Scalability Improvements

- Settlement batching for high-volume corridors
- Regional bridge hubs for reduced latency
- Off-chain state channels for micro-payments

---

**Document Version**: 1.0
**Effective Date**: February 24, 2026
**Review Cycle**: Quarterly
**Document Owner**: BTNG Sovereign Governor
