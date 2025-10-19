package identity

import (
	"strings"
	"testing"
)

func TestAddressManagerGenerateBitcoinAddress(t *testing.T) {
	manager := NewAddressManager()

	addr, err := manager.GenerateBitcoinAddress("Test Bitcoin Wallet")
	if err != nil {
		t.Fatalf("Failed to generate Bitcoin address: %v", err)
	}

	if !strings.HasPrefix(addr.Address, "Btu") {
		t.Errorf("Bitcoin address should start with 'Btu', got: %s", addr.Address)
	}

	if addr.Label != "Test Bitcoin Wallet" {
		t.Errorf("Expected label 'Test Bitcoin Wallet', got: %s", addr.Label)
	}

	if addr.PublicKey == "" {
		t.Error("Public key should not be empty")
	}

	if addr.PrivateKey == "" {
		t.Error("Private key should not be empty")
	}

	// Verify address is in manager
	retrieved, err := manager.GetAddress(addr.Address)
	if err != nil {
		t.Errorf("Failed to retrieve address from manager: %v", err)
	}

	if retrieved.Address != addr.Address {
		t.Errorf("Retrieved address mismatch: expected %s, got %s", addr.Address, retrieved.Address)
	}
}

func TestAddressManagerGenerateEthereumAddress(t *testing.T) {
	manager := NewAddressManager()

	addr, err := manager.GenerateEthereumAddress("Test Ethereum Wallet")
	if err != nil {
		t.Fatalf("Failed to generate Ethereum address: %v", err)
	}

	if !strings.HasPrefix(addr.Address, "0x") {
		t.Errorf("Ethereum address should start with '0x', got: %s", addr.Address)
	}

	if len(addr.Address) != 42 {
		t.Errorf("Ethereum address should be 42 characters, got %d", len(addr.Address))
	}

	if addr.Label != "Test Ethereum Wallet" {
		t.Errorf("Expected label 'Test Ethereum Wallet', got: %s", addr.Label)
	}

	// Verify address is in manager
	retrieved, err := manager.GetAddress(addr.Address)
	if err != nil {
		t.Errorf("Failed to retrieve address from manager: %v", err)
	}

	if retrieved.Address != addr.Address {
		t.Errorf("Retrieved address mismatch: expected %s, got %s", addr.Address, retrieved.Address)
	}
}

func TestAddressManagerMultipleAddressTypes(t *testing.T) {
	manager := NewAddressManager()

	// Generate one of each type
	gldAddr, err := manager.GenerateAddress("GLD Wallet")
	if err != nil {
		t.Fatalf("Failed to generate GLD address: %v", err)
	}

	btcAddr, err := manager.GenerateBitcoinAddress("BTC Wallet")
	if err != nil {
		t.Fatalf("Failed to generate Bitcoin address: %v", err)
	}

	ethAddr, err := manager.GenerateEthereumAddress("ETH Wallet")
	if err != nil {
		t.Fatalf("Failed to generate Ethereum address: %v", err)
	}

	// Verify all addresses are in the manager
	allAddrs := manager.ListAddresses()
	if len(allAddrs) != 3 {
		t.Errorf("Expected 3 addresses in manager, got %d", len(allAddrs))
	}

	// Verify we can retrieve each one
	addresses := []string{gldAddr.Address, btcAddr.Address, ethAddr.Address}
	for _, addr := range addresses {
		_, err := manager.GetAddress(addr)
		if err != nil {
			t.Errorf("Failed to retrieve address %s: %v", addr, err)
		}
	}
}

