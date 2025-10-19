package identity

import (
	"crypto/sha256"
	"strings"
	"testing"
)

func TestGenerateBitcoinStyleAddress(t *testing.T) {
	// Test with valid public key
	publicKey := []byte("test_public_key_data_here_123456")
	address, err := GenerateBitcoinStyleAddress(publicKey)
	
	if err != nil {
		t.Fatalf("Failed to generate Bitcoin-style address: %v", err)
	}
	
	// Check prefix
	if !strings.HasPrefix(address, "Btu") {
		t.Errorf("Bitcoin-style address should start with 'Btu', got: %s", address)
	}
	
	// Check length is reasonable
	if len(address) < 29 || len(address) > 40 {
		t.Errorf("Bitcoin-style address length should be between 29-40, got: %d", len(address))
	}
	
	// Verify address is valid
	if err := ValidateBitcoinStyleAddress(address); err != nil {
		t.Errorf("Generated address failed validation: %v", err)
	}
}

func TestGenerateBitcoinStyleAddressEmptyKey(t *testing.T) {
	// Test with empty public key
	_, err := GenerateBitcoinStyleAddress([]byte{})
	
	if err == nil {
		t.Error("Expected error for empty public key, got nil")
	}
}

func TestGenerateEthereumStyleAddress(t *testing.T) {
	// Test with valid public key
	publicKey := []byte("test_public_key_data_here_123456")
	address, err := GenerateEthereumStyleAddress(publicKey)
	
	if err != nil {
		t.Fatalf("Failed to generate Ethereum-style address: %v", err)
	}
	
	// Check prefix
	if !strings.HasPrefix(address, "0x") {
		t.Errorf("Ethereum-style address should start with '0x', got: %s", address)
	}
	
	// Check length (0x + 40 hex chars = 42)
	if len(address) != 42 {
		t.Errorf("Ethereum-style address should be 42 characters, got: %d", len(address))
	}
	
	// Verify address is valid
	if err := ValidateEthereumStyleAddress(address); err != nil {
		t.Errorf("Generated address failed validation: %v", err)
	}
}

func TestGenerateEthereumStyleAddressEmptyKey(t *testing.T) {
	// Test with empty public key
	_, err := GenerateEthereumStyleAddress([]byte{})
	
	if err == nil {
		t.Error("Expected error for empty public key, got nil")
	}
}

func TestValidateBitcoinStyleAddress(t *testing.T) {
	testCases := []struct {
		name        string
		address     string
		shouldError bool
	}{
		{
			name:        "Valid generated address",
			address:     "",  // Will be generated
			shouldError: false,
		},
		{
			name:        "Missing Btu prefix",
			address:     "xyz1234567890123456789012345678",
			shouldError: true,
		},
		{
			name:        "Too short",
			address:     "Btu123",
			shouldError: true,
		},
		{
			name:        "Too long",
			address:     "Btu12345678901234567890123456789012345678901234567890",
			shouldError: true,
		},
		{
			name:        "Invalid base58 character (0)",
			address:     "Btu0000000000000000000000000",
			shouldError: true,
		},
		{
			name:        "Invalid base58 character (O)",
			address:     "BtuOOOOOOOOOOOOOOOOOOOOOOOOOO",
			shouldError: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			address := tc.address
			
			// Generate a valid address for the first test case
			if tc.name == "Valid generated address" {
				publicKey := []byte("test_public_key_for_validation")
				var err error
				address, err = GenerateBitcoinStyleAddress(publicKey)
				if err != nil {
					t.Fatalf("Failed to generate address: %v", err)
				}
			}
			
			err := ValidateBitcoinStyleAddress(address)
			
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for address '%s', got nil", address)
			}
			
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for address '%s', got: %v", address, err)
			}
		})
	}
}

