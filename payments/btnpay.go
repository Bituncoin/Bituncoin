package payments

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
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

    id := fmt.Sprintf("btngpay_%d", time.Now().UnixNano())
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

// --- HTTP Helpers (optional) ---

// CreateInvoiceHandler handles POST /api/btngpay/invoice
func (bp *BtngPay) CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Merchant string  `json:"merchant"`
        Amount   float64 `json:"amount"`
        Memo     string  `json:"memo"`
        TTL      int64   `json:"ttlSeconds"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if req.TTL == 0 {
        req.TTL = 15 * 60 // default 15 minutes
    }

    inv, err := bp.CreateInvoice(req.Merchant, req.Amount, req.Memo, req.TTL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(inv)
}

// GetInvoiceHandler handles GET /api/btngpay/invoice/{id}
func (bp *BtngPay) GetInvoiceHandler(w http.ResponseWriter, r *http.Request) {
    // assume the last path segment is the invoice id
    parts := splitPath(r.URL.Path)
    if len(parts) == 0 {
        http.Error(w, "invoice id required", http.StatusBadRequest)
        return
    }
    id := parts[len(parts)-1]

    inv, err := bp.GetInvoice(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(inv)
}

// PayInvoiceHandler handles POST /api/btngpay/pay
func (bp *BtngPay) PayInvoiceHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        InvoiceID string `json:"invoiceId"`
        From      string `json:"from"`
        TxID      string `json:"txId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if req.InvoiceID == "" || req.TxID == "" {
        http.Error(w, "invoiceId and txId required", http.StatusBadRequest)
        return
    }

    // In a real implementation, verify the transaction on-chain and confirmations
    if err := bp.MarkPaid(req.InvoiceID, req.TxID); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "paid"})
}

// small helper to split path
func splitPath(p string) []string {
    var parts []string
    cur := ""
    for i := 0; i < len(p); i++ {
        c := p[i]
        if c == '/' {
            if cur != "" {
                parts = append(parts, cur)
                cur = ""
            }
            continue
        }
        cur += string(c)
    }
    if cur != "" {
        parts = append(parts, cur)
    }
    return parts
}