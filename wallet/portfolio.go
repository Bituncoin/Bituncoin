package wallet

import (
	"sync"
	"time"
)

// Asset represents a cryptocurrency asset in the portfolio
type Asset struct {
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name"`
	Balance     float64   `json:"balance"`
	USDValue    float64   `json:"usdValue"`
	PriceUSD    float64   `json:"priceUSD"`
	Change24h   float64   `json:"change24h"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// Portfolio manages multiple cryptocurrency assets
type Portfolio struct {
	Assets      map[string]*Asset
	TotalUSD    float64
	LastUpdated time.Time
	mutex       sync.RWMutex
}

// NewPortfolio creates a new portfolio manager
func NewPortfolio() *Portfolio {
	return &Portfolio{
		Assets:      make(map[string]*Asset),
		LastUpdated: time.Now(),
	}
}

// AddAsset adds or updates an asset in the portfolio
func (p *Portfolio) AddAsset(symbol, name string, balance, priceUSD float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	asset := &Asset{
		Symbol:      symbol,
		Name:        name,
		Balance:     balance,
		PriceUSD:    priceUSD,
		USDValue:    balance * priceUSD,
		LastUpdated: time.Now(),
	}

	p.Assets[symbol] = asset
	p.updateTotalValue()
}

// UpdateBalance updates the balance of an asset
func (p *Portfolio) UpdateBalance(symbol string, newBalance float64) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	asset, exists := p.Assets[symbol]
	if !exists {
		return ErrAssetNotFound
	}

	asset.Balance = newBalance
	asset.USDValue = newBalance * asset.PriceUSD
	asset.LastUpdated = time.Now()

	p.updateTotalValue()
	return nil
}

// UpdatePrice updates the price of an asset
func (p *Portfolio) UpdatePrice(symbol string, newPrice, change24h float64) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	asset, exists := p.Assets[symbol]
	if !exists {
		return ErrAssetNotFound
	}

	asset.PriceUSD = newPrice
	asset.Change24h = change24h
	asset.USDValue = asset.Balance * newPrice
	asset.LastUpdated = time.Now()

	p.updateTotalValue()
	return nil
}

// GetAsset retrieves an asset from the portfolio
func (p *Portfolio) GetAsset(symbol string) (*Asset, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	asset, exists := p.Assets[symbol]
	if !exists {
		return nil, ErrAssetNotFound
	}

	return asset, nil
}

// GetAllAssets returns all assets in the portfolio
func (p *Portfolio) GetAllAssets() []*Asset {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	assets := make([]*Asset, 0, len(p.Assets))
	for _, asset := range p.Assets {
		assets = append(assets, asset)
	}

	return assets
}

// GetTotalValue returns the total portfolio value in USD
func (p *Portfolio) GetTotalValue() float64 {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.TotalUSD
}

// updateTotalValue recalculates the total portfolio value (internal, not thread-safe)
func (p *Portfolio) updateTotalValue() {
	total := 0.0
	for _, asset := range p.Assets {
		total += asset.USDValue
	}
	p.TotalUSD = total
	p.LastUpdated = time.Now()
}

// GetPerformance returns portfolio performance metrics
func (p *Portfolio) GetPerformance() map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	totalChange := 0.0
	assetCount := len(p.Assets)

	for _, asset := range p.Assets {
		totalChange += asset.Change24h
	}

	avgChange := 0.0
	if assetCount > 0 {
		avgChange = totalChange / float64(assetCount)
	}

	return map[string]interface{}{
		"totalValueUSD":   p.TotalUSD,
		"assetCount":      assetCount,
		"avgChange24h":    avgChange,
		"lastUpdated":     p.LastUpdated,
	}
}
