// Package zatools provides South African business tool integrations.
package zatools

import (
	"context"
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// CIPC provides Companies and Intellectual Property Commission integration.
type CIPC struct {
	client *lunahttp.Client
	config CIPCConfig
	strict bool
}

// NewCIPC creates a new CIPC service instance.
func NewCIPC(client *lunahttp.Client, config CIPCConfig, strict bool) *CIPC {
	return &CIPC{
		client: client,
		config: config,
		strict: strict,
	}
}

// Lookup searches for a company by registration number.
func (c *CIPC) Lookup(ctx context.Context, registrationNumber string) (*Company, error) {
	cleaned := regexp.MustCompile(`[\s/]`).ReplaceAllString(registrationNumber, "")

	// Strict Validation ("Rust" safety)
	if c.strict && !c.IsValidRegistrationNumber(cleaned) {
		return nil, &ValidationError{"invalid registration number format (strict mode)"}
	}

	resp, err := c.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   "/v1/za/cipc/companies/" + cleaned,
	})
	if err != nil {
		return nil, err
	}

	var company Company
	if err := json.Unmarshal(resp.Data, &company); err != nil {
		return nil, &ValidationError{"failed to parse response"}
	}

	return &company, nil
}

// SearchByName searches for companies by name.
func (c *CIPC) SearchByName(ctx context.Context, name string, limit int) ([]Company, error) {
	if limit == 0 {
		limit = 10
	}
	query := url.Values{}
	query.Set("name", name)
	query.Set("limit", strconv.Itoa(limit))

	resp, err := c.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   "/v1/za/cipc/companies",
		Query:  query,
	})
	if err != nil {
		return []Company{}, err
	}

	var result struct {
		Data []Company `json:"data"`
	}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return []Company{}, nil
	}

	return result.Data, nil
}

// Verify checks if a company is registered and active.
func (c *CIPC) Verify(ctx context.Context, registrationNumber string) (map[string]interface{}, error) {
	company, _ := c.Lookup(ctx, registrationNumber)

	if company == nil {
		return map[string]interface{}{
			"exists":    false,
			"is_active": false,
		}, nil
	}

	return map[string]interface{}{
		"exists":    true,
		"is_active": company.Status == StatusActive,
		"company":   company,
	}, nil
}

// CheckNameAvailability checks if a company name is available.
func (c *CIPC) CheckNameAvailability(ctx context.Context, name string) (map[string]interface{}, error) {
	companies, _ := c.SearchByName(ctx, name, 5)

	exactMatch := false
	var similarNames []string
	for _, comp := range companies {
		if comp.Name == name {
			exactMatch = true
		}
		similarNames = append(similarNames, comp.Name)
	}

	return map[string]interface{}{
		"available":     !exactMatch,
		"similar_names": similarNames,
	}, nil
}

// GetDirectors gets directors for a company.
func (c *CIPC) GetDirectors(ctx context.Context, registrationNumber string) ([]Director, error) {
	company, _ := c.Lookup(ctx, registrationNumber)
	if company == nil {
		return []Director{}, nil
	}
	return company.Directors, nil
}

// IsValidRegistrationNumber validates registration number format.
func (c *CIPC) IsValidRegistrationNumber(regNumber string) bool {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^\d{4}/\d{6}/\d{2}$`),
		regexp.MustCompile(`^\d{12}$`),
		regexp.MustCompile(`^[A-Z]{2}\d{6}$`),
	}

	for _, p := range patterns {
		if p.MatchString(regNumber) {
			return true
		}
	}
	return false
}

// ParseCompanyType parses company type from registration number.
func (c *CIPC) ParseCompanyType(regNumber string) CompanyType {
	match := regexp.MustCompile(`/(\d{2})$`).FindStringSubmatch(regNumber)
	if len(match) < 2 {
		return ""
	}

	typeCode := match[1]
	typeMap := map[string]CompanyType{
		"07": CompanyPTYLTD,
		"06": CompanyLTD,
		"08": CompanyNPC,
		"23": CompanyCC,
		"21": CompanyINC,
	}

	if t, ok := typeMap[typeCode]; ok {
		return t
	}
	return ""
}

// FormatRegistrationNumber formats registration number for display.
func (c *CIPC) FormatRegistrationNumber(regNumber string) string {
	cleaned := regexp.MustCompile(`[\s/]`).ReplaceAllString(regNumber, "")

	if len(cleaned) == 12 {
		return cleaned[:4] + "/" + cleaned[4:10] + "/" + cleaned[10:]
	}

	return regNumber
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
