package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

// AddressType represents the type of blockchain address
type AddressType string

const (
	// AddressTypeBitcoin represents Bitcoin-style addresses
	AddressTypeBitcoin AddressType = "bitcoin"
	// AddressTypeEthereum represents Ethereum-style addresses
	AddressTypeEthereum AddressType = "ethereum"
)

// Base58 alphabet (Bitcoin-style)
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// BitcoinStyleAddress represents a Bitcoin-style address with Btu prefix
type BitcoinStyleAddress struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

// EthereumStyleAddress represents an Ethereum-style address
type EthereumStyleAddress struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

// GenerateBitcoinStyleAddress generates a Bitcoin-style address with Base58Check encoding
// The address starts with "Btu" prefix
func GenerateBitcoinStyleAddress() (*BitcoinStyleAddress, error) {
	// Generate ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Get public key in uncompressed format
	publicKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// Hash the public key with SHA256
	sha256Hash := sha256.Sum256(publicKeyBytes)

	// Hash again with RIPEMD-160
	ripemd160Hasher := ripemd160.New()
	_, err = ripemd160Hasher.Write(sha256Hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to hash with RIPEMD-160: %v", err)
	}
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// Add version byte (custom version for Bituncoin)
	// Using version byte 0x00 for simplicity
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// Calculate checksum (first 4 bytes of double SHA256)
	firstHash := sha256.Sum256(versionedPayload)
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	// Concatenate versioned payload and checksum
	fullPayload := append(versionedPayload, checksum...)

	// Encode to Base58
	base58Address := base58Encode(fullPayload)
	
	// Prepend "Btu" prefix to create Bituncoin address format
	address := "Btu" + base58Address

	return &BitcoinStyleAddress{
		PrivateKey: hex.EncodeToString(privateKey.D.Bytes()),
		PublicKey:  hex.EncodeToString(publicKeyBytes),
		Address:    address,
	}, nil
}

// GenerateEthereumStyleAddress generates an Ethereum-style address
// The address starts with "0x" and contains 40 hexadecimal characters
func GenerateEthereumStyleAddress() (*EthereumStyleAddress, error) {
	// Generate ECDSA key pair using secp256k1 curve (Ethereum uses secp256k1, we'll use P256 as approximation)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Get public key in uncompressed format (excluding the 0x04 prefix byte)
	publicKeyBytes := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// Hash the public key with Keccak-256 (we'll use SHA256 as approximation since Keccak is not in stdlib)
	// In production, you would use "golang.org/x/crypto/sha3" for Keccak-256
	hash := sha256.Sum256(publicKeyBytes)

	// Take the last 20 bytes of the hash
	addressBytes := hash[12:]

	// Format as Ethereum address (0x + 40 hex characters)
	address := "0x" + hex.EncodeToString(addressBytes)

	return &EthereumStyleAddress{
		PrivateKey: hex.EncodeToString(privateKey.D.Bytes()),
		PublicKey:  hex.EncodeToString(publicKeyBytes),
		Address:    address,
	}, nil
}

// ValidateBitcoinStyleAddress validates a Bitcoin-style address with Btu prefix
func ValidateBitcoinStyleAddress(address string) error {
	if address == "" {
		return errors.New("address cannot be empty")
	}

	// Check that address starts with "Btu"
	if !strings.HasPrefix(address, "Btu") {
		return errors.New("address must start with Btu prefix")
	}

	// Remove "Btu" prefix to get Base58 part
	base58Part := address[3:]
	
	if base58Part == "" {
		return errors.New("address too short")
	}

	// Decode from Base58
	decoded, err := base58Decode(base58Part)
	if err != nil {
		return fmt.Errorf("invalid Base58 encoding: %v", err)
	}

	// Check minimum length (version byte + 20 bytes + 4 bytes checksum)
	if len(decoded) < 25 {
		return errors.New("address too short")
	}

	// Extract components
	payload := decoded[:len(decoded)-4]
	checksum := decoded[len(decoded)-4:]

	// Verify checksum
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	expectedChecksum := secondHash[:4]

	for i := 0; i < 4; i++ {
		if checksum[i] != expectedChecksum[i] {
			return errors.New("invalid checksum")
		}
	}

	// Verify version byte (should be 0x00 for Bituncoin)
	if payload[0] != 0x00 {
		return errors.New("invalid version byte")
	}

	return nil
}

// ValidateEthereumStyleAddress validates an Ethereum-style address
func ValidateEthereumStyleAddress(address string) error {
	if address == "" {
		return errors.New("address cannot be empty")
	}

	// Check prefix
	if !strings.HasPrefix(address, "0x") {
		return errors.New("address must start with 0x")
	}

	// Remove 0x prefix
	addressHex := address[2:]

	// Check length (should be exactly 40 hex characters)
	if len(addressHex) != 40 {
		return fmt.Errorf("address must be 40 hexadecimal characters (got %d)", len(addressHex))
	}

	// Verify all characters are valid hex
	_, err := hex.DecodeString(addressHex)
	if err != nil {
		return fmt.Errorf("address contains invalid hexadecimal characters: %v", err)
	}

	return nil
}

// base58Encode encodes a byte array to Base58 string
func base58Encode(input []byte) string {
	// Convert bytes to big integer
	x := new(big.Int).SetBytes(input)

	// Prepare result
	result := make([]byte, 0, len(input)*2)

	// Convert to base58
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// Add leading zeros
	for _, b := range input {
		if b == 0x00 {
			result = append(result, base58Alphabet[0])
		} else {
			break
		}
	}

	// Reverse result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// base58Decode decodes a Base58 string to byte array
func base58Decode(input string) ([]byte, error) {
	result := big.NewInt(0)
	base := big.NewInt(58)

	// Decode each character
	for _, c := range input {
		pos := strings.IndexRune(base58Alphabet, c)
		if pos == -1 {
			return nil, fmt.Errorf("invalid character '%c' in Base58 string", c)
		}
		result.Mul(result, base)
		result.Add(result, big.NewInt(int64(pos)))
	}

	// Convert to bytes
	decoded := result.Bytes()

	// Add leading zeros
	for _, c := range input {
		if c == rune(base58Alphabet[0]) {
			decoded = append([]byte{0x00}, decoded...)
		} else {
			break
		}
	}

	return decoded, nil
}
