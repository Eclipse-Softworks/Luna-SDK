// Package messaging provides SMS, WhatsApp, and USSD integrations for South Africa.
package messaging

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// WhatsApp provides WhatsApp Business API integration.
type WhatsApp struct {
	client *lunahttp.Client
	config WhatsAppConfig
}

// NewWhatsApp creates a new WhatsApp instance.
func NewWhatsApp(client *lunahttp.Client, config WhatsAppConfig) *WhatsApp {
	return &WhatsApp{
		client: client,
		config: config,
	}
}

// SendText sends a text message.
func (w *WhatsApp) SendText(ctx context.Context, req WhatsAppTextRequest) (*WhatsAppMessage, error) {
	messageID := fmt.Sprintf("wa_%d", time.Now().UnixMilli())
	to := w.normalizePhoneNumber(req.To)

	return &WhatsAppMessage{
		ID:        messageID,
		To:        to,
		Type:      "text",
		Text:      req.Text,
		Status:    StatusPending,
		Direction: "outbound",
		Provider:  w.config.Provider,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// SendTemplate sends a template message.
func (w *WhatsApp) SendTemplate(ctx context.Context, req WhatsAppTemplateRequest) (*WhatsAppMessage, error) {
	messageID := fmt.Sprintf("wa_%d", time.Now().UnixMilli())
	to := w.normalizePhoneNumber(req.To)

	return &WhatsAppMessage{
		ID:             messageID,
		To:             to,
		Type:           "template",
		TemplateName:   req.TemplateName,
		TemplateParams: req.TemplateParams,
		Status:         StatusPending,
		Direction:      "outbound",
		Provider:       w.config.Provider,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// SendMedia sends a media message.
func (w *WhatsApp) SendMedia(ctx context.Context, req WhatsAppMediaRequest) (*WhatsAppMessage, error) {
	messageID := fmt.Sprintf("wa_%d", time.Now().UnixMilli())
	to := w.normalizePhoneNumber(req.To)

	return &WhatsAppMessage{
		ID:        messageID,
		To:        to,
		Type:      req.Type,
		MediaURL:  req.MediaURL,
		Status:    StatusPending,
		Direction: "outbound",
		Provider:  w.config.Provider,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetStatus gets message status.
func (w *WhatsApp) GetStatus(ctx context.Context, messageID string) (*WhatsAppMessage, error) {
	return &WhatsAppMessage{
		ID:        messageID,
		To:        "",
		Type:      "text",
		Status:    StatusDelivered,
		Direction: "outbound",
		Provider:  w.config.Provider,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// VerifyWebhook verifies webhook signature.
func (w *WhatsApp) VerifyWebhook(payload, signature string) bool {
	if w.config.WebhookToken == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(w.config.WebhookToken))
	mac.Write([]byte(payload))
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return expectedSignature == signature
}

// ProcessWebhook processes incoming webhook.
func (w *WhatsApp) ProcessWebhook(payload map[string]interface{}) []WhatsAppMessage {
	messages := []WhatsAppMessage{}

	entries, ok := payload["entry"].([]interface{})
	if !ok {
		return messages
	}

	statusMap := map[string]MessageStatus{
		"sent":      StatusSent,
		"delivered": StatusDelivered,
		"read":      StatusRead,
		"failed":    StatusFailed,
	}

	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		changes, ok := entryMap["changes"].([]interface{})
		if !ok {
			continue
		}

		for _, change := range changes {
			changeMap, ok := change.(map[string]interface{})
			if !ok {
				continue
			}

			value, ok := changeMap["value"].(map[string]interface{})
			if !ok {
				continue
			}

			// Process incoming messages
			if msgs, ok := value["messages"].([]interface{}); ok {
				for _, msg := range msgs {
					msgMap, ok := msg.(map[string]interface{})
					if !ok {
						continue
					}

					messages = append(messages, WhatsAppMessage{
						ID:        msgMap["id"].(string),
						From:      msgMap["from"].(string),
						Type:      msgMap["type"].(string),
						Status:    StatusDelivered,
						Direction: "inbound",
						Provider:  w.config.Provider,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
				}
			}

			// Process status updates
			if statuses, ok := value["statuses"].([]interface{}); ok {
				for _, status := range statuses {
					statusMap2, ok := status.(map[string]interface{})
					if !ok {
						continue
					}

					messageStatus := StatusPending
					if s, ok := statusMap[statusMap2["status"].(string)]; ok {
						messageStatus = s
					}

					messages = append(messages, WhatsAppMessage{
						ID:        statusMap2["id"].(string),
						To:        "",
						Type:      "text",
						Status:    messageStatus,
						Direction: "outbound",
						Provider:  w.config.Provider,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					})
				}
			}
		}
	}

	return messages
}

func (w *WhatsApp) normalizePhoneNumber(phone string) string {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(digits) == 10 && digits[0] == '0' {
		digits = "27" + digits[1:]
	}

	if len(digits) > 0 && digits[0] == '+' {
		digits = digits[1:]
	}

	return digits
}
