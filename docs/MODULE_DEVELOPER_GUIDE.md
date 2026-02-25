# Bituncoin Universal Wallet - Module Developer Guide

## Overview

The Bituncoin Universal Wallet supports a plug-and-play add-on module system that allows developers to extend the wallet's functionality without modifying the core codebase. This guide covers the architecture, development process, and best practices for creating custom modules.

## Module Architecture

### Core Concepts
- **Sandboxing**: Modules run in isolated environments with limited permissions
- **Event-Driven**: Modules communicate through a publish-subscribe event system
- **Permission-Based**: Granular permission system controls module capabilities
- **Version Management**: Semantic versioning with automatic update support

### Module Types
1. **UI Modules**: Add new screens, dialogs, or UI components
2. **Service Modules**: Background services for monitoring, alerts, or automation
3. **Integration Modules**: Third-party service integrations (exchanges, DeFi protocols)
4. **Analysis Modules**: Transaction analysis, portfolio optimization, risk assessment

### Built-in Modules
- **Advanced Staking**: Automated staking strategies across multiple protocols
- **DeFi Lending**: Cross-platform lending and borrowing optimization
- **Portfolio Tracker**: Advanced portfolio analytics and reporting
- **Security Scanner**: Real-time security monitoring and threat detection

## Development Environment Setup

### Prerequisites
- Go 1.21 or later
- Node.js 18.x or later
- Git
- Docker (for testing)

### Project Structure
```
module-name/
├── manifest.json          # Module metadata and configuration
├── main.go               # Main module entry point
├── handlers/             # HTTP handlers and API endpoints
│   └── api.go
├── services/             # Background services
│   └── service.go
├── ui/                   # Frontend components (React)
│   ├── components/
│   └── index.js
├── tests/                # Unit and integration tests
│   ├── main_test.go
│   └── integration_test.go
├── docs/                 # Module documentation
│   └── README.md
└── Dockerfile            # Containerization for sandboxing
```

### Development Tools
```bash
# Install module development CLI
go install github.com/bituncoin/module-cli@latest

# Initialize new module
module-cli init my-module

# Test module
module-cli test

# Build module package
module-cli build
```

## Module Manifest

### Basic Structure
```json
{
  "name": "my-custom-module",
  "version": "1.0.0",
  "description": "A custom module for enhanced wallet functionality",
  "author": "Developer Name",
  "license": "MIT",
  "homepage": "https://github.com/username/my-module",

  "module": {
    "type": "service",
    "entryPoint": "main.go",
    "permissions": [
      "wallet:read",
      "transactions:read",
      "network:access"
    ],
    "dependencies": {
      "go": ">=1.21.0",
      "node": ">=18.0.0"
    }
  },

  "ui": {
    "components": ["DashboardWidget", "SettingsPanel"],
    "routes": [
      {
        "path": "/module/my-module",
        "component": "MainView"
      }
    ]
  },

  "api": {
    "endpoints": [
      {
        "path": "/api/module/my-module/data",
        "method": "GET",
        "handler": "GetData"
      }
    ]
  },

  "events": {
    "subscribes": ["wallet:updated", "transaction:new"],
    "publishes": ["module:my-module:alert"]
  }
}
```

### Permission System

#### Available Permissions
- **wallet:read** - Read wallet balances and addresses
- **wallet:write** - Modify wallet settings (no funds access)
- **transactions:read** - Read transaction history
- **transactions:write** - Create transactions (requires user approval)
- **network:access** - Access external APIs and networks
- **storage:read** - Read module-specific data
- **storage:write** - Write module-specific data
- **ui:inject** - Inject UI components
- **notifications:send** - Send user notifications

#### Permission Request Format
```json
{
  "permissions": [
    {
      "name": "wallet:read",
      "reason": "Required to monitor balance changes for alerts"
    },
    {
      "name": "notifications:send",
      "reason": "Send alerts when conditions are met"
    }
  ]
}
```

