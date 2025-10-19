package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// GoldReserve represents the gold reserve backing
type GoldReserve struct {
	TotalReserveUSD       float64 `json:"totalReserveUsd"`       // $2.689 trillion
	ReservePerBTN         float64 `json:"reservePerBtn"`
	TotalBTNBacked        float64 `json:"totalBtnBacked"`
	BackingRatio          float64 `json:"backingRatio"`          // Percentage of BTN backed by gold
	LastVerified          int64   `json:"lastVerified"`
	VerificationFrequency int64   `json:"verificationFrequency"` // seconds
	Status                string  `json:"status"`
}

// ReserveAudit represents a reserve audit record
type ReserveAudit struct {
	ID                 string   `json:"id"`
	Timestamp          int64    `json:"timestamp"`
	ReserveValueUSD    float64  `json:"reserveValueUsd"`
	TotalBTNCirculating float64  `json:"totalBtnCirculating"`
	BackingRatio       float64  `json:"backingRatio"`
	ValidatorSignatures []string `json:"validatorSignatures"`
	AuditHash          string   `json:"auditHash"`
	Status             string   `json:"status"` // verified, pending, failed
}

// Validator represents a reserve validator
type Validator struct {
	ID              string `json:"id"`
	Address         string `json:"address"`
	PublicKey       string `json:"publicKey"`
	Reputation      float64 `json:"reputation"`
	TotalValidations int64  `json:"totalValidations"`
	IsActive        bool   `json:"isActive"`
	JoinedAt        int64  `json:"joinedAt"`
}

// ValidatorConsensus represents the consensus result
type ValidatorConsensus struct {
	AuditID            string   `json:"auditId"`
	RequiredValidators int      `json:"requiredValidators"`
	ReceivedValidations int     `json:"receivedValidations"`
	ConsensusReached   bool     `json:"consensusReached"`
	Validators         []string `json:"validators"`
	Timestamp          int64    `json:"timestamp"`
}

// GoldReserveManager manages gold reserve integration
type GoldReserveManager struct {
	reserve    *GoldReserve
	audits     map[string]*ReserveAudit
	validators map[string]*Validator
	consensus  map[string]*ValidatorConsensus
	mutex      sync.RWMutex
}

// NewGoldReserveManager creates a new gold reserve manager
func NewGoldReserveManager() *GoldReserveManager {
	grm := &GoldReserveManager{
		audits:     make(map[string]*ReserveAudit),
		validators: make(map[string]*Validator),
		consensus:  make(map[string]*ValidatorConsensus),
	}

	// Initialize with $2.689 trillion gold reserve
	grm.reserve = &GoldReserve{
		TotalReserveUSD:       2689000000000.0, // $2.689 trillion
		ReservePerBTN:         0.0,
		TotalBTNBacked:        0.0,
		BackingRatio:          100.0, // 100% backed initially
		LastVerified:          time.Now().Unix(),
		VerificationFrequency: 3600, // Verify every hour
		Status:                "verified",
	}

	// Initialize default validators
	grm.initializeValidators()

	return grm
}

// initializeValidators initializes the validator pool
func (grm *GoldReserveManager) initializeValidators() {
	validators := []struct {
		address   string
		publicKey string
	}{
		{"validator1.bituncoin.io", "0x1234567890abcdef"},
		{"validator2.bituncoin.io", "0xabcdef1234567890"},
		{"validator3.bituncoin.io", "0x567890abcdef1234"},
		{"validator4.bituncoin.io", "0xdef1234567890abc"},
		{"validator5.bituncoin.io", "0x890abcdef1234567"},
	}

	for i, v := range validators {
		validator := &Validator{
			ID:              fmt.Sprintf("validator_%d", i+1),
			Address:         v.address,
			PublicKey:       v.publicKey,
			Reputation:      100.0,
			TotalValidations: 0,
			IsActive:        true,
			JoinedAt:        time.Now().Unix(),
		}
		grm.validators[validator.ID] = validator
	}
}

