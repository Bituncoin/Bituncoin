package wallet

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrInvalidExchangePair = errors.New("invalid exchange pair")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrExchangeNotFound    = errors.New("exchange not found")
)

// ExchangeType represents the type of exchange
type ExchangeType string

const (
	ExchangeCryptoToCrypto ExchangeType = "crypto-to-crypto"
	ExchangeCryptoToFiat   ExchangeType = "crypto-to-fiat"
	ExchangeFiatToCrypto   ExchangeType = "fiat-to-crypto"
)

// ExchangeRate stores exchange rate information
type ExchangeRate struct {
	FromAsset   string    `json:"fromAsset"`
	ToAsset     string    `json:"toAsset"`
	Rate        float64   `json:"rate"`
	Fee         float64   `json:"fee"` // percentage
	LastUpdated time.Time `json:"lastUpdated"`
}

// ExchangeOrder represents an exchange order
type ExchangeOrder struct {
	ID          string       `json:"id"`
	Type        ExchangeType `json:"type"`
	FromAsset   string       `json:"fromAsset"`
	ToAsset     string       `json:"toAsset"`
	FromAmount  float64      `json:"fromAmount"`
	ToAmount    float64      `json:"toAmount"`
	Rate        float64      `json:"rate"`
	Fee         float64      `json:"fee"`
	Status      string       `json:"status"`
	UserAddress string       `json:"userAddress"`
	Timestamp   time.Time    `json:"timestamp"`
}

// Exchange manages cryptocurrency exchange operations
type Exchange struct {
	rates       map[string]*ExchangeRate // "FROM-TO" -> rate
	orders      map[string]*ExchangeOrder
	supportedFiat []string
	mutex       sync.RWMutex
}

// NewExchange creates a new exchange manager
func NewExchange() *Exchange {
	ex := &Exchange{
		rates:       make(map[string]*ExchangeRate),
		orders:      make(map[string]*ExchangeOrder),
		supportedFiat: []string{"USD", "EUR", "GBP", "JPY", "CNY"},
	}

	// Initialize with some default rates (in production, these would be fetched from APIs)
	ex.initializeDefaultRates()

	return ex
}

// initializeDefaultRates sets up initial exchange rates
func (ex *Exchange) initializeDefaultRates() {
	// Crypto-to-crypto pairs
	ex.SetExchangeRate("BTNG", "BTC", 0.00002, 0.5)
	ex.SetExchangeRate("BTNG", "ETH", 0.0003, 0.5)
	ex.SetExchangeRate("BTNG", "USDT", 1.0, 0.3)
	ex.SetExchangeRate("BTNG", "BNB", 0.002, 0.5)
	ex.SetExchangeRate("GLD", "BTC", 0.00002, 0.5)
	ex.SetExchangeRate("GLD", "ETH", 0.0003, 0.5)
	ex.SetExchangeRate("BTC", "ETH", 15.0, 0.5)
	ex.SetExchangeRate("BTC", "USDT", 50000.0, 0.3)
	ex.SetExchangeRate("ETH", "USDT", 3000.0, 0.3)

	// Crypto-to-fiat pairs
	ex.SetExchangeRate("BTNG", "USD", 1.0, 1.0)
	ex.SetExchangeRate("GLD", "USD", 1.0, 1.0)
	ex.SetExchangeRate("BTC", "USD", 50000.0, 1.0)
	ex.SetExchangeRate("ETH", "USD", 3000.0, 1.0)
	ex.SetExchangeRate("USDT", "USD", 1.0, 0.2)
	ex.SetExchangeRate("BNB", "USD", 500.0, 1.0)
}

// SetExchangeRate sets or updates an exchange rate
func (ex *Exchange) SetExchangeRate(fromAsset, toAsset string, rate, fee float64) {
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	key := fmt.Sprintf("%s-%s", fromAsset, toAsset)
	ex.rates[key] = &ExchangeRate{
		FromAsset:   fromAsset,
		ToAsset:     toAsset,
		Rate:        rate,
		Fee:         fee,
		LastUpdated: time.Now(),
	}

	// Set reverse rate
	reverseKey := fmt.Sprintf("%s-%s", toAsset, fromAsset)
	ex.rates[reverseKey] = &ExchangeRate{
		FromAsset:   toAsset,
		ToAsset:     fromAsset,
		Rate:        1.0 / rate,
		Fee:         fee,
		LastUpdated: time.Now(),
	}
}

