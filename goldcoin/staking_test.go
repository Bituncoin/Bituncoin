package goldcoin

import (
	"testing"
	"time"
)

func TestNewStakingPool(t *testing.T) {
	sp := NewStakingPool()
	
	if sp.AnnualReward != 5.0 {
		t.Errorf("Expected annual reward 5.0, got %f", sp.AnnualReward)
	}
	
	if sp.MinStake != 100.0 {
		t.Errorf("Expected min stake 100.0, got %f", sp.MinStake)
	}
	
	if sp.TotalStaked != 0 {
		t.Errorf("Expected total staked 0, got %f", sp.TotalStaked)
	}
}

func TestCreateStake(t *testing.T) {
	sp := NewStakingPool()
	
	err := sp.CreateStake("address1", 500.0)
	if err != nil {
		t.Fatalf("Failed to create stake: %v", err)
	}
	
	if sp.TotalStaked != 500.0 {
		t.Errorf("Expected total staked 500.0, got %f", sp.TotalStaked)
	}
	
	stake, err := sp.GetStakeInfo("address1")
	if err != nil {
		t.Fatalf("Failed to get stake info: %v", err)
	}
	
	if stake.Amount != 500.0 {
		t.Errorf("Expected stake amount 500.0, got %f", stake.Amount)
	}
	
	if !stake.IsActive {
		t.Error("Expected stake to be active")
	}
}

func TestCreateStakeBelowMinimum(t *testing.T) {
	sp := NewStakingPool()
	
	err := sp.CreateStake("address1", 50.0)
	if err == nil {
		t.Error("Expected error for stake below minimum, got nil")
	}
}

func TestCreateStakeDuplicate(t *testing.T) {
	sp := NewStakingPool()
	
	sp.CreateStake("address1", 500.0)
	
	err := sp.CreateStake("address1", 500.0)
	if err == nil {
		t.Error("Expected error for duplicate stake, got nil")
	}
}

func TestCalculateRewards(t *testing.T) {
	sp := NewStakingPool()
	
	sp.CreateStake("address1", 1000.0)
	
	// Wait a moment to simulate time passage
	time.Sleep(10 * time.Millisecond)
	
	rewards, err := sp.CalculateRewards("address1")
	if err != nil {
		t.Fatalf("Failed to calculate rewards: %v", err)
	}
	
	// Rewards should be very small but positive
	if rewards < 0 {
		t.Errorf("Expected positive rewards, got %f", rewards)
	}
}

func TestClaimRewards(t *testing.T) {
	sp := NewStakingPool()
	
	sp.CreateStake("address1", 1000.0)
	
	// Sleep long enough to accumulate measurable rewards
	time.Sleep(100 * time.Millisecond)
	
	rewards, err := sp.ClaimRewards("address1")
	if err != nil {
		t.Fatalf("Failed to claim rewards: %v", err)
	}
	
	if rewards <= 0 {
		t.Errorf("Expected positive rewards, got %f", rewards)
	}
}

func TestUnstakeBeforeLockPeriod(t *testing.T) {
	sp := NewStakingPool()
	sp.LockPeriod = 3600 // 1 hour for testing
	
	sp.CreateStake("address1", 1000.0)
	
	_, _, err := sp.Unstake("address1")
	if err == nil {
		t.Error("Expected error when unstaking before lock period, got nil")
	}
}

func TestIncreaseStake(t *testing.T) {
	sp := NewStakingPool()
	
	sp.CreateStake("address1", 500.0)
	
	err := sp.IncreaseStake("address1", 300.0)
	if err != nil {
		t.Fatalf("Failed to increase stake: %v", err)
	}
	
	if sp.TotalStaked != 800.0 {
		t.Errorf("Expected total staked 800.0, got %f", sp.TotalStaked)
	}
	
	stake, _ := sp.GetStakeInfo("address1")
	if stake.Amount != 800.0 {
		t.Errorf("Expected stake amount 800.0, got %f", stake.Amount)
	}
}

func TestGetPoolInfo(t *testing.T) {
	sp := NewStakingPool()
	
	sp.CreateStake("address1", 500.0)
	sp.CreateStake("address2", 1000.0)
	
	info := sp.GetPoolInfo()
	
	if info["totalStaked"] != 1500.0 {
		t.Errorf("Expected total staked 1500.0 in pool info, got %v", info["totalStaked"])
	}
	
	if info["activeStakes"] != 2 {
		t.Errorf("Expected 2 active stakes in pool info, got %v", info["activeStakes"])
	}
}
