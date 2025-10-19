package bituncoin

import (
	"errors"
	"sync"
	"time"
)

// StakingPool manages BTN staking
type StakingPool struct {
	MinStake     float64
	LockPeriod   int64 // in seconds
	AnnualReward float64
	Stakes       map[string]*Stake
	mutex        sync.RWMutex
}

// Stake represents a user's staking position
type Stake struct {
	Address     string
	Amount      float64
	StakedAt    int64
	ClaimedAt   int64
	UnstakedAt  int64
	IsActive    bool
}

// NewStakingPool creates a new BTN staking pool
func NewStakingPool(minStake float64, lockPeriod int64, annualReward float64) *StakingPool {
	return &StakingPool{
		MinStake:     minStake,
		LockPeriod:   lockPeriod,
		AnnualReward: annualReward,
		Stakes:       make(map[string]*Stake),
	}
}

// CreateStake creates a new stake
func (sp *StakingPool) CreateStake(address string, amount float64) error {
	if amount < sp.MinStake {
		return errors.New("stake amount below minimum")
	}

	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if _, exists := sp.Stakes[address]; exists {
		return errors.New("stake already exists for address")
	}

	stake := &Stake{
		Address:    address,
		Amount:     amount,
		StakedAt:   time.Now().Unix(),
		ClaimedAt:  time.Now().Unix(),
		UnstakedAt: 0,
		IsActive:   true,
	}

	sp.Stakes[address] = stake
	return nil
}

// CalculateRewards calculates rewards for a stake
func (sp *StakingPool) CalculateRewards(address string) (float64, error) {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return 0, errors.New("stake not found")
	}

	if !stake.IsActive {
		return 0, nil
	}

	now := time.Now().Unix()
	stakingDuration := float64(now - stake.ClaimedAt)
	secondsPerYear := float64(365 * 24 * 60 * 60)
	
	rewards := stake.Amount * (sp.AnnualReward / 100.0) * (stakingDuration / secondsPerYear)
	return rewards, nil
}

// ClaimRewards claims accumulated rewards
func (sp *StakingPool) ClaimRewards(address string) (float64, error) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return 0, errors.New("stake not found")
	}

	if !stake.IsActive {
		return 0, errors.New("stake is not active")
	}

	now := time.Now().Unix()
	stakingDuration := float64(now - stake.ClaimedAt)
	secondsPerYear := float64(365 * 24 * 60 * 60)
	
	rewards := stake.Amount * (sp.AnnualReward / 100.0) * (stakingDuration / secondsPerYear)
	
	stake.ClaimedAt = now
	return rewards, nil
}

// Unstake removes a stake after lock period
func (sp *StakingPool) Unstake(address string) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return errors.New("stake not found")
	}

	if !stake.IsActive {
		return errors.New("stake already unstaked")
	}

	now := time.Now().Unix()
	if now-stake.StakedAt < sp.LockPeriod {
		return errors.New("stake still locked")
	}

	stake.IsActive = false
	stake.UnstakedAt = now
	return nil
}

// IncreaseStake adds more BTN to existing stake
func (sp *StakingPool) IncreaseStake(address string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return errors.New("stake not found")
	}

	if !stake.IsActive {
		return errors.New("cannot increase inactive stake")
	}

	stake.Amount += amount
	return nil
}

// GetPoolInfo returns staking pool statistics
func (sp *StakingPool) GetPoolInfo() map[string]interface{} {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()

	totalStaked := 0.0
	activeStakers := 0

	for _, stake := range sp.Stakes {
		if stake.IsActive {
			totalStaked += stake.Amount
			activeStakers++
		}
	}

	return map[string]interface{}{
		"totalStaked":   totalStaked,
		"activeStakers": activeStakers,
		"minStake":      sp.MinStake,
		"lockPeriod":    sp.LockPeriod,
		"annualReward":  sp.AnnualReward,
	}
}
