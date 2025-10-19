# Module Developer Guide

## Overview
This guide provides developers with instructions for creating custom add-on modules for the Bituncoin wallet platform.

## Table of Contents
1. [Module Architecture](#module-architecture)
2. [Creating a Module](#creating-a-module)
3. [Module Interface](#module-interface)
4. [Module Registration](#module-registration)
5. [Best Practices](#best-practices)
6. [Examples](#examples)

## Module Architecture

Add-on modules extend the Bituncoin wallet functionality through a plug-and-play architecture:

```
┌─────────────────────────────────┐
│      Module Registry            │
│  (Central module manager)       │
└────────────┬────────────────────┘
             │
             ├── Module 1 (Staking)
             ├── Module 2 (Lending)
             ├── Module 3 (Trading)
             └── Module N (Custom)
```

### Module States

Modules can be in one of three states:
- **Disabled**: Registered but not running
- **Enabled**: Active and processing requests
- **Error**: Failed to start or encountered an error

## Creating a Module

### Step 1: Implement the Module Interface

Every module must implement the `Module` interface:

```go
type Module interface {
    GetName() string
    GetVersion() string
    GetCategory() ModuleCategory
    GetDescription() string
    Initialize(config map[string]interface{}) error
    Start() error
    Stop() error
    GetStatus() ModuleStatus
    Execute(action string, params map[string]interface{}) (interface{}, error)
}
```

### Step 2: Define Your Module Structure

```go
package addons

import (
    "errors"
    "sync"
)

// MyCustomModule implements a custom feature
type MyCustomModule struct {
    name    string
    version string
    status  ModuleStatus
    config  map[string]interface{}
    data    map[string]interface{}
    mutex   sync.RWMutex
}

// NewMyCustomModule creates a new instance
func NewMyCustomModule() *MyCustomModule {
    return &MyCustomModule{
        name:    "My Custom Module",
        version: "1.0.0",
        status:  StatusDisabled,
        data:    make(map[string]interface{}),
    }
}
```

### Step 3: Implement Required Methods

```go
// GetName returns the module name
func (m *MyCustomModule) GetName() string {
    return m.name
}

// GetVersion returns the module version
func (m *MyCustomModule) GetVersion() string {
    return m.version
}

// GetCategory returns the module category
func (m *MyCustomModule) GetCategory() ModuleCategory {
    return CategoryUtility // or CategoryDeFi, CategoryStaking, etc.
}

// GetDescription returns module description
func (m *MyCustomModule) GetDescription() string {
    return "A custom module that does something awesome"
}

// Initialize initializes the module with configuration
func (m *MyCustomModule) Initialize(config map[string]interface{}) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.config = config
    
    // Perform initialization tasks
    // Validate config, set up resources, etc.
    
    return nil
}

// Start starts the module
func (m *MyCustomModule) Start() error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // Start background tasks, connect to services, etc.
    
    m.status = StatusEnabled
    return nil
}

// Stop stops the module
func (m *MyCustomModule) Stop() error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // Clean up resources, close connections, etc.
    
    m.status = StatusDisabled
    return nil
}

// GetStatus returns the current status
func (m *MyCustomModule) GetStatus() ModuleStatus {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    return m.status
}

// Execute executes a module-specific action
func (m *MyCustomModule) Execute(action string, params map[string]interface{}) (interface{}, error) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    switch action {
    case "my_action":
        return m.myAction(params)
        
    case "another_action":
        return m.anotherAction(params)
        
    default:
        return nil, errors.New("unknown action")
    }
}

// Private action methods
func (m *MyCustomModule) myAction(params map[string]interface{}) (interface{}, error) {
    // Implement your action logic
    return map[string]string{"result": "success"}, nil
}

func (m *MyCustomModule) anotherAction(params map[string]interface{}) (interface{}, error) {
    // Implement another action
    return nil, nil
}
```

## Module Interface

### Module Categories

Choose the appropriate category for your module:

```go
const (
    CategoryDeFi      ModuleCategory = "defi"      // DeFi features
    CategoryStaking   ModuleCategory = "staking"   // Staking features
    CategoryLending   ModuleCategory = "lending"   // Lending/borrowing
    CategoryTrading   ModuleCategory = "trading"   // Trading features
    CategoryPayment   ModuleCategory = "payment"   // Payment processing
    CategoryAnalytics ModuleCategory = "analytics" // Data analysis
    CategorySecurity  ModuleCategory = "security"  // Security features
    CategoryUtility   ModuleCategory = "utility"   // Utility features
)
```

### Execute Actions

The `Execute` method is the primary interface for interacting with your module. Design clear, RESTful-style actions:

```go
// Good action naming
"list_items"      // GET equivalent
"get_item"        // GET by ID
"create_item"     // POST equivalent
"update_item"     // PUT equivalent
"delete_item"     // DELETE equivalent
"calculate_fee"   // Custom action
```

## Module Registration

### In Application Code

```go
package main

import (
    "github.com/Bituncoin/Bituncoin/addons"
)

func main() {
    // Create module registry
    registry := addons.NewModuleRegistry()
    
    // Create and register your module
    myModule := NewMyCustomModule()
    err := registry.Register(myModule, "Your Name")
    if err != nil {
        panic(err)
    }
    
    // Enable the module with configuration
    config := map[string]interface{}{
        "setting1": "value1",
        "setting2": 42,
    }
    
    err = registry.Enable("My Custom Module", config)
    if err != nil {
        panic(err)
    }
}
```

### Via API

```bash
# Enable module
curl -X POST http://localhost:8080/api/addons/enable \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Custom Module",
    "config": {
      "setting1": "value1",
      "setting2": 42
    }
  }'

# Execute module action
curl -X POST http://localhost:8080/api/addons/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Custom Module",
    "action": "my_action",
    "params": {
      "param1": "value1"
    }
  }'
```

## Best Practices

### 1. Thread Safety
Always use mutexes when accessing shared data:

```go
type MyModule struct {
    data  map[string]interface{}
    mutex sync.RWMutex
}

func (m *MyModule) SafeRead() interface{} {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    return m.data["key"]
}

func (m *MyModule) SafeWrite(key string, value interface{}) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.data[key] = value
}
```

### 2. Error Handling
Return descriptive errors:

```go
func (m *MyModule) Execute(action string, params map[string]interface{}) (interface{}, error) {
    if m.status != StatusEnabled {
        return nil, errors.New("module not enabled")
    }
    
    requiredParam, ok := params["required"]
    if !ok {
        return nil, errors.New("required parameter 'required' missing")
    }
    
    // Process...
}
```

### 3. Configuration Validation
Validate configuration in Initialize:

```go
func (m *MyModule) Initialize(config map[string]interface{}) error {
    // Check required config
    if _, ok := config["api_key"]; !ok {
        return errors.New("api_key required in configuration")
    }
    
    // Validate types
    timeout, ok := config["timeout"].(float64)
    if !ok {
        return errors.New("timeout must be a number")
    }
    
    if timeout < 0 {
        return errors.New("timeout must be positive")
    }
    
    m.config = config
    return nil
}
```

### 4. Resource Cleanup
Clean up resources in Stop:

```go
func (m *MyModule) Stop() error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // Close connections
    if m.connection != nil {
        m.connection.Close()
    }
    
    // Cancel background tasks
    if m.cancelFunc != nil {
        m.cancelFunc()
    }
    
    // Clear data
    m.data = make(map[string]interface{})
    
    m.status = StatusDisabled
    return nil
}
```

### 5. Logging
Include meaningful log messages:

```go
import "log"

func (m *MyModule) Start() error {
    log.Printf("[%s] Starting module version %s", m.name, m.version)
    
    // Start logic...
    
    log.Printf("[%s] Module started successfully", m.name)
    m.status = StatusEnabled
    return nil
}
```

## Examples

### Example 1: Analytics Module

```go
package addons

type AnalyticsModule struct {
    name     string
    version  string
    status   ModuleStatus
    metrics  map[string]float64
    mutex    sync.RWMutex
}

func NewAnalyticsModule() *AnalyticsModule {
    return &AnalyticsModule{
        name:    "Analytics Engine",
        version: "1.0.0",
        status:  StatusDisabled,
        metrics: make(map[string]float64),
    }
}

func (a *AnalyticsModule) Execute(action string, params map[string]interface{}) (interface{}, error) {
    switch action {
    case "track_metric":
        name := params["name"].(string)
        value := params["value"].(float64)
        a.mutex.Lock()
        a.metrics[name] = value
        a.mutex.Unlock()
        return map[string]string{"status": "tracked"}, nil
        
    case "get_metrics":
        a.mutex.RLock()
        defer a.mutex.RUnlock()
        return a.metrics, nil
        
    default:
        return nil, errors.New("unknown action")
    }
}
```

### Example 2: Trading Bot Module

```go
package addons

type TradingBotModule struct {
    name      string
    version   string
    status    ModuleStatus
    strategies map[string]TradingStrategy
    mutex     sync.RWMutex
}

type TradingStrategy struct {
    Name       string
    BuySignal  float64
    SellSignal float64
    Active     bool
}

func NewTradingBotModule() *TradingBotModule {
    return &TradingBotModule{
        name:       "Trading Bot",
        version:    "1.0.0",
        status:     StatusDisabled,
        strategies: make(map[string]TradingStrategy),
    }
}

func (t *TradingBotModule) Execute(action string, params map[string]interface{}) (interface{}, error) {
    switch action {
    case "add_strategy":
        strategy := TradingStrategy{
            Name:       params["name"].(string),
            BuySignal:  params["buy_signal"].(float64),
            SellSignal: params["sell_signal"].(float64),
            Active:     true,
        }
        
        t.mutex.Lock()
        t.strategies[strategy.Name] = strategy
        t.mutex.Unlock()
        
        return strategy, nil
        
    case "list_strategies":
        t.mutex.RLock()
        defer t.mutex.RUnlock()
        
        strategies := make([]TradingStrategy, 0, len(t.strategies))
        for _, s := range t.strategies {
            strategies = append(strategies, s)
        }
        return strategies, nil
        
    default:
        return nil, errors.New("unknown action")
    }
}
```

## Testing Your Module

Create comprehensive tests:

```go
package addons

import "testing"

func TestMyModule(t *testing.T) {
    module := NewMyCustomModule()
    
    // Test initialization
    config := map[string]interface{}{"key": "value"}
    err := module.Initialize(config)
    if err != nil {
        t.Fatalf("Initialize failed: %v", err)
    }
    
    // Test start
    err = module.Start()
    if err != nil {
        t.Fatalf("Start failed: %v", err)
    }
    
    if module.GetStatus() != StatusEnabled {
        t.Error("Module should be enabled")
    }
    
    // Test actions
    result, err := module.Execute("my_action", map[string]interface{}{})
    if err != nil {
        t.Fatalf("Execute failed: %v", err)
    }
    
    // Test stop
    err = module.Stop()
    if err != nil {
        t.Fatalf("Stop failed: %v", err)
    }
    
    if module.GetStatus() != StatusDisabled {
        t.Error("Module should be disabled")
    }
}
```

## Publishing Your Module

1. Create a separate repository for your module
2. Include comprehensive documentation
3. Provide examples and tests
4. Submit to the Bituncoin module registry
5. Maintain and update regularly

## Support

For module development support:
- GitHub Discussions
- Developer Discord channel
- Example modules in `/addons` directory

## License

Ensure your module is compatible with the Bituncoin license (GPL-3.0).
