package consensus

import (
	"fmt"
	"sync"
	"time"

	"github.com/Bituncoin/Bituncoin/bituncoin-btn/core"
)

// Validator represents a network validator
type Validator struct {
	Address    string
	Stake      float64
	LastActive time.Time
	Votes      int
}

// ValidatorEcho implements a simple consensus mechanism
type ValidatorEcho struct {
	Validators map[string]*Validator
	Quorum     int
	mu         sync.RWMutex
}

// NewValidatorEcho creates a new validator echo consensus
func NewValidatorEcho(quorum int) *ValidatorEcho {
	return &ValidatorEcho{
		Validators: make(map[string]*Validator),
		Quorum:     quorum,
	}
}

// RegisterValidator registers a new validator
func (ve *ValidatorEcho) RegisterValidator(address string, stake float64) error {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if stake < 1000.0 {
		return fmt.Errorf("minimum stake requirement not met: need 1000, got %.2f", stake)
	}

	ve.Validators[address] = &Validator{
		Address:    address,
		Stake:      stake,
		LastActive: time.Now(),
		Votes:      0,
	}

	fmt.Printf("Validator %s registered with stake %.2f\n", address, stake)
	return nil
}

// ValidateBlock validates a block using validator consensus
func (ve *ValidatorEcho) ValidateBlock(block core.Block) (bool, error) {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if len(ve.Validators) == 0 {
		return false, fmt.Errorf("no validators registered")
	}

	// Echo consensus: validators vote on block validity
	votes := 0
	totalStake := 0.0
	votedStake := 0.0

	// Calculate total stake
	for _, validator := range ve.Validators {
		totalStake += validator.Stake
	}

	// Simulate validator voting (in real implementation, validators would vote)
	for _, validator := range ve.Validators {
		// Check if validator is active
		if time.Since(validator.LastActive) < 5*time.Minute {
			votes++
			votedStake += validator.Stake
			validator.Votes++
		}
	}

	// Check if quorum is reached
	quorumMet := votes >= ve.Quorum
	stakeMajority := votedStake > (totalStake * 0.51)

	if quorumMet && stakeMajority {
		fmt.Printf("Block validated: %d votes, %.2f%% stake\n", votes, (votedStake/totalStake)*100)
		return true, nil
	}

	return false, fmt.Errorf("consensus not reached: votes=%d, stake=%.2f%%", votes, (votedStake/totalStake)*100)
}

// UpdateValidator updates validator's last active time
func (ve *ValidatorEcho) UpdateValidator(address string) {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	if validator, exists := ve.Validators[address]; exists {
		validator.LastActive = time.Now()
	}
}

// GetActiveValidators returns the number of active validators
func (ve *ValidatorEcho) GetActiveValidators() int {
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	active := 0
	for _, validator := range ve.Validators {
		if time.Since(validator.LastActive) < 5*time.Minute {
			active++
		}
	}

	return active
}

// RemoveValidator removes a validator from the network
func (ve *ValidatorEcho) RemoveValidator(address string) {
	ve.mu.Lock()
	defer ve.mu.Unlock()

	delete(ve.Validators, address)
	fmt.Printf("Validator %s removed\n", address)
}

// GetValidatorInfo returns information about a validator
func (ve *ValidatorEcho) GetValidatorInfo(address string) (*Validator, error) {
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	validator, exists := ve.Validators[address]
	if !exists {
		return nil, fmt.Errorf("validator not found")
	}

	return validator, nil
}

// GetAllValidators returns all registered validators
func (ve *ValidatorEcho) GetAllValidators() []*Validator {
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	validators := make([]*Validator, 0, len(ve.Validators))
	for _, v := range ve.Validators {
		validators = append(validators, v)
	}

	return validators
}

// ValidateTransaction validates a transaction before adding to block
func (ve *ValidatorEcho) ValidateTransaction(tx core.Transaction) bool {
	// Basic validation rules
	if tx.Amount <= 0 {
		return false
	}

	if tx.From == "" || tx.To == "" {
		return false
	}

	if tx.Currency == "" {
		return false
	}

	return true
}

// ReachConsensus attempts to reach consensus on pending transactions
func (ve *ValidatorEcho) ReachConsensus(transactions []core.Transaction) ([]core.Transaction, error) {
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	validTxs := []core.Transaction{}

	for _, tx := range transactions {
		if ve.ValidateTransaction(tx) {
			validTxs = append(validTxs, tx)
		}
	}

	// Check if enough validators are active
	activeValidators := 0
	for _, validator := range ve.Validators {
		if time.Since(validator.LastActive) < 5*time.Minute {
			activeValidators++
		}
	}

	if activeValidators < ve.Quorum {
		return nil, fmt.Errorf("not enough active validators: %d/%d", activeValidators, ve.Quorum)
	}

	fmt.Printf("Consensus reached on %d transactions with %d validators\n", len(validTxs), activeValidators)
	return validTxs, nil
}
