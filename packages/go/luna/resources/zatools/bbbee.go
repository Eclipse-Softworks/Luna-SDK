// Package zatools provides South African business tool integrations.
package zatools

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// BBBEE provides B-BBEE compliance verification.
type BBBEE struct {
	client *lunahttp.Client
	config BBBEEConfig
}

// NewBBBEE creates a new B-BBEE instance.
func NewBBBEE(client *lunahttp.Client, config BBBEEConfig) *BBBEE {
	return &BBBEE{
		client: client,
		config: config,
	}
}

// VerifyCertificate verifies a B-BBEE certificate by number.
func (b *BBBEE) VerifyCertificate(ctx context.Context, certificateNumber string) (*BBBEECertificate, error) {
	resp, err := b.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   "/v1/za/bbbee/certificates/" + certificateNumber,
	})
	if err != nil {
		return nil, err
	}

	var cert BBBEECertificate
	if err := json.Unmarshal(resp.Data, &cert); err != nil {
		return nil, err
	}

	return &cert, nil
}

// LookupByCompany looks up B-BBEE status by company registration number.
func (b *BBBEE) LookupByCompany(ctx context.Context, registrationNumber string) (*BBBEECertificate, error) {
	query := url.Values{}
	query.Set("registrationNumber", registrationNumber)

	resp, err := b.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   "/v1/za/bbbee/lookup",
		Query:  query,
	})
	if err != nil {
		return nil, err
	}

	var cert BBBEECertificate
	if err := json.Unmarshal(resp.Data, &cert); err != nil {
		return nil, err
	}

	return &cert, nil
}

// MeetsRequirement checks if a company meets minimum B-BBEE level requirement.
func (b *BBBEE) MeetsRequirement(ctx context.Context, registrationNumber string, minimumLevel int) (map[string]interface{}, error) {
	certificate, _ := b.LookupByCompany(ctx, registrationNumber)

	if certificate == nil || !certificate.IsValid {
		return map[string]interface{}{"meets": false}, nil
	}

	levelNumber := 9
	if level, ok := certificate.Level.(int); ok {
		levelNumber = level
	}

	return map[string]interface{}{
		"meets":        levelNumber <= minimumLevel,
		"actual_level": certificate.Level,
		"certificate":  certificate,
	}, nil
}

// GetRecognitionLevel gets procurement recognition percentage.
func (b *BBBEE) GetRecognitionLevel(level interface{}) int {
	recognitionMap := map[interface{}]int{
		1:               135,
		2:               125,
		3:               110,
		4:               100,
		5:               80,
		6:               60,
		7:               50,
		8:               10,
		"non-compliant": 0,
	}

	if r, ok := recognitionMap[level]; ok {
		return r
	}
	return 0
}

// DetermineEnterpriseCategory calculates EME or QSE status.
func (b *BBBEE) DetermineEnterpriseCategory(annualRevenue float64) map[string]interface{} {
	if annualRevenue <= 10_000_000 {
		return map[string]interface{}{
			"category":        "EME",
			"automatic_level": 4,
			"description":     "Exempted Micro Enterprise - automatic Level 4 (or Level 1 with 51%+ black ownership)",
		}
	} else if annualRevenue <= 50_000_000 {
		return map[string]interface{}{
			"category":        "QSE",
			"automatic_level": nil,
			"description":     "Qualifying Small Enterprise - simplified scorecard applies",
		}
	}
	return map[string]interface{}{
		"category":        "Generic",
		"automatic_level": nil,
		"description":     "Generic Enterprise - full scorecard applies",
	}
}

// CalculateScore calculates scorecard points.
func (b *BBBEE) CalculateScore(scorecard BBBEEScorecard) map[string]interface{} {
	ownership := min(scorecard.Ownership, 25)
	managementControl := min(scorecard.ManagementControl, 19)
	skillsDevelopment := min(scorecard.SkillsDevelopment, 20)
	enterpriseSupplierDevelopment := min(scorecard.EnterpriseSupplierDevelopment, 40)
	socioEconomicDevelopment := min(scorecard.SocioEconomicDevelopment, 5)

	total := ownership + managementControl + skillsDevelopment + enterpriseSupplierDevelopment + socioEconomicDevelopment
	level := b.pointsToLevel(total)

	return map[string]interface{}{
		"total": total,
		"level": level,
	}
}

func (b *BBBEE) pointsToLevel(points float64) interface{} {
	switch {
	case points >= 100:
		return 1
	case points >= 95:
		return 2
	case points >= 90:
		return 3
	case points >= 80:
		return 4
	case points >= 75:
		return 5
	case points >= 70:
		return 6
	case points >= 55:
		return 7
	case points >= 40:
		return 8
	default:
		return "non-compliant"
	}
}

// IsCertificateValid checks if a certificate is still valid.
func (b *BBBEE) IsCertificateValid(certificate BBBEECertificate) bool {
	expiryDate, err := time.Parse(time.RFC3339, certificate.ExpiryDate)
	if err != nil {
		return false
	}
	return certificate.IsValid && expiryDate.After(time.Now())
}

// GetDaysUntilExpiry gets days until certificate expiry.
func (b *BBBEE) GetDaysUntilExpiry(certificate BBBEECertificate) int {
	expiryDate, err := time.Parse(time.RFC3339, certificate.ExpiryDate)
	if err != nil {
		return 0
	}
	return int(expiryDate.Sub(time.Now()).Hours() / 24)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
