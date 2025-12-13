// Package messaging provides SMS, WhatsApp, and USSD integrations for South Africa.
package messaging

import (
	"fmt"
	"strings"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// USSDHandler is a function that handles USSD sessions.
type USSDHandler func(USSDSession) USSDResponse

// USSD provides USSD service integration for South African networks.
type USSD struct {
	client   *lunahttp.Client
	config   USSDConfig
	handlers map[string]USSDHandler
}

// NewUSSD creates a new USSD instance.
func NewUSSD(client *lunahttp.Client, config USSDConfig) *USSD {
	return &USSD{
		client:   client,
		config:   config,
		handlers: make(map[string]USSDHandler),
	}
}

// OnSession registers a handler for USSD sessions.
func (u *USSD) OnSession(handler USSDHandler) {
	u.handlers["default"] = handler
}

// OnMenu registers a handler for a specific menu path.
func (u *USSD) OnMenu(path string, handler USSDHandler) {
	u.handlers[path] = handler
}

// ProcessRequest processes incoming USSD request.
func (u *USSD) ProcessRequest(session USSDSession) USSDResponse {
	handler, ok := u.handlers[session.Text]
	if !ok {
		handler, ok = u.handlers["default"]
	}

	if !ok {
		return USSDResponse{
			Text: "Service temporarily unavailable. Please try again later.",
			End:  true,
		}
	}

	return handler(session)
}

// Menu creates a menu response.
func Menu(title string, options []struct{ Key, Label string }) string {
	lines := []string{title, ""}
	for _, opt := range options {
		lines = append(lines, fmt.Sprintf("%s. %s", opt.Key, opt.Label))
	}
	return strings.Join(lines, "\n")
}

// ParseAfricasTalkingRequest parses Africa's Talking webhook format.
func (u *USSD) ParseAfricasTalkingRequest(sessionID, phoneNumber, serviceCode, text, networkCode string) USSDSession {
	networks := map[string]string{
		"655001": "Vodacom",
		"655002": "Telkom",
		"655007": "Cell C",
		"655010": "MTN",
	}

	network := networkCode
	if n, ok := networks[networkCode]; ok {
		network = n
	}

	return USSDSession{
		ID:          fmt.Sprintf("ussd_%d", time.Now().UnixMilli()),
		SessionID:   sessionID,
		PhoneNumber: phoneNumber,
		ServiceCode: serviceCode,
		Text:        text,
		State:       USSDActive,
		Network:     network,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// FormatAfricasTalkingResponse formats response for Africa's Talking.
func (u *USSD) FormatAfricasTalkingResponse(response USSDResponse) string {
	prefix := "CON "
	if response.End {
		prefix = "END "
	}
	return prefix + response.Text
}

// ParseClickatellRequest parses Clickatell webhook format.
func (u *USSD) ParseClickatellRequest(sessionID, msisdn, request, shortcode string) USSDSession {
	return USSDSession{
		ID:          fmt.Sprintf("ussd_%d", time.Now().UnixMilli()),
		SessionID:   sessionID,
		PhoneNumber: msisdn,
		ServiceCode: shortcode,
		Text:        request,
		State:       USSDActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ServiceCode returns the USSD service code.
func (u *USSD) ServiceCode() string {
	return u.config.ServiceCode
}

// CreateExampleMenu creates an example session flow.
func CreateExampleMenu() USSDHandler {
	return func(session USSDSession) USSDResponse {
		parts := strings.Split(session.Text, "*")
		var filteredParts []string
		for _, p := range parts {
			if p != "" {
				filteredParts = append(filteredParts, p)
			}
		}

		if len(filteredParts) == 0 {
			return USSDResponse{
				Text: Menu("Welcome to Luna SDK", []struct{ Key, Label string }{
					{"1", "Check Balance"},
					{"2", "Send Payment"},
					{"3", "Mini Statement"},
					{"4", "Exit"},
				}),
				End: false,
			}
		}

		selection := filteredParts[0]

		switch selection {
		case "1":
			return USSDResponse{Text: "Your balance is R1,234.56", End: true}
		case "2":
			if len(filteredParts) == 1 {
				return USSDResponse{Text: "Enter phone number to send payment:", End: false}
			}
			if len(filteredParts) == 2 {
				return USSDResponse{Text: "Enter amount (ZAR):", End: false}
			}
			return USSDResponse{
				Text: fmt.Sprintf("Payment of R%s to %s initiated.", filteredParts[2], filteredParts[1]),
				End:  true,
			}
		case "3":
			return USSDResponse{
				Text: "Mini Statement:\n1. Received R500.00\n2. Sent R100.00\n3. Airtime R50.00",
				End:  true,
			}
		case "4":
			return USSDResponse{Text: "Thank you for using Luna SDK. Goodbye!", End: true}
		default:
			return USSDResponse{Text: "Invalid selection. Please try again.", End: true}
		}
	}
}
