package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// AddressType represents the type of blockchain address
type AddressType string

const (
	// BitcoinStyle represents Bitcoin-style addresses with Btu prefix
	BitcoinStyle AddressType = "bitcoin"
	// EthereumStyle represents Ethereum-style addresses with 0x prefix
	EthereumStyle AddressType = "ethereum"
)

// Base58 alphabet (Bitcoin-compatible)
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// GenerateBitcoinStyleAddress generates a Bitcoin-style address with "Btu" prefix
// Format: Btu + Base58Check(publicKeyHash + checksum)
// Example: BtuzP1eP5QGefi2DMPTfTL5SLmv7DivfNa
func GenerateBitcoinStyleAddress(publicKey []byte) (string, error) {
	if len(publicKey) == 0 {
		return "", errors.New("public key cannot be empty")
	}

	// Hash the public key with SHA-256
	hash := sha256.Sum256(publicKey)
	
	// Take first 20 bytes of the hash
	payload := hash[:20]
	
	// Calculate checksum (first 4 bytes of double SHA-256)
	checksum := calculateChecksum(payload)
	
	// Combine payload and checksum
	fullPayload := append(payload, checksum...)
	
	// Encode with Base58
	encoded := base58Encode(fullPayload)
	
	// Add "Btu" prefix
	address := "Btu" + encoded
	
	return address, nil
}

// GenerateEthereumStyleAddress generates an Ethereum-style address
// Format: 0x + 40 hexadecimal characters
// Example: 0x1234567890abcdef1234567890abcdef12345678
// 
// NOTE: This implementation uses SHA-256 instead of Keccak-256 for simplicity.
// For production use with actual Ethereum network, implement Keccak-256 hashing.
// This format is compatible with EVM-based blockchains that accept this addressing scheme.
func GenerateEthereumStyleAddress(publicKey []byte) (string, error) {
	if len(publicKey) == 0 {
		return "", errors.New("public key cannot be empty")
	}

	// Hash the public key with SHA-256
	// NOTE: Ethereum uses Keccak-256, but we use SHA-256 for Bituncoin's simplified implementation
	hash := sha256.Sum256(publicKey)
	
	// Take last 20 bytes (similar to Ethereum's approach)
	addressBytes := hash[12:]
	
	// Convert to hex and add 0x prefix
	address := "0x" + hex.EncodeToString(addressBytes)
	
	return address, nil
}

// ValidateBitcoinStyleAddress validates a Bitcoin-style address
func ValidateBitcoinStyleAddress(address string) error {
	// Check if address starts with "Btu"
	if !strings.HasPrefix(address, "Btu") {
		return errors.New("invalid Bitcoin-style address: must start with 'Btu'")
	}
	
	// Check minimum length (Btu + at least 26 base58 characters)
	if len(address) < 29 {
		return errors.New("invalid Bitcoin-style address: too short")
	}
	
	// Check maximum length
	if len(address) > 40 {
		return errors.New("invalid Bitcoin-style address: too long")
	}
	
	// Extract the base58 part
	base58Part := address[3:]
	
	// Validate base58 characters
	for _, char := range base58Part {
		if !strings.ContainsRune(base58Alphabet, char) {
			return fmt.Errorf("invalid Bitcoin-style address: invalid character '%c'", char)
		}
	}
	
	// Decode base58
	decoded, err := base58Decode(base58Part)
	if err != nil {
		return fmt.Errorf("invalid Bitcoin-style address: %v", err)
	}
	
	// Check decoded length (at least 24 bytes: 20 payload + 4 checksum)
	if len(decoded) < 24 {
		return errors.New("invalid Bitcoin-style address: decoded data too short")
	}
	
	// Extract payload and checksum
	payload := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]
	
	// Verify checksum
	expectedChecksum := calculateChecksum(payload)
	for i := 0; i < 4; i++ {
		if checksum[i] != expectedChecksum[i] {
			return errors.New("invalid Bitcoin-style address: checksum mismatch")
		}
	}
	
	return nil
}

