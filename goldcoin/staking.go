package goldcoin

import (
	"errors"
	"sync"
	"time"
)

// StakingPool manages staking for Gold-Coin
type StakingPool struct {
	Stakes        map[string]*Stake
	TotalStaked   float64
	AnnualReward  float64
	MinStake      float64
	LockPeriod    int64 // in seconds
	mutex         sync.RWMutex
}

// Stake represents a staking position
type Stake struct {
	Address      string
	Amount       float64
	StartTime    int64
	UnlockTime   int64
	RewardsClaimed float64
	IsActive     bool
}

// NewStakingPool creates a new staking pool
func NewStakingPool() *StakingPool {
	return &StakingPool{
		Stakes:       make(map[string]*Stake),
		TotalStaked:  0,
		AnnualReward: 5.0,              // 5% annual reward
		MinStake:     100.0,            // Minimum 100 GLD
		LockPeriod:   30 * 24 * 60 * 60, // 30 days
	}
}

// CreateStake creates a new stake for an address
func (sp *StakingPool) CreateStake(address string, amount float64) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if address == "" {
		return errors.New("invalid address")
	}

	if amount < sp.MinStake {
		return errors.New("amount below minimum stake")
	}

	if _, exists := sp.Stakes[address]; exists {
		return errors.New("stake already exists for this address")
	}

	now := time.Now().Unix()
	stake := &Stake{
		Address:      address,
		Amount:       amount,
		StartTime:    now,
		UnlockTime:   now + sp.LockPeriod,
		RewardsClaimed: 0,
		IsActive:     true,
	}

	sp.Stakes[address] = stake
	sp.TotalStaked += amount

	return nil
}

// CalculateRewards calculates the current rewards for a stake
func (sp *StakingPool) CalculateRewards(address string) (float64, error) {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return 0, errors.New("stake not found")
	}

	if !stake.IsActive {
		return 0, errors.New("stake is not active")
	}

	// Calculate time elapsed in years
	now := time.Now().Unix()
	timeElapsed := float64(now - stake.StartTime)
	years := timeElapsed / (365.25 * 24 * 60 * 60)

	// Calculate rewards
	rewards := stake.Amount * (sp.AnnualReward / 100.0) * years
	totalRewards := rewards - stake.RewardsClaimed

	return totalRewards, nil
}

// ClaimRewards claims the accumulated rewards
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

	// Calculate rewards
	now := time.Now().Unix()
	timeElapsed := float64(now - stake.StartTime)
	years := timeElapsed / (365.25 * 24 * 60 * 60)
	rewards := stake.Amount * (sp.AnnualReward / 100.0) * years
	claimableRewards := rewards - stake.RewardsClaimed

	if claimableRewards <= 0 {
		return 0, errors.New("no rewards to claim")
	}

	stake.RewardsClaimed += claimableRewards

	return claimableRewards, nil
}

// Unstake removes a stake and returns the staked amount plus unclaimed rewards
func (sp *StakingPool) Unstake(address string) (float64, float64, error) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return 0, 0, errors.New("stake not found")
	}

	if !stake.IsActive {
		return 0, 0, errors.New("stake is not active")
	}

	now := time.Now().Unix()
	if now < stake.UnlockTime {
		return 0, 0, errors.New("stake is still locked")
	}

	// Calculate final rewards
	timeElapsed := float64(now - stake.StartTime)
	years := timeElapsed / (365.25 * 24 * 60 * 60)
	rewards := stake.Amount * (sp.AnnualReward / 100.0) * years
	unclaimedRewards := rewards - stake.RewardsClaimed

	stakedAmount := stake.Amount
	
	// Deactivate stake
	stake.IsActive = false
	sp.TotalStaked -= stakedAmount

	return stakedAmount, unclaimedRewards, nil
}

// GetStakeInfo returns information about a stake
func (sp *StakingPool) GetStakeInfo(address string) (*Stake, error) {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()

	stake, exists := sp.Stakes[address]
	if !exists {
		return nil, errors.New("stake not found")
	}

	return stake, nil
}

// GetPoolInfo returns information about the staking pool
func (sp *StakingPool) GetPoolInfo() map[string]interface{} {
	sp.mutex.RLock()
	defer sp.mutex.RUnlock()

	activeStakes := 0
	for _, stake := range sp.Stakes {
		if stake.IsActive {
			activeStakes++
		}
	}

	return map[string]interface{}{
		"totalStaked":   sp.TotalStaked,
		"annualReward":  sp.AnnualReward,
		"minStake":      sp.MinStake,
		"lockPeriod":    sp.LockPeriod,
		"activeStakes":  activeStakes,
		"totalStakers":  len(sp.Stakes),
	}
}

// IncreaseStake adds more coins to an existing stake
func (sp *StakingPool) IncreaseStake(address string, amount float64) error {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	stake, exists := sp.Stakes[address]
	if !exists {
		return errors.New("stake not found")
	}

	if !stake.IsActive {
		return errors.New("stake is not active")
	}

	stake.Amount += amount
	sp.TotalStaked += amount

	return nil
}
