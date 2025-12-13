// Package zatools provides South African business tool integrations.
package zatools

import (
	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// ZATools provides unified access to SA business tools.
type ZATools struct {
	client       *lunahttp.Client
	config       *Config
	cipc         *CIPC
	bbbee        *BBBEE
	idValidation *IDValidation
	address      *AddressUtils
}

// NewZATools creates a new ZATools resource.
func NewZATools(client *lunahttp.Client, config *Config) *ZATools {
	if config == nil {
		config = &Config{}
	}
	return &ZATools{
		client: client,
		config: config,
	}
}

// CIPC returns the CIPC service instance.
func (z *ZATools) CIPC() *CIPC {
	if z.cipc == nil {
		config := z.config.CIPC
		if config == nil {
			config = &CIPCConfig{}
		}
		z.cipc = NewCIPC(z.client, *config, z.config.Strict)
	}
	return z.cipc
}

// BBBEE returns the B-BBEE service instance.
func (z *ZATools) BBBEE() *BBBEE {
	if z.bbbee == nil {
		config := z.config.BBBEE
		if config == nil {
			config = &BBBEEConfig{}
		}
		z.bbbee = NewBBBEE(z.client, *config)
	}
	return z.bbbee
}

// IDValidation returns the ID validation utility.
func (z *ZATools) IDValidation() *IDValidation {
	if z.idValidation == nil {
		z.idValidation = NewIDValidation()
	}
	return z.idValidation
}

// Address returns the address utilities.
func (z *ZATools) Address() *AddressUtils {
	if z.address == nil {
		z.address = NewAddressUtils()
	}
	return z.address
}

// ValidateID is a convenience method to validate SA ID number.
func (z *ZATools) ValidateID(idNumber string) SAIDInfo {
	return z.IDValidation().Validate(idNumber)
}

// List returns available ZA tools.
func (z *ZATools) List() []string {
	return []string{"cipc", "bbbee", "id_validation", "address"}
}
