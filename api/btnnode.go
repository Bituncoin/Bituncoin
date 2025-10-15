package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Node represents a blockchain node API server
type Node struct {
	Port       int
	Host       string
	IsRunning  bool
	mutex      sync.RWMutex
	endpoints  map[string]http.HandlerFunc
}

// NodeInfo represents node information
type NodeInfo struct {
	Version     string `json:"version"`
	Network     string `json:"network"`
	NodeType    string `json:"nodeType"`
	IsRunning   bool   `json:"isRunning"`
	BlockHeight int    `json:"blockHeight"`
}

// NewNode creates a new API node
func NewNode(host string, port int) *Node {
	return &Node{
		Port:      port,
		Host:      host,
		IsRunning: false,
		endpoints: make(map[string]http.HandlerFunc),
	}
}

// Start starts the API node server
func (n *Node) Start() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.IsRunning {
		return fmt.Errorf("node already running")
	}

	// Register default endpoints
	n.registerEndpoints()

	// Start HTTP server
	go func() {
		mux := http.NewServeMux()
		for path, handler := range n.endpoints {
			mux.HandleFunc(path, handler)
		}

		addr := fmt.Sprintf("%s:%d", n.Host, n.Port)
		http.ListenAndServe(addr, mux)
	}()

	n.IsRunning = true
	return nil
}

// Stop stops the API node server
func (n *Node) Stop() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if !n.IsRunning {
		return fmt.Errorf("node not running")
	}

	n.IsRunning = false
	return nil
}

// registerEndpoints registers API endpoints
func (n *Node) registerEndpoints() {
	n.endpoints["/api/info"] = n.handleInfo
	n.endpoints["/api/health"] = n.handleHealth
	n.endpoints["/api/goldcoin/balance"] = n.handleBalance
	n.endpoints["/api/goldcoin/send"] = n.handleSend
	n.endpoints["/api/goldcoin/stake"] = n.handleStake
	n.endpoints["/api/goldcoin/validators"] = n.handleValidators
}

// handleInfo returns node information
func (n *Node) handleInfo(w http.ResponseWriter, r *http.Request) {
	info := NodeInfo{
		Version:     "1.0.0",
		Network:     "bituncoin-mainnet",
		NodeType:    "full-node",
		IsRunning:   n.IsRunning,
		BlockHeight: 0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// handleHealth returns health status
func (n *Node) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"running": n.IsRunning,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleBalance handles balance queries
func (n *Node) handleBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	
	response := map[string]interface{}{
		"address": address,
		"balance": 0.0,
		"staked":  0.0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSend handles transaction creation
func (n *Node) handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var txData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&txData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status":        "pending",
		"transactionId": "tx_123456789",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleStake handles staking operations
func (n *Node) handleStake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var stakeData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&stakeData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"staked": stakeData["amount"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleValidators returns validator information
func (n *Node) handleValidators(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"validators": []map[string]interface{}{
			{
				"address": "GLDvalidator1...",
				"stake":   10000.0,
				"active":  true,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNodeInfo returns current node information
func (n *Node) GetNodeInfo() NodeInfo {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	return NodeInfo{
		Version:     "1.0.0",
		Network:     "bituncoin-mainnet",
		NodeType:    "full-node",
		IsRunning:   n.IsRunning,
		BlockHeight: 0,
	}
}