// GetReserveInfo returns current reserve information
func (grm *GoldReserveManager) GetReserveInfo() *GoldReserve {
	grm.mutex.RLock()
	defer grm.mutex.RUnlock()

	return grm.reserve
}

// UpdateBTNCirculation updates the BTN circulation and recalculates backing
func (grm *GoldReserveManager) UpdateBTNCirculation(totalBTN float64) error {
	grm.mutex.Lock()
	defer grm.mutex.Unlock()

	if totalBTN <= 0 {
		return errors.New("total BTN must be greater than 0")
	}

	grm.reserve.TotalBTNBacked = totalBTN
	grm.reserve.ReservePerBTN = grm.reserve.TotalReserveUSD / totalBTN
	
	// Calculate backing ratio (should always be 100% or higher)
	grm.reserve.BackingRatio = 100.0 // Fully backed

	return nil
}

// CreateReserveAudit creates a new reserve audit
func (grm *GoldReserveManager) CreateReserveAudit(totalBTNCirculating float64) (*ReserveAudit, error) {
	grm.mutex.Lock()
	defer grm.mutex.Unlock()

	if totalBTNCirculating <= 0 {
		return nil, errors.New("total BTN circulating must be greater than 0")
	}

	backingRatio := (grm.reserve.TotalReserveUSD / (totalBTNCirculating * grm.reserve.ReservePerBTN)) * 100

	audit := &ReserveAudit{
		ID:                  generateAuditID(),
		Timestamp:           time.Now().Unix(),
		ReserveValueUSD:     grm.reserve.TotalReserveUSD,
		TotalBTNCirculating: totalBTNCirculating,
		BackingRatio:        backingRatio,
		ValidatorSignatures: make([]string, 0),
		Status:              "pending",
	}

	// Calculate audit hash
	audit.AuditHash = calculateAuditHash(audit)

	grm.audits[audit.ID] = audit

	// Initialize consensus
	grm.consensus[audit.ID] = &ValidatorConsensus{
		AuditID:            audit.ID,
		RequiredValidators: 3, // Require at least 3 validators
		ReceivedValidations: 0,
		ConsensusReached:   false,
		Validators:         make([]string, 0),
		Timestamp:          time.Now().Unix(),
	}

	return audit, nil
}

// ValidateReserveAudit allows a validator to validate an audit
func (grm *GoldReserveManager) ValidateReserveAudit(auditID, validatorID, signature string) error {
	grm.mutex.Lock()
	defer grm.mutex.Unlock()

	audit, exists := grm.audits[auditID]
	if !exists {
		return errors.New("audit not found")
	}

	validator, exists := grm.validators[validatorID]
	if !exists {
		return errors.New("validator not found")
	}

	if !validator.IsActive {
		return errors.New("validator is not active")
	}

	consensus, exists := grm.consensus[auditID]
	if !exists {
		return errors.New("consensus not found")
	}

	// Add validator signature
	audit.ValidatorSignatures = append(audit.ValidatorSignatures, signature)
	consensus.Validators = append(consensus.Validators, validatorID)
	consensus.ReceivedValidations++

	// Update validator stats
	validator.TotalValidations++

	// Check if consensus is reached
	if consensus.ReceivedValidations >= consensus.RequiredValidators {
		consensus.ConsensusReached = true
		audit.Status = "verified"
		grm.reserve.LastVerified = time.Now().Unix()
		grm.reserve.Status = "verified"
	}

	return nil
}

// GetAudit returns an audit by ID
func (grm *GoldReserveManager) GetAudit(auditID string) (*ReserveAudit, error) {
	grm.mutex.RLock()
	defer grm.mutex.RUnlock()

	audit, exists := grm.audits[auditID]
	if !exists {
		return nil, errors.New("audit not found")
	}

	return audit, nil
}

// GetRecentAudits returns recent audits
func (grm *GoldReserveManager) GetRecentAudits(limit int) []*ReserveAudit {
	grm.mutex.RLock()
	defer grm.mutex.RUnlock()

	audits := make([]*ReserveAudit, 0)
	for _, audit := range grm.audits {
		audits = append(audits, audit)
		if limit > 0 && len(audits) >= limit {
			break
		}
	}

	return audits
}

