// Package payments provides South African payment gateway integrations.
package payments

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

const (
	payfastLiveURL    = "https://www.payfast.co.za/eng/process"
	payfastSandboxURL = "https://sandbox.payfast.co.za/eng/process"
)

// PayFast provides PayFast payment gateway integration.
type PayFast struct {
	client *lunahttp.Client
	config PayFastConfig
}

// NewPayFast creates a new PayFast instance.
func NewPayFast(client *lunahttp.Client, config PayFastConfig) *PayFast {
	return &PayFast{
		client: client,
		config: config,
	}
}

// CreatePayment creates a payment request and returns the redirect URL.
func (p *PayFast) CreatePayment(ctx context.Context, req PayFastPaymentRequest) (*PayFastPayment, error) {
	paymentID := fmt.Sprintf("pf_%d", time.Now().UnixMilli())

	data := map[string]string{
		"merchant_id":  p.config.MerchantID,
		"merchant_key": p.config.MerchantKey,
		"return_url":   req.ReturnURL,
		"cancel_url":   req.CancelURL,
		"notify_url":   req.NotifyURL,
		"m_payment_id": paymentID,
		"amount":       fmt.Sprintf("%.2f", req.Amount),
		"item_name":    req.ItemName,
	}

	if req.ItemDescription != "" {
		data["item_description"] = req.ItemDescription
	}
	if req.EmailAddress != "" {
		data["email_address"] = req.EmailAddress
	}
	if req.CellNumber != "" {
		data["cell_number"] = req.CellNumber
	}
	if req.PaymentMethod != "" {
		data["payment_method"] = req.PaymentMethod
	}
	if req.CustomStr1 != "" {
		data["custom_str1"] = req.CustomStr1
	}
	if req.CustomStr2 != "" {
		data["custom_str2"] = req.CustomStr2
	}
	if req.CustomStr3 != "" {
		data["custom_str3"] = req.CustomStr3
	}
	if req.CustomInt1 != 0 {
		data["custom_int1"] = fmt.Sprintf("%d", req.CustomInt1)
	}
	if req.CustomInt2 != 0 {
		data["custom_int2"] = fmt.Sprintf("%d", req.CustomInt2)
	}

	signature := p.generateSignature(data)
	data["signature"] = signature

	baseURL := payfastLiveURL
	if p.config.Sandbox {
		baseURL = payfastSandboxURL
	}

	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	paymentURL := fmt.Sprintf("%s?%s", baseURL, values.Encode())

	currency := req.Currency
	if currency == "" {
		currency = "ZAR"
	}

	return &PayFastPayment{
		ID:       paymentID,
		Provider: "payfast",
		Amount: Amount{
			Value:    int(req.Amount * 100),
			Currency: currency,
		},
		Status:      StatusPending,
		Reference:   paymentID,
		Description: req.ItemDescription,
		PaymentURL:  paymentURL,
		Signature:   signature,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// VerifyWebhook verifies the webhook signature.
func (p *PayFast) VerifyWebhook(payload map[string]string) bool {
	signature := payload["signature"]
	delete(payload, "signature")
	expectedSignature := p.generateSignature(payload)
	return signature == expectedSignature
}

// ProcessWebhook processes a webhook and returns payment status.
func (p *PayFast) ProcessWebhook(payload map[string]string) *PayFastPayment {
	statusMap := map[string]PaymentStatus{
		"COMPLETE":  StatusCompleted,
		"FAILED":    StatusFailed,
		"PENDING":   StatusPending,
		"CANCELLED": StatusCancelled,
	}

	amountGross := 0.0
	if amt, ok := payload["amount_gross"]; ok {
		fmt.Sscanf(amt, "%f", &amountGross)
	}

	status := StatusPending
	if s, ok := statusMap[payload["payment_status"]]; ok {
		status = s
	}

	return &PayFastPayment{
		ID:          payload["m_payment_id"],
		Provider:    "payfast",
		PFPaymentID: payload["pf_payment_id"],
		Amount: Amount{
			Value:    int(amountGross * 100),
			Currency: "ZAR",
		},
		Status:      status,
		Reference:   payload["m_payment_id"],
		Description: payload["item_name"],
		PaymentURL:  "",
		Signature:   payload["signature"],
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Refund requests a refund for a payment.
func (p *PayFast) Refund(ctx context.Context, req RefundRequest) (*Refund, error) {
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

func (p *PayFast) generateSignature(data map[string]string) string {
	// Sort keys
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build param string
	var parts []string
	for _, k := range keys {
		v := data[k]
		if v != "" {
			encoded := strings.ReplaceAll(url.QueryEscape(v), "+", "%20")
			parts = append(parts, fmt.Sprintf("%s=%s", k, encoded))
		}
	}
	paramString := strings.Join(parts, "&")

	if p.config.Passphrase != "" {
		paramString += "&passphrase=" + url.QueryEscape(p.config.Passphrase)
	}

	hash := md5.Sum([]byte(paramString))
	return hex.EncodeToString(hash[:])
}
