package wallet

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrCardNotFound     = errors.New("card not found")
	ErrInvalidCardType  = errors.New("invalid card type")
	ErrCardLimitExceeded = errors.New("card limit exceeded")
)

// CardType represents the type of payment card
type CardType string

const (
	CardTypeVirtual  CardType = "virtual"
	CardTypePhysical CardType = "physical"
)

// CardProvider represents the card network provider
type CardProvider string

const (
	ProviderMasterCard CardProvider = "mastercard"
	ProviderVisa       CardProvider = "visa"
)

// CardStatus represents the status of a card
type CardStatus string

const (
	CardStatusActive   CardStatus = "active"
	CardStatusInactive CardStatus = "inactive"
	CardStatusBlocked  CardStatus = "blocked"
	CardStatusExpired  CardStatus = "expired"
)

// PaymentCard represents a BTN-Pay card
type PaymentCard struct {
	ID            string       `json:"id"`
	UserAddress   string       `json:"userAddress"`
	CardNumber    string       `json:"cardNumber"`
	CardType      CardType     `json:"cardType"`
	Provider      CardProvider `json:"provider"`
	Status        CardStatus   `json:"status"`
	Balance       float64      `json:"balance"`
	DailyLimit    float64      `json:"dailyLimit"`
	DailySpent    float64      `json:"dailySpent"`
	ExpiryDate    string       `json:"expiryDate"`
	CVV           string       `json:"cvv,omitempty"`
	CreatedAt     time.Time    `json:"createdAt"`
	LastUsed      time.Time    `json:"lastUsed,omitempty"`
}

// CardTransaction represents a card transaction
type CardTransaction struct {
	ID          string    `json:"id"`
	CardID      string    `json:"cardId"`
	Merchant    string    `json:"merchant"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Type        string    `json:"type"` // purchase, refund, withdrawal
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

// CardManager manages payment card operations
type CardManager struct {
	cards        map[string]*PaymentCard
	transactions map[string]*CardTransaction
	userCards    map[string][]string // userAddress -> []cardID
	mutex        sync.RWMutex
}

// NewCardManager creates a new card manager
func NewCardManager() *CardManager {
	return &CardManager{
		cards:        make(map[string]*PaymentCard),
		transactions: make(map[string]*CardTransaction),
		userCards:    make(map[string][]string),
	}
}

// CreateCard creates a new payment card
func (cm *CardManager) CreateCard(userAddress string, cardType CardType, provider CardProvider, dailyLimit float64) (*PaymentCard, error) {
	if cardType != CardTypeVirtual && cardType != CardTypePhysical {
		return nil, ErrInvalidCardType
	}

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Generate card details
	cardNumber := cm.generateCardNumber(provider)
	cvv := cm.generateCVV()
	expiryDate := cm.generateExpiryDate()

	card := &PaymentCard{
		ID:          fmt.Sprintf("CARD-%d", time.Now().UnixNano()),
		UserAddress: userAddress,
		CardNumber:  cardNumber,
		CardType:    cardType,
		Provider:    provider,
		Status:      CardStatusActive,
		Balance:     0.0,
		DailyLimit:  dailyLimit,
		DailySpent:  0.0,
		ExpiryDate:  expiryDate,
		CVV:         cvv,
		CreatedAt:   time.Now(),
	}

	cm.cards[card.ID] = card
	cm.userCards[userAddress] = append(cm.userCards[userAddress], card.ID)

	return card, nil
}

// GetCard retrieves a card by ID
func (cm *CardManager) GetCard(cardID string) (*PaymentCard, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return nil, ErrCardNotFound
	}

	return card, nil
}

// GetUserCards retrieves all cards for a user
func (cm *CardManager) GetUserCards(userAddress string) []*PaymentCard {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	cardIDs := cm.userCards[userAddress]
	cards := make([]*PaymentCard, 0, len(cardIDs))

	for _, cardID := range cardIDs {
		if card, exists := cm.cards[cardID]; exists {
			cards = append(cards, card)
		}
	}

	return cards
}

// TopUpCard adds balance to a card
func (cm *CardManager) TopUpCard(cardID string, amount float64) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return ErrCardNotFound
	}

	if card.Status != CardStatusActive {
		return errors.New("card is not active")
	}

	card.Balance += amount
	return nil
}

// ProcessCardTransaction processes a card transaction
func (cm *CardManager) ProcessCardTransaction(cardID, merchant string, amount float64, txType string) (*CardTransaction, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return nil, ErrCardNotFound
	}

	if card.Status != CardStatusActive {
		return nil, errors.New("card is not active")
	}

	// Check daily limit
	if card.DailySpent+amount > card.DailyLimit {
		return nil, ErrCardLimitExceeded
	}

	// Check balance
	if card.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	// Create transaction
	tx := &CardTransaction{
		ID:          fmt.Sprintf("CTX-%d", time.Now().UnixNano()),
		CardID:      cardID,
		Merchant:    merchant,
		Amount:      amount,
		Currency:    "USD",
		Status:      "completed",
		Type:        txType,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Payment to %s", merchant),
	}

	// Update card
	card.Balance -= amount
	card.DailySpent += amount
	card.LastUsed = time.Now()

	cm.transactions[tx.ID] = tx

	return tx, nil
}

// GetCardTransactions retrieves all transactions for a card
func (cm *CardManager) GetCardTransactions(cardID string) []*CardTransaction {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	transactions := make([]*CardTransaction, 0)
	for _, tx := range cm.transactions {
		if tx.CardID == cardID {
			transactions = append(transactions, tx)
		}
	}

	return transactions
}

// UpdateCardStatus updates the status of a card
func (cm *CardManager) UpdateCardStatus(cardID string, status CardStatus) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	card, exists := cm.cards[cardID]
	if !exists {
		return ErrCardNotFound
	}

	card.Status = status
	return nil
}

// ResetDailySpent resets daily spending (should be called daily)
func (cm *CardManager) ResetDailySpent() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for _, card := range cm.cards {
		card.DailySpent = 0.0
	}
}

// Helper functions

func (cm *CardManager) generateCardNumber(provider CardProvider) string {
	prefix := "4" // Visa starts with 4
	if provider == ProviderMasterCard {
		prefix = "5" // MasterCard starts with 5
	}
	
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s%d", prefix, timestamp%10000000000000000)
}

func (cm *CardManager) generateCVV() string {
	return fmt.Sprintf("%03d", time.Now().UnixNano()%1000)
}

func (cm *CardManager) generateExpiryDate() string {
	now := time.Now()
	expiryDate := now.AddDate(3, 0, 0) // 3 years from now
	return expiryDate.Format("01/06")
}
