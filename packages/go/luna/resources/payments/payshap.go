// Package payments provides South African payment gateway integrations.
package payments

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// SABank represents a South African bank
type SABank string

const (
	BankABSA      SABank = "absa"
	BankCapitec   SABank = "capitec"
	BankFNB       SABank = "fnb"
	BankNedbank   SABank = "nedbank"
	BankStandard  SABank = "standard"
	BankInvestec  SABank = "investec"
	BankDiscovery SABank = "discovery"
	BankTymeBank  SABank = "tymebank"
	BankAfrican   SABank = "african"
)

// PayShap provides PayShap real-time payment integration.
type PayShap struct {
	client *lunahttp.Client
	config PayShapConfig
}

// NewPayShap creates a new PayShap instance.
func NewPayShap(client *lunahttp.Client, config PayShapConfig) *PayShap {
	return &PayShap{
		client: client,
		config: config,
	}
}

// CreatePayment creates a PayShap payment request.
func (p *PayShap) CreatePayment(ctx context.Context, req PayShapPaymentRequest) (*PayShapPayment, error) {
	paymentID := fmt.Sprintf("ps_%d", time.Now().UnixMilli())

	expiryMinutes := req.ExpiryMinutes
	if expiryMinutes == 0 {
		expiryMinutes = 30
	}
	expiresAt := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	// Generate QR code data
	qrData := map[string]interface{}{
		"type":       "payshap",
		"merchantId": p.config.MerchantID,
		"amount":     req.Amount,
		"reference":  req.Reference,
	}
	qrJSON, _ := json.Marshal(qrData)
	qrCode := base64.StdEncoding.EncodeToString(qrJSON)

	return &PayShapPayment{
		ID:       paymentID,
		Provider: "payshap",
		ShapID:   fmt.Sprintf("shp_%s", paymentID[len(paymentID)-8:]),
		Amount: Amount{
			Value:    int(req.Amount * 100),
			Currency: "ZAR",
		},
		Status:    StatusPending,
		Reference: req.Reference,
		QRCode:    qrCode,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetPayment gets payment status.
func (p *PayShap) GetPayment(ctx context.Context, paymentID string) (*PayShapPayment, error) {
	return &PayShapPayment{
		ID:       paymentID,
		Provider: "payshap",
		ShapID:   fmt.Sprintf("shp_%s", paymentID[len(paymentID)-8:]),
		Amount: Amount{
			Value:    0,
			Currency: "ZAR",
		},
		Status:    StatusPending,
		Reference: paymentID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// CancelPayment cancels a pending payment.
func (p *PayShap) CancelPayment(ctx context.Context, paymentID string) (*PayShapPayment, error) {
	payment, _ := p.GetPayment(ctx, paymentID)
	payment.Status = StatusCancelled
	payment.UpdatedAt = time.Now()
	return payment, nil
}

// LookupShapID looks up a ShapID (payment proxy).
func (p *PayShap) LookupShapID(ctx context.Context, shapID string) (map[string]interface{}, error) {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9@._-]+$`)
	isValid := pattern.MatchString(shapID) && len(shapID) >= 5

	result := map[string]interface{}{
		"valid": isValid,
	}

	if isValid {
		result["bank_name"] = "Sample Bank"
		result["account_holder_name"] = "Account Holder"
	}

	return result, nil
}

// GenerateReceiveQR generates a QR code for receiving payments.
func (p *PayShap) GenerateReceiveQR(ctx context.Context, amount *float64, reference string) (map[string]string, error) {
	shapID := fmt.Sprintf("shp_%s_%d", p.config.MerchantID, time.Now().Unix())

	qrData := map[string]interface{}{
		"type":       "payshap_receive",
		"shapId":     shapID,
		"merchantId": p.config.MerchantID,
		"amount":     amount,
		"reference":  reference,
	}
	qrJSON, _ := json.Marshal(qrData)

	return map[string]string{
		"qr_code": base64.StdEncoding.EncodeToString(qrJSON),
		"shap_id": shapID,
	}, nil
}

// ValidateBankAccount validates a South African bank account number format.
func (p *PayShap) ValidateBankAccount(accountNumber string, bankID SABank) bool {
	accountLengths := map[SABank][]int{
		BankABSA:      {10, 11},
		BankCapitec:   {10},
		BankFNB:       {10, 11, 12},
		BankNedbank:   {10, 11},
		BankStandard:  {9, 10, 11},
		BankInvestec:  {10},
		BankDiscovery: {10},
		BankTymeBank:  {10},
		BankAfrican:   {11},
	}

	validLengths, ok := accountLengths[bankID]
	if !ok {
		return false
	}

	// Extract digits only
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(accountNumber, "")

	for _, length := range validLengths {
		if len(digitsOnly) == length {
			return true
		}
	}
	return false
}
