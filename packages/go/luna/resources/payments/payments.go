// Package payments provides South African payment gateway integrations.
package payments

import (
	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// Payments provides unified access to SA payment gateways.
type Payments struct {
	client  *lunahttp.Client
	config  *Config
	payfast *PayFast
	ozow    *Ozow
	yoco    *Yoco
	payshap *PayShap
}

// NewPayments creates a new Payments resource.
func NewPayments(client *lunahttp.Client, config *Config) *Payments {
	if config == nil {
		config = &Config{}
	}
	return &Payments{
		client: client,
		config: config,
	}
}

// PayFast returns the PayFast gateway instance.
func (p *Payments) PayFast() *PayFast {
	if p.payfast == nil {
		if p.config.PayFast == nil {
			panic("PayFast not configured. Provide PayFastConfig when initializing LunaClient.")
		}
		p.payfast = NewPayFast(p.client, *p.config.PayFast)
	}
	return p.payfast
}

// Ozow returns the Ozow gateway instance.
func (p *Payments) Ozow() *Ozow {
	if p.ozow == nil {
		if p.config.Ozow == nil {
			panic("Ozow not configured. Provide OzowConfig when initializing LunaClient.")
		}
		p.ozow = NewOzow(p.client, *p.config.Ozow)
	}
	return p.ozow
}

// Yoco returns the Yoco gateway instance.
func (p *Payments) Yoco() *Yoco {
	if p.yoco == nil {
		if p.config.Yoco == nil {
			panic("Yoco not configured. Provide YocoConfig when initializing LunaClient.")
		}
		p.yoco = NewYoco(p.client, *p.config.Yoco)
	}
	return p.yoco
}

// PayShap returns the PayShap gateway instance.
func (p *Payments) PayShap() *PayShap {
	if p.payshap == nil {
		if p.config.PayShap == nil {
			panic("PayShap not configured. Provide PayShapConfig when initializing LunaClient.")
		}
		p.payshap = NewPayShap(p.client, *p.config.PayShap)
	}
	return p.payshap
}

// List returns available payment gateways.
func (p *Payments) List() []string {
	var available []string
	if p.config.PayFast != nil {
		available = append(available, "payfast")
	}
	if p.config.Ozow != nil {
		available = append(available, "ozow")
	}
	if p.config.Yoco != nil {
		available = append(available, "yoco")
	}
	if p.config.PayShap != nil {
		available = append(available, "payshap")
	}
	return available
}
