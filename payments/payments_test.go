package payments

import "testing"

func TestMobileMoneyPayment(t *testing.T) {
	mms := NewMobileMoneyService()

	tx, err := mms.InitiatePayment(MTN, "+233244123456", 100.0, "GHS")
	if err != nil {
		t.Fatalf("Failed to initiate payment: %v", err)
	}

	if tx.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", tx.Status)
	}

	if tx.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %f", tx.Amount)
	}

	if tx.Provider != MTN {
		t.Errorf("Expected provider MTN, got %s", tx.Provider)
	}
}

func TestMobileMoneyInvalidAmount(t *testing.T) {
	mms := NewMobileMoneyService()

	_, err := mms.InitiatePayment(MTN, "+233244123456", -10.0, "GHS")
	if err == nil {
		t.Error("Expected error for negative amount")
	}
}

func TestQRCodeGeneration(t *testing.T) {
	qrs := NewQRCodeService()

	qrData, err := qrs.GeneratePaymentQR("BTNaddr123", 50.0, "BTN", "Test payment")
	if err != nil {
		t.Fatalf("Failed to generate QR: %v", err)
	}

	if qrData == "" {
		t.Error("Expected non-empty QR data")
	}

	// Parse the generated QR
	data, err := qrs.ParsePaymentQR(qrData)
	if err != nil {
		t.Fatalf("Failed to parse QR: %v", err)
	}

	if data.Address != "BTNaddr123" {
		t.Errorf("Expected address 'BTNaddr123', got '%s'", data.Address)
	}

	if data.Amount != 50.0 {
		t.Errorf("Expected amount 50.0, got %f", data.Amount)
	}
}

func TestCardIssuance(t *testing.T) {
	cs := NewCardService()

	card, err := cs.IssueCard("John Doe", "BTNaddr123", MasterCard, Virtual)
	if err != nil {
		t.Fatalf("Failed to issue card: %v", err)
	}

	if card.HolderName != "John Doe" {
		t.Errorf("Expected holder name 'John Doe', got '%s'", card.HolderName)
	}

	if card.Type != MasterCard {
		t.Errorf("Expected card type MasterCard, got %s", card.Type)
	}

	if card.Category != Virtual {
		t.Errorf("Expected category Virtual, got %s", card.Category)
	}

	if card.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", card.Status)
	}
}

func TestCardPayment(t *testing.T) {
	cs := NewCardService()

	// Issue a card
	card, _ := cs.IssueCard("John Doe", "BTNaddr123", Visa, Physical)

	// Load the card
	err := cs.LoadCard(card.ID, 1000.0)
	if err != nil {
		t.Fatalf("Failed to load card: %v", err)
	}

	// Process payment
	tx, err := cs.ProcessPayment(card.ID, 50.0, "Test Merchant")
	if err != nil {
		t.Fatalf("Failed to process payment: %v", err)
	}

	if tx.Amount != 50.0 {
		t.Errorf("Expected amount 50.0, got %f", tx.Amount)
	}

	if tx.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", tx.Status)
	}

	// Check card balance
	updatedCard, _ := cs.GetCard(card.ID)
	if updatedCard.Balance != 950.0 {
		t.Errorf("Expected balance 950.0, got %f", updatedCard.Balance)
	}
}

func TestCardInsufficientBalance(t *testing.T) {
	cs := NewCardService()

	card, _ := cs.IssueCard("John Doe", "BTNaddr123", MasterCard, Virtual)
	cs.LoadCard(card.ID, 50.0)

	_, err := cs.ProcessPayment(card.ID, 100.0, "Test Merchant")
	if err == nil {
		t.Error("Expected error for insufficient balance")
	}
}
