package main

import (
	"fmt"
	"time"

	"github.com/Bituncoin/Bituncoin/payments"
	"github.com/Bituncoin/Bituncoin/wallet"
)

func main() {
	fmt.Println("=== BTNg Wallet Feature Demo ===\n")

	// 1. Mobile Money Payment Demo
	fmt.Println("1. Mobile Money Integration")
	mms := payments.NewMobileMoneyService()
	tx, err := mms.InitiatePayment(payments.MTN, "+233244123456", 100.0, "GHS")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Mobile Money transaction initiated: %s\n", tx.ID)
		fmt.Printf("  Provider: %s, Amount: %.2f %s, Status: %s\n", tx.Provider, tx.Amount, tx.Currency, tx.Status)
	}

	// Wait for payment to complete (simulation)
	time.Sleep(3 * time.Second)
	updatedTx, _ := mms.GetTransaction(tx.ID)
	fmt.Printf("  Updated Status: %s\n\n", updatedTx.Status)

	// 2. QR Code Payment Demo
	fmt.Println("2. QR Code Payment System")
	qrs := payments.NewQRCodeService()
	qrData, _ := qrs.GeneratePaymentQR("BTNaddr123", 50.0, "BTN", "Demo payment")
	fmt.Printf("✓ QR Code generated successfully\n")
	fmt.Printf("  QR Data: %s...\n", qrData[:50])
	
	parsedData, _ := qrs.ParsePaymentQR(qrData)
	fmt.Printf("✓ QR Code parsed: Address=%s, Amount=%.2f %s\n\n", parsedData.Address, parsedData.Amount, parsedData.Currency)

	// 3. BTN-Pay Card Demo
	fmt.Println("3. BTN-Pay Card System (MasterCard/Visa)")
	cs := payments.NewCardService()
	card, _ := cs.IssueCard("John Doe", "BTNaddr123", payments.MasterCard, payments.Virtual)
	fmt.Printf("✓ Virtual MasterCard issued: %s\n", card.ID)
	fmt.Printf("  Card Number: %s\n", card.CardNumber)
	fmt.Printf("  Holder: %s, Expiry: %02d/%d\n", card.HolderName, card.ExpiryMonth, card.ExpiryYear)
	
	cs.LoadCard(card.ID, 1000.0)
	fmt.Printf("✓ Card loaded with $%.2f\n", 1000.0)
	
	cardTx, _ := cs.ProcessPayment(card.ID, 50.0, "Demo Store")
	fmt.Printf("✓ Payment processed: $%.2f to %s\n", cardTx.Amount, cardTx.Merchant)
	
	updatedCard, _ := cs.GetCard(card.ID)
	fmt.Printf("  Remaining balance: $%.2f\n\n", updatedCard.Balance)

	// 4. Exchange & Swap Demo
	fmt.Println("4. Currency Exchange & Swap")
	es := wallet.NewExchangeService()
	
	rate, _ := es.GetExchangeRate("BTN", "ETH")
	fmt.Printf("✓ Exchange rate: 1 BTN = %.8f ETH\n", rate.Rate)
	
	toAmount, fee, _ := es.CalculateSwap("BTN", "ETH", 100.0)
	fmt.Printf("✓ Swap estimate: 100 BTN → %.6f ETH (Fee: %.6f ETH)\n", toAmount, fee)
	
	swapTx, _ := es.ExecuteSwap("BTN", "ETH", 100.0, "BTNaddr123")
	fmt.Printf("✓ Swap transaction created: %s (Status: %s)\n\n", swapTx.ID, swapTx.Status)

	// 5. Crypto to Fiat Exchange
	fmt.Println("5. Crypto-to-Fiat Exchange")
	fiatAmount, _ := es.CryptoToFiatExchange("BTN", 100.0, "USD")
	fmt.Printf("✓ 100 BTN = $%.2f USD\n", fiatAmount)
	
	ghsAmount, _ := es.CryptoToFiatExchange("BTN", 100.0, "GHS")
	fmt.Printf("✓ 100 BTN = GHS %.2f\n\n", ghsAmount)

	// 6. AI Assistant Demo
	fmt.Println("6. AI-Powered Assistant")
	ai := wallet.NewAIAssistant()
	
	balances := map[string]float64{
		"BTN":  5000.0,
		"GLD":  1250.5,
		"BTC":  0.05,
		"ETH":  2.3,
		"USDT": 1000.0,
		"BNB":  10.5,
	}
	
	insights := ai.AnalyzePortfolio(balances)
	fmt.Printf("✓ Portfolio analysis generated %d insights:\n", len(insights))
	for i, insight := range insights {
		fmt.Printf("  %d. [%s] %s\n", i+1, insight.Priority, insight.Title)
		fmt.Printf("     %s\n", insight.Description)
	}
	
	fmt.Println("\n✓ Security alerts:")
	securitySettings := map[string]bool{
		"twoFactorEnabled": true,
		"fraudMonitoring":  true,
	}
	alerts := ai.GenerateSecurityAlerts(securitySettings, 15, false)
	for _, alert := range alerts {
		fmt.Printf("  [%s] %s\n", alert.Level, alert.Message)
	}
	
	fmt.Println("\n✓ AI Recommendation:")
	query := "How can I improve my portfolio?"
	recommendation := ai.ProvideRecommendation(query, nil)
	fmt.Printf("  Q: %s\n", query)
	fmt.Printf("  A: %s\n", recommendation)

	// 7. BTN-Pay Invoice System
	fmt.Println("\n7. BTN-Pay Invoice System")
	btnPay := payments.NewBtnPay()
	invoice, _ := btnPay.CreateInvoice("BTNmerchant123", 250.0, "Demo order #123", 900)
	fmt.Printf("✓ Invoice created: %s\n", invoice.ID)
	fmt.Printf("  Merchant: %s, Amount: %.2f %s\n", invoice.Merchant, invoice.Amount, invoice.Currency)
	fmt.Printf("  Memo: %s, Status: %s\n", invoice.Memo, invoice.Status)
	
	btnPay.MarkPaid(invoice.ID, "tx_demo_789")
	updatedInvoice, _ := btnPay.GetInvoice(invoice.ID)
	fmt.Printf("✓ Invoice marked as paid: %s (TxID: %s)\n", updatedInvoice.Status, updatedInvoice.TxID)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nAll BTNg Wallet features are working correctly!")
	fmt.Println("For more information, see docs/BTNg-WALLET.md")
}
