package wallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// PaymentMethod represents supported payment methods
type PaymentMethod string

const (
	PaymentMethodQRCode      PaymentMethod = "qr_code"
	PaymentMethodNFC         PaymentMethod = "nfc"
	PaymentMethodMTN         PaymentMethod = "mtn_mobile_money"
	PaymentMethodAirtel      PaymentMethod = "airtel_mobile_money"
	PaymentMethodCard        PaymentMethod = "card"
	PaymentMethodCrypto      PaymentMethod = "crypto"
)

// Merchant represents a registered merchant
type Merchant struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	BusinessType    string    `json:"businessType"`
	WalletAddress   string    `json:"walletAddress"`
	APIKey          string    `json:"apiKey"`
	Status          string    `json:"status"`
	PaymentMethods  []PaymentMethod `json:"paymentMethods"`
	TotalRevenue    float64   `json:"totalRevenue"`
	TransactionCount int64    `json:"transactionCount"`
	RegisteredAt    int64     `json:"registeredAt"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	ID              string            `json:"id"`
	MerchantID      string            `json:"merchantId"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	Description     string            `json:"description"`
	PaymentMethods  []PaymentMethod   `json:"paymentMethods"`
	QRCode          string            `json:"qrCode,omitempty"`
	NFCData         string            `json:"nfcData,omitempty"`
	Status          TransactionStatus `json:"status"`
	CreatedAt       int64             `json:"createdAt"`
	ExpiresAt       int64             `json:"expiresAt"`
	PaidAt          int64             `json:"paidAt,omitempty"`
	CustomerID      string            `json:"customerId,omitempty"`
	TransactionHash string            `json:"transactionHash,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// MobileMoneyPayment represents a mobile money payment
type MobileMoneyPayment struct {
	ID            string            `json:"id"`
	Provider      PaymentMethod     `json:"provider"` // MTN or Airtel
	PhoneNumber   string            `json:"phoneNumber"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Reference     string            `json:"reference"`
	Status        TransactionStatus `json:"status"`
	TransactionID string            `json:"transactionId,omitempty"`
	CreatedAt     int64             `json:"createdAt"`
	CompletedAt   int64             `json:"completedAt,omitempty"`
}

// MerchantService manages merchant operations
type MerchantService struct {
	merchants       map[string]*Merchant
	paymentRequests map[string]*PaymentRequest
	mobilePayments  map[string]*MobileMoneyPayment
	mutex           sync.RWMutex
}

// NewMerchantService creates a new merchant service
func NewMerchantService() *MerchantService {
	return &MerchantService{
		merchants:       make(map[string]*Merchant),
		paymentRequests: make(map[string]*PaymentRequest),
		mobilePayments:  make(map[string]*MobileMoneyPayment),
	}
}

// RegisterMerchant registers a new merchant
func (ms *MerchantService) RegisterMerchant(name, email, phone, businessType, walletAddress string, paymentMethods []PaymentMethod) (*Merchant, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if name == "" || email == "" || walletAddress == "" {
		return nil, errors.New("name, email, and wallet address are required")
	}

	merchant := &Merchant{
		ID:               generateMerchantID(),
		Name:             name,
		Email:            email,
		Phone:            phone,
		BusinessType:     businessType,
		WalletAddress:    walletAddress,
		APIKey:           generateAPIKey(),
		Status:           "active",
		PaymentMethods:   paymentMethods,
		TotalRevenue:     0,
		TransactionCount: 0,
		RegisteredAt:     time.Now().Unix(),
	}

	ms.merchants[merchant.ID] = merchant

	return merchant, nil
}

// GetMerchant returns a merchant by ID
func (ms *MerchantService) GetMerchant(merchantID string) (*Merchant, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	merchant, exists := ms.merchants[merchantID]
	if !exists {
		return nil, errors.New("merchant not found")
	}

	return merchant, nil
}

