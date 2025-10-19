package bituncoin

import (
	"testing"
)

func TestNewBituncoin(t *testing.T) {
	btn := NewBituncoin()

	if btn.Name != "Bituncoin" {
		t.Errorf("expected name Bituncoin, got %s", btn.Name)
	}

	if btn.Symbol != "BTN" {
		t.Errorf("expected symbol BTN, got %s", btn.Symbol)
	}

	if btn.MaxSupply != 100000000 {
		t.Errorf("expected max supply 100000000, got %d", btn.MaxSupply)
	}

	if btn.Decimals != 8 {
		t.Errorf("expected 8 decimals, got %d", btn.Decimals)
	}
}

func TestCreateTransaction(t *testing.T) {
	btn := NewBituncoin()
	tx, err := btn.CreateTransaction("BTNaddr1", "BTNaddr2", 100.0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tx.From != "BTNaddr1" {
		t.Errorf("expected from BTNaddr1, got %s", tx.From)
	}

	if tx.To != "BTNaddr2" {
		t.Errorf("expected to BTNaddr2, got %s", tx.To)
	}

	if tx.Amount != 100.0 {
		t.Errorf("expected amount 100.0, got %f", tx.Amount)
	}

	expectedFee := 100.0 * 0.001
	if tx.Fee != expectedFee {
		t.Errorf("expected fee %f, got %f", expectedFee, tx.Fee)
	}
}

func TestCreateTransactionInvalidAmount(t *testing.T) {
	btn := NewBituncoin()
	_, err := btn.CreateTransaction("BTNaddr1", "BTNaddr2", 0)

	if err == nil {
		t.Error("expected error for zero amount")
	}

	_, err = btn.CreateTransaction("BTNaddr1", "BTNaddr2", -10)
	if err == nil {
		t.Error("expected error for negative amount")
	}
}

func TestCreateTransactionInvalidAddresses(t *testing.T) {
	btn := NewBituncoin()
	_, err := btn.CreateTransaction("", "BTNaddr2", 100.0)

	if err == nil {
		t.Error("expected error for empty from address")
	}

	_, err = btn.CreateTransaction("BTNaddr1", "", 100.0)
	if err == nil {
		t.Error("expected error for empty to address")
	}
}

func TestValidateTransaction(t *testing.T) {
	btn := NewBituncoin()
	tx, _ := btn.CreateTransaction("BTNaddr1", "BTNaddr2", 100.0)

	err := btn.ValidateTransaction(tx)
	if err != nil {
		t.Errorf("expected valid transaction, got error: %v", err)
	}
}

func TestMint(t *testing.T) {
	btn := NewBituncoin()
	err := btn.Mint(1000000)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if btn.CircSupply != 1000000 {
		t.Errorf("expected circ supply 1000000, got %d", btn.CircSupply)
	}
}

func TestMintExceedsMaxSupply(t *testing.T) {
	btn := NewBituncoin()
	err := btn.Mint(100000001)

	if err == nil {
		t.Error("expected error when minting exceeds max supply")
	}
}

func TestGetTokenomics(t *testing.T) {
	btn := NewBituncoin()
	tokenomics := btn.GetTokenomics()

	if tokenomics["symbol"] != "BTN" {
		t.Errorf("expected symbol BTN in tokenomics")
	}

	if tokenomics["name"] != "Bituncoin" {
		t.Errorf("expected name Bituncoin in tokenomics")
	}
}
