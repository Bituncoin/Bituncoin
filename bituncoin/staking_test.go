package bituncoin

import (
	"testing"
	"time"
)

func TestNewStakingPool(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)

	if pool.MinStake != 100.0 {
		t.Errorf("expected min stake 100.0, got %f", pool.MinStake)
	}

	if pool.AnnualReward != 5.0 {
		t.Errorf("expected annual reward 5.0, got %f", pool.AnnualReward)
	}
}

func TestCreateStake(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	err := pool.CreateStake("BTNaddr1", 1000.0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stake := pool.Stakes["BTNaddr1"]
	if stake.Amount != 1000.0 {
		t.Errorf("expected stake amount 1000.0, got %f", stake.Amount)
	}

	if !stake.IsActive {
		t.Error("expected stake to be active")
	}
}

func TestCreateStakeBelowMinimum(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	err := pool.CreateStake("BTNaddr1", 50.0)

	if err == nil {
		t.Error("expected error for stake below minimum")
	}
}

func TestCreateStakeDuplicate(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	pool.CreateStake("BTNaddr1", 1000.0)
	err := pool.CreateStake("BTNaddr1", 500.0)

	if err == nil {
		t.Error("expected error for duplicate stake")
	}
}

func TestCalculateRewards(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	pool.CreateStake("BTNaddr1", 1000.0)

	// Simulate 1 year passing
	stake := pool.Stakes["BTNaddr1"]
	stake.ClaimedAt = time.Now().Unix() - (365 * 24 * 60 * 60)

	rewards, err := pool.CalculateRewards("BTNaddr1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRewards := 1000.0 * 0.05 // 5% of 1000
	if rewards < expectedRewards*0.99 || rewards > expectedRewards*1.01 {
		t.Errorf("expected rewards around %f, got %f", expectedRewards, rewards)
	}
}

func TestClaimRewards(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	pool.CreateStake("BTNaddr1", 1000.0)

	// Simulate time passing
	stake := pool.Stakes["BTNaddr1"]
	stake.ClaimedAt = time.Now().Unix() - (365 * 24 * 60 * 60)

	rewards, err := pool.ClaimRewards("BTNaddr1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rewards <= 0 {
		t.Error("expected positive rewards")
	}

	// Check that ClaimedAt was updated to current time
	now := time.Now().Unix()
	if pool.Stakes["BTNaddr1"].ClaimedAt < now-5 || pool.Stakes["BTNaddr1"].ClaimedAt > now+5 {
		t.Error("expected ClaimedAt to be updated to current time")
	}
}

func TestUnstakeBeforeLockPeriod(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	pool.CreateStake("BTNaddr1", 1000.0)

	err := pool.Unstake("BTNaddr1")
	if err == nil {
		t.Error("expected error when unstaking before lock period")
	}
}

func TestIncreaseStake(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	pool.CreateStake("BTNaddr1", 1000.0)

	err := pool.IncreaseStake("BTNaddr1", 500.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stake := pool.Stakes["BTNaddr1"]
	if stake.Amount != 1500.0 {
		t.Errorf("expected stake amount 1500.0, got %f", stake.Amount)
	}
}

func TestGetPoolInfo(t *testing.T) {
	pool := NewStakingPool(100.0, 2592000, 5.0)
	pool.CreateStake("BTNaddr1", 1000.0)
	pool.CreateStake("BTNaddr2", 2000.0)

	info := pool.GetPoolInfo()

	if info["activeStakers"] != 2 {
		t.Errorf("expected 2 active stakers, got %v", info["activeStakers"])
	}

	if info["totalStaked"] != 3000.0 {
		t.Errorf("expected total staked 3000.0, got %v", info["totalStaked"])
	}
}
