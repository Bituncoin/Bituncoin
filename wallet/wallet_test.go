package wallet

import (
	"testing"
	"time"
)

// Portfolio Tests

func TestNewPortfolio(t *testing.T) {
	portfolio := NewPortfolio()
	if portfolio == nil {
		t.Fatal("Expected portfolio to be created")
	}
	if len(portfolio.Assets) != 0 {
		t.Errorf("Expected 0 assets, got %d", len(portfolio.Assets))
	}
}

func TestAddAsset(t *testing.T) {
	portfolio := NewPortfolio()
	portfolio.AddAsset("BTC", "Bitcoin", 1.0, 50000.0)
	
	if len(portfolio.Assets) != 1 {
		t.Errorf("Expected 1 asset, got %d", len(portfolio.Assets))
	}
	
	btc, err := portfolio.GetAsset("BTC")
	if err != nil {
		t.Fatalf("Expected to get BTC asset: %v", err)
	}
	
	if btc.Balance != 1.0 {
		t.Errorf("Expected balance 1.0, got %.2f", btc.Balance)
	}
	
	if btc.USDValue != 50000.0 {
		t.Errorf("Expected USD value 50000.0, got %.2f", btc.USDValue)
	}
}

func TestUpdateBalance(t *testing.T) {
	portfolio := NewPortfolio()
	portfolio.AddAsset("BTC", "Bitcoin", 1.0, 50000.0)
	
	err := portfolio.UpdateBalance("BTC", 2.0)
	if err != nil {
		t.Fatalf("Failed to update balance: %v", err)
	}
	
	btc, _ := portfolio.GetAsset("BTC")
	if btc.Balance != 2.0 {
		t.Errorf("Expected balance 2.0, got %.2f", btc.Balance)
	}
	
	if btc.USDValue != 100000.0 {
		t.Errorf("Expected USD value 100000.0, got %.2f", btc.USDValue)
	}
}

func TestGetTotalValue(t *testing.T) {
	portfolio := NewPortfolio()
	portfolio.AddAsset("BTC", "Bitcoin", 1.0, 50000.0)
	portfolio.AddAsset("ETH", "Ethereum", 10.0, 3000.0)
	
	totalValue := portfolio.GetTotalValue()
	expectedValue := 80000.0
	
	if totalValue != expectedValue {
		t.Errorf("Expected total value %.2f, got %.2f", expectedValue, totalValue)
	}
}

// Exchange Tests

func TestNewExchange(t *testing.T) {
	exchange := NewExchange()
	if exchange == nil {
		t.Fatal("Expected exchange to be created")
	}
	
	pairs := exchange.GetSupportedPairs()
	if len(pairs) == 0 {
		t.Error("Expected supported pairs to be initialized")
	}
}

func TestGetExchangeRate(t *testing.T) {
	exchange := NewExchange()
	
	rate, err := exchange.GetExchangeRate("BTC", "ETH")
	if err != nil {
		t.Fatalf("Failed to get exchange rate: %v", err)
	}
	
	if rate.Rate <= 0 {
		t.Error("Expected positive exchange rate")
	}
	
	if rate.Fee < 0 {
		t.Error("Expected non-negative fee")
	}
}

func TestCalculateExchange(t *testing.T) {
	exchange := NewExchange()
	
	toAmount, fee, err := exchange.CalculateExchange("BTC", "ETH", 1.0)
	if err != nil {
		t.Fatalf("Failed to calculate exchange: %v", err)
	}
	
	if toAmount <= 0 {
		t.Error("Expected positive exchange amount")
	}
	
	if fee < 0 {
		t.Error("Expected non-negative fee")
	}
}

func TestCreateExchangeOrder(t *testing.T) {
	exchange := NewExchange()
	
	order, err := exchange.CreateExchangeOrder("BTNG123", "BTC", "ETH", 0.5)
	if err != nil {
		t.Fatalf("Failed to create exchange order: %v", err)
	}
	
	if order.ID == "" {
		t.Error("Expected order ID to be set")
	}
	
	if order.FromAmount != 0.5 {
		t.Errorf("Expected from amount 0.5, got %.2f", order.FromAmount)
	}
	
	if order.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", order.Status)
	}
}

// Card Manager Tests

