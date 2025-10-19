package wallet

import (
	"sync"
	"time"
)

// SystemStatus represents the overall system status
type SystemStatus string

const (
	StatusHealthy   SystemStatus = "healthy"
	StatusDegraded  SystemStatus = "degraded"
	StatusDown      SystemStatus = "down"
	StatusMaintenance SystemStatus = "maintenance"
)

// ComponentStatus represents the status of a system component
type ComponentStatus struct {
	Name         string       `json:"name"`
	Status       SystemStatus `json:"status"`
	LastChecked  time.Time    `json:"lastChecked"`
	ResponseTime int64        `json:"responseTime"` // milliseconds
	ErrorCount   int          `json:"errorCount"`
	Uptime       float64      `json:"uptime"` // percentage
}

// BlockchainNetwork represents a blockchain network connection
type BlockchainNetwork struct {
	Name        string    `json:"name"`
	ChainID     int       `json:"chainId"`
	RPC         string    `json:"rpc"`
	Connected   bool      `json:"connected"`
	BlockHeight int64     `json:"blockHeight"`
	LastSync    time.Time `json:"lastSync"`
}

// SystemMetrics represents overall system metrics
type SystemMetrics struct {
	TotalUsers        int64     `json:"totalUsers"`
	ActiveWallets     int64     `json:"activeWallets"`
	TotalTransactions int64     `json:"totalTransactions"`
	TotalVolume       float64   `json:"totalVolume"`
	AvgResponseTime   int64     `json:"avgResponseTime"`
	ErrorRate         float64   `json:"errorRate"`
	Uptime            float64   `json:"uptime"`
	LastUpdated       time.Time `json:"lastUpdated"`
}

// Dashboard manages system monitoring and operations
type Dashboard struct {
	components  map[string]*ComponentStatus
	networks    map[string]*BlockchainNetwork
	metrics     *SystemMetrics
	alerts      []string
	updateQueue []string
	mutex       sync.RWMutex
}

// NewDashboard creates a new operations dashboard
func NewDashboard() *Dashboard {
	d := &Dashboard{
		components:  make(map[string]*ComponentStatus),
		networks:    make(map[string]*BlockchainNetwork),
		metrics:     &SystemMetrics{},
		alerts:      make([]string, 0),
		updateQueue: make([]string, 0),
	}

	d.initializeComponents()
	d.initializeNetworks()

	return d
}

// initializeComponents sets up system components
func (d *Dashboard) initializeComponents() {
	components := []string{
		"wallet_service",
		"exchange_service",
		"merchant_service",
		"card_service",
		"security_service",
		"ai_service",
		"api_gateway",
		"database",
		"cache",
	}

	for _, name := range components {
		d.components[name] = &ComponentStatus{
			Name:        name,
			Status:      StatusHealthy,
			LastChecked: time.Now(),
			Uptime:      99.9,
		}
	}
}

// initializeNetworks sets up blockchain network connections
func (d *Dashboard) initializeNetworks() {
	networks := []BlockchainNetwork{
		{
			Name:      "Bituncoin",
			ChainID:   1,
			RPC:       "https://rpc.bituncoin.io",
			Connected: true,
		},
		{
			Name:      "Bitcoin",
			ChainID:   0,
			RPC:       "https://btc-rpc.bituncoin.io",
			Connected: true,
		},
		{
			Name:      "Ethereum",
			ChainID:   1,
			RPC:       "https://eth-rpc.bituncoin.io",
			Connected: true,
		},
		{
			Name:      "Binance Smart Chain",
			ChainID:   56,
			RPC:       "https://bsc-rpc.bituncoin.io",
			Connected: true,
		},
	}

	for _, network := range networks {
		d.networks[network.Name] = &network
	}
}

// GetSystemStatus returns the overall system status
func (d *Dashboard) GetSystemStatus() SystemStatus {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	downCount := 0
	degradedCount := 0

	for _, component := range d.components {
		switch component.Status {
		case StatusDown:
			downCount++
		case StatusDegraded:
			degradedCount++
		}
	}

	// If any critical component is down
	if downCount > 0 {
		return StatusDegraded
	}

	// If multiple components are degraded
	if degradedCount > 2 {
		return StatusDegraded
	}

	return StatusHealthy
}

