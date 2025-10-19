package wallet

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// FraudDetectionRule represents a fraud detection rule
type FraudDetectionRule struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
	Enabled     bool    `json:"enabled"`
	Threshold   float64 `json:"threshold,omitempty"`
}

// SecurityEvent represents a security event
type SecurityEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	WalletID    string    `json:"walletId"`
	IPAddress   string    `json:"ipAddress,omitempty"`
	DeviceID    string    `json:"deviceId,omitempty"`
	Timestamp   int64     `json:"timestamp"`
	Resolved    bool      `json:"resolved"`
	ResolvedAt  int64     `json:"resolvedAt,omitempty"`
}

// DeviceFingerprint represents a device fingerprint
type DeviceFingerprint struct {
	ID            string            `json:"id"`
	WalletID      string            `json:"walletId"`
	DeviceType    string            `json:"deviceType"`
	OS            string            `json:"os"`
	Browser       string            `json:"browser"`
	IPAddress     string            `json:"ipAddress"`
	Location      string            `json:"location,omitempty"`
	FirstSeen     int64             `json:"firstSeen"`
	LastSeen      int64             `json:"lastSeen"`
	IsTrusted     bool              `json:"isTrusted"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// TransactionApproval represents a transaction approval workflow
type TransactionApproval struct {
	ID                string    `json:"id"`
	WalletID          string    `json:"walletId"`
	TransactionID     string    `json:"transactionId"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	RequiresApproval  bool      `json:"requiresApproval"`
	ApprovalMethod    string    `json:"approvalMethod"` // 2fa, biometric, manual
	Status            string    `json:"status"`         // pending, approved, rejected
	CreatedAt         int64     `json:"createdAt"`
	ApprovedAt        int64     `json:"approvedAt,omitempty"`
	ApprovalCode      string    `json:"approvalCode,omitempty"`
}

// AdvancedSecurity manages advanced security features
type AdvancedSecurity struct {
	fraudRules        map[string]*FraudDetectionRule
	securityEvents    map[string]*SecurityEvent
	deviceFingerprints map[string]*DeviceFingerprint
	approvals         map[string]*TransactionApproval
	blockedIPs        map[string]bool
	blockedDevices    map[string]bool
	mutex             sync.RWMutex
}

// NewAdvancedSecurity creates a new advanced security manager
func NewAdvancedSecurity() *AdvancedSecurity {
	as := &AdvancedSecurity{
		fraudRules:        make(map[string]*FraudDetectionRule),
		securityEvents:    make(map[string]*SecurityEvent),
		deviceFingerprints: make(map[string]*DeviceFingerprint),
		approvals:         make(map[string]*TransactionApproval),
		blockedIPs:        make(map[string]bool),
		blockedDevices:    make(map[string]bool),
	}

	// Initialize default fraud detection rules
	as.initializeDefaultRules()

	return as
}

// initializeDefaultRules initializes default fraud detection rules
func (as *AdvancedSecurity) initializeDefaultRules() {
	rules := []*FraudDetectionRule{
		{
			ID:          "rule_high_value_tx",
			Name:        "High Value Transaction",
			Description: "Flag transactions above $10,000",
			Severity:    "medium",
			Enabled:     true,
			Threshold:   10000.0,
		},
		{
			ID:          "rule_rapid_tx",
			Name:        "Rapid Transactions",
			Description: "Flag more than 5 transactions in 1 minute",
			Severity:    "high",
			Enabled:     true,
			Threshold:   5.0,
		},
		{
			ID:          "rule_new_device",
			Name:        "New Device Login",
			Description: "Flag login from unrecognized device",
			Severity:    "medium",
			Enabled:     true,
		},
		{
			ID:          "rule_unusual_location",
			Name:        "Unusual Location",
			Description: "Flag login from unusual geographic location",
			Severity:    "high",
			Enabled:     true,
		},
		{
			ID:          "rule_withdrawal_spike",
			Name:        "Withdrawal Spike",
			Description: "Flag unusual withdrawal patterns",
			Severity:    "high",
			Enabled:     true,
		},
	}

	for _, rule := range rules {
		as.fraudRules[rule.ID] = rule
	}
}

// DetectFraud analyzes a transaction for fraud indicators
func (as *AdvancedSecurity) DetectFraud(walletID string, txAmount float64, currency, deviceID, ipAddress string) (bool, []string, error) {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	fraudIndicators := make([]string, 0)
	isFraud := false

	// Check if IP is blocked
	if as.blockedIPs[ipAddress] {
		fraudIndicators = append(fraudIndicators, "Blocked IP address")
		isFraud = true
	}

	// Check if device is blocked
	if as.blockedDevices[deviceID] {
		fraudIndicators = append(fraudIndicators, "Blocked device")
		isFraud = true
	}

	// Check high value transaction rule
	if rule, exists := as.fraudRules["rule_high_value_tx"]; exists && rule.Enabled {
		if txAmount > rule.Threshold {
			fraudIndicators = append(fraudIndicators, fmt.Sprintf("High value transaction: $%.2f", txAmount))
			if !isFraud {
				// Don't mark as fraud, just flag for approval
			}
		}
	}

	// Check if device is new/untrusted
	deviceTrusted := false
	for _, fp := range as.deviceFingerprints {
		if fp.ID == deviceID && fp.WalletID == walletID && fp.IsTrusted {
			deviceTrusted = true
			break
		}
	}

	if !deviceTrusted && as.fraudRules["rule_new_device"].Enabled {
		fraudIndicators = append(fraudIndicators, "Unrecognized device")
	}

	return isFraud, fraudIndicators, nil
}

