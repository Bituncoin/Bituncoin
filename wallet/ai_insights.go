package wallet

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// AlertType represents types of alerts
type AlertType string

const (
	AlertTypeMarketTrend      AlertType = "market_trend"
	AlertTypePriceChange      AlertType = "price_change"
	AlertTypeStakingReward    AlertType = "staking_reward"
	AlertTypeSecurityWarning  AlertType = "security_warning"
	AlertTypeTransactionAlert AlertType = "transaction_alert"
	AlertTypePortfolioAlert   AlertType = "portfolio_alert"
)

// AlertSeverity represents alert severity levels
type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

// Alert represents an AI-generated alert
type Alert struct {
	ID          string        `json:"id"`
	Type        AlertType     `json:"type"`
	Severity    AlertSeverity `json:"severity"`
	Title       string        `json:"title"`
	Message     string        `json:"message"`
	Timestamp   int64         `json:"timestamp"`
	WalletID    string        `json:"walletId"`
	ActionURL   string        `json:"actionUrl,omitempty"`
	IsRead      bool          `json:"isRead"`
}

// Insight represents an AI-generated insight
type Insight struct {
	ID          string  `json:"id"`
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"` // 0-100
	Impact      string  `json:"impact"`     // low, medium, high
	Timestamp   int64   `json:"timestamp"`
	WalletID    string  `json:"walletId"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// Recommendation represents an AI recommendation
type Recommendation struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Reasoning   string  `json:"reasoning"`
	Priority    int     `json:"priority"` // 1-10
	Confidence  float64 `json:"confidence"` // 0-100
	Timestamp   int64   `json:"timestamp"`
	WalletID    string  `json:"walletId"`
	Actions     []string `json:"actions,omitempty"`
}

// MarketTrend represents a market trend analysis
type MarketTrend struct {
	Currency      string  `json:"currency"`
	Trend         string  `json:"trend"` // up, down, stable
	ChangePercent float64 `json:"changePercent"`
	Volume24h     float64 `json:"volume24h"`
	Prediction    string  `json:"prediction"`
	Confidence    float64 `json:"confidence"`
	Timestamp     int64   `json:"timestamp"`
}

// AIInsightsEngine manages AI-driven insights and recommendations
type AIInsightsEngine struct {
	alerts          map[string]*Alert
	insights        map[string]*Insight
	recommendations map[string]*Recommendation
	marketTrends    map[string]*MarketTrend
	mutex           sync.RWMutex
}

// NewAIInsightsEngine creates a new AI insights engine
func NewAIInsightsEngine() *AIInsightsEngine {
	return &AIInsightsEngine{
		alerts:          make(map[string]*Alert),
		insights:        make(map[string]*Insight),
		recommendations: make(map[string]*Recommendation),
		marketTrends:    make(map[string]*MarketTrend),
	}
}

// AnalyzePortfolio analyzes a portfolio and generates insights
func (ai *AIInsightsEngine) AnalyzePortfolio(walletID string, portfolio *Portfolio) []*Insight {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	insights := make([]*Insight, 0)

	// Diversification analysis
	if insight := ai.analyzeDiversification(walletID, portfolio); insight != nil {
		insights = append(insights, insight)
		ai.insights[insight.ID] = insight
	}

	// Risk analysis
	if insight := ai.analyzeRisk(walletID, portfolio); insight != nil {
		insights = append(insights, insight)
		ai.insights[insight.ID] = insight
	}

	// Performance analysis
	if insight := ai.analyzePerformance(walletID, portfolio); insight != nil {
		insights = append(insights, insight)
		ai.insights[insight.ID] = insight
	}

	return insights
}

// analyzeDiversification analyzes portfolio diversification
func (ai *AIInsightsEngine) analyzeDiversification(walletID string, portfolio *Portfolio) *Insight {
	// Count non-zero balances
	activeAssets := 0
	maxAllocation := 0.0

	for _, balance := range portfolio.Balances {
		if balance.Amount > 0 {
			activeAssets++
			allocation := balance.USDValue / portfolio.TotalUSD
			if allocation > maxAllocation {
				maxAllocation = allocation
			}
		}
	}

	// Generate insight based on diversification
	if activeAssets <= 2 && portfolio.TotalUSD > 1000 {
		return &Insight{
			ID:          generateInsightID(),
			Category:    "diversification",
			Title:       "Low Portfolio Diversification",
			Description: fmt.Sprintf("Your portfolio is concentrated in only %d assets. Consider diversifying to reduce risk.", activeAssets),
			Confidence:  85.0,
			Impact:      "high",
			Timestamp:   time.Now().Unix(),
			WalletID:    walletID,
			Data: map[string]interface{}{
				"activeAssets":   activeAssets,
				"maxAllocation":  maxAllocation * 100,
			},
		}
	}

	return nil
}

// analyzeRisk analyzes portfolio risk
func (ai *AIInsightsEngine) analyzeRisk(walletID string, portfolio *Portfolio) *Insight {
	// Calculate volatility-based risk score
	volatileAssets := []string{"BTC", "ETH", "BNB"}
	volatileAllocation := 0.0

	for _, curr := range volatileAssets {
		if balance, exists := portfolio.Balances[Currency(curr)]; exists && portfolio.TotalUSD > 0 {
			volatileAllocation += balance.USDValue / portfolio.TotalUSD
		}
	}

	if volatileAllocation > 0.7 { // More than 70% in volatile assets
		return &Insight{
			ID:          generateInsightID(),
			Category:    "risk",
			Title:       "High Portfolio Volatility",
			Description: fmt.Sprintf("%.1f%% of your portfolio is in high-volatility assets. Consider balancing with stable assets like USDT.", volatileAllocation*100),
			Confidence:  90.0,
			Impact:      "medium",
			Timestamp:   time.Now().Unix(),
			WalletID:    walletID,
			Data: map[string]interface{}{
				"volatileAllocation": volatileAllocation * 100,
			},
		}
	}

	return nil
}

// analyzePerformance analyzes portfolio performance
func (ai *AIInsightsEngine) analyzePerformance(walletID string, portfolio *Portfolio) *Insight {
	// Simple performance insight based on portfolio value
	if portfolio.TotalUSD > 10000 {
		return &Insight{
			ID:          generateInsightID(),
			Category:    "performance",
			Title:       "Strong Portfolio Performance",
			Description: fmt.Sprintf("Your portfolio value of $%.2f shows strong growth. Consider staking to earn passive income.", portfolio.TotalUSD),
			Confidence:  75.0,
			Impact:      "medium",
			Timestamp:   time.Now().Unix(),
			WalletID:    walletID,
			Data: map[string]interface{}{
				"portfolioValue": portfolio.TotalUSD,
			},
		}
	}

	return nil
}

// GenerateStakingRecommendations generates staking recommendations
func (ai *AIInsightsEngine) GenerateStakingRecommendations(walletID string, portfolio *Portfolio) []*Recommendation {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	recommendations := make([]*Recommendation, 0)

	// Check BTN balance for staking
	if btnBalance, exists := portfolio.Balances[BTN]; exists && btnBalance.Amount >= 100 {
		rec := &Recommendation{
			ID:          generateRecommendationID(),
			Type:        "staking",
			Title:       "Stake BTN for 5% APY",
			Description: fmt.Sprintf("You have %.2f BTN available. Stake it to earn 5%% annual rewards.", btnBalance.Amount),
			Reasoning:   "BTN staking offers competitive yields with low risk. Your balance exceeds the minimum staking requirement of 100 BTN.",
			Priority:    8,
			Confidence:  95.0,
			Timestamp:   time.Now().Unix(),
			WalletID:    walletID,
			Actions:     []string{"stake_btn", "learn_more"},
		}
		recommendations = append(recommendations, rec)
		ai.recommendations[rec.ID] = rec
	}

	return recommendations
}

// CreateAlert creates a new alert
func (ai *AIInsightsEngine) CreateAlert(walletID string, alertType AlertType, severity AlertSeverity, title, message string) *Alert {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	alert := &Alert{
		ID:        generateAlertID(),
		Type:      alertType,
		Severity:  severity,
		Title:     title,
		Message:   message,
		Timestamp: time.Now().Unix(),
		WalletID:  walletID,
		IsRead:    false,
	}

	ai.alerts[alert.ID] = alert

	return alert
}

// GetAlerts returns alerts for a wallet
func (ai *AIInsightsEngine) GetAlerts(walletID string, unreadOnly bool) []*Alert {
	ai.mutex.RLock()
	defer ai.mutex.RUnlock()

	alerts := make([]*Alert, 0)
	for _, alert := range ai.alerts {
		if alert.WalletID == walletID {
			if !unreadOnly || !alert.IsRead {
				alerts = append(alerts, alert)
			}
		}
	}

	return alerts
}

// MarkAlertAsRead marks an alert as read
func (ai *AIInsightsEngine) MarkAlertAsRead(alertID string) error {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	alert, exists := ai.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found")
	}

	alert.IsRead = true
	return nil
}

// AnalyzeMarketTrend analyzes market trends for a currency
func (ai *AIInsightsEngine) AnalyzeMarketTrend(currency string, currentPrice, previousPrice, volume24h float64) *MarketTrend {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	changePercent := ((currentPrice - previousPrice) / previousPrice) * 100
	
	// Determine trend
	var trend string
	if math.Abs(changePercent) < 1 {
		trend = "stable"
	} else if changePercent > 0 {
		trend = "up"
	} else {
		trend = "down"
	}

	// Simple prediction (in production, this would use ML models)
	prediction := "uncertain"
	confidence := 50.0

	if math.Abs(changePercent) > 5 {
		if trend == "up" {
			prediction = "continued_growth"
			confidence = 70.0
		} else {
			prediction = "correction_expected"
			confidence = 70.0
		}
	}

	marketTrend := &MarketTrend{
		Currency:      currency,
		Trend:         trend,
		ChangePercent: changePercent,
		Volume24h:     volume24h,
		Prediction:    prediction,
		Confidence:    confidence,
		Timestamp:     time.Now().Unix(),
	}

	ai.marketTrends[currency] = marketTrend

	return marketTrend
}

// GetMarketTrend returns the market trend for a currency
func (ai *AIInsightsEngine) GetMarketTrend(currency string) (*MarketTrend, error) {
	ai.mutex.RLock()
	defer ai.mutex.RUnlock()

	trend, exists := ai.marketTrends[currency]
	if !exists {
		return nil, fmt.Errorf("market trend not available for %s", currency)
	}

	return trend, nil
}

// OptimizePortfolio generates portfolio optimization recommendations
func (ai *AIInsightsEngine) OptimizePortfolio(walletID string, portfolio *Portfolio, riskTolerance string) []*Recommendation {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()

	recommendations := make([]*Recommendation, 0)

	// Based on risk tolerance, suggest rebalancing
	if riskTolerance == "conservative" {
		// Suggest more stable assets
		rec := &Recommendation{
			ID:          generateRecommendationID(),
			Type:        "rebalancing",
			Title:       "Increase Stable Asset Allocation",
			Description: "For conservative risk profile, consider increasing USDT allocation to 40-50% of portfolio.",
			Reasoning:   "Stable assets reduce portfolio volatility and protect against market downturns.",
			Priority:    7,
			Confidence:  85.0,
			Timestamp:   time.Now().Unix(),
			WalletID:    walletID,
			Actions:     []string{"rebalance_portfolio", "learn_more"},
		}
		recommendations = append(recommendations, rec)
		ai.recommendations[rec.ID] = rec
	}

	return recommendations
}

// Helper functions

func generateInsightID() string {
	return fmt.Sprintf("insight_%d", time.Now().UnixNano())
}

func generateRecommendationID() string {
	return fmt.Sprintf("rec_%d", time.Now().UnixNano())
}

func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}
