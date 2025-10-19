package bituncoin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Bituncoin represents the BTN cryptocurrency
type Bituncoin struct {
	Name          string
	Symbol        string
	MaxSupply     uint64
	CircSupply    uint64
	Decimals      uint8
	StakingReward float64
	TxFee         float64
	Version       string
}

// Transaction represents a BTN transaction
type Transaction struct {
	ID        string
	From      string
	To        string
	Amount    float64
	Fee       float64
	Timestamp int64
	Signature string
}

// NewBituncoin creates a new instance of Bituncoin (BTN) with defined tokenomics
func NewBituncoin() *Bituncoin {
	return &Bituncoin{
		Name:          "Bituncoin",
		Symbol:        "BTN",
		MaxSupply:     100000000, // 100 million coins
		CircSupply:    0,
		Decimals:      8,
		StakingReward: 5.0,  // 5% annual staking reward
		TxFee:         0.001, // 0.1% transaction fee
		Version:       "1.0.0",
	}
}

// CreateTransaction creates a new BTN transaction
func (btn *Bituncoin) CreateTransaction(from, to string, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if from == "" || to == "" {
		return nil, errors.New("from and to addresses required")
	}

	fee := btn.CalculateFee(amount)

	// Create transaction ID using SHA-256 hash
	data := fmt.Sprintf("%s:%s:%.8f:%d", from, to, amount, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	txID := hex.EncodeToString(hash[:])

	tx := &Transaction{
		ID:        txID,
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Timestamp: time.Now().Unix(),
		Signature: "",
	}

	return tx, nil
}

// CalculateFee calculates the transaction fee (0.1% of amount)
func (btn *Bituncoin) CalculateFee(amount float64) float64 {
	return amount * btn.TxFee
}

// ValidateTransaction validates a transaction
func (btn *Bituncoin) ValidateTransaction(tx *Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("invalid transaction amount")
	}

	if tx.From == "" || tx.To == "" {
		return errors.New("invalid addresses")
	}

	expectedFee := btn.CalculateFee(tx.Amount)
	if tx.Fee < expectedFee {
		return errors.New("insufficient transaction fee")
	}

	return nil
}

// Mint creates new BTN coins (up to max supply)
func (btn *Bituncoin) Mint(amount uint64) error {
	if btn.CircSupply+amount > btn.MaxSupply {
		return errors.New("minting would exceed max supply")
	}

	btn.CircSupply += amount
	return nil
}

// GetTokenomics returns the current tokenomics information
func (btn *Bituncoin) GetTokenomics() map[string]interface{} {
	return map[string]interface{}{
		"name":           btn.Name,
		"symbol":         btn.Symbol,
		"maxSupply":      btn.MaxSupply,
		"circSupply":     btn.CircSupply,
		"decimals":       btn.Decimals,
		"stakingReward":  btn.StakingReward,
		"transactionFee": btn.TxFee,
		"version":        btn.Version,
	}
}
