package main

import (
	"fmt"
	"time"

	"github.com/Bituncoin/Bituncoin/wallet"
)

func main() {
	fmt.Println("=== Bituncoin Comprehensive Wallet Demo ===\n")

	// 1. Portfolio Management Demo
	fmt.Println("1. Portfolio Management")
	fmt.Println("------------------------")
	portfolio := wallet.NewPortfolio()
	
	// Add multiple assets
	portfolio.AddAsset("BTN", "Bituncoin", 1000.0, 1.0)
	portfolio.AddAsset("GLD", "Gold-Coin", 500.0, 1.0)
	portfolio.AddAsset("BTC", "Bitcoin", 0.5, 50000.0)
	portfolio.AddAsset("ETH", "Ethereum", 2.0, 3000.0)
	
	fmt.Printf("Total Portfolio Value: $%.2f\n", portfolio.GetTotalValue())
	
	// Get performance metrics
	performance := portfolio.GetPerformance()
	fmt.Printf("Asset Count: %v\n", performance["assetCount"])
	fmt.Printf("Average 24h Change: %.2f%%\n\n", performance["avgChange24h"])

	// 2. Exchange Demo
	fmt.Println("2. Cryptocurrency Exchange")
	fmt.Println("--------------------------")
	exchange := wallet.NewExchange()
	
	// Get exchange rate
	rate, _ := exchange.GetExchangeRate("BTC", "ETH")
	fmt.Printf("BTC/ETH Rate: %.4f\n", rate.Rate)
	fmt.Printf("Exchange Fee: %.2f%%\n", rate.Fee)
	
	// Calculate exchange
	toAmount, fee, _ := exchange.CalculateExchange("BTC", "ETH", 0.1)
	fmt.Printf("0.1 BTC = %.4f ETH (fee: %.6f ETH)\n", toAmount, fee)
	
	// Create exchange order
	order, _ := exchange.CreateExchangeOrder("BTN123...", "BTC", "ETH", 0.1)
	fmt.Printf("Order Created: %s\n", order.ID)
	fmt.Printf("Status: %s\n\n", order.Status)

	// 3. Payment Cards Demo
	fmt.Println("3. Payment Cards (BTN-Pay)")
	fmt.Println("--------------------------")
	cardManager := wallet.NewCardManager()
	
	// Create virtual card
	card, _ := cardManager.CreateCard(
		"BTN123...",
		wallet.CardTypeVirtual,
		wallet.ProviderVisa,
		1000.0, // daily limit
	)
	fmt.Printf("Card Created: %s\n", card.ID)
	fmt.Printf("Card Number: %s\n", card.CardNumber)
	fmt.Printf("Card Type: %s\n", card.CardType)
	fmt.Printf("Provider: %s\n", card.Provider)
	fmt.Printf("Daily Limit: $%.2f\n", card.DailyLimit)
	
	// Top up card
	cardManager.TopUpCard(card.ID, 500.0)
	fmt.Printf("Card Balance: $%.2f\n", card.Balance)
	
	// Process transaction
	tx, _ := cardManager.ProcessCardTransaction(card.ID, "Amazon Store", 99.99, "purchase")
	fmt.Printf("Transaction: %s at %s for $%.2f\n\n", tx.Type, tx.Merchant, tx.Amount)

	// 4. Merchant Services Demo
	fmt.Println("4. Merchant Services")
	fmt.Println("--------------------")
	merchantService := wallet.NewMerchantService()
	
	// Register merchant
	merchant, _ := merchantService.RegisterMerchant(
		"Coffee Shop",
		"GLD789...",
		"shop@example.com",
		"retail",
	)
	fmt.Printf("Merchant Registered: %s\n", merchant.Name)
	fmt.Printf("Merchant ID: %s\n", merchant.ID)
	
	// Create payment request with QR code
	payment, _ := merchantService.CreatePaymentRequest(
		merchant.ID,
		25.50,
		"GLD",
		wallet.PaymentQRCode,
		"Coffee and pastry",
	)
	fmt.Printf("Payment Request: %s\n", payment.ID)
	fmt.Printf("Amount: %.2f %s\n", payment.Amount, payment.Asset)
	fmt.Printf("QR Code: %s\n", payment.QRCode[:20]+"...")
	
	// Process mobile money payment
	mobilePayment, _ := merchantService.ProcessMobileMoneyPayment(
		merchant.ID,
		wallet.ProviderMTN,
		"+233123456789",
		50.0,
		"GHS",
	)
	fmt.Printf("Mobile Payment: %s\n", mobilePayment.ID)
	fmt.Printf("Status: %s\n\n", mobilePayment.Status)

	// 5. AI Wallet Manager Demo
	fmt.Println("5. AI-Driven Insights")
	fmt.Println("---------------------")
	aiManager := wallet.NewAIWalletManager()
	
	// Create sample transactions for analysis
	transactions := []*wallet.Transaction{
		{
			ID:        "TX1",
			Type:      wallet.TypeSent,
			Amount:    50.0,
			Asset:     "GLD",
			Timestamp: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        "TX2",
			Type:      wallet.TypeSent,
			Amount:    75.0,
			Asset:     "GLD",
			Timestamp: time.Now().Add(-12 * time.Hour),
		},
	}
	
	// Analyze spending
	insights := aiManager.AnalyzeSpending("BTN123...", transactions)
	if len(insights) > 0 {
		fmt.Printf("Spending Insight: %s\n", insights[0].Title)
		fmt.Printf("Description: %s\n", insights[0].Description)
	}
	
	// Generate staking recommendation
	stakingRec := aiManager.GenerateStakingRecommendation("GLD", 1000.0, 5.0)
	if stakingRec != nil {
		fmt.Printf("\nStaking Recommendation:\n")
		fmt.Printf("Type: %s\n", stakingRec.Type)
		fmt.Printf("Asset: %s\n", stakingRec.Asset)
		fmt.Printf("Amount: %.2f\n", stakingRec.Amount)
		fmt.Printf("Reason: %s\n", stakingRec.Reason)
		fmt.Printf("Confidence: %.0f%%\n", stakingRec.Confidence*100)
		fmt.Printf("Potential ROI: $%.2f\n", stakingRec.PotentialROI)
	}
	
	// Create market alert
	alert := aiManager.CreateMarketAlert("BTC", 52000.0, 50000.0)
	fmt.Printf("\nMarket Alert: %s\n", alert.Message)
	fmt.Printf("Priority: %s\n\n", alert.Priority)

	// 6. Platform Configuration Demo
	fmt.Println("6. Platform Detection")
	fmt.Println("---------------------")
	platformConfig := wallet.NewPlatformConfig()
	
	info := platformConfig.GetPlatformInfo()
	fmt.Printf("Platform: %v\n", info["platform"])
	fmt.Printf("Architecture: %v\n", info["arch"])
	fmt.Printf("CPUs: %v\n", info["numCPU"])
	
	fmt.Printf("Biometric Support: %v\n", platformConfig.Capabilities.SupportsBiometric)
	fmt.Printf("NFC Support: %v\n", platformConfig.Capabilities.SupportsNFC)
	
	if platformConfig.IsFeatureEnabled("multi_currency") {
		fmt.Println("Multi-currency support: Enabled")
	}
	fmt.Println()

	// 7. Operations Dashboard Demo
	fmt.Println("7. Operations Dashboard")
	fmt.Println("-----------------------")
	dashboard := wallet.NewDashboard()
	
	// System status
	status := dashboard.GetSystemStatus()
	fmt.Printf("System Status: %s\n", status)
	
	// Health check
	components := dashboard.HealthCheck()
	fmt.Printf("Total Components: %d\n", len(components))
	
	// Network connections
	networks := dashboard.CheckNetworkConnections()
	for name, network := range networks {
		fmt.Printf("Network %s: Connected=%v, Height=%d\n", 
			name, network.Connected, network.BlockHeight)
	}
	
	// Update metrics
	dashboard.UpdateMetrics(1000, 500, 5000, 250000.0)
	metrics := dashboard.GetMetrics()
	fmt.Printf("\nMetrics:\n")
	fmt.Printf("Total Users: %d\n", metrics.TotalUsers)
	fmt.Printf("Active Wallets: %d\n", metrics.ActiveWallets)
	fmt.Printf("Total Transactions: %d\n", metrics.TotalTransactions)
	fmt.Printf("Total Volume: $%.2f\n", metrics.TotalVolume)
	fmt.Printf("Uptime: %.2f%%\n\n", metrics.Uptime)

	// 8. Security & Fraud Detection Demo
	fmt.Println("8. Security & Fraud Detection")
	fmt.Println("------------------------------")
	security := wallet.NewSecurity()
	fraudDetector := wallet.NewFraudDetector()
	alertSystem := wallet.NewAlertSystem()
	
	// Enable security features
	security.EnableTwoFactor("secret123")
	security.EnableBiometric("fingerprint")
	
	secStatus := security.GetSecurityStatus()
	fmt.Printf("2FA Enabled: %v\n", secStatus["twoFactorEnabled"])
	fmt.Printf("Biometric Enabled: %v\n", secStatus["biometricEnabled"])
	fmt.Printf("Encryption: %v\n", secStatus["encryptionType"])
	
	// Check transaction for fraud
	isSuspicious, reason := fraudDetector.CheckTransaction("BTN123...", "BTN456...", 15000.0)
	if isSuspicious {
		fmt.Printf("\nFraud Alert: %s\n", reason)
		alertSystem.SendAlert("fraud", "high", reason, "BTN123...")
	}
	
	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nThe Bituncoin Comprehensive Wallet provides:")
	fmt.Println("✓ Multi-currency portfolio management")
	fmt.Println("✓ Built-in cryptocurrency exchange")
	fmt.Println("✓ Payment card integration (Visa/MasterCard)")
	fmt.Println("✓ Merchant services with QR/NFC/mobile money")
	fmt.Println("✓ AI-driven insights and recommendations")
	fmt.Println("✓ Multi-platform support (iOS/Android/Windows/Mac/Linux/Web)")
	fmt.Println("✓ Advanced security and fraud detection")
	fmt.Println("✓ Unified operations dashboard")
}
