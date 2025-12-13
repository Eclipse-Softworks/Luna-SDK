// Package zatools provides South African business tool integrations.
package zatools

import (
	"regexp"
	"strings"
)

// AddressUtils provides South African address utilities.
type AddressUtils struct{}

// NewAddressUtils creates a new address utils instance.
func NewAddressUtils() *AddressUtils {
	return &AddressUtils{}
}

// Postal code ranges by province
var postalCodeRanges = map[SAProvince][][2]int{
	ProvinceEC:  {{5000, 5999}, {6000, 6499}},
	ProvinceFS:  {{9300, 9999}},
	ProvinceGP:  {{1, 299}, {1000, 2199}},
	ProvinceKZN: {{3000, 4499}},
	ProvinceLP:  {{500, 999}},
	ProvinceMP:  {{1100, 1399}, {2200, 2499}},
	ProvinceNC:  {{8000, 8999}},
	ProvinceNW:  {{2500, 2999}, {300, 499}},
	ProvinceWC:  {{6500, 8099}},
}

// SA cities to province mapping
var saCities = map[string]SAProvince{
	"johannesburg":   ProvinceGP,
	"pretoria":       ProvinceGP,
	"cape town":      ProvinceWC,
	"durban":         ProvinceKZN,
	"port elizabeth": ProvinceEC,
	"gqeberha":       ProvinceEC,
	"bloemfontein":   ProvinceFS,
	"polokwane":      ProvinceLP,
	"nelspruit":      ProvinceMP,
	"mbombela":       ProvinceMP,
	"kimberley":      ProvinceNC,
	"rustenburg":     ProvinceNW,
	"soweto":         ProvinceGP,
	"sandton":        ProvinceGP,
	"centurion":      ProvinceGP,
}

// Province names
var provinceNames = map[SAProvince]string{
	ProvinceEC:  "Eastern Cape",
	ProvinceFS:  "Free State",
	ProvinceGP:  "Gauteng",
	ProvinceKZN: "KwaZulu-Natal",
	ProvinceLP:  "Limpopo",
	ProvinceMP:  "Mpumalanga",
	ProvinceNC:  "Northern Cape",
	ProvinceNW:  "North West",
	ProvinceWC:  "Western Cape",
}

// Validate validates a South African address.
func (a *AddressUtils) Validate(address SAAddress) map[string]interface{} {
	var errors []string
	var warnings []string

	// Validate postal code
	if address.PostalCode != "" {
		if !regexp.MustCompile(`^\d{4}$`).MatchString(address.PostalCode) {
			errors = append(errors, "Postal code must be 4 digits")
		} else {
			detectedProvince := a.GetProvinceFromPostalCode(address.PostalCode)
			if address.Province != "" && detectedProvince != "" && address.Province != detectedProvince {
				warnings = append(warnings, "Postal code province mismatch")
			}
		}
	}

	// Validate province
	if address.Province != "" {
		if _, ok := provinceNames[address.Province]; !ok {
			errors = append(errors, "Invalid province code")
		}
	}

	return map[string]interface{}{
		"valid":    len(errors) == 0,
		"errors":   errors,
		"warnings": warnings,
	}
}

// Format formats address for display.
func (a *AddressUtils) Format(address SAAddress, multiline bool) string {
	var parts []string

	if address.Street != "" {
		parts = append(parts, address.Street)
	}
	if address.Suburb != "" {
		parts = append(parts, address.Suburb)
	}
	if address.City != "" {
		parts = append(parts, address.City)
	}
	if address.Province != "" {
		if name, ok := provinceNames[address.Province]; ok {
			parts = append(parts, name)
		} else {
			parts = append(parts, string(address.Province))
		}
	}
	if address.PostalCode != "" {
		parts = append(parts, address.PostalCode)
	}
	parts = append(parts, address.Country)

	separator := ", "
	if multiline {
		separator = "\n"
	}
	return strings.Join(parts, separator)
}

// GetProvinceFromPostalCode detects province from postal code.
func (a *AddressUtils) GetProvinceFromPostalCode(postalCode string) SAProvince {
	if !regexp.MustCompile(`^\d{4}$`).MatchString(postalCode) {
		return ""
	}

	var code int
	for _, c := range postalCode {
		code = code*10 + int(c-'0')
	}

	for province, ranges := range postalCodeRanges {
		for _, r := range ranges {
			if code >= r[0] && code <= r[1] {
				return province
			}
		}
	}

	return ""
}

// GetProvinceName gets full province name from code.
func (a *AddressUtils) GetProvinceName(code SAProvince) string {
	if name, ok := provinceNames[code]; ok {
		return name
	}
	return string(code)
}

// LookupPostalCode looks up postal code information.
func (a *AddressUtils) LookupPostalCode(postalCode string) map[string]interface{} {
	if !regexp.MustCompile(`^\d{4}$`).MatchString(postalCode) {
		return nil
	}

	province := a.GetProvinceFromPostalCode(postalCode)
	if province == "" {
		return nil
	}

	return map[string]interface{}{
		"postal_code":   postalCode,
		"province":      province,
		"province_name": provinceNames[province],
	}
}

// Parse parses a free-form address string.
func (a *AddressUtils) Parse(addressString string) SAAddress {
	result := SAAddress{Country: "ZA"}

	// Try to extract postal code
	postalMatch := regexp.MustCompile(`\b(\d{4})\b`).FindStringSubmatch(addressString)
	if len(postalMatch) >= 2 {
		result.PostalCode = postalMatch[1]
		result.Province = a.GetProvinceFromPostalCode(postalMatch[1])
	}

	// Try to extract city
	lower := strings.ToLower(addressString)
	for city, province := range saCities {
		if strings.Contains(lower, city) {
			result.City = strings.Title(city)
			if result.Province == "" {
				result.Province = province
			}
			break
		}
	}

	return result
}
