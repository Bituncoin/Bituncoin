package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Bituncoin/Bituncoin/addons"
	"github.com/Bituncoin/Bituncoin/auth"
	"github.com/Bituncoin/Bituncoin/payments"
)

// Node represents a blockchain node API server
type Node struct {
	Port       int
	Host       string
	IsRunning  bool
	mutex      sync.RWMutex
	endpoints  map[string]http.HandlerFunc
	payments   *payments.BtnPay
	accounts   *auth.AccountManager
	addons     *addons.ModuleRegistry
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
		payments:  payments.NewBtnPay(),
		accounts:  auth.NewAccountManager(),
		addons:    addons.NewModuleRegistry(),
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

	// BTN-PAY endpoints
	n.endpoints["/api/btnpay/invoice"] = n.payments.CreateInvoiceHandler
	// register a path prefix for invoice lookups — the handler extracts the last segment as ID
	n.endpoints["/api/btnpay/invoice/"] = n.payments.GetInvoiceHandler
	n.endpoints["/api/btnpay/pay"] = n.payments.PayInvoiceHandler
	
	// Authentication endpoints
	n.endpoints["/api/auth/register"] = n.handleRegister
	n.endpoints["/api/auth/login"] = n.handleLogin
	n.endpoints["/api/auth/logout"] = n.handleLogout
	n.endpoints["/api/auth/validate"] = n.handleValidateSession
	
	// User management endpoints (admin only)
	n.endpoints["/api/users/list"] = n.handleListUsers
	n.endpoints["/api/users/update-role"] = n.handleUpdateUserRole
	n.endpoints["/api/users/deactivate"] = n.handleDeactivateUser
	
	// Add-on module endpoints
	n.endpoints["/api/addons/list"] = n.handleListAddons
	n.endpoints["/api/addons/enable"] = n.handleEnableAddon
	n.endpoints["/api/addons/disable"] = n.handleDisableAddon
	n.endpoints["/api/addons/execute"] = n.handleExecuteAddon
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
			{
				"address": "GLDvalidator2...",
				"stake":   20000.0,
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

// handleRegister handles user registration
func (n *Node) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Username string    `json:"username"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Role     auth.Role `json:"role"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	user, err := n.accounts.CreateUser(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// handleLogin handles user authentication
func (n *Node) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	session, err := n.accounts.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// handleLogout handles user logout
func (n *Node) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}
	
	if err := n.accounts.Logout(sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleValidateSession validates a user session
func (n *Node) handleValidateSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}
	
	user, err := n.accounts.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// handleListUsers lists all users (admin only)
func (n *Node) handleListUsers(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	user, err := n.accounts.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	if !n.accounts.HasPermission(user.ID, auth.PermissionManageUsers) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	users := n.accounts.ListUsers()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// handleUpdateUserRole updates a user's role (admin only)
func (n *Node) handleUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	user, err := n.accounts.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	if !n.accounts.HasPermission(user.ID, auth.PermissionManageUsers) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	var req struct {
		UserID  string    `json:"userId"`
		NewRole auth.Role `json:"newRole"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := n.accounts.UpdateUserRole(req.UserID, req.NewRole); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleDeactivateUser deactivates a user account (admin only)
func (n *Node) handleDeactivateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	user, err := n.accounts.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	if !n.accounts.HasPermission(user.ID, auth.PermissionManageUsers) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	
	var req struct {
		UserID string `json:"userId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := n.accounts.DeactivateUser(req.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleListAddons lists all registered add-ons
func (n *Node) handleListAddons(w http.ResponseWriter, r *http.Request) {
	modules := n.addons.ListModules()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

// handleEnableAddon enables an add-on module
func (n *Node) handleEnableAddon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Name   string                 `json:"name"`
		Config map[string]interface{} `json:"config"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := n.addons.Enable(req.Name, req.Config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleDisableAddon disables an add-on module
func (n *Node) handleDisableAddon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := n.addons.Disable(req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleExecuteAddon executes an add-on module action
func (n *Node) handleExecuteAddon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		Name   string                 `json:"name"`
		Action string                 `json:"action"`
		Params map[string]interface{} `json:"params"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result, err := n.addons.Execute(req.Name, req.Action, req.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