func TestValidateEthereumStyleAddress(t *testing.T) {
	testCases := []struct {
		name        string
		address     string
		shouldError bool
	}{
		{
			name:        "Valid address",
			address:     "0x1234567890abcdef1234567890abcdef12345678",
			shouldError: false,
		},
		{
			name:        "Valid address uppercase",
			address:     "0x1234567890ABCDEF1234567890ABCDEF12345678",
			shouldError: false,
		},
		{
			name:        "Valid address mixed case",
			address:     "0x1234567890AbCdEf1234567890aBcDeF12345678",
			shouldError: false,
		},
		{
			name:        "Missing 0x prefix",
			address:     "1234567890abcdef1234567890abcdef12345678",
			shouldError: true,
		},
		{
			name:        "Too short",
			address:     "0x123456789",
			shouldError: true,
		},
		{
			name:        "Too long",
			address:     "0x1234567890abcdef1234567890abcdef123456789",
			shouldError: true,
		},
		{
			name:        "Invalid hex character (g)",
			address:     "0x1234567890abcdef1234567890abcdefg2345678",
			shouldError: true,
		},
		{
			name:        "Invalid hex character (space)",
			address:     "0x1234567890abcdef 234567890abcdef12345678",
			shouldError: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEthereumStyleAddress(tc.address)
			
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for address '%s', got nil", tc.address)
			}
			
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for address '%s', got: %v", tc.address, err)
			}
		})
	}
}

func TestBase58EncodeDecodRoundtrip(t *testing.T) {
	testCases := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{0xFF, 0xFE, 0xFD, 0xFC},
		[]byte("Hello World"),
		make([]byte, 20), // All zeros
	}
	
	for _, original := range testCases {
		encoded := base58Encode(original)
		decoded, err := base58Decode(encoded)
		
		if err != nil {
			t.Errorf("Failed to decode base58: %v", err)
			continue
		}
		
		if len(decoded) != len(original) {
			t.Errorf("Decoded length mismatch: expected %d, got %d", len(original), len(decoded))
			continue
		}
		
		for i := range original {
			if decoded[i] != original[i] {
				t.Errorf("Byte mismatch at position %d: expected %d, got %d", i, original[i], decoded[i])
				break
			}
		}
	}
}

func TestCalculateChecksum(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	checksum1 := calculateChecksum(payload)
	checksum2 := calculateChecksum(payload)
	
	// Checksums should be deterministic
	if len(checksum1) != 4 {
		t.Errorf("Checksum should be 4 bytes, got %d", len(checksum1))
	}
	
	for i := 0; i < 4; i++ {
		if checksum1[i] != checksum2[i] {
			t.Error("Checksums should be deterministic")
			break
		}
	}
	
	// Different payload should give different checksum
	payload2 := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x06}
	checksum3 := calculateChecksum(payload2)
	
	same := true
	for i := 0; i < 4; i++ {
		if checksum1[i] != checksum3[i] {
			same = false
			break
		}
	}
	
	if same {
		t.Error("Different payloads should produce different checksums")
	}
}

func TestGenerateAddressWithTypeBitcoin(t *testing.T) {
	address, privateKey, err := GenerateAddressWithType(BitcoinStyle)
	
	if err != nil {
		t.Fatalf("Failed to generate Bitcoin-style address: %v", err)
	}
	
	if !strings.HasPrefix(address, "Btu") {
		t.Errorf("Expected Bitcoin-style address to start with 'Btu', got: %s", address)
	}
	
	if len(privateKey) != 32 {
		t.Errorf("Expected private key to be 32 bytes, got: %d", len(privateKey))
	}
	
	// Verify generated address is valid
	if err := ValidateBitcoinStyleAddress(address); err != nil {
		t.Errorf("Generated Bitcoin-style address is invalid: %v", err)
	}
}

func TestGenerateAddressWithTypeEthereum(t *testing.T) {
	address, privateKey, err := GenerateAddressWithType(EthereumStyle)
	
	if err != nil {
		t.Fatalf("Failed to generate Ethereum-style address: %v", err)
	}
	
	if !strings.HasPrefix(address, "0x") {
		t.Errorf("Expected Ethereum-style address to start with '0x', got: %s", address)
	}
	
	if len(privateKey) != 32 {
		t.Errorf("Expected private key to be 32 bytes, got: %d", len(privateKey))
	}
	
	// Verify generated address is valid
	if err := ValidateEthereumStyleAddress(address); err != nil {
		t.Errorf("Generated Ethereum-style address is invalid: %v", err)
	}
}

