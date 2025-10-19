package wallet

import (
	"errors"
	"sync"
	"time"
)

// Currency represents supported cryptocurrencies
type Currency string

const (
	BTN  Currency = "BTN"  // Bituncoin
	BTC  Currency = "BTC"  // Bitcoin
	ETH  Currency = "ETH"  // Ethereum
	USDT Currency = "USDT" // Tether
	BNB  Currency = "BNB"  // Binance Coin
	GLD  Currency = "GLD"  // Gold-Coin (legacy support)
)

// Balance represents a currency balance
type Balance struct {
	Currency Currency `json:"currency"`
	Amount   float64  `json:"amount"`
	USDValue float64  `json:"usdValue"`
	LastSync int64    `json:"lastSync"`
}

// Portfolio represents the complete wallet portfolio
type Portfolio struct {
	Balances     map[Currency]*Balance `json:"balances"`
	TotalUSD     float64               `json:"totalUsd"`
	LastUpdate   int64                 `json:"lastUpdate"`
	mutex        sync.RWMutex
}

// TransactionType represents transaction types
type TransactionType string

const (
	TxTypeSend     TransactionType = "send"
	TxTypeReceive  TransactionType = "receive"
	TxTypeStake    TransactionType = "stake"
	TxTypeSwap     TransactionType = "swap"
	TxTypeExchange TransactionType = "exchange"
	TxTypeCard     TransactionType = "card"
)

// TransactionStatus represents transaction status
type TransactionStatus string

const (
	TxStatusPending   TransactionStatus = "pending"
	TxStatusCompleted TransactionStatus = "completed"
	TxStatusFailed    TransactionStatus = "failed"
	TxStatusCancelled TransactionStatus = "cancelled"
)

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID        string            `json:"id"`
	Type      TransactionType   `json:"type"`
	Currency  Currency          `json:"currency"`
	Amount    float64           `json:"amount"`
	Fee       float64           `json:"fee"`
	From      string            `json:"from"`
	To        string            `json:"to"`
	Status    TransactionStatus `json:"status"`
	Timestamp int64             `json:"timestamp"`
	TxHash    string            `json:"txHash,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// UniversalWallet represents the universal multi-currency wallet
type UniversalWallet struct {
	ID           string                   `json:"id"`
	Owner        string                   `json:"owner"`
	Portfolio    *Portfolio               `json:"portfolio"`
	Transactions []*WalletTransaction     `json:"transactions"`
	Security     *Security                `json:"security"`
	CreatedAt    int64                    `json:"createdAt"`
	mutex        sync.RWMutex
}

// NewUniversalWallet creates a new universal wallet
func NewUniversalWallet(owner string) *UniversalWallet {
	return &UniversalWallet{
		ID:        generateWalletID(),
		Owner:     owner,
		Portfolio: NewPortfolio(),
		Transactions: make([]*WalletTransaction, 0),
		Security: NewSecurity(),
		CreatedAt: time.Now().Unix(),
	}
}

// NewPortfolio creates a new portfolio
func NewPortfolio() *Portfolio {
	portfolio := &Portfolio{
		Balances:   make(map[Currency]*Balance),
		TotalUSD:   0,
		LastUpdate: time.Now().Unix(),
	}
	
	// Initialize all supported currencies
	currencies := []Currency{BTN, BTC, ETH, USDT, BNB, GLD}
	for _, curr := range currencies {
		portfolio.Balances[curr] = &Balance{
			Currency: curr,
			Amount:   0,
			USDValue: 0,
			LastSync: time.Now().Unix(),
		}
	}
	
	return portfolio
}

// GetBalance returns the balance for a specific currency
func (uw *UniversalWallet) GetBalance(currency Currency) (*Balance, error) {
	uw.mutex.RLock()
	defer uw.mutex.RUnlock()
	
	balance, exists := uw.Portfolio.Balances[currency]
	if !exists {
		return nil, errors.New("currency not supported")
	}
	
	return balance, nil
}

// UpdateBalance updates the balance for a specific currency
func (uw *UniversalWallet) UpdateBalance(currency Currency, amount float64, usdValue float64) error {
	uw.mutex.Lock()
	defer uw.mutex.Unlock()
	
	balance, exists := uw.Portfolio.Balances[currency]
	if !exists {
		return errors.New("currency not supported")
	}
	
	balance.Amount = amount
	balance.USDValue = usdValue
	balance.LastSync = time.Now().Unix()
	
	// Recalculate total USD value
	uw.Portfolio.calculateTotalUSD()
	uw.Portfolio.LastUpdate = time.Now().Unix()
	
	return nil
}

// AddTransaction adds a transaction to the wallet
func (uw *UniversalWallet) AddTransaction(tx *WalletTransaction) error {
	uw.mutex.Lock()
	defer uw.mutex.Unlock()
	
	if tx == nil {
		return errors.New("transaction cannot be nil")
	}
	
	uw.Transactions = append(uw.Transactions, tx)
	return nil
}

// GetTransactionHistory returns transaction history with optional filters
func (uw *UniversalWallet) GetTransactionHistory(currency *Currency, txType *TransactionType, limit int) []*WalletTransaction {
	uw.mutex.RLock()
	defer uw.mutex.RUnlock()
	
	filtered := make([]*WalletTransaction, 0)
	
	for _, tx := range uw.Transactions {
		// Apply filters
		if currency != nil && tx.Currency != *currency {
			continue
		}
		if txType != nil && tx.Type != *txType {
			continue
		}
		
		filtered = append(filtered, tx)
		
		// Apply limit
		if limit > 0 && len(filtered) >= limit {
			break
		}
	}
	
	return filtered
}

// GetPortfolioSummary returns a summary of the portfolio
func (uw *UniversalWallet) GetPortfolioSummary() map[string]interface{} {
	uw.mutex.RLock()
	defer uw.mutex.RUnlock()
	
	balances := make(map[string]interface{})
	for currency, balance := range uw.Portfolio.Balances {
		balances[string(currency)] = map[string]interface{}{
			"amount":   balance.Amount,
			"usdValue": balance.USDValue,
			"lastSync": balance.LastSync,
		}
	}
	
	return map[string]interface{}{
		"walletId":    uw.ID,
		"owner":       uw.Owner,
		"balances":    balances,
		"totalUSD":    uw.Portfolio.TotalUSD,
		"lastUpdate":  uw.Portfolio.LastUpdate,
		"txCount":     len(uw.Transactions),
	}
}

// calculateTotalUSD calculates the total USD value of the portfolio
func (p *Portfolio) calculateTotalUSD() {
	total := 0.0
	for _, balance := range p.Balances {
		total += balance.USDValue
	}
	p.TotalUSD = total
}

// generateWalletID generates a unique wallet ID
func generateWalletID() string {
	return "wallet_" + time.Now().Format("20060102150405")
}
