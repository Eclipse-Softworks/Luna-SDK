// Package utils provides utility functions for the Luna SDK.
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateRequestID generates a unique request ID.
func GenerateRequestID() string {
	timestamp := fmt.Sprintf("%x", time.Now().UnixMilli())
	random := make([]byte, 4)
	_, _ = rand.Read(random)
	return fmt.Sprintf("req_%s%s", timestamp, hex.EncodeToString(random))
}

// MaskSensitive masks a sensitive string, showing only start and end characters.
func MaskSensitive(value string, visibleStart, visibleEnd int) string {
	if len(value) <= visibleStart+visibleEnd {
		return strings.Repeat("*", len(value))
	}
	return value[:visibleStart] + "****" + value[len(value)-visibleEnd:]
}

// ValidateID validates an ID has the correct format.
func ValidateID(id, prefix, name string) error {
	if id == "" {
		return fmt.Errorf("%s is required", name)
	}

	pattern := regexp.MustCompile(fmt.Sprintf(`^%s_[a-zA-Z0-9]+$`, prefix))
	if !pattern.MatchString(id) {
		return fmt.Errorf("invalid %s format, expected: %s_<id>", name, prefix)
	}

	return nil
}

// IsRetryableStatus checks if an HTTP status code is retryable.
func IsRetryableStatus(status int) bool {
	switch status {
	case 408, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// DeepMerge merges two maps deeply.
func DeepMerge(base, override map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range base {
		result[k] = v
	}

	for k, v := range override {
		if baseVal, ok := result[k]; ok {
			if baseMap, ok := baseVal.(map[string]interface{}); ok {
				if overrideMap, ok := v.(map[string]interface{}); ok {
					result[k] = DeepMerge(baseMap, overrideMap)
					continue
				}
			}
		}
		result[k] = v
	}

	return result
}

// StringPtr returns a pointer to a string.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int.
func IntPtr(i int) *int {
	return &i
}