// CreatePaymentRequest creates a new payment request
func (ms *MerchantService) CreatePaymentRequest(merchantID string, amount float64, currency, description string, paymentMethods []PaymentMethod, expiryMinutes int) (*PaymentRequest, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	merchant, exists := ms.merchants[merchantID]
	if !exists {
		return nil, errors.New("merchant not found")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	now := time.Now().Unix()
	expiresAt := now + int64(expiryMinutes*60)

	pr := &PaymentRequest{
		ID:             generatePaymentRequestID(),
		MerchantID:     merchantID,
		Amount:         amount,
		Currency:       currency,
		Description:    description,
		PaymentMethods: paymentMethods,
		Status:         TxStatusPending,
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
		Metadata:       make(map[string]string),
	}

	// Generate QR code data
	if containsPaymentMethod(paymentMethods, PaymentMethodQRCode) {
		pr.QRCode = generateQRCode(merchant.WalletAddress, amount, currency, pr.ID)
	}

	// Generate NFC data
	if containsPaymentMethod(paymentMethods, PaymentMethodNFC) {
		pr.NFCData = generateNFCData(merchant.WalletAddress, amount, currency, pr.ID)
	}

	ms.paymentRequests[pr.ID] = pr

	return pr, nil
}

// GetPaymentRequest returns a payment request
func (ms *MerchantService) GetPaymentRequest(requestID string) (*PaymentRequest, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	pr, exists := ms.paymentRequests[requestID]
	if !exists {
		return nil, errors.New("payment request not found")
	}

	// Check if expired
	if time.Now().Unix() > pr.ExpiresAt && pr.Status == TxStatusPending {
		pr.Status = TxStatusFailed
	}

	return pr, nil
}

// ProcessPayment processes a payment for a payment request
func (ms *MerchantService) ProcessPayment(requestID, customerID, transactionHash string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	pr, exists := ms.paymentRequests[requestID]
	if !exists {
		return errors.New("payment request not found")
	}

	if pr.Status != TxStatusPending {
		return errors.New("payment request is not pending")
	}

	if time.Now().Unix() > pr.ExpiresAt {
		pr.Status = TxStatusFailed
		return errors.New("payment request has expired")
	}

	pr.Status = TxStatusCompleted
	pr.PaidAt = time.Now().Unix()
	pr.CustomerID = customerID
	pr.TransactionHash = transactionHash

	// Update merchant statistics
	merchant := ms.merchants[pr.MerchantID]
	merchant.TotalRevenue += pr.Amount
	merchant.TransactionCount++

	return nil
}

// InitiateMobileMoneyPayment initiates a mobile money payment
func (ms *MerchantService) InitiateMobileMoneyPayment(provider PaymentMethod, phoneNumber string, amount float64, currency, reference string) (*MobileMoneyPayment, error) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if provider != PaymentMethodMTN && provider != PaymentMethodAirtel {
		return nil, errors.New("invalid mobile money provider")
	}

	if phoneNumber == "" || amount <= 0 {
		return nil, errors.New("phone number and amount are required")
	}

	payment := &MobileMoneyPayment{
		ID:          generateMobilePaymentID(),
		Provider:    provider,
		PhoneNumber: phoneNumber,
		Amount:      amount,
		Currency:    currency,
		Reference:   reference,
		Status:      TxStatusPending,
		CreatedAt:   time.Now().Unix(),
	}

	ms.mobilePayments[payment.ID] = payment

	// In production, this would:
	// 1. Call MTN/Airtel API to initiate payment
	// 2. Return payment confirmation details
	// 3. Set up webhook for payment status updates

	return payment, nil
}

// CompleteMobileMoneyPayment completes a mobile money payment
func (ms *MerchantService) CompleteMobileMoneyPayment(paymentID, transactionID string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	payment, exists := ms.mobilePayments[paymentID]
	if !exists {
		return errors.New("mobile money payment not found")
	}

	if payment.Status != TxStatusPending {
		return errors.New("payment is not pending")
	}

	payment.Status = TxStatusCompleted
	payment.TransactionID = transactionID
	payment.CompletedAt = time.Now().Unix()

	return nil
}

// GetMobileMoneyPayment returns a mobile money payment
func (ms *MerchantService) GetMobileMoneyPayment(paymentID string) (*MobileMoneyPayment, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	payment, exists := ms.mobilePayments[paymentID]
	if !exists {
		return nil, errors.New("mobile money payment not found")
	}

	return payment, nil
}

// GetMerchantPaymentRequests returns all payment requests for a merchant
func (ms *MerchantService) GetMerchantPaymentRequests(merchantID string, status *TransactionStatus) []*PaymentRequest {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	requests := make([]*PaymentRequest, 0)
	for _, pr := range ms.paymentRequests {
		if pr.MerchantID == merchantID {
			if status == nil || pr.Status == *status {
				requests = append(requests, pr)
			}
		}
	}

	return requests
}

// Helper functions

func generateMerchantID() string {
	return "merchant_" + time.Now().Format("20060102150405")
}

func generatePaymentRequestID() string {
	return "pr_" + time.Now().Format("20060102150405")
}

func generateMobilePaymentID() string {
	return "mp_" + time.Now().Format("20060102150405")
}

func generateAPIKey() string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return "sk_" + hex.EncodeToString(hash[:16])
}

func generateQRCode(walletAddress string, amount float64, currency, requestID string) string {
	// QR code data format: wallet:address:amount:currency:requestId
	data := fmt.Sprintf("%s:%s:%.2f:%s:%s", "btn", walletAddress, amount, currency, requestID)
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded
}

func generateNFCData(walletAddress string, amount float64, currency, requestID string) string {
	// NFC data format (NDEF)
	nfcData := map[string]interface{}{
		"type":      "payment",
		"wallet":    walletAddress,
		"amount":    amount,
		"currency":  currency,
		"requestId": requestID,
		"timestamp": time.Now().Unix(),
	}

	jsonData, _ := json.Marshal(nfcData)
	return base64.StdEncoding.EncodeToString(jsonData)
}

func containsPaymentMethod(methods []PaymentMethod, method PaymentMethod) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}
