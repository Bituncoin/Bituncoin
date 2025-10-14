package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Bituncoin/Bituncoin/bituncoin-btn/core"
	"github.com/Bituncoin/Bituncoin/bituncoin-btn/identity"
	"github.com/Bituncoin/Bituncoin/bituncoin-btn/storage"
)

// Node represents a Bituncoin network node
type Node struct {
	Blockchain *core.Blockchain
	Storage    *storage.LevelDB
	Port       int
	Peers      []string
	mu         sync.RWMutex
}

// WalletRequest represents a wallet creation request
type WalletRequest struct {
	Enable2FA       bool   `json:"enable_2fa"`
	EnableBiometric bool   `json:"enable_biometric"`
	BiometricData   string `json:"biometric_data,omitempty"`
}

// TransactionRequest represents a transaction request
type TransactionRequest struct {
	From         string  `json:"from"`
	To           string  `json:"to"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	CrossChain   bool    `json:"cross_chain"`
	TargetChain  string  `json:"target_chain,omitempty"`
	TwoFactorOTP string  `json:"two_factor_otp,omitempty"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewNode creates a new Bituncoin node
func NewNode(port int, dataDir string) (*Node, error) {
	blockchain := core.NewBlockchain()
	
	db, err := storage.NewLevelDB(dataDir)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Blockchain: blockchain,
		Storage:    db,
		Port:       port,
		Peers:      []string{},
	}

	return node, nil
}

// Start starts the node HTTP server
func (n *Node) Start() error {
	http.HandleFunc("/api/wallet/create", n.handleCreateWallet)
	http.HandleFunc("/api/wallet/balance", n.handleGetBalance)
	http.HandleFunc("/api/transaction/send", n.handleSendTransaction)
	http.HandleFunc("/api/transaction/history", n.handleTransactionHistory)
	http.HandleFunc("/api/blockchain/info", n.handleBlockchainInfo)
	http.HandleFunc("/api/mine", n.handleMineBlock)
	http.HandleFunc("/api/currencies", n.handleGetCurrencies)
	http.HandleFunc("/api/crosschain/bridge", n.handleCrossChainBridge)
	
	fmt.Printf("Bituncoin Node starting on port %d...\n", n.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", n.Port), nil)
}

// handleCreateWallet creates a new wallet
func (n *Node) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Create identity
	identity, err := identity.CreateIdentity()
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to create identity",
		})
		return
	}

	// Enable security features if requested
	if req.Enable2FA {
		identity.EnableTwoFactorAuth("secret-" + identity.Address)
	}

	if req.EnableBiometric && req.BiometricData != "" {
		identity.SetBiometricHash(req.BiometricData)
	}

	// Save wallet to storage
	walletData := &storage.WalletData{
		Address:       identity.Address,
		Balances:      make(map[string]float64),
		Transactions:  []string{},
		EncryptedKeys: identity.PrivateKey,
		AuthConfig: map[string]string{
			"2fa_enabled":       fmt.Sprintf("%t", req.Enable2FA),
			"biometric_enabled": fmt.Sprintf("%t", req.EnableBiometric),
		},
	}

	if err := n.Storage.SaveWallet(identity.Address, walletData); err != nil {
		sendResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Message: "Failed to save wallet",
		})
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Wallet created successfully",
		Data: map[string]interface{}{
			"address":    identity.Address,
			"public_key": identity.PublicKey,
			"private_key": identity.PrivateKey,
			"created_at": identity.CreatedAt,
		},
	})
}

// handleGetBalance returns wallet balance
func (n *Node) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Address parameter required",
		})
		return
	}

	balances := make(map[string]float64)
	for _, currency := range n.Blockchain.SupportedChains {
		balance := n.Blockchain.GetBalance(address, currency)
		balances[currency] = balance
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Balance retrieved",
		Data:    balances,
	})
}

// handleSendTransaction processes a transaction
func (n *Node) handleSendTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Create transaction
	tx := core.Transaction{
		ID:          fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		From:        req.From,
		To:          req.To,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Timestamp:   time.Now().Unix(),
		CrossChain:  req.CrossChain,
		TargetChain: req.TargetChain,
	}

	// Add to blockchain
	n.Blockchain.AddTransaction(tx)

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Transaction added to pending pool",
		Data: map[string]interface{}{
			"transaction_id": tx.ID,
			"status":         "pending",
		},
	})
}

// handleTransactionHistory returns transaction history
func (n *Node) handleTransactionHistory(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Address parameter required",
		})
		return
	}

	var transactions []core.Transaction
	for _, block := range n.Blockchain.Blocks {
		for _, tx := range block.Transactions {
			if tx.From == address || tx.To == address {
				transactions = append(transactions, tx)
			}
		}
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Transaction history retrieved",
		Data:    transactions,
	})
}

// handleBlockchainInfo returns blockchain information
func (n *Node) handleBlockchainInfo(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Blockchain info retrieved",
		Data: map[string]interface{}{
			"blocks":            len(n.Blockchain.Blocks),
			"pending_txs":       len(n.Blockchain.PendingTxs),
			"difficulty":        n.Blockchain.Difficulty,
			"mining_reward":     n.Blockchain.MiningReward,
			"supported_chains":  n.Blockchain.SupportedChains,
			"chain_valid":       n.Blockchain.IsChainValid(),
		},
	})
}

// handleMineBlock mines a new block
func (n *Node) handleMineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	minerAddress := r.URL.Query().Get("miner")
	if minerAddress == "" {
		sendResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Miner address required",
		})
		return
	}

	block := n.Blockchain.MineBlock(minerAddress)

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Block mined successfully",
		Data:    block,
	})
}

// handleGetCurrencies returns supported currencies
func (n *Node) handleGetCurrencies(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Supported currencies",
		Data:    n.Blockchain.SupportedChains,
	})
}

// handleCrossChainBridge handles cross-chain transactions
func (n *Node) handleCrossChainBridge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendResponse(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	sendResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Cross-chain bridge transaction initiated",
		Data: map[string]interface{}{
			"status": "bridging",
			"info":   "Transaction will be processed across chains",
		},
	})
}

// sendResponse sends a JSON response
func sendResponse(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
