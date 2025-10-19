package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"
)

// Security manages wallet security features
type Security struct {
	TwoFactorEnabled   bool
	BiometricEnabled   bool
	EncryptionKey      []byte
	BackupEncrypted    bool
	LastBackupTime     int64
	mutex              sync.RWMutex
}

// TwoFactorAuth represents 2FA settings
type TwoFactorAuth struct {
	Secret    string
	Enabled   bool
	BackupCodes []string
}

// BiometricAuth represents biometric settings
type BiometricAuth struct {
	Enabled      bool
	Type         string // fingerprint, face, etc.
	LastUsed     int64
}

// NewSecurity creates a new security manager
func NewSecurity() *Security {
	return &Security{
		TwoFactorEnabled: false,
		BiometricEnabled: false,
		EncryptionKey:    generateEncryptionKey(),
		BackupEncrypted:  true,
		LastBackupTime:   0,
	}
}

// generateEncryptionKey generates a random encryption key
func generateEncryptionKey() []byte {
	key := make([]byte, 32) // AES-256
	rand.Read(key)
	return key
}

// EnableTwoFactor enables two-factor authentication
func (s *Security) EnableTwoFactor(secret string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if secret == "" {
		return errors.New("secret cannot be empty")
	}

	s.TwoFactorEnabled = true
	return nil
}

// DisableTwoFactor disables two-factor authentication
func (s *Security) DisableTwoFactor() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.TwoFactorEnabled = false
	return nil
}

// EnableBiometric enables biometric authentication
func (s *Security) EnableBiometric(biometricType string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if biometricType == "" {
		return errors.New("biometric type cannot be empty")
	}

	s.BiometricEnabled = true
	return nil
}

// DisableBiometric disables biometric authentication
func (s *Security) DisableBiometric() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.BiometricEnabled = false
	return nil
}

// Encrypt encrypts data using AES-256
func (s *Security) Encrypt(plaintext []byte) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	block, err := aes.NewCipher(s.EncryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts AES-256 encrypted data
func (s *Security) Decrypt(ciphertext string) ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(s.EncryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// CreateBackup creates an encrypted wallet backup
func (s *Security) CreateBackup(walletData []byte) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	encrypted, err := s.Encrypt(walletData)
	if err != nil {
		return "", err
	}

	s.LastBackupTime = time.Now().Unix()
	return encrypted, nil
}

// RestoreBackup restores wallet from encrypted backup
func (s *Security) RestoreBackup(encryptedBackup string) ([]byte, error) {
	return s.Decrypt(encryptedBackup)
}

// HashPassword hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// GenerateRecoveryPhrase generates a recovery phrase
func GenerateRecoveryPhrase() ([]string, error) {
	words := []string{
		"abandon", "ability", "able", "about", "above", "absent",
		"absorb", "abstract", "absurd", "abuse", "access", "accident",
	}

	phrase := make([]string, 12)
	for i := 0; i < 12; i++ {
		randomIndex := make([]byte, 1)
		if _, err := rand.Read(randomIndex); err != nil {
			return nil, err
		}
		phrase[i] = words[int(randomIndex[0])%len(words)]
	}

	return phrase, nil
}

// GetSecurityStatus returns current security settings
func (s *Security) GetSecurityStatus() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"twoFactorEnabled": s.TwoFactorEnabled,
		"biometricEnabled": s.BiometricEnabled,
		"encryptionType":   "AES-256",
		"backupEncrypted":  s.BackupEncrypted,
		"lastBackup":       s.LastBackupTime,
	}
}

// FraudDetector manages fraud detection
type FraudDetector struct {
	suspiciousActivities map[string][]string // address -> activities
	blockedAddresses     map[string]bool
	alertThreshold       int
	mutex                sync.RWMutex
}

// NewFraudDetector creates a new fraud detector
func NewFraudDetector() *FraudDetector {
	return &FraudDetector{
		suspiciousActivities: make(map[string][]string),
		blockedAddresses:     make(map[string]bool),
		alertThreshold:       3,
	}
}

