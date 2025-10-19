package wallet

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrMerchantNotFound = errors.New("merchant not found")
	ErrInvalidQRCode    = errors.New("invalid qr code")
)

// PaymentMethod represents the method of payment
type PaymentMethod string

const (
	PaymentQRCode      PaymentMethod = "qr_code"
	PaymentNFC         PaymentMethod = "nfc"
	PaymentWallet      PaymentMethod = "wallet_transfer"
	PaymentMobileMoney PaymentMethod = "mobile_money"
)

// MobileMoneyProvider represents mobile money service providers
type MobileMoneyProvider string

const (
	ProviderMTN      MobileMoneyProvider = "mtn_mobile_money"
	ProviderVodafone MobileMoneyProvider = "vodafone_cash"
	ProviderAirtel   MobileMoneyProvider = "airtel_money"
	ProviderTigo     MobileMoneyProvider = "tigo_cash"
)

// Merchant represents a registered merchant
type Merchant struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	WalletAddress    string    `json:"walletAddress"`
	Email            string    `json:"email"`
	BusinessType     string    `json:"businessType"`
	AcceptedMethods  []PaymentMethod `json:"acceptedMethods"`
	AcceptedAssets   []string  `json:"acceptedAssets"`
	RegistrationDate time.Time `json:"registrationDate"`
	Status           string    `json:"status"`
	TotalReceived    float64   `json:"totalReceived"`
}

// PaymentRequest represents a merchant payment request
type PaymentRequest struct {
	ID             string        `json:"id"`
	MerchantID     string        `json:"merchantId"`
	Amount         float64       `json:"amount"`
	Asset          string        `json:"asset"`
	PaymentMethod  PaymentMethod `json:"paymentMethod"`
	QRCode         string        `json:"qrCode,omitempty"`
	NFCToken       string        `json:"nfcToken,omitempty"`
	Status         string        `json:"status"`
	Description    string        `json:"description"`
	CreatedAt      time.Time     `json:"createdAt"`
	ExpiresAt      time.Time     `json:"expiresAt"`
	CompletedAt    time.Time     `json:"completedAt,omitempty"`
	CustomerAddress string       `json:"customerAddress,omitempty"`
	TxHash         string        `json:"txHash,omitempty"`
}

// MobileMoneyPayment represents a mobile money payment
type MobileMoneyPayment struct {
	ID              string              `json:"id"`
	Provider        MobileMoneyProvider `json:"provider"`
	PhoneNumber     string              `json:"phoneNumber"`
	Amount          float64             `json:"amount"`
	Currency        string              `json:"currency"`
	MerchantID      string              `json:"merchantId"`
	Status          string              `json:"status"`
	TransactionRef  string              `json:"transactionRef"`
	Timestamp       time.Time           `json:"timestamp"`
}