func TestGenerateAddressWithTypeInvalid(t *testing.T) {
	_, _, err := GenerateAddressWithType("invalid_type")
	
	if err == nil {
		t.Error("Expected error for invalid address type, got nil")
	}
}

func TestValidateAddressTypeBitcoin(t *testing.T) {
	// Generate a Bitcoin-style address
	publicKey := []byte("test_public_key")
	address, _ := GenerateBitcoinStyleAddress(publicKey)
	
	addressType, err := ValidateAddressType(address)
	
	if err != nil {
		t.Fatalf("Failed to validate address type: %v", err)
	}
	
	if addressType != BitcoinStyle {
		t.Errorf("Expected Bitcoin-style address type, got: %s", addressType)
	}
}

func TestValidateAddressTypeEthereum(t *testing.T) {
	address := "0x1234567890abcdef1234567890abcdef12345678"
	
	addressType, err := ValidateAddressType(address)
	
	if err != nil {
		t.Fatalf("Failed to validate address type: %v", err)
	}
	
	if addressType != EthereumStyle {
		t.Errorf("Expected Ethereum-style address type, got: %s", addressType)
	}
}

func TestValidateAddressTypeInvalid(t *testing.T) {
	invalidAddresses := []string{
		"invalid_address",
		"",
		"xyz123",
		"0xinvalid",
	}
	
	for _, address := range invalidAddresses {
		_, err := ValidateAddressType(address)
		
		if err == nil {
			t.Errorf("Expected error for invalid address '%s', got nil", address)
		}
	}
}

func TestBitcoinStyleAddressUniqueness(t *testing.T) {
	// Generate multiple addresses and ensure they're unique
	addresses := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		publicKey := make([]byte, 32)
		copy(publicKey, []byte{byte(i)})
		hash := sha256.Sum256(publicKey)
		
		address, err := GenerateBitcoinStyleAddress(hash[:])
		if err != nil {
			t.Fatalf("Failed to generate address: %v", err)
		}
		
		if addresses[address] {
			t.Errorf("Duplicate address generated: %s", address)
		}
		addresses[address] = true
	}
}

func TestEthereumStyleAddressUniqueness(t *testing.T) {
	// Generate multiple addresses and ensure they're unique
	addresses := make(map[string]bool)
	
	for i := 0; i < 100; i++ {
		publicKey := make([]byte, 32)
		copy(publicKey, []byte{byte(i)})
		hash := sha256.Sum256(publicKey)
		
		address, err := GenerateEthereumStyleAddress(hash[:])
		if err != nil {
			t.Fatalf("Failed to generate address: %v", err)
		}
		
		if addresses[address] {
			t.Errorf("Duplicate address generated: %s", address)
		}
		addresses[address] = true
	}
}

func TestBitcoinStyleAddressDeterministic(t *testing.T) {
	// Same public key should always generate the same address
	publicKey := []byte("consistent_test_key")
	
	address1, err := GenerateBitcoinStyleAddress(publicKey)
	if err != nil {
		t.Fatalf("Failed to generate address 1: %v", err)
	}
	
	address2, err := GenerateBitcoinStyleAddress(publicKey)
	if err != nil {
		t.Fatalf("Failed to generate address 2: %v", err)
	}
	
	if address1 != address2 {
		t.Errorf("Addresses should be deterministic. Got %s and %s", address1, address2)
	}
}

func TestEthereumStyleAddressDeterministic(t *testing.T) {
	// Same public key should always generate the same address
	publicKey := []byte("consistent_test_key")
	
	address1, err := GenerateEthereumStyleAddress(publicKey)
	if err != nil {
		t.Fatalf("Failed to generate address 1: %v", err)
	}
	
	address2, err := GenerateEthereumStyleAddress(publicKey)
	if err != nil {
		t.Fatalf("Failed to generate address 2: %v", err)
	}
	
	if address1 != address2 {
		t.Errorf("Addresses should be deterministic. Got %s and %s", address1, address2)
	}
}
