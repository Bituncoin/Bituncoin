package wallet

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ExchangeType represents the type of exchange
type ExchangeType string

const (
	ExchangeCryptoToCrypto ExchangeType = "crypto_to_crypto"
	ExchangeCryptoToFiat   ExchangeType = "crypto_to_fiat"
	ExchangeFiatToCrypto   ExchangeType = "fiat_to_crypto"
)

// ExchangeRate represents the exchange rate between two currencies
type ExchangeRate struct {
	FromCurrency string  `json:"fromCurrency"`
	ToCurrency   string  `json:"toCurrency"`
	Rate         float64 `json:"rate"`
	LastUpdate   int64   `json:"lastUpdate"`
}

// ExchangeOrder represents an exchange order
type ExchangeOrder struct {
	ID           string            `json:"id"`
	WalletID     string            `json:"walletId"`
	Type         ExchangeType      `json:"type"`
	FromCurrency string            `json:"fromCurrency"`
	ToCurrency   string            `json:"toCurrency"`
	FromAmount   float64           `json:"fromAmount"`
	ToAmount     float64           `json:"toAmount"`
	Rate         float64           `json:"rate"`
	Fee          float64           `json:"fee"`
	Status       TransactionStatus `json:"status"`
	CreatedAt    int64             `json:"createdAt"`
	CompletedAt  int64             `json:"completedAt,omitempty"`
	TxHash       string            `json:"txHash,omitempty"`
}

// ExchangeSystem manages cryptocurrency exchanges
type ExchangeSystem struct {
	rates  map[string]*ExchangeRate
	orders map[string]*ExchangeOrder
	mutex  sync.RWMutex
}

// NewExchangeSystem creates a new exchange system
func NewExchangeSystem() *ExchangeSystem {
	es := &ExchangeSystem{
		rates:  make(map[string]*ExchangeRate),
		orders: make(map[string]*ExchangeOrder),
	}

	// Initialize with default rates (in production, these would come from live APIs)
	es.initializeRates()

	return es
}

// initializeRates initializes exchange rates
func (es *ExchangeSystem) initializeRates() {
	// Mock rates for demonstration (real implementation would use live APIs)
	rates := []struct {
		from, to string
		rate     float64
	}{
		// BTN rates
		{"BTN", "USD", 10.0},
		{"BTN", "BTC", 0.00022},
		{"BTN", "ETH", 0.0033},
		{"BTN", "USDT", 10.0},
		{"BTN", "BNB", 0.016},

		// BTC rates
		{"BTC", "USD", 45000.0},
		{"BTC", "BTN", 4500.0},
		{"BTC", "ETH", 15.0},
		{"BTC", "USDT", 45000.0},
		{"BTC", "BNB", 75.0},

		// ETH rates
		{"ETH", "USD", 3000.0},
		{"ETH", "BTN", 300.0},
		{"ETH", "BTC", 0.067},
		{"ETH", "USDT", 3000.0},
		{"ETH", "BNB", 5.0},

		// USDT rates
		{"USDT", "USD", 1.0},
		{"USDT", "BTN", 0.1},
		{"USDT", "BTC", 0.000022},
		{"USDT", "ETH", 0.00033},
		{"USDT", "BNB", 0.0016},

		// BNB rates
		{"BNB", "USD", 600.0},
		{"BNB", "BTN", 60.0},
		{"BNB", "BTC", 0.013},
		{"BNB", "ETH", 0.2},
		{"BNB", "USDT", 600.0},
	}

	now := time.Now().Unix()
	for _, r := range rates {
		key := fmt.Sprintf("%s_%s", r.from, r.to)
		es.rates[key] = &ExchangeRate{
			FromCurrency: r.from,
			ToCurrency:   r.to,
			Rate:         r.rate,
			LastUpdate:   now,
		}
	}
}

// GetExchangeRate returns the exchange rate between two currencies
func (es *ExchangeSystem) GetExchangeRate(fromCurrency, toCurrency string) (*ExchangeRate, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	key := fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
	rate, exists := es.rates[key]
	if !exists {
		return nil, errors.New("exchange rate not available")
	}

	return rate, nil
}

// UpdateExchangeRate updates the exchange rate
func (es *ExchangeSystem) UpdateExchangeRate(fromCurrency, toCurrency string, rate float64) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	key := fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
	es.rates[key] = &ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
		LastUpdate:   time.Now().Unix(),
	}

	return nil
}

