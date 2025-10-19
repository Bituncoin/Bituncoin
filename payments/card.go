package payments

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CardType represents the type of payment card
type CardType string

const (
	MasterCard CardType = "MasterCard"
	Visa       CardType = "Visa"
)

// CardCategory represents virtual or physical card
type CardCategory string

const (
	Virtual  CardCategory = "Virtual"
	Physical CardCategory = "Physical"
)

// Card represents a BTN-Pay payment card
type Card struct {
	ID            string       `json:"id"`
	CardNumber    string       `json:"cardNumber"`
	Type          CardType     `json:"type"`
	Category      CardCategory `json:"category"`
	HolderName    string       `json:"holderName"`
	ExpiryMonth   int          `json:"expiryMonth"`
	ExpiryYear    int          `json:"expiryYear"`
	CVV           string       `json:"cvv"`
	Balance       float64      `json:"balance"`
	Currency      string       `json:"currency"`
	Status        string       `json:"status"` // active, frozen, expired
	LinkedAddress string       `json:"linkedAddress"`
	CreatedAt     int64        `json:"createdAt"`
	DailyLimit    float64      `json:"dailyLimit"`
	MonthlyLimit  float64      `json:"monthlyLimit"`
}

// CardTransaction represents a card payment transaction
type CardTransaction struct {
	ID            string    `json:"id"`
	CardID        string    `json:"cardId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Merchant      string    `json:"merchant"`
	MerchantCode  string    `json:"merchantCode"`
	Status        string    `json:"status"`
	TransactionAt int64     `json:"transactionAt"`
	Description   string    `json:"description"`
}

// CardService manages BTN-Pay payment cards
type CardService struct {
	cards        map[string]*Card
	transactions map[string]*CardTransaction
	mutex        sync.RWMutex
}

// NewCardService creates a new card service
func NewCardService() *CardService {
	return &CardService{
		cards:        make(map[string]*Card),
		transactions: make(map[string]*CardTransaction),
	}
}

// IssueCard issues a new BTN-Pay card
func (cs *CardService) IssueCard(holderName, linkedAddress string, cardType CardType, category CardCategory) (*Card, error) {
	if holderName == "" {
		return nil, errors.New("holder name required")
	}
	if linkedAddress == "" {
		return nil, errors.New("linked address required")
	}

	// Generate card number (mock)
	cardNumber := generateCardNumber(cardType)
	cvv := generateCVV()
	
	// Set expiry date (3 years from now)
	now := time.Now()
	expiryMonth := int(now.Month())
	expiryYear := now.Year() + 3

	card := &Card{
		ID:            fmt.Sprintf("card_%d", time.Now().UnixNano()),
		CardNumber:    cardNumber,
		Type:          cardType,
		Category:      category,
		HolderName:    holderName,
		ExpiryMonth:   expiryMonth,
		ExpiryYear:    expiryYear,
		CVV:           cvv,
		Balance:       0,
		Currency:      "USD",
		Status:        "active",
		LinkedAddress: linkedAddress,
		CreatedAt:     time.Now().Unix(),
		DailyLimit:    5000.0,
		MonthlyLimit:  50000.0,
	}

	cs.mutex.Lock()
	cs.cards[card.ID] = card
	cs.mutex.Unlock()

	return card, nil
}

// GetCard retrieves a card by ID
func (cs *CardService) GetCard(cardID string) (*Card, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	card, ok := cs.cards[cardID]
	if !ok {
		return nil, errors.New("card not found")
	}

	return card, nil
}

// LoadCard loads funds onto a card from wallet
func (cs *CardService) LoadCard(cardID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	card, ok := cs.cards[cardID]
	if !ok {
		return errors.New("card not found")
	}

	if card.Status != "active" {
		return errors.New("card is not active")
	}

	// In production, this would deduct from wallet balance
	card.Balance += amount

	return nil
}

// ProcessPayment processes a card payment
func (cs *CardService) ProcessPayment(cardID string, amount float64, merchant string) (*CardTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	card, ok := cs.cards[cardID]
	if !ok {
		return nil, errors.New("card not found")
	}

	if card.Status != "active" {
		return nil, errors.New("card is not active")
	}

	if card.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Check daily limit
	dailySpent := cs.getDailySpent(cardID)
	if dailySpent+amount > card.DailyLimit {
		return nil, errors.New("daily limit exceeded")
	}

	// Process payment
	card.Balance -= amount

	tx := &CardTransaction{
		ID:            fmt.Sprintf("tx_%d", time.Now().UnixNano()),
		CardID:        cardID,
		Amount:        amount,
		Currency:      card.Currency,
		Merchant:      merchant,
		MerchantCode:  generateMerchantCode(),
		Status:        "completed",
		TransactionAt: time.Now().Unix(),
		Description:   fmt.Sprintf("Payment to %s", merchant),
	}

	cs.transactions[tx.ID] = tx

	return tx, nil
}

// FreezeCard freezes a card
func (cs *CardService) FreezeCard(cardID string) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	card, ok := cs.cards[cardID]
	if !ok {
		return errors.New("card not found")
	}

	card.Status = "frozen"
	return nil
}

// UnfreezeCard unfreezes a card
func (cs *CardService) UnfreezeCard(cardID string) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	card, ok := cs.cards[cardID]
	if !ok {
		return errors.New("card not found")
	}

	if card.Status == "expired" {
		return errors.New("cannot unfreeze expired card")
	}

	card.Status = "active"
	return nil
}

// GetCardTransactions retrieves transactions for a card
func (cs *CardService) GetCardTransactions(cardID string) ([]*CardTransaction, error) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	if _, ok := cs.cards[cardID]; !ok {
		return nil, errors.New("card not found")
	}

	var transactions []*CardTransaction
	for _, tx := range cs.transactions {
		if tx.CardID == cardID {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// getDailySpent calculates total spent today
func (cs *CardService) getDailySpent(cardID string) float64 {
	today := time.Now().Truncate(24 * time.Hour).Unix()
	spent := 0.0

	for _, tx := range cs.transactions {
		if tx.CardID == cardID && tx.TransactionAt >= today && tx.Status == "completed" {
			spent += tx.Amount
		}
	}

	return spent
}

// Helper functions

func generateCardNumber(cardType CardType) string {
	// Generate mock card number
	prefix := "5"
	if cardType == Visa {
		prefix = "4"
	}
	
	// Mock card number (not a real algorithm)
	return fmt.Sprintf("%s%d%d%d %d%d%d%d %d%d%d%d %d%d%d%d",
		prefix,
		time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10,
		time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10,
		time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10,
		time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10, time.Now().Unix()%10,
	)
}

func generateCVV() string {
	return fmt.Sprintf("%03d", time.Now().Unix()%1000)
}

func generateMerchantCode() string {
	return fmt.Sprintf("MRC%d", time.Now().Unix()%1000000)
}

// UpdateCardLimits updates card spending limits
func (cs *CardService) UpdateCardLimits(cardID string, dailyLimit, monthlyLimit float64) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	card, ok := cs.cards[cardID]
	if !ok {
		return errors.New("card not found")
	}

	if dailyLimit > 0 {
		card.DailyLimit = dailyLimit
	}
	if monthlyLimit > 0 {
		card.MonthlyLimit = monthlyLimit
	}

	return nil
}
