package testutil

import (
	"encoding/json"
	"time"
)

// User fixtures
var MockUser = map[string]interface{}{
	"id":         "usr_123456789",
	"name":       "John Doe",
	"email":      "john@example.com",
	"avatar":     "https://example.com/avatar.jpg",
	"created_at": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	"updated_at": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
}

var MockUsers = []map[string]interface{}{
	MockUser,
	{
		"id":         "usr_987654321",
		"name":       "Jane Smith",
		"email":      "jane@example.com",
		"avatar":     nil,
		"created_at": time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		"updated_at": time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	},
}

var MockUserCreate = map[string]string{
	"name":  "New User",
	"email": "newuser@example.com",
}

// Project fixtures
var MockProject = map[string]interface{}{
	"id":          "prj_123456789",
	"name":        "Test Project",
	"description": "A test project for unit tests",
	"owner_id":    "usr_123456789",
	"created_at":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	"updated_at":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
}

var MockProjects = []map[string]interface{}{
	MockProject,
	{
		"id":          "prj_987654321",
		"name":        "Another Project",
		"description": "Another test project",
		"owner_id":    "usr_987654321",
		"created_at":  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		"updated_at":  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	},
}

// Storage fixtures
var MockBucket = map[string]interface{}{
	"id":         "bkt_123456789",
	"name":       "test-bucket",
	"public":     false,
	"created_at": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
}

var MockBuckets = []map[string]interface{}{
	MockBucket,
	{
		"id":         "bkt_987654321",
		"name":       "public-bucket",
		"public":     true,
		"created_at": time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	},
}

var MockFile = map[string]interface{}{
	"id":         "file_123456789",
	"name":       "test-file.pdf",
	"bucket_id":  "bkt_123456789",
	"size":       1024,
	"mime_type":  "application/pdf",
	"created_at": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
}

// Helper functions
func MockListResponse(data interface{}, hasMore bool, nextCursor string) map[string]interface{} {
	response := map[string]interface{}{
		"data":     data,
		"has_more": hasMore,
	}
	if nextCursor != "" {
		response["next_cursor"] = nextCursor
	}
	return response
}

// Error response fixtures
var MockErrorNotFound = map[string]interface{}{
	"error": map[string]interface{}{
		"message": "Resource not found",
		"code":    "NOT_FOUND",
		"status":  404,
	},
}

var MockErrorValidation = map[string]interface{}{
	"error": map[string]interface{}{
		"message": "Validation failed",
		"code":    "VALIDATION_ERROR",
		"status":  400,
		"details": []map[string]string{
			{"field": "email", "message": "Invalid email format"},
		},
	},
}

var MockErrorRateLimit = map[string]interface{}{
	"error": map[string]interface{}{
		"message":     "Rate limit exceeded",
		"code":        "RATE_LIMIT_EXCEEDED",
		"status":      429,
		"retry_after": 60,
	},
}

var MockErrorAuth = map[string]interface{}{
	"error": map[string]interface{}{
		"message": "Invalid API key",
		"code":    "AUTHENTICATION_ERROR",
		"status":  401,
	},
}

var MockErrorServer = map[string]interface{}{
	"error": map[string]interface{}{
		"message": "Internal server error",
		"code":    "SERVER_ERROR",
		"status":  500,
	},
}

// ToJSON converts a fixture to JSON bytes
func ToJSON(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}
