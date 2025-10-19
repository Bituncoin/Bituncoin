package wallet

import (
	"testing"
)

func TestNewCardManager(t *testing.T) {
	cm := NewCardManager()

	if cm == nil {
		t.Fatal("Expected card manager to be created")
	}

	if cm.cards == nil || cm.transactions == nil {
		t.Fatal("Expected maps to be initialized")
	}
}

func TestCreateCard(t *testing.T) {
	cm := NewCardManager()

	card, err := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if card.CardholderName != "John Doe" {
		t.Errorf("Expected cardholder name 'John Doe', got '%s'", card.CardholderName)
	}

	if card.CardType != CardTypeVirtualVisa {
		t.Errorf("Expected card type VirtualVisa, got %v", card.CardType)
	}

	if card.Status != CardStatusPending {
		t.Errorf("Expected status Pending, got %v", card.Status)
	}

	if card.SpendingLimit != 5000.0 {
		t.Errorf("Expected spending limit 5000.0, got %f", card.SpendingLimit)
	}
}

func TestActivateCard(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)

	err := cm.ActivateCard(card.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	activatedCard, _ := cm.GetCard(card.ID)
	if activatedCard.Status != CardStatusActive {
		t.Errorf("Expected status Active, got %v", activatedCard.Status)
	}

	if activatedCard.ActivatedAt == 0 {
		t.Error("Expected ActivatedAt to be set")
	}
}

func TestSuspendCard(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	cm.ActivateCard(card.ID)

	err := cm.SuspendCard(card.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	suspendedCard, _ := cm.GetCard(card.ID)
	if suspendedCard.Status != CardStatusSuspended {
		t.Errorf("Expected status Suspended, got %v", suspendedCard.Status)
	}
}

func TestCancelCard(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)

	err := cm.CancelCard(card.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	cancelledCard, _ := cm.GetCard(card.ID)
	if cancelledCard.Status != CardStatusCancelled {
		t.Errorf("Expected status Cancelled, got %v", cancelledCard.Status)
	}
}

func TestProcessTransaction(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	cm.ActivateCard(card.ID)

	tx, err := cm.ProcessTransaction(card.ID, 100.0, "USD", "Test Merchant", "MCC_1234")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tx.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %f", tx.Amount)
	}

	if tx.Status != TxStatusCompleted {
		t.Errorf("Expected status Completed, got %v", tx.Status)
	}

	// Verify card spending was updated
	updatedCard, _ := cm.GetCard(card.ID)
	if updatedCard.SpentToday != 100.0 {
		t.Errorf("Expected SpentToday 100.0, got %f", updatedCard.SpentToday)
	}

	if updatedCard.SpentThisMonth != 100.0 {
		t.Errorf("Expected SpentThisMonth 100.0, got %f", updatedCard.SpentThisMonth)
	}
}

func TestProcessTransactionExceedsLimit(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 100.0, 50.0, 500.0)
	cm.ActivateCard(card.ID)

	// Try to process transaction exceeding spending limit
	_, err := cm.ProcessTransaction(card.ID, 150.0, "USD", "Test Merchant", "MCC_1234")
	if err == nil {
		t.Error("Expected error when exceeding spending limit")
	}
}

func TestProcessTransactionExceedsDailyLimit(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 100.0, 10000.0)
	cm.ActivateCard(card.ID)

	// Process first transaction
	cm.ProcessTransaction(card.ID, 60.0, "USD", "Merchant 1", "MCC_1234")

	// Try to process transaction exceeding daily limit
	_, err := cm.ProcessTransaction(card.ID, 50.0, "USD", "Merchant 2", "MCC_1234")
	if err == nil {
		t.Error("Expected error when exceeding daily limit")
	}
}