## Module Interface

### Core Interface
```go
type Module interface {
    // Initialize the module
    Init(ctx context.Context, config ModuleConfig) error

    // Start the module
    Start() error

    // Stop the module
    Stop() error

    // Get module information
    Info() ModuleInfo

    // Handle events
    HandleEvent(event Event) error

    // Handle API requests
    HandleAPIRequest(req APIRequest) (APIResponse, error)
}
```

### Module Configuration
```go
type ModuleConfig struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Version     string                 `json:"version"`
    Settings    map[string]interface{} `json:"settings"`
    Permissions []string               `json:"permissions"`
    DataDir     string                 `json:"dataDir"`
}
```

### Event Handling
```go
type Event struct {
    Type      string                 `json:"type"`
    Source    string                 `json:"source"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}

// Example event handler
func (m *MyModule) HandleEvent(event Event) error {
    switch event.Type {
    case "wallet:updated":
        return m.handleWalletUpdate(event.Data)
    case "transaction:new":
        return m.handleNewTransaction(event.Data)
    default:
        return nil
    }
}
```

### API Request Handling
```go
type APIRequest struct {
    Method  string                 `json:"method"`
    Path    string                 `json:"path"`
    Headers map[string]string      `json:"headers"`
    Body    map[string]interface{} `json:"body"`
    UserID  string                 `json:"userId"`
}

type APIResponse struct {
    StatusCode int                    `json:"statusCode"`
    Headers    map[string]string      `json:"headers"`
    Body       map[string]interface{} `json:"body"`
}

// Example API handler
func (m *MyModule) HandleAPIRequest(req APIRequest) (APIResponse, error) {
    switch req.Path {
    case "/api/module/my-module/data":
        return m.handleGetData(req)
    case "/api/module/my-module/settings":
        return m.handleUpdateSettings(req)
    default:
        return APIResponse{StatusCode: 404}, nil
    }
}
```

## UI Components

### React Component Structure
```jsx
import React, { useState, useEffect } from 'react';
import { useWallet } from '@bituncoin/wallet-hooks';

const MyModuleWidget = () => {
  const { balance, transactions } = useWallet();
  const [moduleData, setModuleData] = useState(null);

  useEffect(() => {
    // Fetch module-specific data
    fetchModuleData();
  }, []);

  const fetchModuleData = async () => {
    try {
      const response = await fetch('/api/module/my-module/data');
      const data = await response.json();
      setModuleData(data);
    } catch (error) {
      console.error('Failed to fetch module data:', error);
    }
  };

  return (
    <div className="module-widget">
      <h3>My Custom Module</h3>
      <div className="balance-info">
        Current Balance: {balance} BTC
      </div>
      <div className="module-data">
        {moduleData && (
          <pre>{JSON.stringify(moduleData, null, 2)}</pre>
        )}
      </div>
    </div>
  );
};

export default MyModuleWidget;
```

### UI Registration
```javascript
// ui/index.js
import MyModuleWidget from './components/MyModuleWidget';
import SettingsPanel from './components/SettingsPanel';

export const components = {
  MyModuleWidget,
  SettingsPanel
};

export const routes = [
  {
    path: '/module/my-module',
    component: 'MyModuleWidget'
  },
  {
    path: '/module/my-module/settings',
    component: 'SettingsPanel'
  }
];
```

## Testing

### Unit Tests
```go
package mymodule

import (
    "testing"
    "github.com/bituncoin/module-sdk/testing"
)

func TestModuleInit(t *testing.T) {
    config := ModuleConfig{
        ID: "test-module",
        Name: "Test Module",
        Version: "1.0.0",
    }

    module := &MyModule{}
    err := module.Init(context.Background(), config)

    if err != nil {
        t.Fatalf("Failed to initialize module: %v", err)
    }

    info := module.Info()
    if info.Name != "Test Module" {
        t.Errorf("Expected module name 'Test Module', got '%s'", info.Name)
    }
}