func TestValidateAddressUniversal(t *testing.T) {
	// Generate valid addresses for testing
	btcAddr, _ := GenerateBitcoinStyleAddress()
	ethAddr, _ := GenerateEthereumStyleAddress()

	testCases := []struct {
		address string
		valid   bool
		name    string
	}{
		{"GLD5386547d51fe30e25ffd6616319524e486b17f27", true, "valid GLD address"},
		{btcAddr.Address, true, "valid Bitcoin-style address"},
		{ethAddr.Address, true, "valid Ethereum address"},
		{"0x123", false, "too short Ethereum address"},
		{"InvalidAddress", false, "invalid prefix"},
		{"", false, "empty address"},
		{"Btu123", false, "invalid Bitcoin address (too short)"},
	}

	for _, tc := range testCases {
		err := ValidateAddress(tc.address)
		if tc.valid && err != nil {
			t.Errorf("%s: expected valid, got error: %v", tc.name, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("%s: expected invalid, but validation passed", tc.name)
		}
	}
}

func TestAddressManagerDeleteAddress(t *testing.T) {
	manager := NewAddressManager()

	// Generate an address
	addr, err := manager.GenerateBitcoinAddress("Test Wallet")
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Verify it exists
	_, err = manager.GetAddress(addr.Address)
	if err != nil {
		t.Errorf("Address should exist: %v", err)
	}

	// Delete it
	err = manager.DeleteAddress(addr.Address)
	if err != nil {
		t.Errorf("Failed to delete address: %v", err)
	}

	// Verify it's gone
	_, err = manager.GetAddress(addr.Address)
	if err == nil {
		t.Error("Address should not exist after deletion")
	}
}

func TestAddressManagerSignMessage(t *testing.T) {
	manager := NewAddressManager()

	// Generate a Bitcoin-style address
	addr, err := manager.GenerateBitcoinAddress("Signing Wallet")
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Sign a message
	message := "Test message"
	signature, err := manager.SignMessage(addr.Address, message)
	if err != nil {
		t.Errorf("Failed to sign message: %v", err)
	}

	if signature == "" {
		t.Error("Signature should not be empty")
	}

	// Verify signature
	valid := VerifySignature(addr.Address, message, signature)
	if !valid {
		t.Error("Signature verification failed")
	}
}

func TestAddressManagerUniqueAddresses(t *testing.T) {
	manager := NewAddressManager()

	// Generate multiple Bitcoin addresses
	addresses := make(map[string]bool)
	for i := 0; i < 10; i++ {
		addr, err := manager.GenerateBitcoinAddress("Wallet")
		if err != nil {
			t.Fatalf("Failed to generate address %d: %v", i, err)
		}

		if addresses[addr.Address] {
			t.Errorf("Duplicate address generated: %s", addr.Address)
		}
		addresses[addr.Address] = true
	}

	// Generate multiple Ethereum addresses
	for i := 0; i < 10; i++ {
		addr, err := manager.GenerateEthereumAddress("Wallet")
		if err != nil {
			t.Fatalf("Failed to generate address %d: %v", i, err)
		}

		if addresses[addr.Address] {
			t.Errorf("Duplicate address generated: %s", addr.Address)
		}
		addresses[addr.Address] = true
	}

	if len(addresses) != 20 {
		t.Errorf("Expected 20 unique addresses, got %d", len(addresses))
	}
}

func TestAddressManagerCreatedAt(t *testing.T) {
	manager := NewAddressManager()

	// Test Bitcoin-style address
	btcAddr, err := manager.GenerateBitcoinAddress("Test Wallet")
	if err != nil {
		t.Fatalf("Failed to generate Bitcoin address: %v", err)
	}

	if btcAddr.CreatedAt == 0 {
		t.Error("CreatedAt timestamp should not be 0")
	}

	// Test Ethereum-style address
	ethAddr, err := manager.GenerateEthereumAddress("Test Wallet")
	if err != nil {
		t.Fatalf("Failed to generate Ethereum address: %v", err)
	}

	if ethAddr.CreatedAt == 0 {
		t.Error("CreatedAt timestamp should not be 0")
	}

	// Test legacy GLD address
	gldAddr, err := manager.GenerateAddress("Test Wallet")
	if err != nil {
		t.Fatalf("Failed to generate GLD address: %v", err)
	}

	if gldAddr.CreatedAt == 0 {
		t.Error("CreatedAt timestamp should not be 0")
	}
}

func TestValidateAddressEdgeCases(t *testing.T) {
	testCases := []struct {
		address string
		name    string
		valid   bool
	}{
		{"0x", "0x prefix only", false},
		{"Bt", "Bt prefix (too short)", false},
		{"GL", "GL prefix (too short)", false},
		{"", "empty string", false},
		{"a", "single character", false},
	}

	for _, tc := range testCases {
		err := ValidateAddress(tc.address)
		if tc.valid && err != nil {
			t.Errorf("%s: expected valid, got error: %v", tc.name, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("%s: expected invalid, but validation passed", tc.name)
		}
	}
}