func TestGetCardsByWallet(t *testing.T) {
	cm := NewCardManager()

	// Create multiple cards
	card1, err1 := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	if err1 != nil {
		t.Fatalf("Failed to create card 1: %v", err1)
	}
	t.Logf("Created card1: ID=%s, WalletID=%s", card1.ID, card1.WalletID)
	
	card2, err2 := cm.CreateCard("wallet_123", "John Doe", CardTypePhysicalMastercard, 3000.0, 500.0, 5000.0)
	if err2 != nil {
		t.Fatalf("Failed to create card 2: %v", err2)
	}
	t.Logf("Created card2: ID=%s, WalletID=%s", card2.ID, card2.WalletID)
	
	card3, err3 := cm.CreateCard("wallet_456", "Jane Smith", CardTypeVirtualVisa, 2000.0, 400.0, 4000.0)
	if err3 != nil {
		t.Fatalf("Failed to create card 3: %v", err3)
	}
	t.Logf("Created card3: ID=%s, WalletID=%s", card3.ID, card3.WalletID)

	// Verify cards were stored
	if card1.WalletID != "wallet_123" {
		t.Errorf("Card 1 wallet ID mismatch: %s", card1.WalletID)
	}
	if card2.WalletID != "wallet_123" {
		t.Errorf("Card 2 wallet ID mismatch: %s", card2.WalletID)
	}
	if card3.WalletID != "wallet_456" {
		t.Errorf("Card 3 wallet ID mismatch: %s", card3.WalletID)
	}

	cards := cm.GetCardsByWallet("wallet_123")
	t.Logf("GetCardsByWallet returned %d cards", len(cards))
	for i, card := range cards {
		t.Logf("  Card %d: ID=%s, WalletID=%s", i, card.ID, card.WalletID)
	}
	
	if len(cards) != 2 {
		t.Errorf("Expected 2 cards for wallet_123, got %d", len(cards))
	}
}

func TestGetCardTransactions(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	cm.ActivateCard(card.ID)

	// Process multiple transactions
	cm.ProcessTransaction(card.ID, 100.0, "USD", "Merchant 1", "MCC_1234")
	cm.ProcessTransaction(card.ID, 50.0, "USD", "Merchant 2", "MCC_5678")

	txs, err := cm.GetCardTransactions(card.ID, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(txs) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(txs))
	}
}

func TestUpdateSpendingLimits(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)

	err := cm.UpdateSpendingLimits(card.ID, 6000.0, 1500.0, 15000.0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	updatedCard, _ := cm.GetCard(card.ID)
	if updatedCard.SpendingLimit != 6000.0 {
		t.Errorf("Expected spending limit 6000.0, got %f", updatedCard.SpendingLimit)
	}

	if updatedCard.DailyLimit != 1500.0 {
		t.Errorf("Expected daily limit 1500.0, got %f", updatedCard.DailyLimit)
	}
}

func TestResetDailySpending(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	cm.ActivateCard(card.ID)
	cm.ProcessTransaction(card.ID, 100.0, "USD", "Merchant", "MCC_1234")

	cm.ResetDailySpending()

	resetCard, _ := cm.GetCard(card.ID)
	if resetCard.SpentToday != 0 {
		t.Errorf("Expected SpentToday to be 0 after reset, got %f", resetCard.SpentToday)
	}

	if resetCard.SpentThisMonth != 100.0 {
		t.Errorf("Expected SpentThisMonth to remain 100.0, got %f", resetCard.SpentThisMonth)
	}
}

func TestResetMonthlySpending(t *testing.T) {
	cm := NewCardManager()
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	cm.ActivateCard(card.ID)
	cm.ProcessTransaction(card.ID, 100.0, "USD", "Merchant", "MCC_1234")

	cm.ResetMonthlySpending()

	resetCard, _ := cm.GetCard(card.ID)
	if resetCard.SpentThisMonth != 0 {
		t.Errorf("Expected SpentThisMonth to be 0 after reset, got %f", resetCard.SpentThisMonth)
	}
}

func TestInvalidCardOperations(t *testing.T) {
	cm := NewCardManager()

	// Test activating non-existent card
	err := cm.ActivateCard("invalid_card")
	if err == nil {
		t.Error("Expected error when activating non-existent card")
	}

	// Test getting non-existent card
	_, err = cm.GetCard("invalid_card")
	if err == nil {
		t.Error("Expected error when getting non-existent card")
	}

	// Test processing transaction with inactive card
	card, _ := cm.CreateCard("wallet_123", "John Doe", CardTypeVirtualVisa, 5000.0, 1000.0, 10000.0)
	_, err = cm.ProcessTransaction(card.ID, 100.0, "USD", "Merchant", "MCC_1234")
	if err == nil {
		t.Error("Expected error when processing transaction with inactive card")
	}
}
