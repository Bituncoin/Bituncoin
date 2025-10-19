package identity

import (
	"strings"
	"testing"
)

func TestGenerateBitcoinStyleAddress(t *testing.T) {
	addr, err := GenerateBitcoinStyleAddress()
	if err != nil {
		t.Fatalf("Failed to generate Bitcoin-style address: %v", err)
	}

	// Check that address is not empty
	if addr.Address == "" {
		t.Error("Generated address is empty")
	}

	// Check that address starts with "Btu"
	if !strings.HasPrefix(addr.Address, "Btu") {
		t.Errorf("Address should start with 'Btu', got: %s", addr.Address)
	}

	// Check that private key is not empty
	if addr.PrivateKey == "" {
		t.Error("Private key is empty")
	}

	// Check that public key is not empty
	if addr.PublicKey == "" {
		t.Error("Public key is empty")
	}

	// Check address format (example: BtuzP1eP5QGefi2DMPTfTL5SLmv7DivfNa)
	// Should be Base58 encoded with "Btu" prefix
	if len(addr.Address) < 26 || len(addr.Address) > 40 {
		t.Errorf("Address length should be between 26 and 40 characters, got %d", len(addr.Address))
	}
}

func TestGenerateMultipleBitcoinStyleAddresses(t *testing.T) {
	addresses := make(map[string]bool)

	// Generate multiple addresses and ensure they're unique
	for i := 0; i < 10; i++ {
		addr, err := GenerateBitcoinStyleAddress()
		if err != nil {
			t.Fatalf("Failed to generate address %d: %v", i, err)
		}

		if addresses[addr.Address] {
			t.Errorf("Duplicate address generated: %s", addr.Address)
		}
		addresses[addr.Address] = true

		// Ensure all start with Btu
		if !strings.HasPrefix(addr.Address, "Btu") {
			t.Errorf("Address %d does not start with Btu: %s", i, addr.Address)
		}
	}

	if len(addresses) != 10 {
		t.Errorf("Expected 10 unique addresses, got %d", len(addresses))
	}
}

func TestGenerateEthereumStyleAddress(t *testing.T) {
	addr, err := GenerateEthereumStyleAddress()
	if err != nil {
		t.Fatalf("Failed to generate Ethereum-style address: %v", err)
	}

	// Check that address is not empty
	if addr.Address == "" {
		t.Error("Generated address is empty")
	}

	// Check that address starts with "0x"
	if !strings.HasPrefix(addr.Address, "0x") {
		t.Errorf("Address should start with '0x', got: %s", addr.Address)
	}

	// Check that private key is not empty
	if addr.PrivateKey == "" {
		t.Error("Private key is empty")
	}

	// Check that public key is not empty
	if addr.PublicKey == "" {
		t.Error("Public key is empty")
	}

	// Check address format (0x + 40 hex characters)
	if len(addr.Address) != 42 {
		t.Errorf("Ethereum address should be 42 characters (0x + 40 hex), got %d", len(addr.Address))
	}

	// Check that the rest is hexadecimal
	addressHex := addr.Address[2:]
	for _, c := range addressHex {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			t.Errorf("Address contains non-hexadecimal character: %c", c)
		}
	}
}

func TestGenerateMultipleEthereumStyleAddresses(t *testing.T) {
	addresses := make(map[string]bool)

	// Generate multiple addresses and ensure they're unique
	for i := 0; i < 10; i++ {
		addr, err := GenerateEthereumStyleAddress()
		if err != nil {
			t.Fatalf("Failed to generate address %d: %v", i, err)
		}

		if addresses[addr.Address] {
			t.Errorf("Duplicate address generated: %s", addr.Address)
		}
		addresses[addr.Address] = true

		// Ensure all start with 0x and are correct length
		if !strings.HasPrefix(addr.Address, "0x") || len(addr.Address) != 42 {
			t.Errorf("Invalid Ethereum address format: %s", addr.Address)
		}
	}

	if len(addresses) != 10 {
		t.Errorf("Expected 10 unique addresses, got %d", len(addresses))
	}
}

func TestValidateBitcoinStyleAddress(t *testing.T) {
	// Generate a valid address
	addr, err := GenerateBitcoinStyleAddress()
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Validate it
	err = ValidateBitcoinStyleAddress(addr.Address)
	if err != nil {
		t.Errorf("Valid address failed validation: %v", err)
	}
}

func TestValidateBitcoinStyleAddressEmpty(t *testing.T) {
	err := ValidateBitcoinStyleAddress("")
	if err == nil {
		t.Error("Expected error for empty address, got nil")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected error about empty address, got: %v", err)
	}
}

func TestValidateBitcoinStyleAddressInvalidPrefix(t *testing.T) {
	err := ValidateBitcoinStyleAddress("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	if err == nil {
		t.Error("Expected error for address without Btu prefix, got nil")
	}
}

func TestValidateBitcoinStyleAddressInvalidBase58(t *testing.T) {
	err := ValidateBitcoinStyleAddress("Btu0OIl")
	if err == nil {
		t.Error("Expected error for invalid Base58 characters, got nil")
	}
}

func TestValidateBitcoinStyleAddressTooShort(t *testing.T) {
	err := ValidateBitcoinStyleAddress("Btu123")
	if err == nil {
		t.Error("Expected error for too short address, got nil")
	}
}

