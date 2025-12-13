// Package messaging provides SMS, WhatsApp, and USSD integrations for South Africa.
package messaging

import "time"

// Provider types
type SMSProvider string
type WhatsAppProvider string
type USSDProvider string

const (
	SMSClickatell     SMSProvider = "clickatell"
	SMSAfricasTalking SMSProvider = "africastalking"
	SMSTwilio         SMSProvider = "twilio"

	WhatsAppCloudAPI   WhatsAppProvider = "cloud_api"
	WhatsAppClickatell WhatsAppProvider = "clickatell"
	WhatsAppWati       WhatsAppProvider = "wati"
	WhatsAppInfobip    WhatsAppProvider = "infobip"

	USSDClickatell     USSDProvider = "clickatell"
	USSDAfricasTalking USSDProvider = "africastalking"
)

// MessageStatus represents message delivery status
type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

// USSDState represents USSD session state
type USSDState string

const (
	USSDNew       USSDState = "new"
	USSDActive    USSDState = "active"
	USSDCompleted USSDState = "completed"
	USSDTimeout   USSDState = "timeout"
)

// SMSConfig holds SMS configuration
type SMSConfig struct {
	Provider SMSProvider `json:"provider"`
	APIKey   string      `json:"api_key"`
	Username string      `json:"username,omitempty"`
	SenderID string      `json:"sender_id,omitempty"`
	Sandbox  bool        `json:"sandbox"`
}

// WhatsAppConfig holds WhatsApp configuration
type WhatsAppConfig struct {
	Provider      WhatsAppProvider `json:"provider"`
	APIKey        string           `json:"api_key"`
	PhoneNumberID string           `json:"phone_number_id,omitempty"`
	WebhookToken  string           `json:"webhook_token,omitempty"`
	Sandbox       bool             `json:"sandbox"`
}

// USSDConfig holds USSD configuration
type USSDConfig struct {
	Provider    USSDProvider `json:"provider"`
	APIKey      string       `json:"api_key"`
	ServiceCode string       `json:"service_code"`
	Sandbox     bool         `json:"sandbox"`
}

// Config holds all messaging configuration
type Config struct {
	SMS      *SMSConfig      `json:"sms,omitempty"`
	WhatsApp *WhatsAppConfig `json:"whatsapp,omitempty"`
	USSD     *USSDConfig     `json:"ussd,omitempty"`
}

// SMSSendRequest represents an SMS send request
type SMSSendRequest struct {
	To          []string               `json:"to"`
	Body        string                 `json:"body"`
	From        string                 `json:"from,omitempty"`
	CallbackURL string                 `json:"callback_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SMSMessage represents an SMS message
type SMSMessage struct {
	ID        string                 `json:"id"`
	To        string                 `json:"to"`
	From      string                 `json:"from,omitempty"`
	Body      string                 `json:"body"`
	Status    MessageStatus          `json:"status"`
	Direction string                 `json:"direction"` // inbound/outbound
	Provider  SMSProvider            `json:"provider,omitempty"`
	Parts     int                    `json:"parts,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// SMSBulkResult represents bulk SMS result
type SMSBulkResult struct {
	Successful []SMSMessage `json:"successful"`
	Failed     []struct {
		To    string `json:"to"`
		Error string `json:"error"`
	} `json:"failed"`
}

// WhatsAppTextRequest represents a WhatsApp text message request
type WhatsAppTextRequest struct {
	To   string `json:"to"`
	Text string `json:"text"`
}

// WhatsAppTemplateRequest represents a WhatsApp template message request
type WhatsAppTemplateRequest struct {
	To             string                 `json:"to"`
	TemplateName   string                 `json:"template_name"`
	TemplateParams map[string]interface{} `json:"template_params,omitempty"`
	Language       string                 `json:"language,omitempty"`
}

// WhatsAppMediaRequest represents a WhatsApp media message request
type WhatsAppMediaRequest struct {
	To       string `json:"to"`
	Type     string `json:"type"` // image, document, audio, video
	MediaURL string `json:"media_url"`
	Caption  string `json:"caption,omitempty"`
}

// WhatsAppMessage represents a WhatsApp message
type WhatsAppMessage struct {
	ID             string                 `json:"id"`
	To             string                 `json:"to"`
	From           string                 `json:"from,omitempty"`
	Type           string                 `json:"type"`
	Text           string                 `json:"text,omitempty"`
	TemplateName   string                 `json:"template_name,omitempty"`
	TemplateParams map[string]interface{} `json:"template_params,omitempty"`
	MediaURL       string                 `json:"media_url,omitempty"`
	Status         MessageStatus          `json:"status"`
	Direction      string                 `json:"direction"`
	Provider       WhatsAppProvider       `json:"provider,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// USSDSession represents a USSD session
type USSDSession struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	PhoneNumber string    `json:"phone_number"`
	ServiceCode string    `json:"service_code"`
	Text        string    `json:"text"`
	State       USSDState `json:"state"`
	Network     string    `json:"network,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// USSDResponse represents a USSD response
type USSDResponse struct {
	Text string `json:"text"`
	End  bool   `json:"end"`
}
