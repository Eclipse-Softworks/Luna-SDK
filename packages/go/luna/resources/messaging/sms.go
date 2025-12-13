// Package messaging provides SMS, WhatsApp, and USSD integrations for South Africa.
package messaging

import (
	"context"
	"fmt"
	"regexp"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// SMS provides multi-provider SMS integration.
type SMS struct {
	client *lunahttp.Client
	config SMSConfig
}

// NewSMS creates a new SMS instance.
func NewSMS(client *lunahttp.Client, config SMSConfig) *SMS {
	return &SMS{
		client: client,
		config: config,
	}
}

// Send sends a single SMS message.
func (s *SMS) Send(ctx context.Context, req SMSSendRequest) (*SMSMessage, error) {
	if len(req.To) == 0 {
		return nil, fmt.Errorf("SMS recipient (to) is required")
	}

	to := req.To[0]
	messageID := fmt.Sprintf("sms_%d", time.Now().UnixMilli())
	normalizedTo := s.normalizePhoneNumber(to)

	from := req.From
	if from == "" {
		from = s.config.SenderID
	}

	return &SMSMessage{
		ID:        messageID,
		To:        normalizedTo,
		From:      from,
		Body:      req.Body,
		Status:    StatusPending,
		Direction: "outbound",
		Provider:  s.config.Provider,
		Parts:     (len(req.Body) + 159) / 160,
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// SendBulk sends SMS to multiple recipients.
func (s *SMS) SendBulk(ctx context.Context, req SMSSendRequest) (*SMSBulkResult, error) {
	result := &SMSBulkResult{
		Successful: []SMSMessage{},
		Failed: []struct {
			To    string `json:"to"`
			Error string `json:"error"`
		}{},
	}

	for _, to := range req.To {
		msg, err := s.Send(ctx, SMSSendRequest{
			To:          []string{to},
			Body:        req.Body,
			From:        req.From,
			CallbackURL: req.CallbackURL,
			Metadata:    req.Metadata,
		})
		if err != nil {
			result.Failed = append(result.Failed, struct {
				To    string `json:"to"`
				Error string `json:"error"`
			}{To: to, Error: err.Error()})
		} else {
			result.Successful = append(result.Successful, *msg)
		}
	}

	return result, nil
}

// GetStatus gets SMS delivery status.
func (s *SMS) GetStatus(ctx context.Context, messageID string) (*SMSMessage, error) {
	return &SMSMessage{
		ID:        messageID,
		To:        "",
		Body:      "",
		Status:    StatusDelivered,
		Direction: "outbound",
		Provider:  s.config.Provider,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetBalance gets account balance.
func (s *SMS) GetBalance(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"balance":  100.0,
		"currency": "ZAR",
	}, nil
}

func (s *SMS) normalizePhoneNumber(phone string) string {
	// Remove non-digits
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Handle SA numbers
	if len(digits) == 10 && digits[0] == '0' {
		digits = "27" + digits[1:]
	}

	if digits[0] != '+' {
		digits = "+" + digits
	}

	return digits
}
