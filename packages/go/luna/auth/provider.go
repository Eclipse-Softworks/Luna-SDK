// Package auth provides authentication providers for the Luna SDK.
package auth

// Provider is the interface for authentication providers
type Provider interface {
	// GetHeaders returns authorization headers for a request
	GetHeaders() (map[string]string, error)

	// NeedsRefresh returns true if credentials need refresh
	NeedsRefresh() bool

	// Refresh refreshes the credentials
	Refresh() error
}
