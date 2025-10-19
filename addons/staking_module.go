package addons

import (
	"errors"
	"sync"
)

// StakingModule implements advanced staking features as an add-on
type StakingModule struct {
	name        string
	version     string
	status      ModuleStatus
	config      map[string]interface{}
	stakePools  map[string]*AdvancedStakePool
	mutex       sync.RWMutex
}

// AdvancedStakePool represents an advanced staking pool
type AdvancedStakePool struct {
	ID              string
	Name            string
	TokenSymbol     string
	APY             float64
	MinStake        float64
	LockPeriod      int64
	TotalStaked     float64
	MaxCapacity     float64
	AutoCompound    bool
	PoolRewards     float64
}

// NewStakingModule creates a new staking add-on module
func NewStakingModule() *StakingModule {
	return &StakingModule{
		name:       "Advanced Staking",
		version:    "1.0.0",
		status:     StatusDisabled,
		stakePools: make(map[string]*AdvancedStakePool),
	}
}

// GetName returns the module name
func (sm *StakingModule) GetName() string {
	return sm.name
}

// GetVersion returns the module version
func (sm *StakingModule) GetVersion() string {
	return sm.version
}

// GetCategory returns the module category
func (sm *StakingModule) GetCategory() ModuleCategory {
	return CategoryStaking
}

// GetDescription returns module description
func (sm *StakingModule) GetDescription() string {
	return "Advanced staking with multiple pools, auto-compounding, and flexible lock periods"
}

// Initialize initializes the module
func (sm *StakingModule) Initialize(config map[string]interface{}) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.config = config
	
	// Create default staking pools
	sm.stakePools["gld-flexible"] = &AdvancedStakePool{
		ID:           "gld-flexible",
		Name:         "GLD Flexible",
		TokenSymbol:  "GLD",
		APY:          3.5,
		MinStake:     10.0,
		LockPeriod:   0,
		MaxCapacity:  1000000.0,
		AutoCompound: false,
	}
	
	sm.stakePools["gld-locked"] = &AdvancedStakePool{
		ID:           "gld-locked",
		Name:         "GLD Locked 90 Days",
		TokenSymbol:  "GLD",
		APY:          7.5,
		MinStake:     100.0,
		LockPeriod:   90 * 24 * 60 * 60,
		MaxCapacity:  5000000.0,
		AutoCompound: true,
	}
	
	return nil
}

// Start starts the module
func (sm *StakingModule) Start() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.status = StatusEnabled
	return nil
}

// Stop stops the module
func (sm *StakingModule) Stop() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.status = StatusDisabled
	return nil
}

// GetStatus returns the current status
func (sm *StakingModule) GetStatus() ModuleStatus {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	return sm.status
}

// Execute executes a module-specific action
func (sm *StakingModule) Execute(action string, params map[string]interface{}) (interface{}, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	switch action {
	case "list_pools":
		return sm.listPools(), nil
		
	case "get_pool":
		poolID, ok := params["pool_id"].(string)
		if !ok {
			return nil, errors.New("pool_id required")
		}
		return sm.getPool(poolID)
		
	case "create_pool":
		return sm.createPool(params)
		
	default:
		return nil, errors.New("unknown action")
	}
}

func (sm *StakingModule) listPools() []AdvancedStakePool {
	pools := make([]AdvancedStakePool, 0, len(sm.stakePools))
	for _, pool := range sm.stakePools {
		pools = append(pools, *pool)
	}
	return pools
}

func (sm *StakingModule) getPool(poolID string) (*AdvancedStakePool, error) {
	pool, exists := sm.stakePools[poolID]
	if !exists {
		return nil, errors.New("pool not found")
	}
	return pool, nil
}

func (sm *StakingModule) createPool(params map[string]interface{}) (*AdvancedStakePool, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, errors.New("id required")
	}
	
	name, ok := params["name"].(string)
	if !ok {
		return nil, errors.New("name required")
	}
	
	pool := &AdvancedStakePool{
		ID:          id,
		Name:        name,
		TokenSymbol: "GLD",
		APY:         5.0,
		MinStake:    100.0,
		LockPeriod:  30 * 24 * 60 * 60,
	}
	
	sm.stakePools[id] = pool
	return pool, nil
}