func TestHandleEvent(t *testing.T) {
    module := &MyModule{}
    // Initialize module...

    event := Event{
        Type: "wallet:updated",
        Data: map[string]interface{}{
            "balance": 1.5,
            "currency": "BTC",
        },
    }

    err := module.HandleEvent(event)
    if err != nil {
        t.Fatalf("Failed to handle event: %v", err)
    }

    // Assert expected behavior
}
```

### Integration Tests
```go
func TestModuleIntegration(t *testing.T) {
    // Start test wallet instance
    wallet := testing.NewTestWallet()

    // Load and start module
    module, err := testing.LoadModule("my-module", wallet)
    if err != nil {
        t.Fatalf("Failed to load module: %v", err)
    }

    // Test module functionality
    balance, err := module.GetBalance()
    if err != nil {
        t.Fatalf("Failed to get balance: %v", err)
    }

    if balance != 0.0 {
        t.Errorf("Expected initial balance 0.0, got %f", balance)
    }
}
```

## Security Best Practices

### Input Validation
```go
func validateInput(input string) error {
    if len(input) > 1000 {
        return errors.New("input too long")
    }

    // Use allowlist validation
    matched, err := regexp.MatchString(`^[a-zA-Z0-9\s\-_.]+$`, input)
    if err != nil {
        return err
    }
    if !matched {
        return errors.New("invalid characters in input")
    }

    return nil
}
```

### Secure API Calls
```go
func makeSecureAPICall(url string, data interface{}) error {
    // Validate URL
    parsedURL, err := url.Parse(url)
    if err != nil {
        return err
    }

    // Only allow HTTPS
    if parsedURL.Scheme != "https" {
        return errors.New("only HTTPS URLs allowed")
    }

    // Set reasonable timeout
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    // Make request with proper headers
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "Bituncoin-Module/1.0")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}
```

### Data Encryption
```go
func encryptData(data []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}
```

## Deployment

### Module Packaging
```bash
# Build module
module-cli build

# This creates module-name.zip with:
# - Compiled binaries
# - UI assets
# - Manifest file
# - Documentation
```

### Installation
```bash
# Install from registry
wallet-cli module install my-module

# Install from file
wallet-cli module install ./my-module.zip

# Verify installation
wallet-cli module list
```

### Update Process
```bash
# Check for updates
wallet-cli module check-updates

# Update specific module
wallet-cli module update my-module

# Update all modules
wallet-cli module update-all
```

## Distribution

### Module Registry
Publish your module to the official Bituncoin Module Registry:

```bash
# Login to registry
module-cli registry login

# Publish module
module-cli registry publish

# View module stats
module-cli registry stats my-module
```

### Private Distribution
For private or enterprise modules:

```bash
# Create private registry
module-cli registry create-private

# Publish to private registry
module-cli registry publish --registry my-private-registry
```

## Support and Community

### Resources
- **Documentation**: https://docs.bituncoin.com/modules
- **API Reference**: https://api.bituncoin.com/modules
- **Community Forum**: https://forum.bituncoin.com/c/modules
- **GitHub Issues**: https://github.com/bituncoin/modules/issues

### Getting Help
1. Check the documentation
2. Search existing issues
3. Ask in the community forum
4. Create a GitHub issue for bugs
5. Request features via GitHub discussions

### Contributing
We welcome contributions to the module ecosystem:

1. Fork the module SDK
2. Create your feature branch
3. Add tests for new functionality
4. Submit a pull request
5. Participate in code review

## Version Information

- Module SDK Version: 1.0.0
- Compatible Wallet Versions: 1.0.0+
- Go Version Requirements: 1.21+
- Node.js Version Requirements: 18+

For the latest updates and breaking changes, see the [changelog](https://github.com/bituncoin/module-sdk/blob/main/CHANGELOG.md).
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