func TestCreateCard(t *testing.T) {
	cardManager := NewCardManager()
	
	card, err := cardManager.CreateCard("BTNG123", CardTypeVirtual, ProviderVisa, 1000.0)
	if err != nil {
		t.Fatalf("Failed to create card: %v", err)
	}
	
	if card.ID == "" {
		t.Error("Expected card ID to be set")
	}
	
	if card.CardType != CardTypeVirtual {
		t.Errorf("Expected card type %s, got %s", CardTypeVirtual, card.CardType)
	}
	
	if card.Status != CardStatusActive {
		t.Errorf("Expected status %s, got %s", CardStatusActive, card.Status)
	}
}

func TestTopUpCard(t *testing.T) {
	cardManager := NewCardManager()
	card, _ := cardManager.CreateCard("BTNG123", CardTypeVirtual, ProviderVisa, 1000.0)
	
	err := cardManager.TopUpCard(card.ID, 500.0)
	if err != nil {
		t.Fatalf("Failed to top up card: %v", err)
	}
	
	updatedCard, _ := cardManager.GetCard(card.ID)
	if updatedCard.Balance != 500.0 {
		t.Errorf("Expected balance 500.0, got %.2f", updatedCard.Balance)
	}
}

func TestProcessCardTransaction(t *testing.T) {
	cardManager := NewCardManager()
	card, _ := cardManager.CreateCard("BTNG123", CardTypeVirtual, ProviderVisa, 1000.0)
	cardManager.TopUpCard(card.ID, 500.0)
	
	tx, err := cardManager.ProcessCardTransaction(card.ID, "Test Store", 100.0, "purchase")
	if err != nil {
		t.Fatalf("Failed to process transaction: %v", err)
	}
	
	if tx.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %.2f", tx.Amount)
	}
	
	updatedCard, _ := cardManager.GetCard(card.ID)
	if updatedCard.Balance != 400.0 {
		t.Errorf("Expected balance 400.0, got %.2f", updatedCard.Balance)
	}
}

// Merchant Service Tests

func TestRegisterMerchant(t *testing.T) {
	merchantService := NewMerchantService()
	
	merchant, err := merchantService.RegisterMerchant("Test Shop", "GLD123", "test@example.com", "retail")
	if err != nil {
		t.Fatalf("Failed to register merchant: %v", err)
	}
	
	if merchant.ID == "" {
		t.Error("Expected merchant ID to be set")
	}
	
	if merchant.Status != "active" {
		t.Errorf("Expected status 'active', got %s", merchant.Status)
	}
}

func TestCreatePaymentRequest(t *testing.T) {
	merchantService := NewMerchantService()
	merchant, _ := merchantService.RegisterMerchant("Test Shop", "GLD123", "test@example.com", "retail")
	
	payment, err := merchantService.CreatePaymentRequest(merchant.ID, 50.0, "GLD", PaymentQRCode, "Test payment")
	if err != nil {
		t.Fatalf("Failed to create payment request: %v", err)
	}
	
	if payment.ID == "" {
		t.Error("Expected payment ID to be set")
	}
	
	if payment.QRCode == "" {
		t.Error("Expected QR code to be generated")
	}
	
	if payment.Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", payment.Status)
	}
}

// Transaction History Tests

func TestAddTransaction(t *testing.T) {
	history := NewTransactionHistory()
	
	tx := &Transaction{
		ID:        "TX123",
		Type:      TypeSent,
		From:      "BTNG123",
		To:        "BTNG456",
		Amount:    100.0,
		Asset:     "GLD",
		Timestamp: time.Now(),
	}
	
	err := history.AddTransaction(tx)
	if err != nil {
		t.Fatalf("Failed to add transaction: %v", err)
	}
	
	retrieved, err := history.GetTransaction("TX123")
	if err != nil {
		t.Fatalf("Failed to get transaction: %v", err)
	}
	
	if retrieved.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %.2f", retrieved.Amount)
	}
}

