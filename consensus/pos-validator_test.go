package consensus

import (
	"testing"
)

func TestNewProofOfStake(t *testing.T) {
	pos := NewProofOfStake()
	
	if pos.MinStake != 1000.0 {
		t.Errorf("Expected min stake 1000.0, got %f", pos.MinStake)
	}
	
	if pos.BlockTime != 10 {
		t.Errorf("Expected block time 10, got %d", pos.BlockTime)
	}
	
	if pos.RewardPerBlock != 2.0 {
		t.Errorf("Expected reward per block 2.0, got %f", pos.RewardPerBlock)
	}
}

func TestRegisterValidator(t *testing.T) {
	pos := NewProofOfStake()
	
	err := pos.RegisterValidator("validator1", 2000.0)
	if err != nil {
		t.Fatalf("Failed to register validator: %v", err)
	}
	
	validator, err := pos.GetValidatorInfo("validator1")
	if err != nil {
		t.Fatalf("Failed to get validator info: %v", err)
	}
	
	if validator.StakedAmount != 2000.0 {
		t.Errorf("Expected staked amount 2000.0, got %f", validator.StakedAmount)
	}
	
	if !validator.IsActive {
		t.Error("Expected validator to be active")
	}
}

func TestRegisterValidatorBelowMinStake(t *testing.T) {
	pos := NewProofOfStake()
	
	err := pos.RegisterValidator("validator1", 500.0)
	if err == nil {
		t.Error("Expected error for stake below minimum, got nil")
	}
}

func TestRegisterValidatorDuplicate(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	
	err := pos.RegisterValidator("validator1", 2000.0)
	if err == nil {
		t.Error("Expected error for duplicate validator, got nil")
	}
}

func TestSelectValidator(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	pos.RegisterValidator("validator2", 3000.0)
	
	validator, err := pos.SelectValidator()
	if err != nil {
		t.Fatalf("Failed to select validator: %v", err)
	}
	
	if validator.Address != "validator1" && validator.Address != "validator2" {
		t.Errorf("Selected validator should be validator1 or validator2, got %s", validator.Address)
	}
}

func TestSelectValidatorNoValidators(t *testing.T) {
	pos := NewProofOfStake()
	
	_, err := pos.SelectValidator()
	if err == nil {
		t.Error("Expected error when no validators available, got nil")
	}
}

func TestCreateBlock(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	
	transactions := []string{"tx1", "tx2", "tx3"}
	prevHash := "0000000000"
	
	block, err := pos.CreateBlock(transactions, prevHash)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}
	
	if block.Index != 1 {
		t.Errorf("Expected block index 1, got %d", block.Index)
	}
	
	if block.PrevHash != prevHash {
		t.Errorf("Expected prev hash %s, got %s", prevHash, block.PrevHash)
	}
	
	if len(block.Transactions) != 3 {
		t.Errorf("Expected 3 transactions, got %d", len(block.Transactions))
	}
	
	if block.Hash == "" {
		t.Error("Block hash should not be empty")
	}
}

func TestValidateBlock(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	
	block, _ := pos.CreateBlock([]string{"tx1"}, "0000000000")
	
	err := pos.ValidateBlock(block)
	if err != nil {
		t.Errorf("Block validation failed: %v", err)
	}
}

func TestValidateBlockInvalidHash(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	
	block, _ := pos.CreateBlock([]string{"tx1"}, "0000000000")
	block.Hash = "invalid_hash"
	
	err := pos.ValidateBlock(block)
	if err == nil {
		t.Error("Expected error for invalid block hash, got nil")
	}
}

func TestUnstakeValidator(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	
	stakedAmount, err := pos.UnstakeValidator("validator1")
	if err != nil {
		t.Fatalf("Failed to unstake validator: %v", err)
	}
	
	if stakedAmount != 2000.0 {
		t.Errorf("Expected staked amount 2000.0, got %f", stakedAmount)
	}
	
	_, err = pos.GetValidatorInfo("validator1")
	if err == nil {
		t.Error("Expected error when getting unstaked validator, got nil")
	}
}

func TestGetAllValidators(t *testing.T) {
	pos := NewProofOfStake()
	
	pos.RegisterValidator("validator1", 2000.0)
	pos.RegisterValidator("validator2", 3000.0)
	pos.RegisterValidator("validator3", 1500.0)
	
	validators := pos.GetAllValidators()
	
	if len(validators) != 3 {
		t.Errorf("Expected 3 validators, got %d", len(validators))
	}
}
