package wallet

import (
	"testing"
)

func TestNewUniversalWallet(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	if wallet == nil {
		t.Fatal("Expected wallet to be created")
	}

	if wallet.Owner != "test_user" {
		t.Errorf("Expected owner to be 'test_user', got '%s'", wallet.Owner)
	}

	if wallet.Portfolio == nil {
		t.Fatal("Expected portfolio to be initialized")
	}

	if len(wallet.Portfolio.Balances) != 6 {
		t.Errorf("Expected 6 currencies, got %d", len(wallet.Portfolio.Balances))
	}
}

func TestGetBalance(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	balance, err := wallet.GetBalance(BTN)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if balance.Currency != BTN {
		t.Errorf("Expected currency BTN, got %v", balance.Currency)
	}

	if balance.Amount != 0 {
		t.Errorf("Expected initial balance to be 0, got %f", balance.Amount)
	}
}

func TestUpdateBalance(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	err := wallet.UpdateBalance(BTN, 1000.0, 10000.0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	balance, _ := wallet.GetBalance(BTN)
	if balance.Amount != 1000.0 {
		t.Errorf("Expected balance to be 1000.0, got %f", balance.Amount)
	}

	if balance.USDValue != 10000.0 {
		t.Errorf("Expected USD value to be 10000.0, got %f", balance.USDValue)
	}

	if wallet.Portfolio.TotalUSD != 10000.0 {
		t.Errorf("Expected total USD to be 10000.0, got %f", wallet.Portfolio.TotalUSD)
	}
}

func TestAddTransaction(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	tx := &WalletTransaction{
		ID:       "tx_123",
		Type:     TxTypeSend,
		Currency: BTN,
		Amount:   100.0,
		Status:   TxStatusCompleted,
	}

	err := wallet.AddTransaction(tx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(wallet.Transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(wallet.Transactions))
	}
}

func TestGetTransactionHistory(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	// Add multiple transactions
	transactions := []*WalletTransaction{
		{ID: "tx_1", Type: TxTypeSend, Currency: BTN, Amount: 100.0},
		{ID: "tx_2", Type: TxTypeReceive, Currency: BTC, Amount: 0.5},
		{ID: "tx_3", Type: TxTypeSwap, Currency: ETH, Amount: 2.0},
	}

	for _, tx := range transactions {
		wallet.AddTransaction(tx)
	}

	// Test getting all transactions
	history := wallet.GetTransactionHistory(nil, nil, 0)
	if len(history) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(history))
	}

	// Test filtering by currency
	btnCurrency := BTN
	btnHistory := wallet.GetTransactionHistory(&btnCurrency, nil, 0)
	if len(btnHistory) != 1 {
		t.Errorf("Expected 1 BTN transaction, got %d", len(btnHistory))
	}

	// Test filtering by type
	sendType := TxTypeSend
	sendHistory := wallet.GetTransactionHistory(nil, &sendType, 0)
	if len(sendHistory) != 1 {
		t.Errorf("Expected 1 send transaction, got %d", len(sendHistory))
	}

	// Test limit
	limitedHistory := wallet.GetTransactionHistory(nil, nil, 2)
	if len(limitedHistory) != 2 {
		t.Errorf("Expected 2 transactions with limit, got %d", len(limitedHistory))
	}
}

func TestGetPortfolioSummary(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	// Update some balances
	wallet.UpdateBalance(BTN, 1000.0, 10000.0)
	wallet.UpdateBalance(BTC, 0.5, 22500.0)

	summary := wallet.GetPortfolioSummary()

	if summary["walletId"] != wallet.ID {
		t.Errorf("Expected wallet ID to match")
	}

	if summary["owner"] != "test_user" {
		t.Errorf("Expected owner to be 'test_user'")
	}

	if summary["totalUSD"].(float64) != 32500.0 {
		t.Errorf("Expected total USD to be 32500.0, got %f", summary["totalUSD"].(float64))
	}
}

func TestMultipleCurrencies(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	currencies := []struct {
		currency Currency
		amount   float64
		usdValue float64
	}{
		{BTN, 1000.0, 10000.0},
		{BTC, 0.5, 22500.0},
		{ETH, 5.0, 15000.0},
		{USDT, 5000.0, 5000.0},
		{BNB, 10.0, 6000.0},
	}

	for _, c := range currencies {
		err := wallet.UpdateBalance(c.currency, c.amount, c.usdValue)
		if err != nil {
			t.Fatalf("Error updating balance for %s: %v", c.currency, err)
		}
	}

	expectedTotal := 58500.0
	if wallet.Portfolio.TotalUSD != expectedTotal {
		t.Errorf("Expected total USD to be %f, got %f", expectedTotal, wallet.Portfolio.TotalUSD)
	}
}

func TestInvalidCurrency(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	_, err := wallet.GetBalance(Currency("INVALID"))
	if err == nil {
		t.Error("Expected error for invalid currency")
	}
}

func TestAddNilTransaction(t *testing.T) {
	wallet := NewUniversalWallet("test_user")

	err := wallet.AddTransaction(nil)
	if err == nil {
		t.Error("Expected error when adding nil transaction")
	}
}
