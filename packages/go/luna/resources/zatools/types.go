// Package zatools provides South African business tool integrations.
package zatools

import "time"

// CompanyType represents CIPC company types
type CompanyType string

const (
	CompanyPTYLTD CompanyType = "PTY_LTD"
	CompanyLTD    CompanyType = "LTD"
	CompanyNPC    CompanyType = "NPC"
	CompanyCC     CompanyType = "CC"
	CompanyINC    CompanyType = "INC"
)

// CompanyStatus represents company registration status
type CompanyStatus string

const (
	StatusActive         CompanyStatus = "active"
	StatusBusinessRescue CompanyStatus = "in_business_rescue"
	StatusDeregistered   CompanyStatus = "deregistered"
	StatusLiquidated     CompanyStatus = "liquidated"
	StatusDissolved      CompanyStatus = "dissolved"
)

// SAProvince represents South African provinces
type SAProvince string

const (
	ProvinceEC  SAProvince = "EC"  // Eastern Cape
	ProvinceFS  SAProvince = "FS"  // Free State
	ProvinceGP  SAProvince = "GP"  // Gauteng
	ProvinceKZN SAProvince = "KZN" // KwaZulu-Natal
	ProvinceLP  SAProvince = "LP"  // Limpopo
	ProvinceMP  SAProvince = "MP"  // Mpumalanga
	ProvinceNC  SAProvince = "NC"  // Northern Cape
	ProvinceNW  SAProvince = "NW"  // North West
	ProvinceWC  SAProvince = "WC"  // Western Cape
)

// BBBEELevel represents B-BBEE compliance levels
type BBBEELevel interface{}

// Director represents a company director
type Director struct {
	Name          string `json:"name"`
	IDNumber      string `json:"id_number,omitempty"`
	Role          string `json:"role"`
	AppointedDate string `json:"appointed_date,omitempty"`
	ResignedDate  string `json:"resigned_date,omitempty"`
}

// Company represents CIPC company information
type Company struct {
	RegistrationNumber string        `json:"registration_number"`
	Name               string        `json:"name"`
	Type               CompanyType   `json:"type"`
	Status             CompanyStatus `json:"status"`
	RegistrationDate   string        `json:"registration_date,omitempty"`
	FinancialYearEnd   string        `json:"financial_year_end,omitempty"`
	Directors          []Director    `json:"directors"`
	PhysicalAddress    string        `json:"physical_address,omitempty"`
	PostalAddress      string        `json:"postal_address,omitempty"`
}

// BBBEECertificate represents a B-BBEE certificate
type BBBEECertificate struct {
	CertificateNumber  string                 `json:"certificate_number"`
	CompanyName        string                 `json:"company_name"`
	RegistrationNumber string                 `json:"registration_number"`
	Level              BBBEELevel             `json:"level"`
	VerificationAgency string                 `json:"verification_agency"`
	IssueDate          string                 `json:"issue_date"`
	ExpiryDate         string                 `json:"expiry_date"`
	IsValid            bool                   `json:"is_valid"`
	Scorecard          map[string]interface{} `json:"scorecard,omitempty"`
}

// BBBEEScorecard represents B-BBEE scorecard breakdown
type BBBEEScorecard struct {
	Ownership                     float64 `json:"ownership"`
	ManagementControl             float64 `json:"management_control"`
	SkillsDevelopment             float64 `json:"skills_development"`
	EnterpriseSupplierDevelopment float64 `json:"enterprise_supplier_development"`
	SocioEconomicDevelopment      float64 `json:"socio_economic_development"`
}

// SAIDInfo represents South African ID information
type SAIDInfo struct {
	IDNumber      string    `json:"id_number"`
	IsValid       bool      `json:"is_valid"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Gender        string    `json:"gender"` // male or female
	IsSACitizen   bool      `json:"is_sa_citizen"`
	ChecksumValid bool      `json:"checksum_valid"`
}

// SAAddress represents a South African address
type SAAddress struct {
	Street     string     `json:"street,omitempty"`
	Suburb     string     `json:"suburb,omitempty"`
	City       string     `json:"city,omitempty"`
	Province   SAProvince `json:"province,omitempty"`
	PostalCode string     `json:"postal_code,omitempty"`
	Country    string     `json:"country"`
}

// CIPCConfig holds CIPC configuration
type CIPCConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	Sandbox bool   `json:"sandbox"`
}

// BBBEEConfig holds B-BBEE configuration
type BBBEEConfig struct {
	APIKey  string `json:"api_key,omitempty"`
	Sandbox bool   `json:"sandbox"`
}

// Config holds configuration for ZA tools.
type Config struct {
	CIPC   *CIPCConfig
	BBBEE  *BBBEEConfig
	Strict bool
}