// ValidateEthereumStyleAddress validates an Ethereum-style address
func ValidateEthereumStyleAddress(address string) error {
	// Check if address starts with "0x"
	if !strings.HasPrefix(address, "0x") {
		return errors.New("invalid Ethereum-style address: must start with '0x'")
	}
	
	// Check length (0x + 40 hex characters = 42 total)
	if len(address) != 42 {
		return fmt.Errorf("invalid Ethereum-style address: must be 42 characters, got %d", len(address))
	}
	
	// Extract hex part
	hexPart := address[2:]
	
	// Validate hex characters
	matched, err := regexp.MatchString("^[0-9a-fA-F]{40}$", hexPart)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid Ethereum-style address: must contain only hexadecimal characters")
	}
	
	return nil
}

// calculateChecksum calculates a 4-byte checksum using double SHA-256
func calculateChecksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:4]
}

// base58Encode encodes bytes to base58 string
func base58Encode(input []byte) string {
	// Handle empty input
	if len(input) == 0 {
		return ""
	}
	
	// Count leading zeros
	leadingZeros := 0
	for _, b := range input {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}
	
	// Convert bytes to a big integer (represented as []byte)
	// We'll use a simple implementation with byte slices
	result := make([]byte, 0, len(input)*2)
	
	for _, b := range input {
		carry := int(b)
		for i := len(result) - 1; i >= 0; i-- {
			carry += int(result[i]) * 256
			result[i] = byte(carry % 58)
			carry /= 58
		}
		for carry > 0 {
			result = append([]byte{byte(carry % 58)}, result...)
			carry /= 58
		}
	}
	
	// Convert to base58 string
	var encoded strings.Builder
	
	// Add leading '1's for leading zero bytes
	for i := 0; i < leadingZeros; i++ {
		encoded.WriteByte(base58Alphabet[0])
	}
	
	// Add the rest
	for _, val := range result {
		encoded.WriteByte(base58Alphabet[val])
	}
	
	return encoded.String()
}

// base58Decode decodes a base58 string to bytes
func base58Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return []byte{}, nil
	}
	
	// Count leading '1's
	leadingOnes := 0
	for _, char := range input {
		if char == rune(base58Alphabet[0]) {
			leadingOnes++
		} else {
			break
		}
	}
	
	// Convert base58 string to bytes
	result := make([]byte, 0)
	
	for _, char := range input {
		pos := strings.IndexRune(base58Alphabet, char)
		if pos == -1 {
			return nil, fmt.Errorf("invalid base58 character: %c", char)
		}
		
		carry := pos
		for i := len(result) - 1; i >= 0; i-- {
			carry += int(result[i]) * 58
			result[i] = byte(carry % 256)
			carry /= 256
		}
		for carry > 0 {
			result = append([]byte{byte(carry % 256)}, result...)
			carry /= 256
		}
	}
	
	// Add leading zero bytes for leading '1's
	for i := 0; i < leadingOnes; i++ {
		result = append([]byte{0}, result...)
	}
	
	return result, nil
}

// GenerateAddressWithType generates an address of the specified type
// NOTE: This is a simplified implementation suitable for Bituncoin's blockchain.
// For production use with Bitcoin or Ethereum networks:
// - Use secp256k1 elliptic curve for proper public key derivation
// - Use Keccak-256 for Ethereum addresses instead of SHA-256
func GenerateAddressWithType(addressType AddressType) (string, []byte, error) {
	// Generate random private key (32 bytes)
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)
	if err != nil {
		return "", nil, err
	}
	
	// Derive public key from private key (simplified - using hash)
	// PRODUCTION NOTE: Use proper elliptic curve point multiplication
	publicKeyHash := sha256.Sum256(privateKey)
	publicKey := publicKeyHash[:]
	
	var address string
	switch addressType {
	case BitcoinStyle:
		address, err = GenerateBitcoinStyleAddress(publicKey)
	case EthereumStyle:
		address, err = GenerateEthereumStyleAddress(publicKey)
	default:
		return "", nil, fmt.Errorf("unsupported address type: %s", addressType)
	}
	
	if err != nil {
		return "", nil, err
	}
	
	return address, privateKey, nil
}

// ValidateAddressType validates an address and returns its type
func ValidateAddressType(address string) (AddressType, error) {
	// Try Bitcoin-style validation
	if err := ValidateBitcoinStyleAddress(address); err == nil {
		return BitcoinStyle, nil
	}
	
	// Try Ethereum-style validation
	if err := ValidateEthereumStyleAddress(address); err == nil {
		return EthereumStyle, nil
	}
	
	return "", errors.New("unknown address type: address does not match Bitcoin-style or Ethereum-style format")
}
