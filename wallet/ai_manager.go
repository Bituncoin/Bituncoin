package wallet

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

// InsightType represents the type of AI insight
type InsightType string

const (
	InsightSpending     InsightType = "spending"
	InsightStaking      InsightType = "staking"
	InsightTrading      InsightType = "trading"
	InsightMarket       InsightType = "market"
	InsightSecurity     InsightType = "security"
	InsightOptimization InsightType = "optimization"
)

// AlertPriority represents the priority level of an alert
type AlertPriority string

const (
	PriorityLow      AlertPriority = "low"
	PriorityMedium   AlertPriority = "medium"
	PriorityHigh     AlertPriority = "high"
	PriorityCritical AlertPriority = "critical"
)

// AIInsight represents an AI-generated insight
type AIInsight struct {
	ID          string        `json:"id"`
	Type        InsightType   `json:"type"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Priority    AlertPriority `json:"priority"`
	Data        map[string]interface{} `json:"data"`
	CreatedAt   time.Time     `json:"createdAt"`
	ActionItems []string      `json:"actionItems,omitempty"`
}

// MarketAlert represents a market trend alert
type MarketAlert struct {
	ID          string        `json:"id"`
	Asset       string        `json:"asset"`
	AlertType   string        `json:"alertType"` // price_increase, price_decrease, volatility
	Message     string        `json:"message"`
	Priority    AlertPriority `json:"priority"`
	CurrentPrice float64      `json:"currentPrice"`
	TargetPrice  float64      `json:"targetPrice,omitempty"`
	PercentChange float64     `json:"percentChange"`
	CreatedAt   time.Time     `json:"createdAt"`
	IsActive    bool          `json:"isActive"`
}

// Recommendation represents an AI recommendation
type Recommendation struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // stake, swap, hold, sell
	Asset       string    `json:"asset"`
	Amount      float64   `json:"amount"`
	Reason      string    `json:"reason"`
	Confidence  float64   `json:"confidence"` // 0-1
	PotentialROI float64  `json:"potentialRoi,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// AIWalletManager provides AI-driven wallet management
type AIWalletManager struct {
	insights        map[string]*AIInsight
	alerts          map[string]*MarketAlert
	recommendations map[string]*Recommendation
	spendingPatterns map[string][]float64 // asset -> daily spending
	mutex           sync.RWMutex
}

// NewAIWalletManager creates a new AI wallet manager
func NewAIWalletManager() *AIWalletManager {
	return &AIWalletManager{
		insights:        make(map[string]*AIInsight),
		alerts:          make(map[string]*MarketAlert),
		recommendations: make(map[string]*Recommendation),
		spendingPatterns: make(map[string][]float64),
	}
}

// AnalyzeSpending analyzes spending patterns and generates insights
func (ai *AIWalletManager) AnalyzeSpending(userAddress string, transactions []*Transaction) []*AIInsight {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	insights := make([]*AIInsight, 0)

	// Calculate daily spending
	dailySpending := make(map[string]float64)
	totalSpent := 0.0

	for _, tx := range transactions {
		if tx.Type == TypeSent {
			date := tx.Timestamp.Format("2006-01-02")
			dailySpending[date] += tx.Amount
			totalSpent += tx.Amount
		}
	}

	// Calculate average daily spending
	avgSpending := 0.0
	if len(dailySpending) > 0 {
		avgSpending = totalSpent / float64(len(dailySpending))
	}

	// Generate spending insight
	insight := &AIInsight{
		ID:          fmt.Sprintf("INS-SP-%d", time.Now().UnixNano()),
		Type:        InsightSpending,
		Title:       "Spending Pattern Analysis",
		Description: fmt.Sprintf("Your average daily spending is $%.2f", avgSpending),
		Priority:    PriorityMedium,
		Data: map[string]interface{}{
			"totalSpent":    totalSpent,
			"avgDaily":      avgSpending,
			"activeDays":    len(dailySpending),
		},
		CreatedAt: time.Now(),
		ActionItems: []string{
			"Consider setting up spending limits",
			"Review recurring transactions",
		},
	}

	ai.insights[insight.ID] = insight
	insights = append(insights, insight)

	return insights
}

// GenerateStakingRecommendation generates staking recommendations
func (ai *AIWalletManager) GenerateStakingRecommendation(asset string, balance float64, currentAPY float64) *Recommendation {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	// Simple recommendation logic: if balance > 100 and APY > 3%, recommend staking
	if balance > 100 && currentAPY > 3.0 {
		rec := &Recommendation{
			ID:          fmt.Sprintf("REC-STK-%d", time.Now().UnixNano()),
			Type:        "stake",
			Asset:       asset,
			Amount:      balance * 0.7, // Recommend staking 70%
			Reason:      fmt.Sprintf("Current APY of %.2f%% provides good returns. Staking 70%% maintains liquidity.", currentAPY),
			Confidence:  0.85,
			PotentialROI: (balance * 0.7 * currentAPY / 100),
			Timestamp:   time.Now(),
		}

		ai.recommendations[rec.ID] = rec
		return rec
	}

	return nil
}

// CreateMarketAlert creates a market trend alert
func (ai *AIWalletManager) CreateMarketAlert(asset string, currentPrice, previousPrice float64) *MarketAlert {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	percentChange := ((currentPrice - previousPrice) / previousPrice) * 100

	// Determine alert priority based on percent change
	priority := PriorityLow
	if math.Abs(percentChange) > 10 {
		priority = PriorityHigh
	} else if math.Abs(percentChange) > 5 {
		priority = PriorityMedium
	}

	alertType := "price_increase"
	message := fmt.Sprintf("%s increased by %.2f%%", asset, percentChange)
	if percentChange < 0 {
		alertType = "price_decrease"
		message = fmt.Sprintf("%s decreased by %.2f%%", asset, math.Abs(percentChange))
	}

	alert := &MarketAlert{
		ID:           fmt.Sprintf("ALERT-%d", time.Now().UnixNano()),
		Asset:        asset,
		AlertType:    alertType,
		Message:      message,
		Priority:     priority,
		CurrentPrice: currentPrice,
		PercentChange: percentChange,
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	ai.alerts[alert.ID] = alert
	return alert
}

// GenerateTradingRecommendation generates trading recommendations based on market analysis
func (ai *AIWalletManager) GenerateTradingRecommendation(asset string, price, avgPrice float64, volatility float64) *Recommendation {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	// Simple trading logic
	priceDeviation := ((price - avgPrice) / avgPrice) * 100

	var recType string
	var reason string
	confidence := 0.7

	if priceDeviation < -5 && volatility < 15 {
		recType = "buy"
		reason = fmt.Sprintf("%s is %.2f%% below average with low volatility. Good buying opportunity.", asset, math.Abs(priceDeviation))
		confidence = 0.8
	} else if priceDeviation > 10 {
		recType = "sell"
		reason = fmt.Sprintf("%s is %.2f%% above average. Consider taking profits.", asset, priceDeviation)
		confidence = 0.75
	} else {
		recType = "hold"
		reason = fmt.Sprintf("%s price is stable. Continue holding.", asset)
		confidence = 0.6
	}

	rec := &Recommendation{
		ID:         fmt.Sprintf("REC-TRD-%d", time.Now().UnixNano()),
		Type:       recType,
		Asset:      asset,
		Reason:     reason,
		Confidence: confidence,
		Timestamp:  time.Now(),
	}

	ai.recommendations[rec.ID] = rec
	return rec
}

// GetActiveAlerts retrieves all active alerts
func (ai *AIWalletManager) GetActiveAlerts() []*MarketAlert {
	ai.mutex.RLock()
	defer ai.mutex.RUnlock()

	alerts := make([]*MarketAlert, 0)
	for _, alert := range ai.alerts {
		if alert.IsActive {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// GetInsights retrieves all insights
func (ai *AIWalletManager) GetInsights() []*AIInsight {
	ai.mutex.RLock()
	defer ai.mutex.RUnlock()

	insights := make([]*AIInsight, 0, len(ai.insights))
	for _, insight := range ai.insights {
		insights = append(insights, insight)
	}

	return insights
}

// GetRecommendations retrieves all recommendations
func (ai *AIWalletManager) GetRecommendations() []*Recommendation {
	ai.mutex.RLock()
	defer ai.mutex.RUnlock()

	recommendations := make([]*Recommendation, 0, len(ai.recommendations))
	for _, rec := range ai.recommendations {
		recommendations = append(recommendations, rec)
	}

	return recommendations
}

// DismissAlert dismisses an active alert
func (ai *AIWalletManager) DismissAlert(alertID string) error {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	alert, exists := ai.alerts[alertID]
	if !exists {
		return errors.New("alert not found")
	}

	alert.IsActive = false
	return nil
}

// GenerateOptimizationInsight generates portfolio optimization insights
func (ai *AIWalletManager) GenerateOptimizationInsight(portfolio *Portfolio) *AIInsight {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	assets := portfolio.GetAllAssets()
	
	// Calculate portfolio diversity
	diversity := len(assets)
	
	// Calculate risk level based on volatility
	highVolatilityCount := 0
	for _, asset := range assets {
		if math.Abs(asset.Change24h) > 5.0 {
			highVolatilityCount++
		}
	}

	actionItems := make([]string, 0)
	priority := PriorityLow

	if diversity < 3 {
		actionItems = append(actionItems, "Consider diversifying into more assets")
		priority = PriorityMedium
	}

	if highVolatilityCount > diversity/2 {
		actionItems = append(actionItems, "High portfolio volatility detected. Consider adding stable assets")
		priority = PriorityHigh
	}

	insight := &AIInsight{
		ID:          fmt.Sprintf("INS-OPT-%d", time.Now().UnixNano()),
		Type:        InsightOptimization,
		Title:       "Portfolio Optimization",
		Description: fmt.Sprintf("Your portfolio has %d assets with %d high-volatility holdings", diversity, highVolatilityCount),
		Priority:    priority,
		Data: map[string]interface{}{
			"diversity":     diversity,
			"volatileAssets": highVolatilityCount,
			"totalValue":    portfolio.GetTotalValue(),
		},
		CreatedAt:   time.Now(),
		ActionItems: actionItems,
	}

	ai.insights[insight.ID] = insight
	return insight
}
