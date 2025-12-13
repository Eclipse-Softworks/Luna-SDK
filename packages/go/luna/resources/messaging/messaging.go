// Package messaging provides SMS, WhatsApp, and USSD integrations for South Africa.
package messaging

import (
	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// Messaging provides unified access to messaging channels.
type Messaging struct {
	client   *lunahttp.Client
	config   *Config
	sms      *SMS
	whatsapp *WhatsApp
	ussd     *USSD
}

// NewMessaging creates a new Messaging resource.
func NewMessaging(client *lunahttp.Client, config *Config) *Messaging {
	if config == nil {
		config = &Config{}
	}
	return &Messaging{
		client: client,
		config: config,
	}
}

// SMS returns the SMS instance.
func (m *Messaging) SMS() *SMS {
	if m.sms == nil {
		if m.config.SMS == nil {
			panic("SMS not configured. Provide SMSConfig when initializing LunaClient.")
		}
		m.sms = NewSMS(m.client, *m.config.SMS)
	}
	return m.sms
}

// WhatsApp returns the WhatsApp instance.
func (m *Messaging) WhatsApp() *WhatsApp {
	if m.whatsapp == nil {
		if m.config.WhatsApp == nil {
			panic("WhatsApp not configured. Provide WhatsAppConfig when initializing LunaClient.")
		}
		m.whatsapp = NewWhatsApp(m.client, *m.config.WhatsApp)
	}
	return m.whatsapp
}

// USSD returns the USSD instance.
func (m *Messaging) USSD() *USSD {
	if m.ussd == nil {
		if m.config.USSD == nil {
			panic("USSD not configured. Provide USSDConfig when initializing LunaClient.")
		}
		m.ussd = NewUSSD(m.client, *m.config.USSD)
	}
	return m.ussd
}

// List returns available messaging channels.
func (m *Messaging) List() []string {
	var available []string
	if m.config.SMS != nil {
		available = append(available, "sms")
	}
	if m.config.WhatsApp != nil {
		available = append(available, "whatsapp")
	}
	if m.config.USSD != nil {
		available = append(available, "ussd")
	}
	return available
}
