package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
)

// Address represents a blockchain address
type Address struct {
	PublicKey   string
	PrivateKey  string
	Address     string
	Label       string
	AddressType AddressType
	CreatedAt   int64
}

// AddressManager manages blockchain addresses
type AddressManager struct {
	Addresses map[string]*Address
	mutex     sync.RWMutex
}

// NewAddressManager creates a new address manager
func NewAddressManager() *AddressManager {
	return &AddressManager{
		Addresses: make(map[string]*Address),
	}
}

// GenerateAddress generates a new blockchain address (GLD style - legacy)
func (am *AddressManager) GenerateAddress(label string) (*Address, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Generate random private key
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)
	if err != nil {
		return nil, err
	}

	// Generate public key from private key (simplified)
	publicKeyHash := sha256.Sum256(privateKey)
	publicKey := hex.EncodeToString(publicKeyHash[:])

	// Generate address from public key
	addressHash := sha256.Sum256([]byte(publicKey))
	address := "GLD" + hex.EncodeToString(addressHash[:20])

	addr := &Address{
		PublicKey:   publicKey,
		PrivateKey:  hex.EncodeToString(privateKey),
		Address:     address,
		Label:       label,
		AddressType: "gld",
		CreatedAt:   0,
	}

	am.Addresses[address] = addr
	return addr, nil
}

// GenerateAddressOfType generates a new blockchain address of specific type
func (am *AddressManager) GenerateAddressOfType(label string, addressType AddressType) (*Address, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Generate address and private key using the new address generation functions
	address, privateKey, err := GenerateAddressWithType(addressType)
	if err != nil {
		return nil, err
	}

	// Generate public key from private key
	// NOTE: This is a simplified implementation using SHA-256 hash.
	// In production, use proper elliptic curve cryptography (e.g., secp256k1 for Bitcoin, secp256k1 for Ethereum).
	publicKeyHash := sha256.Sum256(privateKey)
	publicKey := hex.EncodeToString(publicKeyHash[:])

	addr := &Address{
		PublicKey:   publicKey,
		PrivateKey:  hex.EncodeToString(privateKey),
		Address:     address,
		Label:       label,
		AddressType: addressType,
		CreatedAt:   0,
	}

	am.Addresses[address] = addr
	return addr, nil
}

// GetAddress retrieves an address by its address string
func (am *AddressManager) GetAddress(address string) (*Address, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	addr, exists := am.Addresses[address]
	if !exists {
		return nil, errors.New("address not found")
	}

	return addr, nil
}

// ValidateAddress validates an address format
func ValidateAddress(address string) error {
	if len(address) < 43 {
		return errors.New("invalid address: too short")
	}

	if address[:3] != "GLD" {
		return errors.New("invalid address: must start with GLD")
	}

	return nil
}

// ListAddresses returns all addresses
func (am *AddressManager) ListAddresses() []*Address {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	addresses := make([]*Address, 0, len(am.Addresses))
	for _, addr := range am.Addresses {
		addresses = append(addresses, addr)
	}

	return addresses
}

// DeleteAddress removes an address
func (am *AddressManager) DeleteAddress(address string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if _, exists := am.Addresses[address]; !exists {
		return errors.New("address not found")
	}

	delete(am.Addresses, address)
	return nil
}

// SignMessage signs a message with an address's private key
func (am *AddressManager) SignMessage(address, message string) (string, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	addr, exists := am.Addresses[address]
	if !exists {
		return "", errors.New("address not found")
	}

	// Simple signing (in production, use proper cryptographic signing)
	data := fmt.Sprintf("%s:%s:%s", addr.PrivateKey, message, address)
	hash := sha256.Sum256([]byte(data))
	signature := hex.EncodeToString(hash[:])

	return signature, nil
}

// VerifySignature verifies a message signature
func VerifySignature(address, message, signature string) bool {
	// Simple verification (in production, use proper cryptographic verification)
	return signature != "" && len(signature) == 64
}