func TestValidateEthereumStyleAddress(t *testing.T) {
	// Generate a valid address
	addr, err := GenerateEthereumStyleAddress()
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Validate it
	err = ValidateEthereumStyleAddress(addr.Address)
	if err != nil {
		t.Errorf("Valid address failed validation: %v", err)
	}
}

func TestValidateEthereumStyleAddressEmpty(t *testing.T) {
	err := ValidateEthereumStyleAddress("")
	if err == nil {
		t.Error("Expected error for empty address, got nil")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected error about empty address, got: %v", err)
	}
}

func TestValidateEthereumStyleAddressInvalidPrefix(t *testing.T) {
	err := ValidateEthereumStyleAddress("1234567890abcdef1234567890abcdef12345678")
	if err == nil {
		t.Error("Expected error for address without 0x prefix, got nil")
	}
	if !strings.Contains(err.Error(), "0x") {
		t.Errorf("Expected error about 0x prefix, got: %v", err)
	}
}

func TestValidateEthereumStyleAddressInvalidLength(t *testing.T) {
	testCases := []struct {
		address string
		name    string
	}{
		{"0x123", "too short"},
		{"0x1234567890abcdef1234567890abcdef123456789", "too long"},
		{"0x1234567890abcdef1234567890abcdef1234567", "one char short"},
	}

	for _, tc := range testCases {
		err := ValidateEthereumStyleAddress(tc.address)
		if err == nil {
			t.Errorf("Expected error for %s address, got nil", tc.name)
		}
	}
}

func TestValidateEthereumStyleAddressInvalidHex(t *testing.T) {
	err := ValidateEthereumStyleAddress("0x1234567890abcdef1234567890abcdefghijklmn")
	if err == nil {
		t.Error("Expected error for invalid hex characters, got nil")
	}
	if !strings.Contains(err.Error(), "hexadecimal") {
		t.Errorf("Expected error about hexadecimal characters, got: %v", err)
	}
}

func TestValidateEthereumStyleAddressValidExample(t *testing.T) {
	// Test with a valid example format
	validAddress := "0x1234567890abcdef1234567890abcdef12345678"
	err := ValidateEthereumStyleAddress(validAddress)
	if err != nil {
		t.Errorf("Valid example address failed validation: %v", err)
	}
}

func TestBase58EncodeDecode(t *testing.T) {
	testCases := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{0xff, 0xfe, 0xfd, 0xfc},
		{0x12, 0x34, 0x56, 0x78, 0x90},
	}

	for _, tc := range testCases {
		encoded := base58Encode(tc)
		decoded, err := base58Decode(encoded)
		if err != nil {
			t.Errorf("Failed to decode Base58: %v", err)
		}

		if len(decoded) != len(tc) {
			t.Errorf("Decoded length mismatch: expected %d, got %d", len(tc), len(decoded))
		}

		for i := range tc {
			if decoded[i] != tc[i] {
				t.Errorf("Decoded byte mismatch at index %d: expected %x, got %x", i, tc[i], decoded[i])
			}
		}
	}
}

func TestBase58DecodeInvalidCharacter(t *testing.T) {
	// Base58 alphabet doesn't include 0, O, I, l
	invalidStrings := []string{
		"0",     // Zero is not in Base58
		"O",     // Letter O is not in Base58
		"I",     // Letter I is not in Base58
		"l",     // Letter l is not in Base58
		"Btu0",  // Contains invalid character
	}

	for _, invalid := range invalidStrings {
		_, err := base58Decode(invalid)
		if err == nil {
			t.Errorf("Expected error for invalid Base58 string '%s', got nil", invalid)
		}
	}
}

func TestBitcoinStyleAddressExampleFormat(t *testing.T) {
	// Generate an address and verify it matches expected format
	addr, err := GenerateBitcoinStyleAddress()
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Example format: BtuzP1eP5QGefi2DMPTfTL5SLmv7DivfNa
	// The address should:
	// 1. Start with "Btu"
	// 2. Be Base58 encoded
	// 3. Have checksum for error detection
	// 4. Be roughly 26-35 characters long

	if !strings.HasPrefix(addr.Address, "Btu") {
		t.Errorf("Address doesn't match expected format, should start with 'Btu': %s", addr.Address)
	}

	// Validate the address using our validation function
	if err := ValidateBitcoinStyleAddress(addr.Address); err != nil {
		t.Errorf("Generated address failed validation: %v", err)
	}
}

func TestEthereumStyleAddressExampleFormat(t *testing.T) {
	// Generate an address and verify it matches expected format
	addr, err := GenerateEthereumStyleAddress()
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Example format: 0x1234567890abcdef1234567890abcdef12345678
	// The address should:
	// 1. Start with "0x"
	// 2. Have exactly 40 hexadecimal characters after "0x"

	if !strings.HasPrefix(addr.Address, "0x") {
		t.Errorf("Address doesn't match expected format, should start with '0x': %s", addr.Address)
	}

	if len(addr.Address) != 42 {
		t.Errorf("Address length should be 42 (0x + 40 hex chars), got %d: %s", len(addr.Address), addr.Address)
	}

	// Validate the address using our validation function
	if err := ValidateEthereumStyleAddress(addr.Address); err != nil {
		t.Errorf("Generated address failed validation: %v", err)
	}
}
