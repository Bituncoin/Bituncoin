package consensus

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Validator represents a PoS validator
type Validator struct {
	Address     string
	StakedAmount float64
	RewardRate   float64
	IsActive     bool
	JoinedAt     int64
}

// Block represents a blockchain block
type Block struct {
	Index        int
	Timestamp    int64
	Transactions []string
	PrevHash     string
	Hash         string
	Validator    string
	Nonce        int
}

// ProofOfStake implements the Proof-of-Stake consensus mechanism
type ProofOfStake struct {
	Validators       map[string]*Validator
	MinStake         float64
	BlockTime        int64
	RewardPerBlock   float64
	mutex            sync.RWMutex
	currentBlockIndex int
}

// NewProofOfStake creates a new PoS consensus instance
func NewProofOfStake() *ProofOfStake {
	return &ProofOfStake{
		Validators:       make(map[string]*Validator),
		MinStake:         1000.0,  // Minimum 1000 GLD to stake
		BlockTime:        10,       // 10 seconds block time
		RewardPerBlock:   2.0,      // 2 GLD per block
		currentBlockIndex: 0,
	}
}

// RegisterValidator registers a new validator
func (pos *ProofOfStake) RegisterValidator(address string, stakedAmount float64) error {
	pos.mutex.Lock()
	defer pos.mutex.Unlock()

	if address == "" {
		return errors.New("invalid address")
	}

	if stakedAmount < pos.MinStake {
		return fmt.Errorf("insufficient stake: minimum %f GLD required", pos.MinStake)
	}

	if _, exists := pos.Validators[address]; exists {
		return errors.New("validator already registered")
	}

	validator := &Validator{
		Address:      address,
		StakedAmount: stakedAmount,
		RewardRate:   5.0, // 5% annual reward
		IsActive:     true,
		JoinedAt:     time.Now().Unix(),
	}

	pos.Validators[address] = validator
	return nil
}

// SelectValidator selects a validator based on stake weight
func (pos *ProofOfStake) SelectValidator() (*Validator, error) {
	pos.mutex.RLock()
	defer pos.mutex.RUnlock()

	if len(pos.Validators) == 0 {
		return nil, errors.New("no validators available")
	}

	// Calculate total stake
	var totalStake float64
	activeValidators := make([]*Validator, 0)
	
	for _, v := range pos.Validators {
		if v.IsActive {
			totalStake += v.StakedAmount
			activeValidators = append(activeValidators, v)
		}
	}

	if len(activeValidators) == 0 {
		return nil, errors.New("no active validators")
	}

	// Weighted random selection based on stake
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Float64() * totalStake
	
	var cumulative float64
	for _, v := range activeValidators {
		cumulative += v.StakedAmount
		if cumulative >= randomValue {
			return v, nil
		}
	}

	// Fallback to first active validator
	return activeValidators[0], nil
}

// CreateBlock creates a new block with the selected validator
func (pos *ProofOfStake) CreateBlock(transactions []string, prevHash string) (*Block, error) {
	validator, err := pos.SelectValidator()
	if err != nil {
		return nil, err
	}

	pos.mutex.Lock()
	pos.currentBlockIndex++
	blockIndex := pos.currentBlockIndex
	pos.mutex.Unlock()

	block := &Block{
		Index:        blockIndex,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Validator:    validator.Address,
		Nonce:        0,
	}

	block.Hash = block.calculateHash()

	// Reward the validator
	pos.rewardValidator(validator.Address)

	return block, nil
}

// calculateHash calculates the block hash
func (b *Block) calculateHash() string {
	data := fmt.Sprintf("%d%d%v%s%s%d",
		b.Index, b.Timestamp, b.Transactions, b.PrevHash, b.Validator, b.Nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// rewardValidator rewards a validator for creating a block
func (pos *ProofOfStake) rewardValidator(address string) {
	pos.mutex.Lock()
	defer pos.mutex.Unlock()

	if validator, exists := pos.Validators[address]; exists {
		validator.StakedAmount += pos.RewardPerBlock
	}
}

// GetValidatorInfo returns information about a validator
func (pos *ProofOfStake) GetValidatorInfo(address string) (*Validator, error) {
	pos.mutex.RLock()
	defer pos.mutex.RUnlock()

	validator, exists := pos.Validators[address]
	if !exists {
		return nil, errors.New("validator not found")
	}

	return validator, nil
}

// UnstakeValidator removes a validator and returns their stake
func (pos *ProofOfStake) UnstakeValidator(address string) (float64, error) {
	pos.mutex.Lock()
	defer pos.mutex.Unlock()

	validator, exists := pos.Validators[address]
	if !exists {
		return 0, errors.New("validator not found")
	}

	stakedAmount := validator.StakedAmount
	delete(pos.Validators, address)
	
	return stakedAmount, nil
}

// GetAllValidators returns all validators
func (pos *ProofOfStake) GetAllValidators() []*Validator {
	pos.mutex.RLock()
	defer pos.mutex.RUnlock()

	validators := make([]*Validator, 0, len(pos.Validators))
	for _, v := range pos.Validators {
		validators = append(validators, v)
	}

	return validators
}

// ValidateBlock validates a block
func (pos *ProofOfStake) ValidateBlock(block *Block) error {
	if block == nil {
		return errors.New("block is nil")
	}

	// Verify hash
	expectedHash := block.calculateHash()
	if block.Hash != expectedHash {
		return errors.New("invalid block hash")
	}

	// Verify validator exists and is active
	pos.mutex.RLock()
	validator, exists := pos.Validators[block.Validator]
	pos.mutex.RUnlock()

	if !exists {
		return errors.New("validator not found")
	}

	if !validator.IsActive {
		return errors.New("validator is not active")
	}

	return nil
}
