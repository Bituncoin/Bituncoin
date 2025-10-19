package wallet

import (
	"testing"
)

func TestAIPortfolioAnalysis(t *testing.T) {
	ai := NewAIAssistant()

	balances := map[string]float64{
		"BTN":  5000.0,
		"GLD":  1250.5,
		"BTC":  0.05,
		"ETH":  2.3,
		"USDT": 1000.0,
		"BNB":  10.5,
	}

	insights := ai.AnalyzePortfolio(balances)

	if len(insights) == 0 {
		t.Error("Expected at least one insight")
	}

	// Check that insights have proper structure
	for _, insight := range insights {
		if insight.Title == "" {
			t.Error("Insight title should not be empty")
		}
		if insight.Description == "" {
			t.Error("Insight description should not be empty")
		}
		if insight.Priority == "" {
			t.Error("Insight priority should not be empty")
		}
	}
}

func TestAISecurityAlerts(t *testing.T) {
	ai := NewAIAssistant()

	securitySettings := map[string]bool{
		"twoFactorEnabled": true,
		"fraudMonitoring":  true,
	}

	alerts := ai.GenerateSecurityAlerts(securitySettings, 15, false)

	if len(alerts) == 0 {
		t.Error("Expected at least one alert")
	}

	// Check for 2FA enabled alert
	found2FAAlert := false
	for _, alert := range alerts {
		if alert.Level == "info" && alert.Message != "" {
			found2FAAlert = true
			break
		}
	}

	if !found2FAAlert {
		t.Error("Expected to find 2FA status alert")
	}
}

func TestAIRecommendations(t *testing.T) {
	ai := NewAIAssistant()
	context := make(map[string]interface{})

	tests := []struct {
		query    string
		expected string // substring to check
	}{
		{"How can I improve my portfolio?", "portfolio"},
		{"What are the best staking strategies?", "staking"},
		{"How do I secure my wallet?", "security"},
		{"Should I swap tokens?", "swap"},
	}

	for _, tt := range tests {
		recommendation := ai.ProvideRecommendation(tt.query, context)
		if recommendation == "" {
			t.Errorf("Expected non-empty recommendation for query: %s", tt.query)
		}
	}
}

func TestExchangeRates(t *testing.T) {
	es := NewExchangeService()

	rate, err := es.GetExchangeRate("BTN", "BTC")
	if err != nil {
		t.Fatalf("Failed to get exchange rate: %v", err)
	}

	if rate.Rate <= 0 {
		t.Error("Exchange rate should be positive")
	}

	if rate.From != "BTN" || rate.To != "BTC" {
		t.Errorf("Expected BTN to BTC, got %s to %s", rate.From, rate.To)
	}
}

func TestSwapCalculation(t *testing.T) {
	es := NewExchangeService()

	toAmount, fee, err := es.CalculateSwap("BTN", "ETH", 100.0)
	if err != nil {
		t.Fatalf("Failed to calculate swap: %v", err)
	}

	if toAmount <= 0 {
		t.Error("To amount should be positive")
	}

	if fee <= 0 {
		t.Error("Fee should be positive")
	}

	// Fee should be approximately 0.5% of the exchange amount
	expectedFee := toAmount * 0.005 / 0.995 // Approximate reverse calculation
	if fee < expectedFee*0.9 || fee > expectedFee*1.1 {
		t.Errorf("Fee seems incorrect: got %f, expected around %f", fee, expectedFee)
	}
}

func TestSwapExecution(t *testing.T) {
	es := NewExchangeService()

	tx, err := es.ExecuteSwap("BTN", "ETH", 100.0, "BTNaddr123")
	if err != nil {
		t.Fatalf("Failed to execute swap: %v", err)
	}

	if tx.FromCurrency != "BTN" {
		t.Errorf("Expected from currency BTN, got %s", tx.FromCurrency)
	}

	if tx.ToCurrency != "ETH" {
		t.Errorf("Expected to currency ETH, got %s", tx.ToCurrency)
	}

	if tx.FromAmount != 100.0 {
		t.Errorf("Expected from amount 100.0, got %f", tx.FromAmount)
	}

	if tx.Status != "pending" {
		t.Errorf("Expected initial status 'pending', got '%s'", tx.Status)
	}
}

func TestCryptoToFiatExchange(t *testing.T) {
	es := NewExchangeService()

	fiatAmount, err := es.CryptoToFiatExchange("BTN", 100.0, "USD")
	if err != nil {
		t.Fatalf("Failed to convert crypto to fiat: %v", err)
	}

	if fiatAmount <= 0 {
		t.Error("Fiat amount should be positive")
	}

	// BTN is $15, so 100 BTN should be approximately $1500
	if fiatAmount < 1400 || fiatAmount > 1600 {
		t.Errorf("Expected fiat amount around 1500, got %f", fiatAmount)
	}
}

func TestDiversificationCheck(t *testing.T) {
	ai := NewAIAssistant()

	// Diversified portfolio
	diversified := map[string]float64{
		"BTN": 200.0,  // 200 * 15 = 3000
		"GLD": 200.0,  // 200 * 10 = 2000
		"BTC": 0.02,   // 0.02 * 45000 = 900
		"ETH": 1.0,    // 1 * 3000 = 3000
	}
	// Total: 8900, max is 3000/8900 = 33.7% < 50%

	if !ai.isDiversified(diversified) {
		t.Error("Expected portfolio to be diversified")
	}

	// Not diversified (single asset)
	notDiversified := map[string]float64{
		"BTN": 10000.0,
		"GLD": 0.0,
		"BTC": 0.0,
	}

	if ai.isDiversified(notDiversified) {
		t.Error("Expected portfolio to not be diversified")
	}
}