// CheckTransaction checks if a transaction is suspicious
func (fd *FraudDetector) CheckTransaction(from, to string, amount float64) (bool, string) {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	// Check if address is blocked
	if fd.blockedAddresses[from] || fd.blockedAddresses[to] {
		return true, "Transaction involves blocked address"
	}

	// Check for unusually large transaction
	if amount > 10000 {
		fd.logSuspiciousActivity(from, "Large transaction amount")
		return true, "Unusually large transaction amount"
	}

	// Check for rapid transactions (simplified)
	activities := fd.suspiciousActivities[from]
	if len(activities) > 5 {
		return true, "High frequency of transactions detected"
	}

	return false, ""
}

// logSuspiciousActivity logs suspicious activity
func (fd *FraudDetector) logSuspiciousActivity(address, activity string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := timestamp + " - " + activity
	fd.suspiciousActivities[address] = append(fd.suspiciousActivities[address], entry)

	// Auto-block if threshold exceeded
	if len(fd.suspiciousActivities[address]) >= fd.alertThreshold {
		fd.blockedAddresses[address] = true
	}
}

// BlockAddress manually blocks an address
func (fd *FraudDetector) BlockAddress(address string) {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	fd.blockedAddresses[address] = true
}

// UnblockAddress unblocks an address
func (fd *FraudDetector) UnblockAddress(address string) {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	delete(fd.blockedAddresses, address)
}

// GetSuspiciousActivities returns suspicious activities for an address
func (fd *FraudDetector) GetSuspiciousActivities(address string) []string {
	fd.mutex.RLock()
	defer fd.mutex.RUnlock()

	return fd.suspiciousActivities[address]
}

// IsAddressBlocked checks if an address is blocked
func (fd *FraudDetector) IsAddressBlocked(address string) bool {
	fd.mutex.RLock()
	defer fd.mutex.RUnlock()

	return fd.blockedAddresses[address]
}

// AlertSystem manages real-time security alerts
type AlertSystem struct {
	alerts      []SecurityAlert
	subscribers map[string]chan SecurityAlert // address -> alert channel
	mutex       sync.RWMutex
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // fraud, suspicious, security
	Severity  string    `json:"severity"` // low, medium, high, critical
	Message   string    `json:"message"`
	Address   string    `json:"address"`
	Timestamp time.Time `json:"timestamp"`
	Resolved  bool      `json:"resolved"`
}

// NewAlertSystem creates a new alert system
func NewAlertSystem() *AlertSystem {
	return &AlertSystem{
		alerts:      make([]SecurityAlert, 0),
		subscribers: make(map[string]chan SecurityAlert),
	}
}

// SendAlert sends a security alert
func (as *AlertSystem) SendAlert(alertType, severity, message, address string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	alert := SecurityAlert{
		ID:        generateAlertID(),
		Type:      alertType,
		Severity:  severity,
		Message:   message,
		Address:   address,
		Timestamp: time.Now(),
		Resolved:  false,
	}

	as.alerts = append(as.alerts, alert)

	// Send to subscribers
	if ch, exists := as.subscribers[address]; exists {
		select {
		case ch <- alert:
		default:
			// Channel full, skip
		}
	}
}

// Subscribe subscribes to alerts for an address
func (as *AlertSystem) Subscribe(address string) chan SecurityAlert {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	ch := make(chan SecurityAlert, 100)
	as.subscribers[address] = ch
	return ch
}

// Unsubscribe unsubscribes from alerts
func (as *AlertSystem) Unsubscribe(address string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	if ch, exists := as.subscribers[address]; exists {
		close(ch)
		delete(as.subscribers, address)
	}
}

// GetAlerts returns alerts for an address
func (as *AlertSystem) GetAlerts(address string) []SecurityAlert {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	alerts := make([]SecurityAlert, 0)
	for _, alert := range as.alerts {
		if alert.Address == address && !alert.Resolved {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// ResolveAlert marks an alert as resolved
func (as *AlertSystem) ResolveAlert(alertID string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	for i := range as.alerts {
		if as.alerts[i].ID == alertID {
			as.alerts[i].Resolved = true
			break
		}
	}
}

func generateAlertID() string {
	return "ALERT-" + time.Now().Format("20060102150405")
}
