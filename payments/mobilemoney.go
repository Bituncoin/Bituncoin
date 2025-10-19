package payments

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// MobileMoneyProvider represents supported mobile money services
type MobileMoneyProvider string

const (
	MTN        MobileMoneyProvider = "MTN"
	AirtelTigo MobileMoneyProvider = "AirtelTigo"
	Vodafone   MobileMoneyProvider = "Vodafone"
)

// MobileMoneyTransaction represents a mobile money transaction
type MobileMoneyTransaction struct {
	ID          string              `json:"id"`
	Provider    MobileMoneyProvider `json:"provider"`
	PhoneNumber string              `json:"phoneNumber"`
	Amount      float64             `json:"amount"`
	Currency    string              `json:"currency"`
	Status      string              `json:"status"`
	Reference   string              `json:"reference"`
	CreatedAt   int64               `json:"createdAt"`
	CompletedAt int64               `json:"completedAt,omitempty"`
}

// MobileMoneyService handles mobile money integrations
type MobileMoneyService struct {
	transactions map[string]*MobileMoneyTransaction
	mutex        sync.RWMutex
}

// NewMobileMoneyService creates a new mobile money service
func NewMobileMoneyService() *MobileMoneyService {
	return &MobileMoneyService{
		transactions: make(map[string]*MobileMoneyTransaction),
	}
}

// InitiatePayment initiates a mobile money payment
func (mms *MobileMoneyService) InitiatePayment(provider MobileMoneyProvider, phoneNumber string, amount float64, currency string) (*MobileMoneyTransaction, error) {
	if phoneNumber == "" {
		return nil, errors.New("phone number required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if !isValidProvider(provider) {
		return nil, errors.New("invalid mobile money provider")
	}

	id := fmt.Sprintf("momo_%d", time.Now().UnixNano())
	tx := &MobileMoneyTransaction{
		ID:          id,
		Provider:    provider,
		PhoneNumber: phoneNumber,
		Amount:      amount,
		Currency:    currency,
		Status:      "pending",
		Reference:   generateReference(),
		CreatedAt:   time.Now().Unix(),
	}

	mms.mutex.Lock()
	mms.transactions[id] = tx
	mms.mutex.Unlock()

	// In production, this would trigger actual mobile money API calls
	// For now, we simulate the payment flow
	go mms.simulatePayment(id)

	return tx, nil
}

// GetTransaction retrieves a transaction by ID
func (mms *MobileMoneyService) GetTransaction(id string) (*MobileMoneyTransaction, error) {
	mms.mutex.RLock()
	defer mms.mutex.RUnlock()

	tx, ok := mms.transactions[id]
	if !ok {
		return nil, errors.New("transaction not found")
	}

	return tx, nil
}

// simulatePayment simulates the payment processing
func (mms *MobileMoneyService) simulatePayment(id string) {
	// Simulate processing time
	time.Sleep(2 * time.Second)

	mms.mutex.Lock()
	defer mms.mutex.Unlock()

	if tx, ok := mms.transactions[id]; ok {
		tx.Status = "completed"
		tx.CompletedAt = time.Now().Unix()
	}
}

// isValidProvider checks if provider is valid
func isValidProvider(provider MobileMoneyProvider) bool {
	switch provider {
	case MTN, AirtelTigo, Vodafone:
		return true
	default:
		return false
	}
}

// generateReference generates a unique reference for the transaction
func generateReference() string {
	return fmt.Sprintf("REF%d", time.Now().UnixNano()%1000000000)
}

// WithdrawToMobileMoney withdraws crypto to mobile money
func (mms *MobileMoneyService) WithdrawToMobileMoney(provider MobileMoneyProvider, phoneNumber string, cryptoAmount float64, cryptoCurrency string) (*MobileMoneyTransaction, error) {
	// Convert crypto to fiat (simplified conversion)
	conversionRate := getConversionRate(cryptoCurrency)
	fiatAmount := cryptoAmount * conversionRate

	return mms.InitiatePayment(provider, phoneNumber, fiatAmount, "GHS")
}

// getConversionRate returns mock conversion rates
func getConversionRate(currency string) float64 {
	rates := map[string]float64{
		"BTN":  15.0,  // 1 BTN = 15 GHS
		"GLD":  10.0,  // 1 GLD = 10 GHS
		"BTC":  450000.0,
		"ETH":  30000.0,
		"USDT": 10.0,
		"BNB":  3000.0,
	}
	if rate, ok := rates[currency]; ok {
		return rate
	}
	return 1.0
}
