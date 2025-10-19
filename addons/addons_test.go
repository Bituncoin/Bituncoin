package addons

import (
	"testing"
)

func TestRegisterModule(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewStakingModule()
	
	err := registry.Register(module, "Bituncoin Team")
	if err != nil {
		t.Fatalf("Failed to register module: %v", err)
	}
	
	// Try to register the same module again
	err = registry.Register(module, "Bituncoin Team")
	if err == nil {
		t.Error("Expected error when registering duplicate module")
	}
}

func TestEnableModule(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewStakingModule()
	
	registry.Register(module, "Bituncoin Team")
	
	config := make(map[string]interface{})
	err := registry.Enable("Advanced Staking", config)
	if err != nil {
		t.Fatalf("Failed to enable module: %v", err)
	}
	
	info, err := registry.GetModuleInfo("Advanced Staking")
	if err != nil {
		t.Fatalf("Failed to get module info: %v", err)
	}
	
	if info.Status != StatusEnabled {
		t.Errorf("Expected module status 'enabled', got '%s'", info.Status)
	}
}

func TestDisableModule(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewStakingModule()
	
	registry.Register(module, "Bituncoin Team")
	registry.Enable("Advanced Staking", make(map[string]interface{}))
	
	err := registry.Disable("Advanced Staking")
	if err != nil {
		t.Fatalf("Failed to disable module: %v", err)
	}
	
	info, err := registry.GetModuleInfo("Advanced Staking")
	if err != nil {
		t.Fatalf("Failed to get module info: %v", err)
	}
	
	if info.Status != StatusDisabled {
		t.Errorf("Expected module status 'disabled', got '%s'", info.Status)
	}
}

func TestListModules(t *testing.T) {
	registry := NewModuleRegistry()
	
	staking := NewStakingModule()
	lending := NewLendingModule()
	
	registry.Register(staking, "Bituncoin Team")
	registry.Register(lending, "Bituncoin Team")
	
	modules := registry.ListModules()
	
	if len(modules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(modules))
	}
}

func TestListModulesByCategory(t *testing.T) {
	registry := NewModuleRegistry()
	
	staking := NewStakingModule()
	lending := NewLendingModule()
	
	registry.Register(staking, "Bituncoin Team")
	registry.Register(lending, "Bituncoin Team")
	
	stakingModules := registry.ListModulesByCategory(CategoryStaking)
	if len(stakingModules) != 1 {
		t.Errorf("Expected 1 staking module, got %d", len(stakingModules))
	}
	
	lendingModules := registry.ListModulesByCategory(CategoryLending)
	if len(lendingModules) != 1 {
		t.Errorf("Expected 1 lending module, got %d", len(lendingModules))
	}
}

func TestExecuteModuleAction(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewStakingModule()
	
	registry.Register(module, "Bituncoin Team")
	registry.Enable("Advanced Staking", make(map[string]interface{}))
	
	// Test listing pools
	result, err := registry.Execute("Advanced Staking", "list_pools", nil)
	if err != nil {
		t.Fatalf("Failed to execute action: %v", err)
	}
	
	pools, ok := result.([]AdvancedStakePool)
	if !ok {
		t.Error("Expected result to be []AdvancedStakePool")
	}
	
	if len(pools) < 2 {
		t.Errorf("Expected at least 2 pools, got %d", len(pools))
	}
}

func TestUnregisterModule(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewStakingModule()
	
	registry.Register(module, "Bituncoin Team")
	
	err := registry.Unregister("Advanced Staking")
	if err != nil {
		t.Fatalf("Failed to unregister module: %v", err)
	}
	
	_, err = registry.GetModuleInfo("Advanced Staking")
	if err == nil {
		t.Error("Expected error when getting unregistered module")
	}
}

func TestStakingModuleCreatePool(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewStakingModule()
	
	registry.Register(module, "Bituncoin Team")
	registry.Enable("Advanced Staking", make(map[string]interface{}))
	
	params := map[string]interface{}{
		"id":   "custom-pool",
		"name": "Custom Staking Pool",
	}
	
	result, err := registry.Execute("Advanced Staking", "create_pool", params)
	if err != nil {
		t.Fatalf("Failed to create pool: %v", err)
	}
	
	pool, ok := result.(*AdvancedStakePool)
	if !ok {
		t.Error("Expected result to be *AdvancedStakePool")
	}
	
	if pool.ID != "custom-pool" {
		t.Errorf("Expected pool ID 'custom-pool', got '%s'", pool.ID)
	}
}

func TestLendingModuleCreateOffer(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewLendingModule()
	
	registry.Register(module, "Bituncoin Team")
	registry.Enable("DeFi Lending", make(map[string]interface{}))
	
	params := map[string]interface{}{
		"lender":        "user123",
		"amount":        1000.0,
		"interest_rate": 7.5,
	}
	
	result, err := registry.Execute("DeFi Lending", "create_offer", params)
	if err != nil {
		t.Fatalf("Failed to create offer: %v", err)
	}
	
	offer, ok := result.(*LendingOffer)
	if !ok {
		t.Error("Expected result to be *LendingOffer")
	}
	
	if offer.Lender != "user123" {
		t.Errorf("Expected lender 'user123', got '%s'", offer.Lender)
	}
	
	if offer.Amount != 1000.0 {
		t.Errorf("Expected amount 1000.0, got %f", offer.Amount)
	}
}

func TestLendingModuleCreateLoan(t *testing.T) {
	registry := NewModuleRegistry()
	module := NewLendingModule()
	
	registry.Register(module, "Bituncoin Team")
	registry.Enable("DeFi Lending", make(map[string]interface{}))
	
	// First create an offer
	offerParams := map[string]interface{}{
		"lender":        "lender123",
		"amount":        1000.0,
		"interest_rate": 5.0,
	}
	
	offerResult, err := registry.Execute("DeFi Lending", "create_offer", offerParams)
	if err != nil {
		t.Fatalf("Failed to create offer: %v", err)
	}
	
	offer, _ := offerResult.(*LendingOffer)
	
	// Now create a loan from the offer
	loanParams := map[string]interface{}{
		"borrower":   "borrower123",
		"offer_id":   offer.ID,
		"collateral": 1600.0, // 160% collateral
	}
	
	loanResult, err := registry.Execute("DeFi Lending", "create_loan", loanParams)
	if err != nil {
		t.Fatalf("Failed to create loan: %v", err)
	}
	
	loan, ok := loanResult.(*Loan)
	if !ok {
		t.Error("Expected result to be *Loan")
	}
	
	if loan.Borrower != "borrower123" {
		t.Errorf("Expected borrower 'borrower123', got '%s'", loan.Borrower)
	}
	
	if loan.Status != "active" {
		t.Errorf("Expected loan status 'active', got '%s'", loan.Status)
	}
}
