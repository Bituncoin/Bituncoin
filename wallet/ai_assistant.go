package wallet

import (
	"fmt"
	"math"
	"time"
)

// AIInsight represents an AI-generated insight
type AIInsight struct {
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"` // low, medium, high
	CreatedAt   time.Time `json:"createdAt"`
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	Level       string    `json:"level"` // info, warning, critical
	Message     string    `json:"message"`
	ActionItems []string  `json:"actionItems,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// AIAssistant provides AI-powered wallet insights and recommendations
type AIAssistant struct {
	// In production, this would integrate with ML models and analytics
}

// NewAIAssistant creates a new AI assistant
func NewAIAssistant() *AIAssistant {
	return &AIAssistant{}
}

// AnalyzePortfolio analyzes the wallet portfolio and provides insights
func (ai *AIAssistant) AnalyzePortfolio(balances map[string]float64) []AIInsight {
	insights := []AIInsight{}

	totalValue := ai.calculateTotalValue(balances)
	
	// Check for diversification
	if ai.isDiversified(balances) {
		insights = append(insights, AIInsight{
			Type:        "portfolio",
			Title:       "Well-Diversified Portfolio",
			Description: "Your portfolio is well-diversified across multiple assets, reducing overall risk.",
			Priority:    "low",
			CreatedAt:   time.Now(),
		})
	} else {
		insights = append(insights, AIInsight{
			Type:        "portfolio",
			Title:       "Diversification Opportunity",
			Description: "Consider diversifying your portfolio across more assets to reduce risk.",
			Priority:    "medium",
			CreatedAt:   time.Now(),
		})
	}

	// Check for staking opportunities
	if balances["BTN"] > 100 || balances["GLD"] > 100 {
		availableForStaking := math.Max(balances["BTN"], balances["GLD"])
		potentialRewards := availableForStaking * 0.05 // 5% APY
		
		insights = append(insights, AIInsight{
			Type:        "staking",
			Title:       "Staking Opportunity",
			Description: fmt.Sprintf("You have %.2f tokens available for staking. Potential annual rewards: %.2f", availableForStaking, potentialRewards),
			Priority:    "medium",
			CreatedAt:   time.Now(),
		})
	}

	// Market trend insight (mock)
	insights = append(insights, AIInsight{
		Type:        "market",
		Title:       "Market Trend Analysis",
		Description: "BTN is showing strong upward momentum. Consider holding or accumulating.",
		Priority:    "low",
		CreatedAt:   time.Now(),
	})

	return insights
}

// GenerateSecurityAlerts generates security alerts based on wallet state
func (ai *AIAssistant) GenerateSecurityAlerts(securitySettings map[string]bool, lastBackupDays int, hasRecentSuspiciousActivity bool) []SecurityAlert {
	alerts := []SecurityAlert{}

	// Check 2FA status
	if !securitySettings["twoFactorEnabled"] {
		alerts = append(alerts, SecurityAlert{
			Level:   "warning",
			Message: "Two-Factor Authentication is disabled. Enable 2FA to secure your wallet.",
			ActionItems: []string{
				"Go to Security settings",
				"Enable Two-Factor Authentication",
			},
			CreatedAt: time.Now(),
		})
	} else {
		alerts = append(alerts, SecurityAlert{
			Level:     "info",
			Message:   "Your 2FA is enabled. Your wallet is secure.",
			CreatedAt: time.Now(),
		})
	}

	// Check backup status
	if lastBackupDays > 30 {
		alerts = append(alerts, SecurityAlert{
			Level:   "warning",
			Message: fmt.Sprintf("Last backup was %d days ago. Consider creating a new backup.", lastBackupDays),
			ActionItems: []string{
				"Create encrypted backup",
				"Store backup in secure location",
			},
			CreatedAt: time.Now(),
		})
	}

	// Check for suspicious activity
	if hasRecentSuspiciousActivity {
		alerts = append(alerts, SecurityAlert{
			Level:   "critical",
			Message: "Suspicious activity detected on your wallet. Please review recent transactions.",
			ActionItems: []string{
				"Review transaction history",
				"Change password immediately",
				"Enable fraud monitoring",
			},
			CreatedAt: time.Now(),
		})
	} else {
		alerts = append(alerts, SecurityAlert{
			Level:     "info",
			Message:   "No suspicious activity detected in the last 30 days.",
			CreatedAt: time.Now(),
		})
	}

	// Fraud monitoring status
	if securitySettings["fraudMonitoring"] {
		alerts = append(alerts, SecurityAlert{
			Level:     "info",
			Message:   "Real-time fraud monitoring is active.",
			CreatedAt: time.Now(),
		})
	}

	return alerts
}

// ProvideRecommendation provides a recommendation based on query
func (ai *AIAssistant) ProvideRecommendation(query string, context map[string]interface{}) string {
	// Simple keyword-based recommendations (in production, use NLP/ML)
	switch {
	case containsKeyword(query, "portfolio", "improve", "diversify"):
		return "To improve your portfolio: 1) Diversify across multiple assets (BTN, GLD, BTC, ETH, USDT, BNB), 2) Consider staking for passive income, 3) Rebalance periodically based on market conditions."
	
	case containsKeyword(query, "staking", "stake", "rewards"):
		return "Staking strategies: 1) Stake at least 100 BTN/GLD to earn 5% annual rewards, 2) Lock tokens for 30 days minimum, 3) Claim rewards regularly to compound earnings, 4) Monitor validator performance."
	
	case containsKeyword(query, "security", "secure", "protect"):
		return "Wallet security best practices: 1) Enable 2FA and biometric authentication, 2) Create encrypted backups regularly, 3) Never share recovery phrases, 4) Enable fraud monitoring, 5) Use strong passwords, 6) Verify transaction addresses carefully."
	
	case containsKeyword(query, "swap", "exchange", "trade"):
		return "For cross-chain swaps: 1) Use the built-in swap feature for lowest fees, 2) Check current exchange rates, 3) Consider transaction fees (1% base + network fee), 4) Verify destination address, 5) Start with small amounts for testing."
	
	default:
		return "I can help you with portfolio management, staking strategies, wallet security, and transaction guidance. Please ask a specific question about these topics."
	}
}

// calculateTotalValue calculates total portfolio value in USD
func (ai *AIAssistant) calculateTotalValue(balances map[string]float64) float64 {
	rates := map[string]float64{
		"BTN":  15.0,
		"GLD":  10.0,
		"BTC":  45000.0,
		"ETH":  3000.0,
		"USDT": 1.0,
		"BNB":  300.0,
	}

	total := 0.0
	for currency, balance := range balances {
		if rate, ok := rates[currency]; ok {
			total += balance * rate
		}
	}
	return total
}

// isDiversified checks if portfolio is well-diversified
func (ai *AIAssistant) isDiversified(balances map[string]float64) bool {
	// Portfolio is diversified if no single asset represents > 50% of total value
	totalValue := ai.calculateTotalValue(balances)
	if totalValue == 0 {
		return false
	}

	rates := map[string]float64{
		"BTN":  15.0,
		"GLD":  10.0,
		"BTC":  45000.0,
		"ETH":  3000.0,
		"USDT": 1.0,
		"BNB":  300.0,
	}

	for currency, balance := range balances {
		if rate, ok := rates[currency]; ok {
			assetValue := balance * rate
			if assetValue/totalValue > 0.5 {
				return false
			}
		}
	}

	// Also check if at least 3 different assets are held
	activeAssets := 0
	for _, balance := range balances {
		if balance > 0 {
			activeAssets++
		}
	}

	return activeAssets >= 3
}

// containsKeyword checks if text contains any of the given keywords
func containsKeyword(text string, keywords ...string) bool {
	textLower := toLower(text)
	for _, keyword := range keywords {
		if contains(textLower, toLower(keyword)) {
			return true
		}
	}
	return false
}

// Helper functions
func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
