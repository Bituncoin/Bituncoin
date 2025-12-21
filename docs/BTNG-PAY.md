# BTNG-PAY - Bituncoin Gold Payment Protocol

This document describes BTNG-PAY, a lightweight on-chain/off-chain payment flow for Bituncoin (Gold-Coin, GLD). It includes a specification, API endpoints, example Go server integration, and a simple client workflow.

## Goals

- Provide a merchant-friendly payment flow (invoice creation, payment, status callback).
- Support on-chain GLD payments and off-chain invoice settlement for convenience.
- Simple API that can be integrated with the existing API node and wallet.

## Terminology

- Invoice: A payment request created by a merchant. Contains ID, amount (GLD), memo, expiry, status.
- Payer/Customer: The user who pays the invoice using GLD.
- Merchant: The receiver of the payment who creates invoices.

## Endpoints (suggested)

- POST /api/btnpay/invoice
  - Create a new invoice
  - Body: { "merchant": "addr", "amount": 123.45, "currency": "GLD", "memo": "Order #123" }
  - Returns: { "invoiceId": "btngpay_...", "expiresAt": 169... }

- GET /api/btngpay/invoice/{id}
  - Get invoice details/status

- POST /api/btngpay/pay
  - Submit proof-of-payment or on-chain tx id
  - Body: { "invoiceId": "...", "from": "addr", "txId": "tx_..." }
  - Returns: { "status": "pending|paid|failed" }

- POST /api/btngpay/callback
  - (Optional) Merchant webhook to notify of invoice status changes

## Recommended flow

1. Merchant creates invoice using /api/btngpay/invoice.
2. Wallet fetches invoice and displays amount + QR code with payee address and invoiceId.
3. Customer pays on-chain using GLD to merchant address (including invoiceId in metadata) or uses off-chain channel.
4. Wallet submits /api/btngpay/pay with txId or proof.
5. BTNGPAY service validates on-chain transaction (confirmations) or marks paid after proof is verified; merchant receives webhook callback.

## Example: Minimal Go implementation (in-memory)

The payments module can be added under `payments/btngpay.go` and wired into `api/btngnode.go` to expose endpoints. Below is a small example that can be used as a starting point.

```go
package payments

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

// InvoiceStatus enumerates invoice states

type InvoiceStatus string

const (
    StatusPending InvoiceStatus = "pending"
    StatusPaid    InvoiceStatus = "paid"
    StatusExpired InvoiceStatus = "expired"
    StatusFailed  InvoiceStatus = "failed"
)

// Invoice represents a payment request

type Invoice struct {
    ID        string        `json:"id"`
    Merchant  string        `json:"merchant"`
    Amount    float64       `json:"amount"`
    Currency  string        `json:"currency"`
    Memo      string        `json:"memo"`
    CreatedAt int64         `json:"createdAt"`
    ExpiresAt int64         `json:"expiresAt"`
    Status    InvoiceStatus `json:"status"`
    TxID      string        `json:"txId,omitempty"`
}

// BtngPay is a minimal in-memory BTNG-PAY service

type BtngPay struct {
    invoices map[string]*Invoice
    mutex    sync.RWMutex
}

// NewBtngPay creates a new service
func NewBtngPay() *BtngPay {
    return &BtngPay{
        invoices: make(map[string]*Invoice),
    }
}

// CreateInvoice creates a new invoice
func (bp *BtngPay) CreateInvoice(merchant string, amount float64, memo string, ttlSeconds int64) (*Invoice, error) {
    if merchant == "" {
        return nil, errors.New("merchant required")
    }
    if amount <= 0 {
        return nil, errors.New("amount must be > 0")
    }

    id := fmt.Sprintf("btnggpay_%d", time.Now().UnixNano())
    now := time.Now().Unix()
    inv := &Invoice{
        ID:        id,
        Merchant:  merchant,
        Amount:    amount,
        Currency:  "GLD",
        Memo:      memo,
        CreatedAt: now,
        ExpiresAt: now + ttlSeconds,
        Status:    StatusPending,
    }

    bp.mutex.Lock()
    bp.invoices[id] = inv
    bp.mutex.Unlock()

    return inv, nil
}

// GetInvoice retrieves an invoice
func (bp *BtngPay) GetInvoice(id string) (*Invoice, error) {
    bp.mutex.RLock()
    defer bp.mutex.RUnlock()

    inv, ok := bp.invoices[id]
    if !ok {
        return nil, errors.New("invoice not found")
    }

    // expire if past expiration
    if time.Now().Unix() > inv.ExpiresAt && inv.Status == StatusPending {
        inv.Status = StatusExpired
    }

    return inv, nil
}

// MarkPaid marks an invoice as paid (typically after verifying tx)
func (bp *BtngPay) MarkPaid(id, txId string) error {
    bp.mutex.Lock()
    defer bp.mutex.Unlock()

    inv, ok := bp.invoices[id]
    if !ok {
        return errors.New("invoice not found")
    }
    if inv.Status != StatusPending {
        return errors.New("invoice not pending")
    }

    inv.Status = StatusPaid
    inv.TxID = txId
    return nil
}

```

## Integration notes

- To integrate with existing api/btngnode.go, add routes:
  - POST /api/btngpay/invoice -> create invoice
  - GET /api/btngpay/invoice/{id} -> get invoice
  - POST /api/btngpay/pay -> submit payment proof

- Validation of on-chain txs can reuse existing Gold-Coin logic (e.g., checking mempool or block confirmations) — a simple approach is to require the payer to POST the txId and the service verifies a few confirmations before marking paid.

- Persistence: switch the in-memory map to storage/LevelDB (or another DB) for production. Use payments.PutJSON/GetJSON helper wrappers.

## Wallet UI

- Wallet should call POST /api/btngpay/invoice to create an invoice (merchant side), or retrieve invoice and display a QR with params:
  - gld:<merchantAddress>?amount=123.45&memo=btngpay_... 

- After sending funds, wallet calls POST /api/btngpay/pay with invoiceId and txId, then polls GET /api/btngpay/invoice/{id} for status changes.

## Security

- Require authentication for merchant invoice creation (API key / JWT).
- Rate-limit invoice creation to prevent abuse.
- Sanitize memo fields and enforce maximum amounts.

## Example curl

Create invoice:

```bash
curl -X POST "http://localhost:8080/api/btngpay/invoice" \
  -H "Content-Type: application/json" \
  -d '{"merchant":"GLDmerchantAddr","amount":12.5,"memo":"Order #100"}'
```

Pay invoice (wallet submits tx id after sending funds):

```bash
curl -X POST "http://localhost:8080/api/btngpay/pay" \
  -H "Content-Type: application/json" \
  -d '{"invoiceId":"btngpay_...","from":"GLDfromAddr","txId":"tx_123"}'
```

---

This file is intended as a specification + starter implementation for BTNG-PAY. You can ask me to:
- Add the minimal Go implementation file under payments/btngpay.go and wire API handlers in api/btngnode.go.
- Add persistence using storage/leveldb.go.
- Add tests and examples.