func TestGetUserTransactions(t *testing.T) {
	history := NewTransactionHistory()
	
	tx1 := &Transaction{
		ID:        "TX1",
		Type:      TypeSent,
		From:      "BTNG123",
		To:        "BTNG456",
		Amount:    100.0,
		Asset:     "GLD",
		Timestamp: time.Now(),
	}
	
	tx2 := &Transaction{
		ID:        "TX2",
		Type:      TypeReceived,
		From:      "BTNG789",
		To:        "BTNG123",
		Amount:    50.0,
		Asset:     "GLD",
		Timestamp: time.Now(),
	}
	
	history.AddTransaction(tx1)
	history.AddTransaction(tx2)
	
	transactions := history.GetUserTransactions("BTNG123")
	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}
}

// AI Manager Tests

func TestCreateMarketAlert(t *testing.T) {
	aiManager := NewAIWalletManager()
	
	alert := aiManager.CreateMarketAlert("BTC", 52000.0, 50000.0)
	if alert == nil {
		t.Fatal("Expected alert to be created")
	}
	
	if alert.Asset != "BTC" {
		t.Errorf("Expected asset BTC, got %s", alert.Asset)
	}
	
	if alert.PercentChange == 0 {
		t.Error("Expected non-zero percent change")
	}
}

func TestGenerateStakingRecommendation(t *testing.T) {
	aiManager := NewAIWalletManager()
	
	rec := aiManager.GenerateStakingRecommendation("GLD", 1000.0, 5.0)
	if rec == nil {
		t.Fatal("Expected recommendation to be created")
	}
	
	if rec.Type != "stake" {
		t.Errorf("Expected type 'stake', got %s", rec.Type)
	}
	
	if rec.Amount <= 0 {
		t.Error("Expected positive recommendation amount")
	}
}

// Platform Config Tests

func TestNewPlatformConfig(t *testing.T) {
	config := NewPlatformConfig()
	if config == nil {
		t.Fatal("Expected platform config to be created")
	}
	
	if config.Platform == PlatformUnknown {
		t.Error("Expected platform to be detected")
	}
}

func TestIsFeatureEnabled(t *testing.T) {
	config := NewPlatformConfig()
	
	if !config.IsFeatureEnabled("multi_currency") {
		t.Error("Expected multi_currency feature to be enabled")
	}
}

// Dashboard Tests

func TestNewDashboard(t *testing.T) {
	dashboard := NewDashboard()
	if dashboard == nil {
		t.Fatal("Expected dashboard to be created")
	}
	
	status := dashboard.GetSystemStatus()
	if status == "" {
		t.Error("Expected system status to be set")
	}
}

func TestHealthCheck(t *testing.T) {
	dashboard := NewDashboard()
	
	components := dashboard.HealthCheck()
	if len(components) == 0 {
		t.Error("Expected components to be checked")
	}
}

func TestUpdateMetrics(t *testing.T) {
	dashboard := NewDashboard()
	
	dashboard.UpdateMetrics(1000, 500, 5000, 250000.0)
	metrics := dashboard.GetMetrics()
	
	if metrics.TotalUsers != 1000 {
		t.Errorf("Expected 1000 users, got %d", metrics.TotalUsers)
	}
	
	if metrics.TotalVolume != 250000.0 {
		t.Errorf("Expected volume 250000.0, got %.2f", metrics.TotalVolume)
	}
}

// Security & Fraud Detection Tests

func TestFraudDetector(t *testing.T) {
	detector := NewFraudDetector()
	
	// Test large transaction
	isSuspicious, reason := detector.CheckTransaction("BTNG123", "BTNG456", 15000.0)
	if !isSuspicious {
		t.Error("Expected large transaction to be flagged")
	}
	
	if reason == "" {
		t.Error("Expected reason to be provided")
	}
}

func TestBlockAddress(t *testing.T) {
	detector := NewFraudDetector()
	
	detector.BlockAddress("BTNG123")
	
	if !detector.IsAddressBlocked("BTNG123") {
		t.Error("Expected address to be blocked")
	}
}

func TestAlertSystem(t *testing.T) {
	alertSystem := NewAlertSystem()
	
	alertSystem.SendAlert("security", "high", "Test alert", "BTNG123")
	
	alerts := alertSystem.GetAlerts("BTNG123")
	if len(alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(alerts))
	}
	
	if alerts[0].Message != "Test alert" {
		t.Errorf("Expected message 'Test alert', got %s", alerts[0].Message)
	}
}
