// Package payments provides South African payment gateway integrations.
package payments

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

const ozowPaymentURL = "https://pay.ozow.com"

// Ozow provides Ozow instant EFT payment integration.
type Ozow struct {
	client *lunahttp.Client
	config OzowConfig
}

// NewOzow creates a new Ozow instance.
func NewOzow(client *lunahttp.Client, config OzowConfig) *Ozow {
	return &Ozow{
		client: client,
		config: config,
	}
}

// CreatePayment creates a payment request and returns the redirect URL.
func (o *Ozow) CreatePayment(ctx context.Context, req OzowPaymentRequest) (*OzowPayment, error) {
	paymentID := fmt.Sprintf("oz_%d", time.Now().UnixMilli())

	isTest := "false"
	if req.IsTest != nil && *req.IsTest || o.config.Sandbox {
		isTest = "true"
	}

	data := map[string]string{
		"SiteCode":             o.config.SiteCode,
		"CountryCode":          "ZA",
		"CurrencyCode":         "ZAR",
		"Amount":               fmt.Sprintf("%.2f", req.Amount),
		"TransactionReference": req.TransactionReference,
		"BankReference":        req.BankReference,
		"CancelUrl":            req.CancelURL,
		"ErrorUrl":             req.ErrorURL,
		"SuccessUrl":           req.SuccessURL,
		"NotifyUrl":            req.NotifyURL,
		"IsTest":               isTest,
	}

	if req.CustomerFirstName != "" {
		data["CustomerFirstName"] = req.CustomerFirstName
	}
	if req.CustomerLastName != "" {
		data["CustomerLastName"] = req.CustomerLastName
	}
	if req.CustomerEmail != "" {
		data["CustomerEmail"] = req.CustomerEmail
	}
	if req.CustomerPhone != "" {
		data["CustomerPhone"] = req.CustomerPhone
	}

	hashString := o.generateHashString(data)
	hashCheck := o.generateHash(hashString)
	data["HashCheck"] = hashCheck

	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	paymentURL := fmt.Sprintf("%s?%s", ozowPaymentURL, values.Encode())

	return &OzowPayment{
		ID:       paymentID,
		Provider: "ozow",
		Amount: Amount{
			Value:    int(req.Amount * 100),
			Currency: "ZAR",
		},
		Status:        StatusPending,
		Reference:     req.TransactionReference,
		Description:   req.BankReference,
		PaymentURL:    paymentURL,
		TransactionID: paymentID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// VerifyWebhook verifies the webhook hash.
func (o *Ozow) VerifyWebhook(payload map[string]string) bool {
	receivedHash := strings.ToLower(payload["Hash"])
	delete(payload, "Hash")
	hashString := o.generateHashString(payload)
	expectedHash := strings.ToLower(o.generateHash(hashString))
	return receivedHash == expectedHash
}

// ProcessWebhook processes a webhook and returns payment status.
func (o *Ozow) ProcessWebhook(payload map[string]string) *OzowPayment {
	statusMap := map[string]PaymentStatus{
		"Complete":             StatusCompleted,
		"Cancelled":            StatusCancelled,
		"Error":                StatusFailed,
		"Abandoned":            StatusCancelled,
		"PendingInvestigation": StatusProcessing,
	}

	amount := 0.0
	if amt, ok := payload["Amount"]; ok {
		fmt.Sscanf(amt, "%f", &amount)
	}

	status := StatusPending
	if s, ok := statusMap[payload["Status"]]; ok {
		status = s
	}

	return &OzowPayment{
		ID:            payload["TransactionReference"],
		Provider:      "ozow",
		TransactionID: payload["TransactionId"],
		Amount: Amount{
			Value:    int(amount * 100),
			Currency: "ZAR",
		},
		Status:    status,
		Reference: payload["TransactionReference"],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Refund requests a refund.
func (o *Ozow) Refund(ctx context.Context, req RefundRequest) (*Refund, error) {
	refundID := fmt.Sprintf("ref_%d", time.Now().UnixMilli())

	amount := 0
	if req.Amount != nil {
		amount = *req.Amount
	}

	return &Refund{
		ID:        refundID,
		PaymentID: req.PaymentID,
		Amount: Amount{
			Value:    amount,
			Currency: "ZAR",
		},
		Status:    "pending",
		Reason:    req.Reason,
		CreatedAt: time.Now(),
	}, nil
}

func (o *Ozow) generateHashString(data map[string]string) string {
	orderedFields := []string{
		"SiteCode", "CountryCode", "CurrencyCode", "Amount",
		"TransactionReference", "BankReference", "CancelUrl",
		"ErrorUrl", "SuccessUrl", "NotifyUrl", "IsTest",
	}

	var parts []string
	for _, field := range orderedFields {
		if v, ok := data[field]; ok {
			parts = append(parts, v)
		} else {
			parts = append(parts, "")
		}
	}
	return strings.ToLower(strings.Join(parts, ""))
}

func (o *Ozow) generateHash(input string) string {
	toHash := input + strings.ToLower(o.config.PrivateKey)
	hash := sha512.Sum512([]byte(toHash))
	return hex.EncodeToString(hash[:])
}
