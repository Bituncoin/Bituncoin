package goldcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// GoldCoin represents the Gold-Coin cryptocurrency
type GoldCoin struct {
	Name          string
	Symbol        string
	MaxSupply     uint64
	CircSupply    uint64
	Decimals      uint8
	StakingReward float64
	TxFee         float64
	Version       string
}

// Transaction represents a Gold-Coin transaction
type Transaction struct {
	ID        string
	From      string
	To        string
	Amount    float64
	Fee       float64
	Timestamp int64
	Signature string
}

// NewGoldCoin creates a new instance of Gold-Coin with defined tokenomics
func NewGoldCoin() *GoldCoin {
	return &GoldCoin{
		Name:          "Gold-Coin",
		Symbol:        "GLD",
		MaxSupply:     100000000, // 100 million coins
		CircSupply:    0,
		Decimals:      8,
		StakingReward: 5.0,  // 5% annual staking reward
		TxFee:         0.001, // 0.1% transaction fee
		Version:       "1.0.0",
	}
}

// CreateTransaction creates a new transaction
func (gc *GoldCoin) CreateTransaction(from, to string, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount: must be greater than 0")
	}

	if from == "" || to == "" {
		return nil, errors.New("invalid addresses: from and to cannot be empty")
	}

	fee := amount * gc.TxFee

	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
	}

	// Generate transaction ID
	tx.ID = tx.generateID()

	return tx, nil
}

// generateID generates a unique transaction ID using SHA-256
func (tx *Transaction) generateID() string {
	data := fmt.Sprintf("%s%s%f%d", tx.From, tx.To, tx.Amount, tx.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ValidateTransaction validates a transaction
func (gc *GoldCoin) ValidateTransaction(tx *Transaction) error {
	if tx == nil {
		return errors.New("transaction is nil")
	}

	if tx.Amount <= 0 {
		return errors.New("invalid amount")
	}

	if tx.From == "" || tx.To == "" {
		return errors.New("invalid addresses")
	}

	// Verify transaction ID
	expectedID := tx.generateID()
	if tx.ID != expectedID {
		return errors.New("invalid transaction ID")
	}

	return nil
}

// Mint creates new coins (only up to max supply)
func (gc *GoldCoin) Mint(amount uint64) error {
	if gc.CircSupply+amount > gc.MaxSupply {
		return errors.New("cannot mint: would exceed max supply")
	}
	gc.CircSupply += amount
	return nil
}

// GetTokenomics returns the current tokenomics information
func (gc *GoldCoin) GetTokenomics() map[string]interface{} {
	return map[string]interface{}{
		"name":           gc.Name,
		"symbol":         gc.Symbol,
		"maxSupply":      gc.MaxSupply,
		"circSupply":     gc.CircSupply,
		"decimals":       gc.Decimals,
		"stakingReward":  gc.StakingReward,
		"transactionFee": gc.TxFee,
		"version":        gc.Version,
	}
}