// HealthCheck performs a health check on all components
func (d *Dashboard) HealthCheck() map[string]*ComponentStatus {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, component := range d.components {
		// Simulate health check
		component.LastChecked = time.Now()
		component.ResponseTime = int64(10 + time.Now().UnixNano()%100)
		
		// Update uptime
		if component.ErrorCount > 10 {
			component.Status = StatusDegraded
			component.Uptime = 95.0
		} else {
			component.Status = StatusHealthy
			component.Uptime = 99.9
		}
	}

	return d.components
}

// CheckNetworkConnections checks blockchain network connections
func (d *Dashboard) CheckNetworkConnections() map[string]*BlockchainNetwork {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, network := range d.networks {
		// Simulate network check
		network.LastSync = time.Now()
		network.BlockHeight++
		
		// Randomly simulate connection status (in production, this would be actual checks)
		if time.Now().Unix()%100 > 95 {
			network.Connected = false
		} else {
			network.Connected = true
		}
	}

	return d.networks
}

// GetMetrics returns current system metrics
func (d *Dashboard) GetMetrics() *SystemMetrics {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.metrics
}

// UpdateMetrics updates system metrics
func (d *Dashboard) UpdateMetrics(users, activeWallets, transactions int64, volume float64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.metrics.TotalUsers = users
	d.metrics.ActiveWallets = activeWallets
	d.metrics.TotalTransactions = transactions
	d.metrics.TotalVolume = volume
	d.metrics.LastUpdated = time.Now()

	// Calculate average response time
	totalResponseTime := int64(0)
	count := 0
	for _, component := range d.components {
		totalResponseTime += component.ResponseTime
		count++
	}
	if count > 0 {
		d.metrics.AvgResponseTime = totalResponseTime / int64(count)
	}

	// Calculate error rate
	totalErrors := 0
	for _, component := range d.components {
		totalErrors += component.ErrorCount
	}
	if transactions > 0 {
		d.metrics.ErrorRate = (float64(totalErrors) / float64(transactions)) * 100
	}

	// Calculate overall uptime
	totalUptime := 0.0
	for _, component := range d.components {
		totalUptime += component.Uptime
	}
	d.metrics.Uptime = totalUptime / float64(len(d.components))
}

// AddAlert adds a system alert
func (d *Dashboard) AddAlert(message string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	alert := time.Now().Format("2006-01-02 15:04:05") + " - " + message
	d.alerts = append(d.alerts, alert)

	// Keep only last 100 alerts
	if len(d.alerts) > 100 {
		d.alerts = d.alerts[len(d.alerts)-100:]
	}
}

// GetAlerts returns recent system alerts
func (d *Dashboard) GetAlerts() []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.alerts
}

// ScheduleUpdate schedules a system update
func (d *Dashboard) ScheduleUpdate(component string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.updateQueue = append(d.updateQueue, component)
	d.AddAlert("Update scheduled for " + component)
}

// GetUpdateQueue returns pending updates
func (d *Dashboard) GetUpdateQueue() []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.updateQueue
}

// ExecuteUpdate executes a pending update
func (d *Dashboard) ExecuteUpdate(component string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Mark component as under maintenance
	if comp, exists := d.components[component]; exists {
		comp.Status = StatusMaintenance
	}

	// Remove from update queue
	for i, name := range d.updateQueue {
		if name == component {
			d.updateQueue = append(d.updateQueue[:i], d.updateQueue[i+1:]...)
			break
		}
	}

	d.AddAlert("Update completed for " + component)

	// Restore component status
	if comp, exists := d.components[component]; exists {
		comp.Status = StatusHealthy
	}

	return nil
}

// GetDashboardSummary returns a complete dashboard summary
func (d *Dashboard) GetDashboardSummary() map[string]interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return map[string]interface{}{
		"systemStatus": d.GetSystemStatus(),
		"components":   d.components,
		"networks":     d.networks,
		"metrics":      d.metrics,
		"alerts":       d.alerts[max(0, len(d.alerts)-10):], // Last 10 alerts
		"updateQueue":  d.updateQueue,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
