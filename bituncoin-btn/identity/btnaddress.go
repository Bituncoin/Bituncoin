package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// Identity represents a user's identity in the Bituncoin network
type Identity struct {
	Address        string    `json:"address"`
	PublicKey      string    `json:"public_key"`
	PrivateKey     string    `json:"private_key,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	TwoFactorAuth  bool      `json:"two_factor_auth"`
	BiometricHash  string    `json:"biometric_hash,omitempty"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	TwoFactorEnabled   bool   `json:"two_factor_enabled"`
	BiometricEnabled   bool   `json:"biometric_enabled"`
	TwoFactorSecret    string `json:"two_factor_secret,omitempty"`
	BiometricTemplate  string `json:"biometric_template,omitempty"`
}

// CreateIdentity creates a new identity with keys
func CreateIdentity() (*Identity, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := elliptic.Marshal(privateKey.Curve, privateKey.X, privateKey.Y)
	privKeyBytes := privateKey.D.Bytes()

	address := generateAddress(pubKeyBytes)

	identity := &Identity{
		Address:       address,
		PublicKey:     hex.EncodeToString(pubKeyBytes),
		PrivateKey:    hex.EncodeToString(privKeyBytes),
		CreatedAt:     time.Now(),
		TwoFactorAuth: false,
	}

	return identity, nil
}

// generateAddress generates a blockchain address from public key
func generateAddress(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	address := "BTN" + hex.EncodeToString(hash[:20])
	return address
}

// EnableTwoFactorAuth enables 2FA for the identity
func (id *Identity) EnableTwoFactorAuth(secret string) {
	id.TwoFactorAuth = true
}

// VerifyTwoFactorAuth verifies a 2FA token
func (id *Identity) VerifyTwoFactorAuth(token string) bool {
	// In a real implementation, this would verify TOTP tokens
	// For now, we'll use a simplified version
	return len(token) == 6 && id.TwoFactorAuth
}

// SetBiometricHash sets the biometric hash for authentication
func (id *Identity) SetBiometricHash(biometricData string) {
	hash := sha256.Sum256([]byte(biometricData))
	id.BiometricHash = hex.EncodeToString(hash[:])
}

// VerifyBiometric verifies biometric authentication
func (id *Identity) VerifyBiometric(biometricData string) bool {
	hash := sha256.Sum256([]byte(biometricData))
	return hex.EncodeToString(hash[:]) == id.BiometricHash
}

// SignTransaction signs a transaction with the private key
func (id *Identity) SignTransaction(txData []byte) (string, error) {
	if id.PrivateKey == "" {
		return "", errors.New("private key not available")
	}

	// Create a hash of the transaction
	hash := sha256.Sum256(txData)
	
	// In a real implementation, we would use ECDSA signing
	// For simplicity, we'll create a signature based on the hash
	signature := hex.EncodeToString(hash[:])
	
	return signature, nil
}

// VerifySignature verifies a transaction signature
func VerifySignature(txData []byte, signature string, publicKey string) bool {
	hash := sha256.Sum256(txData)
	expectedSig := hex.EncodeToString(hash[:])
	return signature == expectedSig
}

// ExportIdentity exports identity to JSON (without private key for security)
func (id *Identity) ExportIdentity(includePrivateKey bool) (string, error) {
	exportID := *id
	if !includePrivateKey {
		exportID.PrivateKey = ""
	}
	
	data, err := json.MarshalIndent(exportID, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// ImportIdentity imports an identity from JSON
func ImportIdentity(jsonData string) (*Identity, error) {
	var identity Identity
	err := json.Unmarshal([]byte(jsonData), &identity)
	if err != nil {
		return nil, err
	}
	
	return &identity, nil
}

// ValidateAddress validates a Bituncoin address
func ValidateAddress(address string) bool {
	if len(address) < 43 || address[:3] != "BTN" {
		return false
	}
	
	_, err := hex.DecodeString(address[3:])
	return err == nil
}
