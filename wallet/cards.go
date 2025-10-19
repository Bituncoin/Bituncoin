package wallet

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// CardType represents the type of payment card
type CardType string

const (
	CardTypeVirtualVisa       CardType = "virtual_visa"
	CardTypeVirtualMastercard CardType = "virtual_mastercard"
	CardTypePhysicalVisa      CardType = "physical_visa"
	CardTypePhysicalMastercard CardType = "physical_mastercard"
)

// CardStatus represents the status of a card
type CardStatus string

const (
	CardStatusActive    CardStatus = "active"
	CardStatusSuspended CardStatus = "suspended"
	CardStatusCancelled CardStatus = "cancelled"
	CardStatusPending   CardStatus = "pending"
)

// PaymentCard represents a BTN-Pay card
type PaymentCard struct {
	ID              string     `json:"id"`
	WalletID        string     `json:"walletId"`
	CardNumber      string     `json:"cardNumber"`
	CardType        CardType   `json:"cardType"`
	ExpiryMonth     int        `json:"expiryMonth"`
	ExpiryYear      int        `json:"expiryYear"`
	CVV             string     `json:"cvv"`
	CardholderName  string     `json:"cardholderName"`
	Status          CardStatus `json:"status"`
	SpendingLimit   float64    `json:"spendingLimit"`
	DailyLimit      float64    `json:"dailyLimit"`
	MonthlyLimit    float64    `json:"monthlyLimit"`
	SpentToday      float64    `json:"spentToday"`
	SpentThisMonth  float64    `json:"spentThisMonth"`
	ActivatedAt     int64      `json:"activatedAt,omitempty"`
	CreatedAt       int64      `json:"createdAt"`
	LastUsed        int64      `json:"lastUsed,omitempty"`
}

// CardTransaction represents a card transaction
type CardTransaction struct {
	ID           string            `json:"id"`
	CardID       string            `json:"cardId"`
	Amount       float64           `json:"amount"`
	Currency     string            `json:"currency"`
	Merchant     string            `json:"merchant"`
	MerchantCode string            `json:"merchantCode"`
	Status       TransactionStatus `json:"status"`
	Timestamp    int64             `json:"timestamp"`
	Location     string            `json:"location,omitempty"`
	Description  string            `json:"description,omitempty"`
}

// CardManager manages payment cards
type CardManager struct {
	cards        map[string]*PaymentCard
	transactions map[string][]*CardTransaction
	mutex        sync.RWMutex
}

// NewCardManager creates a new card manager
func NewCardManager() *CardManager {
	return &CardManager{
		cards:        make(map[string]*PaymentCard),
		transactions: make(map[string][]*CardTransaction),
	}
}

// CreateCard creates a new payment card
func (cm *CardManager) CreateCard(walletID, cardholderName string, cardType CardType, spendingLimit, dailyLimit, monthlyLimit float64) (*PaymentCard, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if walletID == "" || cardholderName == "" {
		return nil, errors.New("wallet ID and cardholder name are required")
	}

	cardNumber := generateCardNumber()
	cvv := generateCVV()
	expiryDate := time.Now().AddDate(3, 0, 0) // 3 years from now

	card := &PaymentCard{
		ID:             generateCardID(),
		WalletID:       walletID,
		CardNumber:     cardNumber,
		CardType:       cardType,
		ExpiryMonth:    int(expiryDate.Month()),
		ExpiryYear:     expiryDate.Year(),
		CVV:            cvv,
		CardholderName: cardholderName,
		Status:         CardStatusPending,
		SpendingLimit:  spendingLimit,
		DailyLimit:     dailyLimit,
		MonthlyLimit:   monthlyLimit,
		SpentToday:     0,
		SpentThisMonth: 0,
		CreatedAt:      time.Now().Unix(),
	}

	cm.cards[card.ID] = card
	cm.transactions[card.ID] = make([]*CardTransaction, 0)

	return card, nil
}

// ActivateCard activates a card
func (cm *CardManager) ActivateCard(cardID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return errors.New("card not found")
	}

	if card.Status != CardStatusPending {
		return errors.New("card cannot be activated")
	}

	card.Status = CardStatusActive
	card.ActivatedAt = time.Now().Unix()

	return nil
}

