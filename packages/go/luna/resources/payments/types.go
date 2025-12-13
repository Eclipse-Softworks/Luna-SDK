// Package payments provides South African payment gateway integrations.
package payments

import (
	"time"
)

// PaymentProvider represents available payment providers
type PaymentProvider string

const (
	ProviderPayFast PaymentProvider = "payfast"
	ProviderOzow    PaymentProvider = "ozow"
	ProviderYoco    PaymentProvider = "yoco"
	ProviderPayShap PaymentProvider = "payshap"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusCompleted  PaymentStatus = "completed"
	StatusFailed     PaymentStatus = "failed"
	StatusCancelled  PaymentStatus = "cancelled"
	StatusRefunded   PaymentStatus = "refunded"
)

// Amount represents a monetary amount
type Amount struct {
	Value    int    `json:"value"`    // In cents
	Currency string `json:"currency"` // ISO 4217 code
}

// PayFastConfig holds PayFast configuration
type PayFastConfig struct {
	MerchantID  string `json:"merchant_id"`
	MerchantKey string `json:"merchant_key"`
	Passphrase  string `json:"passphrase,omitempty"`
	Sandbox     bool   `json:"sandbox"`
}

// OzowConfig holds Ozow configuration
type OzowConfig struct {
	SiteCode   string `json:"site_code"`
	PrivateKey string `json:"private_key"`
	APIKey     string `json:"api_key,omitempty"`
	Sandbox    bool   `json:"sandbox"`
}

// YocoConfig holds Yoco configuration
type YocoConfig struct {
	SecretKey string `json:"secret_key"`
	PublicKey string `json:"public_key,omitempty"`
	Sandbox   bool   `json:"sandbox"`
}

// PayShapConfig holds PayShap configuration
type PayShapConfig struct {
	MerchantID string `json:"merchant_id"`
	BankID     string `json:"bank_id"`
	APIKey     string `json:"api_key,omitempty"`
	Sandbox    bool   `json:"sandbox"`
}

// Config holds all payments configuration
type Config struct {
	PayFast *PayFastConfig `json:"payfast,omitempty"`
	Ozow    *OzowConfig    `json:"ozow,omitempty"`
	Yoco    *YocoConfig    `json:"yoco,omitempty"`
	PayShap *PayShapConfig `json:"payshap,omitempty"`
}

// PayFastPaymentRequest represents a PayFast payment request
type PayFastPaymentRequest struct {
	Amount          float64 `json:"amount"`
	ItemName        string  `json:"item_name"`
	ReturnURL       string  `json:"return_url"`
	CancelURL       string  `json:"cancel_url"`
	NotifyURL       string  `json:"notify_url"`
	ItemDescription string  `json:"item_description,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	EmailAddress    string  `json:"email_address,omitempty"`
	CellNumber      string  `json:"cell_number,omitempty"`
	PaymentMethod   string  `json:"payment_method,omitempty"`
	CustomStr1      string  `json:"custom_str1,omitempty"`
	CustomStr2      string  `json:"custom_str2,omitempty"`
	CustomStr3      string  `json:"custom_str3,omitempty"`
	CustomInt1      int     `json:"custom_int1,omitempty"`
	CustomInt2      int     `json:"custom_int2,omitempty"`
}

// PayFastPayment represents a PayFast payment
type PayFastPayment struct {
	ID          string        `json:"id"`
	Provider    string        `json:"provider"`
	Amount      Amount        `json:"amount"`
	Status      PaymentStatus `json:"status"`
	Reference   string        `json:"reference,omitempty"`
	Description string        `json:"description,omitempty"`
	PaymentURL  string        `json:"payment_url"`
	Signature   string        `json:"signature,omitempty"`
	PFPaymentID string        `json:"pf_payment_id,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// OzowPaymentRequest represents an Ozow payment request
type OzowPaymentRequest struct {
	Amount               float64 `json:"amount"`
	TransactionReference string  `json:"transaction_reference"`
	BankReference        string  `json:"bank_reference"`
	SuccessURL           string  `json:"success_url"`
	CancelURL            string  `json:"cancel_url"`
	ErrorURL             string  `json:"error_url"`
	NotifyURL            string  `json:"notify_url"`
	IsTest               *bool   `json:"is_test,omitempty"`
	CustomerFirstName    string  `json:"customer_first_name,omitempty"`
	CustomerLastName     string  `json:"customer_last_name,omitempty"`
	CustomerEmail        string  `json:"customer_email,omitempty"`
	CustomerPhone        string  `json:"customer_phone,omitempty"`
}

// OzowPayment represents an Ozow payment
type OzowPayment struct {
	ID            string        `json:"id"`
	Provider      string        `json:"provider"`
	Amount        Amount        `json:"amount"`
	Status        PaymentStatus `json:"status"`
	Reference     string        `json:"reference,omitempty"`
	Description   string        `json:"description,omitempty"`
	PaymentURL    string        `json:"payment_url"`
	TransactionID string        `json:"transaction_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// YocoPaymentRequest represents a Yoco payment request
type YocoPaymentRequest struct {
	Amount     int                    `json:"amount"` // In cents
	SuccessURL string                 `json:"success_url"`
	CancelURL  string                 `json:"cancel_url"`
	FailureURL string                 `json:"failure_url,omitempty"`
	Currency   string                 `json:"currency,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	LineItems  []YocoLineItem         `json:"line_items,omitempty"`
}

// YocoLineItem represents a line item in Yoco checkout
type YocoLineItem struct {
	DisplayName    string `json:"displayName"`
	Quantity       int    `json:"quantity"`
	PricingDetails struct {
		Price int `json:"price"`
	} `json:"pricingDetails"`
}

// YocoPayment represents a Yoco payment
type YocoPayment struct {
	ID          string                 `json:"id"`
	Provider    string                 `json:"provider"`
	CheckoutID  string                 `json:"checkout_id"`
	Amount      Amount                 `json:"amount"`
	Status      PaymentStatus          `json:"status"`
	Reference   string                 `json:"reference,omitempty"`
	RedirectURL string                 `json:"redirect_url"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// PayShapPaymentRequest represents a PayShap payment request
type PayShapPaymentRequest struct {
	Amount           float64 `json:"amount"`
	Reference        string  `json:"reference"`
	ShapID           string  `json:"shap_id,omitempty"`
	RecipientAccount string  `json:"recipient_account,omitempty"`
	RecipientBank    string  `json:"recipient_bank,omitempty"`
	ExpiryMinutes    int     `json:"expiry_minutes,omitempty"`
}

// PayShapPayment represents a PayShap payment
type PayShapPayment struct {
	ID        string        `json:"id"`
	Provider  string        `json:"provider"`
	ShapID    string        `json:"shap_id"`
	Amount    Amount        `json:"amount"`
	Status    PaymentStatus `json:"status"`
	Reference string        `json:"reference,omitempty"`
	QRCode    string        `json:"qr_code,omitempty"`
	ExpiresAt time.Time     `json:"expires_at,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID string `json:"payment_id"`
	Amount    *int   `json:"amount,omitempty"` // Partial refund in cents
	Reason    string `json:"reason,omitempty"`
}

// Refund represents a refund
type Refund struct {
	ID        string    `json:"id"`
	PaymentID string    `json:"payment_id"`
	Amount    Amount    `json:"amount"`
	Status    string    `json:"status"` // pending, completed, failed
	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