// MerchantService manages merchant payment operations
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
func (ms *MerchantService) RegisterMerchant(name, walletAddress, email, businessType string) (*Merchant, error) {
	if name == "" || walletAddress == "" {
		return nil, errors.New("name and wallet address required")
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	merchant := &Merchant{
		ID:               fmt.Sprintf("MERCH-%d", time.Now().UnixNano()),
		Name:             name,
		WalletAddress:    walletAddress,
		Email:            email,
		BusinessType:     businessType,
		AcceptedMethods:  []PaymentMethod{PaymentQRCode, PaymentNFC, PaymentWallet},
		AcceptedAssets:   []string{"BTNG", "GLD", "BTC", "ETH", "USDT"},
		RegistrationDate: time.Now(),
		Status:           "active",
		TotalReceived:    0.0,
	}

	ms.merchants[merchant.ID] = merchant
	return merchant, nil
}

// GetMerchant retrieves a merchant by ID
func (ms *MerchantService) GetMerchant(merchantID string) (*Merchant, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	merchant, exists := ms.merchants[merchantID]
	if !exists {
		return nil, ErrMerchantNotFound
	}

	return merchant, nil
}

// CreatePaymentRequest creates a new payment request
func (ms *MerchantService) CreatePaymentRequest(merchantID string, amount float64, asset string, method PaymentMethod, description string) (*PaymentRequest, error) {
	merchant, err := ms.GetMerchant(merchantID)
	if err != nil {
		return nil, err
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	request := &PaymentRequest{
		ID:            fmt.Sprintf("PAY-%d", time.Now().UnixNano()),
		MerchantID:    merchantID,
		Amount:        amount,
		Asset:         asset,
		PaymentMethod: method,
		Status:        "pending",
		Description:   description,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(15 * time.Minute),
	}

	// Generate QR code data or NFC token
	if method == PaymentQRCode {
		request.QRCode = ms.generateQRCode(merchant.WalletAddress, amount, asset, request.ID)
	} else if method == PaymentNFC {
		request.NFCToken = ms.generateNFCToken()
	}

	ms.paymentRequests[request.ID] = request
	return request, nil
}

// GetPaymentRequest retrieves a payment request
func (ms *MerchantService) GetPaymentRequest(requestID string) (*PaymentRequest, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	request, exists := ms.paymentRequests[requestID]
	if !exists {
		return nil, errors.New("payment request not found")
	}

	// Check expiration
	if time.Now().After(request.ExpiresAt) && request.Status == "pending" {
		request.Status = "expired"
	}

	return request, nil
}

// CompletePaymentRequest marks a payment request as completed
func (ms *MerchantService) CompletePaymentRequest(requestID, customerAddress, txHash string) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	request, exists := ms.paymentRequests[requestID]
	if !exists {
		return errors.New("payment request not found")
	}

	if request.Status != "pending" {
		return errors.New("payment request not pending")
	}

	request.Status = "completed"
	request.CustomerAddress = customerAddress
	request.TxHash = txHash
	request.CompletedAt = time.Now()

	// Update merchant total
	if merchant, exists := ms.merchants[request.MerchantID]; exists {
		merchant.TotalReceived += request.Amount
	}

	return nil
}

// ProcessMobileMoneyPayment processes a mobile money payment
func (ms *MerchantService) ProcessMobileMoneyPayment(merchantID string, provider MobileMoneyProvider, phoneNumber string, amount float64, currency string) (*MobileMoneyPayment, error) {
	merchant, err := ms.GetMerchant(merchantID)
	if err != nil {
		return nil, err
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	payment := &MobileMoneyPayment{
		ID:             fmt.Sprintf("MM-%d", time.Now().UnixNano()),
		Provider:       provider,
		PhoneNumber:    phoneNumber,
		Amount:         amount,
		Currency:       currency,
		MerchantID:     merchantID,
		Status:         "pending",
		TransactionRef: ms.generateTransactionRef(),
		Timestamp:      time.Now(),
	}

	ms.mobilePayments[payment.ID] = payment

	// In production, this would integrate with actual mobile money APIs
	// For now, simulate successful payment
	payment.Status = "completed"
	merchant.TotalReceived += amount

	return payment, nil
}

// GetMerchantPayments retrieves all payments for a merchant
func (ms *MerchantService) GetMerchantPayments(merchantID string) []*PaymentRequest {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	payments := make([]*PaymentRequest, 0)
	for _, request := range ms.paymentRequests {
		if request.MerchantID == merchantID {
			payments = append(payments, request)
		}
	}

	return payments
}

// Helper functions

func (ms *MerchantService) generateQRCode(walletAddress string, amount float64, asset, requestID string) string {
	// Generate QR code data in standard format
	qrData := fmt.Sprintf("btn:%s?amount=%.8f&asset=%s&request=%s", 
		walletAddress, amount, asset, requestID)
	return base64.StdEncoding.EncodeToString([]byte(qrData))
}

func (ms *MerchantService) generateNFCToken() string {
	// Generate secure NFC token
	token := make([]byte, 32)
	rand.Read(token)
	return base64.StdEncoding.EncodeToString(token)
}

func (ms *MerchantService) generateTransactionRef() string {
	ref := make([]byte, 16)
	rand.Read(ref)
	return fmt.Sprintf("TXN-%s", base64.URLEncoding.EncodeToString(ref)[:16])
}

// UpdateMerchantPaymentMethods updates accepted payment methods for a merchant
func (ms *MerchantService) UpdateMerchantPaymentMethods(merchantID string, methods []PaymentMethod) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	merchant, exists := ms.merchants[merchantID]
	if !exists {
		return ErrMerchantNotFound
	}

	merchant.AcceptedMethods = methods
	return nil
}