// SuspendCard suspends a card
func (cm *CardManager) SuspendCard(cardID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return errors.New("card not found")
	}

	if card.Status != CardStatusActive {
		return errors.New("only active cards can be suspended")
	}

	card.Status = CardStatusSuspended

	return nil
}

// CancelCard cancels a card
func (cm *CardManager) CancelCard(cardID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return errors.New("card not found")
	}

	card.Status = CardStatusCancelled

	return nil
}

// ProcessTransaction processes a card transaction
func (cm *CardManager) ProcessTransaction(cardID string, amount float64, currency, merchant, merchantCode string) (*CardTransaction, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return nil, errors.New("card not found")
	}

	if card.Status != CardStatusActive {
		return nil, errors.New("card is not active")
	}

	// Check spending limits
	if card.SpendingLimit > 0 && amount > card.SpendingLimit {
		return nil, errors.New("transaction exceeds spending limit")
	}

	if card.DailyLimit > 0 && (card.SpentToday+amount) > card.DailyLimit {
		return nil, errors.New("transaction exceeds daily limit")
	}

	if card.MonthlyLimit > 0 && (card.SpentThisMonth+amount) > card.MonthlyLimit {
		return nil, errors.New("transaction exceeds monthly limit")
	}

	// Create transaction
	tx := &CardTransaction{
		ID:           generateTransactionID(),
		CardID:       cardID,
		Amount:       amount,
		Currency:     currency,
		Merchant:     merchant,
		MerchantCode: merchantCode,
		Status:       TxStatusCompleted,
		Timestamp:    time.Now().Unix(),
	}

	// Update card spending
	card.SpentToday += amount
	card.SpentThisMonth += amount
	card.LastUsed = time.Now().Unix()

	// Store transaction
	cm.transactions[cardID] = append(cm.transactions[cardID], tx)

	return tx, nil
}

// GetCard returns a card by ID
func (cm *CardManager) GetCard(cardID string) (*PaymentCard, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return nil, errors.New("card not found")
	}

	return card, nil
}

// GetCardsByWallet returns all cards for a wallet
func (cm *CardManager) GetCardsByWallet(walletID string) []*PaymentCard {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	cards := make([]*PaymentCard, 0)
	for _, card := range cm.cards {
		if card.WalletID == walletID {
			cards = append(cards, card)
		}
	}

	return cards
}

// GetCardTransactions returns transactions for a card
func (cm *CardManager) GetCardTransactions(cardID string, limit int) ([]*CardTransaction, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	txs, exists := cm.transactions[cardID]
	if !exists {
		return nil, errors.New("card not found")
	}

	if limit > 0 && len(txs) > limit {
		return txs[:limit], nil
	}

	return txs, nil
}

// UpdateSpendingLimits updates card spending limits
func (cm *CardManager) UpdateSpendingLimits(cardID string, spendingLimit, dailyLimit, monthlyLimit float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return errors.New("card not found")
	}

	if spendingLimit > 0 {
		card.SpendingLimit = spendingLimit
	}
	if dailyLimit > 0 {
		card.DailyLimit = dailyLimit
	}
	if monthlyLimit > 0 {
		card.MonthlyLimit = monthlyLimit
	}

	return nil
}

// ResetDailySpending resets daily spending (should be called daily)
func (cm *CardManager) ResetDailySpending() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, card := range cm.cards {
		card.SpentToday = 0
	}
}

// ResetMonthlySpending resets monthly spending (should be called monthly)
func (cm *CardManager) ResetMonthlySpending() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, card := range cm.cards {
		card.SpentThisMonth = 0
	}
}

// Helper functions

func generateCardNumber() string {
	// Generate a 16-digit card number (simplified)
	prefix := "5399" // BTN-Pay prefix
	randomBytes := make([]byte, 6)
	rand.Read(randomBytes)
	suffix := hex.EncodeToString(randomBytes)
	return prefix + suffix
}

func generateCVV() string {
	randomBytes := make([]byte, 2)
	rand.Read(randomBytes)
	cvv := fmt.Sprintf("%03d", int(randomBytes[0])%1000)
	return cvv
}

func generateCardID() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	return "card_" + hex.EncodeToString(randomBytes)
}

func generateTransactionID() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	return "tx_" + hex.EncodeToString(randomBytes)
}