// CreateExchangeOrder creates a new exchange order
func (es *ExchangeSystem) CreateExchangeOrder(walletID, fromCurrency, toCurrency string, fromAmount float64) (*ExchangeOrder, error) {
	if fromAmount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get exchange rate
	rate, err := es.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}

	// Calculate amounts
	toAmount := fromAmount * rate.Rate
	fee := fromAmount * 0.001 // 0.1% fee

	// Determine exchange type
	exchangeType := ExchangeCryptoToCrypto
	if toCurrency == "USD" || toCurrency == "EUR" || toCurrency == "GBP" {
		exchangeType = ExchangeCryptoToFiat
	} else if fromCurrency == "USD" || fromCurrency == "EUR" || fromCurrency == "GBP" {
		exchangeType = ExchangeFiatToCrypto
	}

	order := &ExchangeOrder{
		ID:           generateExchangeID(),
		WalletID:     walletID,
		Type:         exchangeType,
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   fromAmount,
		ToAmount:     toAmount,
		Rate:         rate.Rate,
		Fee:          fee,
		Status:       TxStatusPending,
		CreatedAt:    time.Now().Unix(),
	}

	es.mutex.Lock()
	es.orders[order.ID] = order
	es.mutex.Unlock()

	return order, nil
}

// ExecuteExchangeOrder executes an exchange order
func (es *ExchangeSystem) ExecuteExchangeOrder(orderID string) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	order, exists := es.orders[orderID]
	if !exists {
		return errors.New("order not found")
	}

	if order.Status != TxStatusPending {
		return errors.New("order is not pending")
	}

	// In a real implementation, this would:
	// 1. Verify wallet has sufficient balance
	// 2. Lock the funds
	// 3. Execute the exchange
	// 4. Update balances
	// 5. Generate transaction hash

	order.Status = TxStatusCompleted
	order.CompletedAt = time.Now().Unix()
	order.TxHash = generateTxHash()

	return nil
}

// CancelExchangeOrder cancels an exchange order
func (es *ExchangeSystem) CancelExchangeOrder(orderID string) error {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	order, exists := es.orders[orderID]
	if !exists {
		return errors.New("order not found")
	}

	if order.Status != TxStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	order.Status = TxStatusCancelled

	return nil
}

// GetExchangeOrder returns an exchange order
func (es *ExchangeSystem) GetExchangeOrder(orderID string) (*ExchangeOrder, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	order, exists := es.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}

	return order, nil
}

// GetExchangeOrdersByWallet returns all exchange orders for a wallet
func (es *ExchangeSystem) GetExchangeOrdersByWallet(walletID string) []*ExchangeOrder {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	orders := make([]*ExchangeOrder, 0)
	for _, order := range es.orders {
		if order.WalletID == walletID {
			orders = append(orders, order)
		}
	}

	return orders
}

// CalculateExchange calculates the exchange amount without creating an order
func (es *ExchangeSystem) CalculateExchange(fromCurrency, toCurrency string, fromAmount float64) (float64, float64, error) {
	rate, err := es.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, 0, err
	}

	toAmount := fromAmount * rate.Rate
	fee := fromAmount * 0.001 // 0.1% fee

	return toAmount, fee, nil
}

// GetAllRates returns all available exchange rates
func (es *ExchangeSystem) GetAllRates() map[string]*ExchangeRate {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	rates := make(map[string]*ExchangeRate)
	for key, rate := range es.rates {
		rates[key] = rate
	}

	return rates
}

// RefreshRates refreshes all exchange rates (in production, this would fetch from APIs)
func (es *ExchangeSystem) RefreshRates() error {
	// In production, this would:
	// 1. Fetch rates from external APIs
	// 2. Update internal rates
	// 3. Notify subscribers of rate changes

	// For now, just update timestamps
	es.mutex.Lock()
	defer es.mutex.Unlock()

	now := time.Now().Unix()
	for _, rate := range es.rates {
		rate.LastUpdate = now
	}

	return nil
}

func generateExchangeID() string {
	return "exchange_" + time.Now().Format("20060102150405")
}

func generateTxHash() string {
	return fmt.Sprintf("0x%x", time.Now().UnixNano())
}
