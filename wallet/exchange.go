package wallet

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ExchangeRate represents exchange rate between two currencies
type ExchangeRate struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Rate      float64   `json:"rate"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SwapTransaction represents a crypto swap transaction
type SwapTransaction struct {
	ID          string    `json:"id"`
	FromCurrency string   `json:"fromCurrency"`
	ToCurrency   string   `json:"toCurrency"`
	FromAmount   float64  `json:"fromAmount"`
	ToAmount     float64  `json:"toAmount"`
	ExchangeRate float64  `json:"exchangeRate"`
	Fee          float64  `json:"fee"`
	Status       string   `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	CompletedAt  time.Time `json:"completedAt,omitempty"`
}

// ExchangeService provides crypto-to-crypto and crypto-to-fiat exchange
type ExchangeService struct {
	rates        map[string]ExchangeRate
	transactions map[string]*SwapTransaction
	mutex        sync.RWMutex
}

// NewExchangeService creates a new exchange service
func NewExchangeService() *ExchangeService {
	es := &ExchangeService{
		rates:        make(map[string]ExchangeRate),
		transactions: make(map[string]*SwapTransaction),
	}
	
	// Initialize mock exchange rates
	es.initializeRates()
	
	return es
}

// initializeRates initializes exchange rates (mock data)
func (es *ExchangeService) initializeRates() {
	// Base rates in USD
	baseRates := map[string]float64{
		"BTN":  15.0,
		"GLD":  10.0,
		"BTC":  45000.0,
		"ETH":  3000.0,
		"USDT": 1.0,
		"BNB":  300.0,
	}

	now := time.Now()
	
	// Generate all currency pair rates
	currencies := []string{"BTN", "GLD", "BTC", "ETH", "USDT", "BNB"}
	for _, from := range currencies {
		for _, to := range currencies {
			if from != to {
				key := fmt.Sprintf("%s-%s", from, to)
				rate := baseRates[to] / baseRates[from]
				es.rates[key] = ExchangeRate{
					From:      from,
					To:        to,
					Rate:      rate,
					UpdatedAt: now,
				}
			}
		}
	}
}

// GetExchangeRate retrieves the exchange rate between two currencies
func (es *ExchangeService) GetExchangeRate(from, to string) (*ExchangeRate, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	key := fmt.Sprintf("%s-%s", from, to)
	rate, ok := es.rates[key]
	if !ok {
		return nil, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return &rate, nil
}

// CalculateSwap calculates the swap amount and fees
func (es *ExchangeService) CalculateSwap(fromCurrency, toCurrency string, amount float64) (toAmount float64, fee float64, err error) {
	if amount <= 0 {
		return 0, 0, errors.New("amount must be positive")
	}

	rate, err := es.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		return 0, 0, err
	}

	// Calculate base amount
	baseAmount := amount * rate.Rate

	// Calculate fee (0.5% of base amount)
	fee = baseAmount * 0.005

	// Final amount after fee
	toAmount = baseAmount - fee

	return toAmount, fee, nil
}

// ExecuteSwap executes a currency swap
func (es *ExchangeService) ExecuteSwap(fromCurrency, toCurrency string, amount float64, userAddress string) (*SwapTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	toAmount, fee, err := es.CalculateSwap(fromCurrency, toCurrency, amount)
	if err != nil {
		return nil, err
	}

	rate, _ := es.GetExchangeRate(fromCurrency, toCurrency)

	// Create swap transaction
	tx := &SwapTransaction{
		ID:           fmt.Sprintf("swap_%d", time.Now().UnixNano()),
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		FromAmount:   amount,
		ToAmount:     toAmount,
		ExchangeRate: rate.Rate,
		Fee:          fee,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	es.mutex.Lock()
	es.transactions[tx.ID] = tx
	es.mutex.Unlock()

	// Simulate swap processing
	go es.processSwap(tx.ID)

	return tx, nil
}

// processSwap simulates swap processing
func (es *ExchangeService) processSwap(txID string) {
	// Simulate processing time
	time.Sleep(2 * time.Second)

	es.mutex.Lock()
	defer es.mutex.Unlock()

	if tx, ok := es.transactions[txID]; ok {
		tx.Status = "completed"
		tx.CompletedAt = time.Now()
	}
}

// GetSwapTransaction retrieves a swap transaction by ID
func (es *ExchangeService) GetSwapTransaction(txID string) (*SwapTransaction, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	tx, ok := es.transactions[txID]
	if !ok {
		return nil, errors.New("transaction not found")
	}

	return tx, nil
}

// GetSwapHistory retrieves swap transaction history
func (es *ExchangeService) GetSwapHistory(userAddress string, limit int) ([]*SwapTransaction, error) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	var history []*SwapTransaction
	count := 0

	// In production, filter by userAddress
	for _, tx := range es.transactions {
		if limit > 0 && count >= limit {
			break
		}
		history = append(history, tx)
		count++
	}

	return history, nil
}

