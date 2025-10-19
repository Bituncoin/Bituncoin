package payments

import (
	"encoding/json"
	"errors"
	"fmt"
)

// QRPaymentData represents payment information encoded in QR code
type QRPaymentData struct {
	Address    string  `json:"address"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Memo       string  `json:"memo,omitempty"`
	InvoiceID  string  `json:"invoiceId,omitempty"`
	MerchantID string  `json:"merchantId,omitempty"`
}

// QRCodeService handles QR code payment generation and parsing
type QRCodeService struct {
	// In production, this would integrate with actual QR code libraries
}

// NewQRCodeService creates a new QR code service
func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

// GeneratePaymentQR generates QR code data for payment
func (qrs *QRCodeService) GeneratePaymentQR(address string, amount float64, currency string, memo string) (string, error) {
	if address == "" {
		return "", errors.New("address required")
	}
	if amount <= 0 {
		return "", errors.New("amount must be positive")
	}

	data := QRPaymentData{
		Address:  address,
		Amount:   amount,
		Currency: currency,
		Memo:     memo,
	}

	// Convert to JSON for QR encoding
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// In production, generate actual QR code image
	// For now, return the data string that would be encoded
	return string(jsonData), nil
}

// GenerateMerchantQR generates QR code for merchant payment
func (qrs *QRCodeService) GenerateMerchantQR(merchantID string, address string, amount float64, currency string, invoiceID string) (string, error) {
	if merchantID == "" {
		return "", errors.New("merchant ID required")
	}
	if address == "" {
		return "", errors.New("address required")
	}

	data := QRPaymentData{
		Address:    address,
		Amount:     amount,
		Currency:   currency,
		MerchantID: merchantID,
		InvoiceID:  invoiceID,
		Memo:       fmt.Sprintf("Payment to merchant %s", merchantID),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// ParsePaymentQR parses QR code data
func (qrs *QRCodeService) ParsePaymentQR(qrData string) (*QRPaymentData, error) {
	var data QRPaymentData
	err := json.Unmarshal([]byte(qrData), &data)
	if err != nil {
		return nil, errors.New("invalid QR code data")
	}

	// Validate parsed data
	if data.Address == "" {
		return nil, errors.New("invalid payment address in QR code")
	}
	if data.Amount <= 0 {
		return nil, errors.New("invalid amount in QR code")
	}

	return &data, nil
}

// GenerateReceiveQR generates QR code for receiving payments
func (qrs *QRCodeService) GenerateReceiveQR(address string, currency string) (string, error) {
	if address == "" {
		return "", errors.New("address required")
	}

	data := QRPaymentData{
		Address:  address,
		Amount:   0, // Amount to be specified by sender
		Currency: currency,
		Memo:     "Receive payment",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// ValidateQRPayment validates a QR payment before processing
func (qrs *QRCodeService) ValidateQRPayment(qrData string, expectedCurrency string) error {
	data, err := qrs.ParsePaymentQR(qrData)
	if err != nil {
		return err
	}

	if data.Currency != expectedCurrency {
		return fmt.Errorf("currency mismatch: expected %s, got %s", expectedCurrency, data.Currency)
	}

	return nil
}
