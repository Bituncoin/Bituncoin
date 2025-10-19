package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Bituncoin/Bituncoin/payments"
	"github.com/Bituncoin/Bituncoin/wallet"
)

// Node represents a blockchain node API server
type Node struct {
	Port          int
	Host          string
	IsRunning     bool
	mutex         sync.RWMutex
	endpoints     map[string]http.HandlerFunc
	payments      *payments.BtnPay
	mobileMoney   *payments.MobileMoneyService
	qrCode        *payments.QRCodeService
	cards         *payments.CardService
	aiAssistant   *wallet.AIAssistant
	exchange      *wallet.ExchangeService
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
		Port:        port,
		Host:        host,
		IsRunning:   false,
		endpoints:   make(map[string]http.HandlerFunc),
		payments:    payments.NewBtnPay(),
		mobileMoney: payments.NewMobileMoneyService(),
		qrCode:      payments.NewQRCodeService(),
		cards:       payments.NewCardService(),
		aiAssistant: wallet.NewAIAssistant(),
		exchange:    wallet.NewExchangeService(),
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
	// Core endpoints
	n.endpoints["/api/info"] = n.handleInfo
	n.endpoints["/api/health"] = n.handleHealth
	n.endpoints["/api/goldcoin/balance"] = n.handleBalance
	n.endpoints["/api/goldcoin/send"] = n.handleSend
	n.endpoints["/api/goldcoin/stake"] = n.handleStake
	n.endpoints["/api/goldcoin/validators"] = n.handleValidators

	// BTN-PAY endpoints
	n.endpoints["/api/btnpay/invoice"] = n.payments.CreateInvoiceHandler
	n.endpoints["/api/btnpay/invoice/"] = n.payments.GetInvoiceHandler
	n.endpoints["/api/btnpay/pay"] = n.payments.PayInvoiceHandler

	// Mobile Money endpoints
	n.endpoints["/api/mobilemoney/pay"] = n.handleMobileMoneyPayment
	n.endpoints["/api/mobilemoney/status"] = n.handleMobileMoneyStatus

	// QR Code endpoints
	n.endpoints["/api/qrcode/generate"] = n.handleGenerateQR
	n.endpoints["/api/qrcode/parse"] = n.handleParseQR

	// Card endpoints
	n.endpoints["/api/card/issue"] = n.handleIssueCard
	n.endpoints["/api/card/load"] = n.handleLoadCard
	n.endpoints["/api/card/payment"] = n.handleCardPayment
	n.endpoints["/api/card/transactions"] = n.handleCardTransactions

	// Exchange endpoints
	n.endpoints["/api/exchange/rates"] = n.handleExchangeRates
	n.endpoints["/api/exchange/swap"] = n.handleSwap
	n.endpoints["/api/exchange/estimate"] = n.handleSwapEstimate

	// AI Assistant endpoints
	n.endpoints["/api/ai/insights"] = n.handleAIInsights
	n.endpoints["/api/ai/alerts"] = n.handleAIAlerts
	n.endpoints["/api/ai/ask"] = n.handleAIAsk
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
// Mobile Money handlers

func (n *Node) handleMobileMoneyPayment(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
Provider    string  `json:"provider"`
PhoneNumber string  `json:"phoneNumber"`
Amount      float64 `json:"amount"`
Currency    string  `json:"currency"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

tx, err := n.mobileMoney.InitiatePayment(
payments.MobileMoneyProvider(req.Provider),
req.PhoneNumber,
req.Amount,
req.Currency,
)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(tx)
}

func (n *Node) handleMobileMoneyStatus(w http.ResponseWriter, r *http.Request) {
txID := r.URL.Query().Get("id")
if txID == "" {
http.Error(w, "transaction ID required", http.StatusBadRequest)
return
}

tx, err := n.mobileMoney.GetTransaction(txID)
if err != nil {
http.Error(w, err.Error(), http.StatusNotFound)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(tx)
}

// QR Code handlers

func (n *Node) handleGenerateQR(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
Address    string  `json:"address"`
Amount     float64 `json:"amount"`
Currency   string  `json:"currency"`
Memo       string  `json:"memo"`
MerchantID string  `json:"merchantId"`
InvoiceID  string  `json:"invoiceId"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

var qrData string
var err error

if req.MerchantID != "" {
qrData, err = n.qrCode.GenerateMerchantQR(req.MerchantID, req.Address, req.Amount, req.Currency, req.InvoiceID)
} else if req.Amount > 0 {
qrData, err = n.qrCode.GeneratePaymentQR(req.Address, req.Amount, req.Currency, req.Memo)
} else {
qrData, err = n.qrCode.GenerateReceiveQR(req.Address, req.Currency)
}

if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

response := map[string]string{
"qrData": qrData,
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
}

func (n *Node) handleParseQR(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
QRData string `json:"qrData"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

data, err := n.qrCode.ParsePaymentQR(req.QRData)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
}

// Card handlers

func (n *Node) handleIssueCard(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
HolderName    string `json:"holderName"`
LinkedAddress string `json:"linkedAddress"`
CardType      string `json:"cardType"`
Category      string `json:"category"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

card, err := n.cards.IssueCard(
req.HolderName,
req.LinkedAddress,
payments.CardType(req.CardType),
payments.CardCategory(req.Category),
)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(card)
}

func (n *Node) handleLoadCard(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
CardID string  `json:"cardId"`
Amount float64 `json:"amount"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

err := n.cards.LoadCard(req.CardID, req.Amount)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

response := map[string]string{
"status": "success",
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
}

func (n *Node) handleCardPayment(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
CardID   string  `json:"cardId"`
Amount   float64 `json:"amount"`
Merchant string  `json:"merchant"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

tx, err := n.cards.ProcessPayment(req.CardID, req.Amount, req.Merchant)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(tx)
}

func (n *Node) handleCardTransactions(w http.ResponseWriter, r *http.Request) {
cardID := r.URL.Query().Get("cardId")
if cardID == "" {
http.Error(w, "card ID required", http.StatusBadRequest)
return
}

transactions, err := n.cards.GetCardTransactions(cardID)
if err != nil {
http.Error(w, err.Error(), http.StatusNotFound)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(transactions)
}

// Exchange handlers

func (n *Node) handleExchangeRates(w http.ResponseWriter, r *http.Request) {
from := r.URL.Query().Get("from")
to := r.URL.Query().Get("to")

if from == "" || to == "" {
http.Error(w, "from and to currencies required", http.StatusBadRequest)
return
}

rate, err := n.exchange.GetExchangeRate(from, to)
if err != nil {
http.Error(w, err.Error(), http.StatusNotFound)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(rate)
}

func (n *Node) handleSwap(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
From        string  `json:"from"`
To          string  `json:"to"`
Amount      float64 `json:"amount"`
UserAddress string  `json:"userAddress"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

tx, err := n.exchange.ExecuteSwap(req.From, req.To, req.Amount, req.UserAddress)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(tx)
}

func (n *Node) handleSwapEstimate(w http.ResponseWriter, r *http.Request) {
from := r.URL.Query().Get("from")
to := r.URL.Query().Get("to")
amount := 0.0
fmt.Sscanf(r.URL.Query().Get("amount"), "%f", &amount)

if from == "" || to == "" || amount <= 0 {
http.Error(w, "from, to, and amount required", http.StatusBadRequest)
return
}

toAmount, fee, err := n.exchange.CalculateSwap(from, to, amount)
if err != nil {
http.Error(w, err.Error(), http.StatusBadRequest)
return
}

response := map[string]interface{}{
"from":     from,
"to":       to,
"amount":   amount,
"toAmount": toAmount,
"fee":      fee,
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
}

// AI Assistant handlers

func (n *Node) handleAIInsights(w http.ResponseWriter, r *http.Request) {
// In production, get balances from user account
balances := map[string]float64{
"BTN":  5000.0,
"GLD":  1250.5,
"BTC":  0.05,
"ETH":  2.3,
"USDT": 1000.0,
"BNB":  10.5,
}

insights := n.aiAssistant.AnalyzePortfolio(balances)

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(insights)
}

func (n *Node) handleAIAlerts(w http.ResponseWriter, r *http.Request) {
// In production, get actual security settings
securitySettings := map[string]bool{
"twoFactorEnabled": true,
"fraudMonitoring":  true,
}

alerts := n.aiAssistant.GenerateSecurityAlerts(securitySettings, 15, false)

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(alerts)
}

func (n *Node) handleAIAsk(w http.ResponseWriter, r *http.Request) {
if r.Method != http.MethodPost {
http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
return
}

var req struct {
Query string `json:"query"`
}

if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "Invalid request body", http.StatusBadRequest)
return
}

context := make(map[string]interface{})
recommendation := n.aiAssistant.ProvideRecommendation(req.Query, context)

response := map[string]string{
"query":          req.Query,
"recommendation": recommendation,
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
}