// CryptoToFiatExchange converts crypto to fiat currency
func (es *ExchangeService) CryptoToFiatExchange(cryptoCurrency string, amount float64, fiatCurrency string) (float64, error) {
	// Base rates in USD (how much 1 unit is worth in USD)
	baseRates := map[string]float64{
		"BTN":  15.0,
		"GLD":  10.0,
		"BTC":  45000.0,
		"ETH":  3000.0,
		"USDT": 1.0,
		"BNB":  300.0,
	}

	// Get crypto value in USD
	cryptoRate, ok := baseRates[cryptoCurrency]
	if !ok {
		return 0, fmt.Errorf("cryptocurrency %s not supported", cryptoCurrency)
	}
	
	usdAmount := amount * cryptoRate

	// Convert USD to target fiat (mock rates)
	fiatRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.73,
		"GHS": 12.0,  // Ghanaian Cedi
		"NGN": 770.0, // Nigerian Naira
		"KES": 130.0, // Kenyan Shilling
	}

	conversionRate, ok := fiatRates[fiatCurrency]
	if !ok {
		return 0, fmt.Errorf("fiat currency %s not supported", fiatCurrency)
	}

	return usdAmount * conversionRate, nil
}

// FiatToCryptoExchange converts fiat to crypto currency
func (es *ExchangeService) FiatToCryptoExchange(fiatCurrency string, amount float64, cryptoCurrency string) (float64, error) {
	// Convert fiat to USD first
	fiatRates := map[string]float64{
		"USD": 1.0,
		"EUR": 0.85,
		"GBP": 0.73,
		"GHS": 12.0,
		"NGN": 770.0,
		"KES": 130.0,
	}

	conversionRate, ok := fiatRates[fiatCurrency]
	if !ok {
		return 0, fmt.Errorf("fiat currency %s not supported", fiatCurrency)
	}

	usdAmount := amount / conversionRate

	// Get USD to crypto rate
	rate, err := es.GetExchangeRate("USDT", cryptoCurrency)
	if err != nil {
		return 0, err
	}

	return usdAmount * rate.Rate, nil
}

// UpdateExchangeRates updates exchange rates (simulates price feed)
func (es *ExchangeService) UpdateExchangeRates() {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	// In production, fetch from external price feeds
	// For now, just update timestamps
	now := time.Now()
	for key, rate := range es.rates {
		rate.UpdatedAt = now
		es.rates[key] = rate
	}
}

// GetSupportedCurrencies returns list of supported currencies
func (es *ExchangeService) GetSupportedCurrencies() []string {
	return []string{"BTN", "GLD", "BTC", "ETH", "USDT", "BNB"}
}

// EstimateSwapTime estimates time for swap completion
func (es *ExchangeService) EstimateSwapTime(fromCurrency, toCurrency string) time.Duration {
	// Mock estimation based on currency type
	if fromCurrency == "BTC" || toCurrency == "BTC" {
		return 30 * time.Minute // BTC transactions take longer
	}
	if fromCurrency == "ETH" || toCurrency == "ETH" {
		return 5 * time.Minute
	}
	return 2 * time.Minute // Fast transactions for native tokens
}
