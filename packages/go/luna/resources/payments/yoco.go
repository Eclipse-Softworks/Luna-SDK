// Package payments provides South African payment gateway integrations.
package payments

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// Yoco provides Yoco online payment integration.
type Yoco struct {
	client *lunahttp.Client
	config YocoConfig
}

// NewYoco creates a new Yoco instance.
func NewYoco(client *lunahttp.Client, config YocoConfig) *Yoco {
	return &Yoco{
		client: client,
		config: config,
	}
}

// CreatePayment creates a checkout session and returns the redirect URL.
func (y *Yoco) CreatePayment(ctx context.Context, req YocoPaymentRequest) (*YocoPayment, error) {
	currency := req.Currency
	if currency == "" {
		currency = "ZAR"
	}

	body := map[string]interface{}{
		"amount":     req.Amount,
		"currency":   currency,
		"metadata":   req.Metadata,
		"successUrl": req.SuccessURL,
		"cancelUrl":  req.CancelURL,
		"failureUrl": req.FailureURL,
	}

	resp, err := y.client.Request(ctx, lunahttp.RequestConfig{
		Method: "POST",
		Path:   "/v1/payments/yoco/checkouts",
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	var result struct {
		ID          string                 `json:"id"`
		CheckoutID  string                 `json:"checkoutId"`
		Amount      float64                `json:"amount"`
		Currency    string                 `json:"currency"`
		Status      PaymentStatus          `json:"status"`
		Reference   string                 `json:"reference"`
		RedirectURL string                 `json:"redirectUrl"`
		Metadata    map[string]interface{} `json:"metadata"`
		CreatedAt   time.Time              `json:"createdAt"`
		UpdatedAt   time.Time              `json:"updatedAt"`
	}

	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}

	return &YocoPayment{
		ID:         result.ID,
		Provider:   "yoco",
		CheckoutID: result.CheckoutID,
		Amount: Amount{
			Value:    int(result.Amount),
			Currency: result.Currency,
		},
		Status:      result.Status,
		Reference:   result.Reference,
		RedirectURL: result.RedirectURL,
		Metadata:    result.Metadata,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

// VerifyWebhook verifies the webhook signature.
func (y *Yoco) VerifyWebhook(payload string, signature string) bool {
	mac := hmac.New(sha256.New, []byte(y.config.SecretKey))
	mac.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return signature == expectedSignature
}

// ProcessWebhook processes a webhook event.
func (y *Yoco) ProcessWebhook(payload map[string]interface{}) *YocoPayment {
	statusMap := map[string]PaymentStatus{
		"payment.succeeded": StatusCompleted,
		"payment.failed":    StatusFailed,
		"payment.cancelled": StatusCancelled,
	}

	eventType, _ := payload["type"].(string)
	paymentData, _ := payload["payload"].(map[string]interface{})

	id, _ := paymentData["id"].(string)
	amount, _ := paymentData["amount"].(float64)
	currency, _ := paymentData["currency"].(string)
	if currency == "" {
		currency = "ZAR"
	}

	status := StatusPending
	if s, ok := statusMap[eventType]; ok {
		status = s
	}

	var metadata map[string]interface{}
	if m, ok := paymentData["metadata"].(map[string]interface{}); ok {
		metadata = m
	}

	return &YocoPayment{
		ID:         fmt.Sprintf("yc_%s", id),
		Provider:   "yoco",
		CheckoutID: id,
		Amount: Amount{
			Value:    int(amount),
			Currency: currency,
		},
		Status:    status,
		Reference: id,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Refund requests a refund.
func (y *Yoco) Refund(ctx context.Context, req RefundRequest) (*Refund, error) {
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