// GetExchangeRate retrieves the exchange rate for a pair
func (ex *Exchange) GetExchangeRate(fromAsset, toAsset string) (*ExchangeRate, error) {
	ex.mutex.RLock()
	defer ex.mutex.RUnlock()

	key := fmt.Sprintf("%s-%s", fromAsset, toAsset)
	rate, exists := ex.rates[key]
	if !exists {
		return nil, ErrInvalidExchangePair
	}

	return rate, nil
}

// CalculateExchange calculates the exchange output amount including fees
func (ex *Exchange) CalculateExchange(fromAsset, toAsset string, fromAmount float64) (float64, float64, error) {
	rate, err := ex.GetExchangeRate(fromAsset, toAsset)
	if err != nil {
		return 0, 0, err
	}

	// Calculate base amount
	toAmount := fromAmount * rate.Rate

	// Calculate fee
	fee := toAmount * (rate.Fee / 100.0)

	// Final amount after fee
	finalAmount := toAmount - fee

	return finalAmount, fee, nil
}

// CreateExchangeOrder creates a new exchange order
func (ex *Exchange) CreateExchangeOrder(userAddress, fromAsset, toAsset string, fromAmount float64) (*ExchangeOrder, error) {
	toAmount, fee, err := ex.CalculateExchange(fromAsset, toAsset, fromAmount)
	if err != nil {
		return nil, err
	}

	rate, _ := ex.GetExchangeRate(fromAsset, toAsset)

	// Determine exchange type
	exchangeType := ExchangeCryptoToCrypto
	if ex.isFiat(toAsset) {
		exchangeType = ExchangeCryptoToFiat
	} else if ex.isFiat(fromAsset) {
		exchangeType = ExchangeFiatToCrypto
	}

	order := &ExchangeOrder{
		ID:          fmt.Sprintf("EX-%d", time.Now().UnixNano()),
		Type:        exchangeType,
		FromAsset:   fromAsset,
		ToAsset:     toAsset,
		FromAmount:  fromAmount,
		ToAmount:    toAmount,
		Rate:        rate.Rate,
		Fee:         fee,
		Status:      "pending",
		UserAddress: userAddress,
		Timestamp:   time.Now(),
	}

	ex.mutex.Lock()
	ex.orders[order.ID] = order
	ex.mutex.Unlock()

	return order, nil
}

// GetExchangeOrder retrieves an exchange order
func (ex *Exchange) GetExchangeOrder(orderID string) (*ExchangeOrder, error) {
	ex.mutex.RLock()
	defer ex.mutex.RUnlock()

	order, exists := ex.orders[orderID]
	if !exists {
		return nil, ErrExchangeNotFound
	}

	return order, nil
}

// UpdateOrderStatus updates the status of an exchange order
func (ex *Exchange) UpdateOrderStatus(orderID, status string) error {
	ex.mutex.Lock()
	defer ex.mutex.Unlock()

	order, exists := ex.orders[orderID]
	if !exists {
		return ErrExchangeNotFound
	}

	order.Status = status
	return nil
}

// GetUserOrders retrieves all orders for a user
func (ex *Exchange) GetUserOrders(userAddress string) []*ExchangeOrder {
	ex.mutex.RLock()
	defer ex.mutex.RUnlock()

	orders := make([]*ExchangeOrder, 0)
	for _, order := range ex.orders {
		if order.UserAddress == userAddress {
			orders = append(orders, order)
		}
	}

	return orders
}

// isFiat checks if an asset is a fiat currency
func (ex *Exchange) isFiat(asset string) bool {
	for _, fiat := range ex.supportedFiat {
		if asset == fiat {
			return true
		}
	}
	return false
}

// GetSupportedPairs returns all supported exchange pairs
func (ex *Exchange) GetSupportedPairs() []string {
	ex.mutex.RLock()
	defer ex.mutex.RUnlock()

	pairs := make([]string, 0, len(ex.rates))
	for key := range ex.rates {
		pairs = append(pairs, key)
	}

	return pairs
}