// RegisterDevice registers a device fingerprint
func (as *AdvancedSecurity) RegisterDevice(walletID, deviceType, os, browser, ipAddress, location string) (*DeviceFingerprint, error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	deviceID := generateDeviceID()

	fp := &DeviceFingerprint{
		ID:         deviceID,
		WalletID:   walletID,
		DeviceType: deviceType,
		OS:         os,
		Browser:    browser,
		IPAddress:  ipAddress,
		Location:   location,
		FirstSeen:  time.Now().Unix(),
		LastSeen:   time.Now().Unix(),
		IsTrusted:  false,
		Metadata:   make(map[string]string),
	}

	as.deviceFingerprints[deviceID] = fp

	return fp, nil
}

// TrustDevice marks a device as trusted
func (as *AdvancedSecurity) TrustDevice(deviceID string) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	fp, exists := as.deviceFingerprints[deviceID]
	if !exists {
		return errors.New("device not found")
	}

	fp.IsTrusted = true

	return nil
}

// BlockIP blocks an IP address
func (as *AdvancedSecurity) BlockIP(ipAddress string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	as.blockedIPs[ipAddress] = true
}

// UnblockIP unblocks an IP address
func (as *AdvancedSecurity) UnblockIP(ipAddress string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	delete(as.blockedIPs, ipAddress)
}

// BlockDevice blocks a device
func (as *AdvancedSecurity) BlockDevice(deviceID string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	as.blockedDevices[deviceID] = true
}

// UnblockDevice unblocks a device
func (as *AdvancedSecurity) UnblockDevice(deviceID string) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	delete(as.blockedDevices, deviceID)
}

// CreateSecurityEvent creates a security event
func (as *AdvancedSecurity) CreateSecurityEvent(walletID, eventType, severity, description, ipAddress, deviceID string) *SecurityEvent {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	event := &SecurityEvent{
		ID:          generateSecurityEventID(),
		Type:        eventType,
		Severity:    severity,
		Description: description,
		WalletID:    walletID,
		IPAddress:   ipAddress,
		DeviceID:    deviceID,
		Timestamp:   time.Now().Unix(),
		Resolved:    false,
	}

	as.securityEvents[event.ID] = event

	return event
}

// GetSecurityEvents returns security events for a wallet
func (as *AdvancedSecurity) GetSecurityEvents(walletID string, unresolvedOnly bool) []*SecurityEvent {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	events := make([]*SecurityEvent, 0)
	for _, event := range as.securityEvents {
		if event.WalletID == walletID {
			if !unresolvedOnly || !event.Resolved {
				events = append(events, event)
			}
		}
	}

	return events
}

// ResolveSecurityEvent resolves a security event
func (as *AdvancedSecurity) ResolveSecurityEvent(eventID string) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	event, exists := as.securityEvents[eventID]
	if !exists {
		return errors.New("security event not found")
	}

	event.Resolved = true
	event.ResolvedAt = time.Now().Unix()

	return nil
}

// CreateTransactionApproval creates a transaction approval request
func (as *AdvancedSecurity) CreateTransactionApproval(walletID, transactionID string, amount float64, currency, approvalMethod string) (*TransactionApproval, error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	approval := &TransactionApproval{
		ID:               generateApprovalID(),
		WalletID:         walletID,
		TransactionID:    transactionID,
		Amount:           amount,
		Currency:         currency,
		RequiresApproval: true,
		ApprovalMethod:   approvalMethod,
		Status:           "pending",
		CreatedAt:        time.Now().Unix(),
		ApprovalCode:     generateApprovalCode(),
	}

	as.approvals[approval.ID] = approval

	return approval, nil
}

// ApproveTransaction approves a transaction
func (as *AdvancedSecurity) ApproveTransaction(approvalID, code string) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	approval, exists := as.approvals[approvalID]
	if !exists {
		return errors.New("approval request not found")
	}

	if approval.Status != "pending" {
		return errors.New("approval request is not pending")
	}

	if approval.ApprovalCode != code {
		return errors.New("invalid approval code")
	}

	approval.Status = "approved"
	approval.ApprovedAt = time.Now().Unix()

	return nil
}

// RejectTransaction rejects a transaction
func (as *AdvancedSecurity) RejectTransaction(approvalID string) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	approval, exists := as.approvals[approvalID]
	if !exists {
		return errors.New("approval request not found")
	}

	if approval.Status != "pending" {
		return errors.New("approval request is not pending")
	}

	approval.Status = "rejected"

	return nil
}

// GetPendingApprovals returns pending approvals for a wallet
func (as *AdvancedSecurity) GetPendingApprovals(walletID string) []*TransactionApproval {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	approvals := make([]*TransactionApproval, 0)
	for _, approval := range as.approvals {
		if approval.WalletID == walletID && approval.Status == "pending" {
			approvals = append(approvals, approval)
		}
	}

	return approvals
}

// Helper functions

func generateDeviceID() string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	return "dev_" + hex.EncodeToString(randomBytes)
}

func generateSecurityEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}

func generateApprovalID() string {
	return fmt.Sprintf("approval_%d", time.Now().UnixNano())
}

func generateApprovalCode() string {
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	return fmt.Sprintf("%06d", int(randomBytes[0])<<16|int(randomBytes[1])<<8|int(randomBytes[2]))
}
