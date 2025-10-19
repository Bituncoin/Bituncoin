package wallet

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrAssetNotFound       = errors.New("asset not found")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidTransaction  = errors.New("invalid transaction")
)

// TransactionType represents different types of transactions
type TransactionType string

const (
	TypeSent     TransactionType = "sent"
	TypeReceived TransactionType = "received"
	TypeStaked   TransactionType = "staked"
	TypeUnstaked TransactionType = "unstaked"
	TypeSwapped  TransactionType = "swapped"
	TypePaid     TransactionType = "paid"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusConfirmed TransactionStatus = "confirmed"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID            string            `json:"id"`
	Type          TransactionType   `json:"type"`
	Status        TransactionStatus `json:"status"`
	From          string            `json:"from"`
	To            string            `json:"to"`
	Asset         string            `json:"asset"`
	Amount        float64           `json:"amount"`
	Fee           float64           `json:"fee"`
	Hash          string            `json:"hash"`
	Confirmations int               `json:"confirmations"`
	Timestamp     time.Time         `json:"timestamp"`
	Memo          string            `json:"memo,omitempty"`
}

// TransactionHistory manages wallet transaction history
type TransactionHistory struct {
	transactions map[string]*Transaction
	userTxs      map[string][]string // userAddress -> []txID
	mutex        sync.RWMutex
}

// NewTransactionHistory creates a new transaction history manager
func NewTransactionHistory() *TransactionHistory {
	return &TransactionHistory{
		transactions: make(map[string]*Transaction),
		userTxs:      make(map[string][]string),
	}
}

// AddTransaction adds a new transaction to history
func (th *TransactionHistory) AddTransaction(tx *Transaction) error {
	if tx == nil || tx.ID == "" {
		return ErrInvalidTransaction
	}

	th.mutex.Lock()
	defer th.mutex.Unlock()

	th.transactions[tx.ID] = tx

	// Index by user addresses
	if tx.From != "" {
		th.userTxs[tx.From] = append(th.userTxs[tx.From], tx.ID)
	}
	if tx.To != "" && tx.To != tx.From {
		th.userTxs[tx.To] = append(th.userTxs[tx.To], tx.ID)
	}

	return nil
}

// GetTransaction retrieves a transaction by ID
func (th *TransactionHistory) GetTransaction(txID string) (*Transaction, error) {
	th.mutex.RLock()
	defer th.mutex.RUnlock()

	tx, exists := th.transactions[txID]
	if !exists {
		return nil, ErrTransactionNotFound
	}

	return tx, nil
}

// GetUserTransactions retrieves all transactions for a user address
func (th *TransactionHistory) GetUserTransactions(address string) []*Transaction {
	th.mutex.RLock()
	defer th.mutex.RUnlock()

	txIDs, exists := th.userTxs[address]
	if !exists {
		return []*Transaction{}
	}

	transactions := make([]*Transaction, 0, len(txIDs))
	for _, txID := range txIDs {
		if tx, exists := th.transactions[txID]; exists {
			transactions = append(transactions, tx)
		}
	}

	return transactions
}

// UpdateTransactionStatus updates the status of a transaction
func (th *TransactionHistory) UpdateTransactionStatus(txID string, status TransactionStatus, confirmations int) error {
	th.mutex.Lock()
	defer th.mutex.Unlock()

	tx, exists := th.transactions[txID]
	if !exists {
		return ErrTransactionNotFound
	}

	tx.Status = status
	tx.Confirmations = confirmations

	return nil
}

// FilterTransactions filters transactions by type and asset
func (th *TransactionHistory) FilterTransactions(address string, txType TransactionType, asset string) []*Transaction {
	th.mutex.RLock()
	defer th.mutex.RUnlock()

	allTxs := th.GetUserTransactions(address)
	filtered := make([]*Transaction, 0)

	for _, tx := range allTxs {
		if (txType == "" || tx.Type == txType) && (asset == "" || tx.Asset == asset) {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

// GetRecentTransactions returns the most recent N transactions for a user
func (th *TransactionHistory) GetRecentTransactions(address string, limit int) []*Transaction {
	allTxs := th.GetUserTransactions(address)

	// Sort by timestamp (newest first)
	for i := 0; i < len(allTxs)-1; i++ {
		for j := i + 1; j < len(allTxs); j++ {
			if allTxs[i].Timestamp.Before(allTxs[j].Timestamp) {
				allTxs[i], allTxs[j] = allTxs[j], allTxs[i]
			}
		}
	}

	if limit > 0 && len(allTxs) > limit {
		return allTxs[:limit]
	}

	return allTxs
}
