package addons

import (
	"errors"
	"sync"
)

// ModuleStatus represents the status of a module
type ModuleStatus string

const (
	StatusEnabled  ModuleStatus = "enabled"
	StatusDisabled ModuleStatus = "disabled"
	StatusError    ModuleStatus = "error"
)

// ModuleCategory represents the category of functionality
type ModuleCategory string

const (
	CategoryDeFi     ModuleCategory = "defi"
	CategoryStaking  ModuleCategory = "staking"
	CategoryLending  ModuleCategory = "lending"
	CategoryTrading  ModuleCategory = "trading"
	CategoryPayment  ModuleCategory = "payment"
	CategoryAnalytics ModuleCategory = "analytics"
	CategorySecurity ModuleCategory = "security"
	CategoryUtility  ModuleCategory = "utility"
)

// Module represents an add-on module interface
type Module interface {
	// GetName returns the module name
	GetName() string
	
	// GetVersion returns the module version
	GetVersion() string
	
	// GetCategory returns the module category
	GetCategory() ModuleCategory
	
	// GetDescription returns module description
	GetDescription() string
	
	// Initialize initializes the module
	Initialize(config map[string]interface{}) error
	
	// Start starts the module
	Start() error
	
	// Stop stops the module
	Stop() error
	
	// GetStatus returns the current status
	GetStatus() ModuleStatus
	
	// Execute executes a module-specific action
	Execute(action string, params map[string]interface{}) (interface{}, error)
}

// ModuleInfo contains metadata about a module
type ModuleInfo struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Category    ModuleCategory `json:"category"`
	Description string         `json:"description"`
	Status      ModuleStatus   `json:"status"`
	Author      string         `json:"author"`
	Config      map[string]interface{} `json:"config"`
}

// ModuleRegistry manages all add-on modules
type ModuleRegistry struct {
	modules map[string]Module
	info    map[string]*ModuleInfo
	mutex   sync.RWMutex
}

// NewModuleRegistry creates a new module registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]Module),
		info:    make(map[string]*ModuleInfo),
	}
}

// Register registers a new module
func (mr *ModuleRegistry) Register(module Module, author string) error {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	
	name := module.GetName()
	if name == "" {
		return errors.New("module name cannot be empty")
	}
	
	if _, exists := mr.modules[name]; exists {
		return errors.New("module already registered")
	}
	
	mr.modules[name] = module
	mr.info[name] = &ModuleInfo{
		Name:        name,
		Version:     module.GetVersion(),
		Category:    module.GetCategory(),
		Description: module.GetDescription(),
		Status:      StatusDisabled,
		Author:      author,
		Config:      make(map[string]interface{}),
	}
	
	return nil
}

// Unregister removes a module from the registry
func (mr *ModuleRegistry) Unregister(name string) error {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	
	module, exists := mr.modules[name]
	if !exists {
		return errors.New("module not found")
	}
	
	// Stop module if running
	if module.GetStatus() == StatusEnabled {
		module.Stop()
	}
	
	delete(mr.modules, name)
	delete(mr.info, name)
	
	return nil
}

// Enable enables a module
func (mr *ModuleRegistry) Enable(name string, config map[string]interface{}) error {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	
	module, exists := mr.modules[name]
	if !exists {
		return errors.New("module not found")
	}
	
	info := mr.info[name]
	
	// Initialize module with config
	if err := module.Initialize(config); err != nil {
		info.Status = StatusError
		return err
	}
	
	// Start module
	if err := module.Start(); err != nil {
		info.Status = StatusError
		return err
	}
	
	info.Status = StatusEnabled
	info.Config = config
	
	return nil
}

// Disable disables a module
func (mr *ModuleRegistry) Disable(name string) error {
	mr.mutex.Lock()
	defer mr.mutex.Unlock()
	
	module, exists := mr.modules[name]
	if !exists {
		return errors.New("module not found")
	}
	
	info := mr.info[name]
	
	if err := module.Stop(); err != nil {
		info.Status = StatusError
		return err
	}
	
	info.Status = StatusDisabled
	
	return nil
}

// Execute executes a module action
func (mr *ModuleRegistry) Execute(name, action string, params map[string]interface{}) (interface{}, error) {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	
	module, exists := mr.modules[name]
	if !exists {
		return nil, errors.New("module not found")
	}
	
	if module.GetStatus() != StatusEnabled {
		return nil, errors.New("module not enabled")
	}
	
	return module.Execute(action, params)
}

// GetModule returns a module by name
func (mr *ModuleRegistry) GetModule(name string) (Module, error) {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	
	module, exists := mr.modules[name]
	if !exists {
		return nil, errors.New("module not found")
	}
	
	return module, nil
}

// ListModules returns all registered modules
func (mr *ModuleRegistry) ListModules() []*ModuleInfo {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	
	modules := make([]*ModuleInfo, 0, len(mr.info))
	for _, info := range mr.info {
		modules = append(modules, info)
	}
	
	return modules
}

// ListModulesByCategory returns modules in a specific category
func (mr *ModuleRegistry) ListModulesByCategory(category ModuleCategory) []*ModuleInfo {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	
	modules := make([]*ModuleInfo, 0)
	for _, info := range mr.info {
		if info.Category == category {
			modules = append(modules, info)
		}
	}
	
	return modules
}

// GetModuleInfo returns module information
func (mr *ModuleRegistry) GetModuleInfo(name string) (*ModuleInfo, error) {
	mr.mutex.RLock()
	defer mr.mutex.RUnlock()
	
	info, exists := mr.info[name]
	if !exists {
		return nil, errors.New("module not found")
	}
	
	return info, nil
}
