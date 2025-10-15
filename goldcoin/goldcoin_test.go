package goldcoin

import (
	"testing"
)

func TestNewGoldCoin(t *testing.T) {
	gc := NewGoldCoin()
	
	if gc.Name != "Gold-Coin" {
		t.Errorf("Expected name Gold-Coin, got %s", gc.Name)
	}
	
	if gc.Symbol != "GLD" {
		t.Errorf("Expected symbol GLD, got %s", gc.Symbol)
	}
	
	if gc.MaxSupply != 100000000 {
		t.Errorf("Expected max supply 100000000, got %d", gc.MaxSupply)
	}
	
	if gc.Decimals != 8 {
		t.Errorf("Expected 8 decimals, got %d", gc.Decimals)
	}
}

func TestCreateTransaction(t *testing.T) {
	gc := NewGoldCoin()
	
	tx, err := gc.CreateTransaction("from_addr", "to_addr", 100.0)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
	
	if tx.Amount != 100.0 {
		t.Errorf("Expected amount 100.0, got %f", tx.Amount)
	}
	
	expectedFee := 100.0 * gc.TxFee
	if tx.Fee != expectedFee {
		t.Errorf("Expected fee %f, got %f", expectedFee, tx.Fee)
	}
	
	if tx.From != "from_addr" {
		t.Errorf("Expected from address from_addr, got %s", tx.From)
	}
	
	if tx.To != "to_addr" {
		t.Errorf("Expected to address to_addr, got %s", tx.To)
	}
}

func TestCreateTransactionInvalidAmount(t *testing.T) {
	gc := NewGoldCoin()
	
	_, err := gc.CreateTransaction("from_addr", "to_addr", 0)
	if err == nil {
		t.Error("Expected error for zero amount, got nil")
	}
	
	_, err = gc.CreateTransaction("from_addr", "to_addr", -100)
	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
}

func TestCreateTransactionInvalidAddresses(t *testing.T) {
	gc := NewGoldCoin()
	
	_, err := gc.CreateTransaction("", "to_addr", 100)
	if err == nil {
		t.Error("Expected error for empty from address, got nil")
	}
	
	_, err = gc.CreateTransaction("from_addr", "", 100)
	if err == nil {
		t.Error("Expected error for empty to address, got nil")
	}
}

func TestValidateTransaction(t *testing.T) {
	gc := NewGoldCoin()
	
	tx, _ := gc.CreateTransaction("from_addr", "to_addr", 100.0)
	
	err := gc.ValidateTransaction(tx)
	if err != nil {
		t.Errorf("Transaction validation failed: %v", err)
	}
}

func TestMint(t *testing.T) {
	gc := NewGoldCoin()
	
	err := gc.Mint(1000000)
	if err != nil {
		t.Fatalf("Failed to mint: %v", err)
	}
	
	if gc.CircSupply != 1000000 {
		t.Errorf("Expected circulating supply 1000000, got %d", gc.CircSupply)
	}
}

func TestMintExceedsMaxSupply(t *testing.T) {
	gc := NewGoldCoin()
	
	err := gc.Mint(gc.MaxSupply + 1)
	if err == nil {
		t.Error("Expected error when minting exceeds max supply, got nil")
	}
}

func TestGetTokenomics(t *testing.T) {
	gc := NewGoldCoin()
	
	tokenomics := gc.GetTokenomics()
	
	if tokenomics["name"] != "Gold-Coin" {
		t.Errorf("Expected name Gold-Coin in tokenomics, got %v", tokenomics["name"])
	}
	
	if tokenomics["symbol"] != "GLD" {
		t.Errorf("Expected symbol GLD in tokenomics, got %v", tokenomics["symbol"])
	}
}