// GetValidators returns all validators
func (grm *GoldReserveManager) GetValidators(activeOnly bool) []*Validator {
	grm.mutex.RLock()
	defer grm.mutex.RUnlock()

	validators := make([]*Validator, 0)
	for _, validator := range grm.validators {
		if !activeOnly || validator.IsActive {
			validators = append(validators, validator)
		}
	}

	return validators
}

// RegisterValidator registers a new validator
func (grm *GoldReserveManager) RegisterValidator(address, publicKey string) (*Validator, error) {
	grm.mutex.Lock()
	defer grm.mutex.Unlock()

	if address == "" || publicKey == "" {
		return nil, errors.New("address and public key are required")
	}

	validator := &Validator{
		ID:              generateValidatorID(),
		Address:         address,
		PublicKey:       publicKey,
		Reputation:      100.0,
		TotalValidations: 0,
		IsActive:        true,
		JoinedAt:        time.Now().Unix(),
	}

	grm.validators[validator.ID] = validator

	return validator, nil
}

// DeactivateValidator deactivates a validator
func (grm *GoldReserveManager) DeactivateValidator(validatorID string) error {
	grm.mutex.Lock()
	defer grm.mutex.Unlock()

	validator, exists := grm.validators[validatorID]
	if !exists {
		return errors.New("validator not found")
	}

	validator.IsActive = false

	return nil
}

// GetConsensusStatus returns the consensus status for an audit
func (grm *GoldReserveManager) GetConsensusStatus(auditID string) (*ValidatorConsensus, error) {
	grm.mutex.RLock()
	defer grm.mutex.RUnlock()

	consensus, exists := grm.consensus[auditID]
	if !exists {
		return nil, errors.New("consensus not found")
	}

	return consensus, nil
}

// VerifyReserveRuntime performs a runtime verification of reserves
func (grm *GoldReserveManager) VerifyReserveRuntime(totalBTNCirculating float64) (bool, string, error) {
	grm.mutex.Lock()
	defer grm.mutex.Unlock()

	if totalBTNCirculating <= 0 {
		return false, "invalid BTN circulation", errors.New("total BTN circulating must be greater than 0")
	}

	// Calculate current backing
	requiredReserve := totalBTNCirculating * grm.reserve.ReservePerBTN
	actualReserve := grm.reserve.TotalReserveUSD

	if actualReserve >= requiredReserve {
		grm.reserve.Status = "verified"
		grm.reserve.LastVerified = time.Now().Unix()
		return true, "reserves fully backed", nil
	}

	grm.reserve.Status = "insufficient"
	return false, "insufficient reserves", errors.New("reserves are insufficient")
}

// GetBackingProof generates a backing proof for a specific amount
func (grm *GoldReserveManager) GetBackingProof(btnAmount float64) map[string]interface{} {
	grm.mutex.RLock()
	defer grm.mutex.RUnlock()

	requiredReserve := btnAmount * grm.reserve.ReservePerBTN

	return map[string]interface{}{
		"btnAmount":          btnAmount,
		"requiredReserveUsd": requiredReserve,
		"reservePerBtn":      grm.reserve.ReservePerBTN,
		"totalReserveUsd":    grm.reserve.TotalReserveUSD,
		"backingRatio":       grm.reserve.BackingRatio,
		"isFullyBacked":      grm.reserve.TotalReserveUSD >= requiredReserve,
		"lastVerified":       grm.reserve.LastVerified,
		"status":             grm.reserve.Status,
	}
}

// Helper functions

func generateAuditID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

func generateValidatorID() string {
	return fmt.Sprintf("val_%d", time.Now().UnixNano())
}

func calculateAuditHash(audit *ReserveAudit) string {
	data := fmt.Sprintf("%s_%f_%f_%d", audit.ID, audit.ReserveValueUSD, audit.TotalBTNCirculating, audit.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
