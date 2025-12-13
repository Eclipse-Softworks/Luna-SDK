// Package zatools provides South African business tool integrations.
package zatools

import (
	"regexp"
	"strconv"
	"time"
)

// IDValidation provides South African ID number validation.
type IDValidation struct{}

// NewIDValidation creates a new ID validation instance.
func NewIDValidation() *IDValidation {
	return &IDValidation{}
}

// Validate validates and parses a South African ID number.
func (v *IDValidation) Validate(idNumber string) SAIDInfo {
	// Remove spaces and dashes
	cleaned := regexp.MustCompile(`[\s-]`).ReplaceAllString(idNumber, "")

	if len(cleaned) != 13 {
		return v.invalidResult(idNumber)
	}

	if !regexp.MustCompile(`^\d{13}$`).MatchString(cleaned) {
		return v.invalidResult(idNumber)
	}

	// Extract components
	year, _ := strconv.Atoi(cleaned[0:2])
	month, _ := strconv.Atoi(cleaned[2:4])
	day, _ := strconv.Atoi(cleaned[4:6])
	genderDigits, _ := strconv.Atoi(cleaned[6:10])
	citizenshipDigit, _ := strconv.Atoi(cleaned[10:11])

	// Determine full year
	fullYear := 1900 + year
	if year <= 30 {
		fullYear = 2000 + year
	}

	// Validate date
	dateOfBirth := time.Date(fullYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if dateOfBirth.Month() != time.Month(month) || dateOfBirth.Day() != day {
		return v.invalidResult(idNumber)
	}

	// Validate checksum
	checksumValid := v.validateLuhn(cleaned)

	gender := "female"
	if genderDigits >= 5000 {
		gender = "male"
	}

	return SAIDInfo{
		IDNumber:      cleaned,
		IsValid:       checksumValid,
		DateOfBirth:   dateOfBirth,
		Gender:        gender,
		IsSACitizen:   citizenshipDigit == 0,
		ChecksumValid: checksumValid,
	}
}

// IsValid returns whether the ID number is valid.
func (v *IDValidation) IsValid(idNumber string) bool {
	return v.Validate(idNumber).IsValid
}

// GetDateOfBirth extracts date of birth from ID number.
func (v *IDValidation) GetDateOfBirth(idNumber string) *time.Time {
	result := v.Validate(idNumber)
	if !result.IsValid {
		return nil
	}
	return &result.DateOfBirth
}

// GetAge gets age from ID number.
func (v *IDValidation) GetAge(idNumber string) *int {
	dob := v.GetDateOfBirth(idNumber)
	if dob == nil {
		return nil
	}

	now := time.Now()
	age := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	return &age
}

// IsOver checks if person is over a specific age.
func (v *IDValidation) IsOver(idNumber string, age int) bool {
	currentAge := v.GetAge(idNumber)
	return currentAge != nil && *currentAge >= age
}

// Generate generates a valid SA ID number.
func (v *IDValidation) Generate(dateOfBirth *time.Time, gender string, isCitizen bool) string {
	dob := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	if dateOfBirth != nil {
		dob = *dateOfBirth
	}

	year := dob.Year() % 100
	month := int(dob.Month())
	day := dob.Day()

	// Gender sequence
	genderMin, genderMax := 5000, 9999
	if gender == "female" {
		genderMin, genderMax = 0, 4999
	}
	genderSeq := genderMin + (time.Now().Nanosecond() % (genderMax - genderMin + 1))

	// Citizenship
	citizen := "0"
	if !isCitizen {
		citizen = "1"
	}

	// A digit
	aDigit := "8"

	// Build partial ID
	partial := sprintf("%02d%02d%02d%04d%s%s", year, month, day, genderSeq, citizen, aDigit)
	checksum := v.calculateLuhnChecksum(partial)

	return partial + checksum
}

func sprintf(format string, args ...interface{}) string {
	result := ""
	argIdx := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) {
			switch format[i+1] {
			case '0':
				// Handle %02d, %04d
				width := 0
				j := i + 1
				for j < len(format) && format[j] >= '0' && format[j] <= '9' {
					width = width*10 + int(format[j]-'0')
					j++
				}
				if j < len(format) && format[j] == 'd' {
					val := args[argIdx].(int)
					strVal := strconv.Itoa(val)
					for len(strVal) < width {
						strVal = "0" + strVal
					}
					result += strVal
					argIdx++
					i = j
					continue
				}
			case 's':
				result += args[argIdx].(string)
				argIdx++
				i++
				continue
			case 'd':
				result += strconv.Itoa(args[argIdx].(int))
				argIdx++
				i++
				continue
			}
		}
		result += string(format[i])
	}
	return result
}

func (v *IDValidation) invalidResult(idNumber string) SAIDInfo {
	return SAIDInfo{
		IDNumber:      idNumber,
		IsValid:       false,
		DateOfBirth:   time.Time{},
		Gender:        "male",
		IsSACitizen:   false,
		ChecksumValid: false,
	}
}

func (v *IDValidation) validateLuhn(number string) bool {
	sum := 0
	isEven := false

	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))

		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}

func (v *IDValidation) calculateLuhnChecksum(partial string) string {
	for i := 0; i <= 9; i++ {
		if v.validateLuhn(partial + strconv.Itoa(i)) {
			return strconv.Itoa(i)
		}
	}
	return "0"
}
