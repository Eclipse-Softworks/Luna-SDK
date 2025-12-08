package utils_test

import (
	"strings"
	"testing"

	"github.com/eclipse-softworks/luna-sdk-go/luna/utils"
)

func TestGenerateRequestID(t *testing.T) {
	t.Run("generates unique IDs", func(t *testing.T) {
		id1 := utils.GenerateRequestID()
		id2 := utils.GenerateRequestID()

		if !strings.HasPrefix(id1, "req_") {
			t.Errorf("expected id to start with req_, got %s", id1)
		}
		if !strings.HasPrefix(id2, "req_") {
			t.Errorf("expected id to start with req_, got %s", id2)
		}
		if id1 == id2 {
			t.Error("expected unique IDs")
		}
	})
}

func TestMaskSensitive(t *testing.T) {
	t.Run("masks long strings", func(t *testing.T) {
		result := utils.MaskSensitive("lk_test_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 7, 4)

		if !strings.HasPrefix(result, "lk_test") {
			t.Errorf("expected to start with lk_test, got %s", result)
		}
		if !strings.HasSuffix(result, "aaaa") {
			t.Errorf("expected to end with aaaa, got %s", result)
		}
		if !strings.Contains(result, "****") {
			t.Errorf("expected to contain ****, got %s", result)
		}
	})

	t.Run("masks short strings completely", func(t *testing.T) {
		result := utils.MaskSensitive("short", 7, 4)
		if result != "*****" {
			t.Errorf("expected *****, got %s", result)
		}
	})
}

func TestValidateID(t *testing.T) {
	t.Run("accepts valid user ID", func(t *testing.T) {
		err := utils.ValidateID("usr_abc123", "usr", "user ID")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("accepts valid project ID", func(t *testing.T) {
		err := utils.ValidateID("prj_abc123", "prj", "project ID")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("rejects empty ID", func(t *testing.T) {
		err := utils.ValidateID("", "usr", "user ID")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("rejects invalid format", func(t *testing.T) {
		err := utils.ValidateID("invalid", "usr", "user ID")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestIsRetryableStatus(t *testing.T) {
	t.Run("identifies retryable statuses", func(t *testing.T) {
		retryable := []int{408, 429, 500, 502, 503, 504}
		for _, status := range retryable {
			if !utils.IsRetryableStatus(status) {
				t.Errorf("expected %d to be retryable", status)
			}
		}
	})

	t.Run("identifies non-retryable statuses", func(t *testing.T) {
		nonRetryable := []int{400, 401, 403, 404}
		for _, status := range nonRetryable {
			if utils.IsRetryableStatus(status) {
				t.Errorf("expected %d to not be retryable", status)
			}
		}
	})
}

func TestDeepMerge(t *testing.T) {
	t.Run("merges flat maps", func(t *testing.T) {
		base := map[string]interface{}{"a": 1, "b": 2}
		override := map[string]interface{}{"b": 3, "c": 4}

		result := utils.DeepMerge(base, override)

		if result["a"] != 1 {
			t.Errorf("expected a=1, got %v", result["a"])
		}
		if result["b"] != 3 {
			t.Errorf("expected b=3, got %v", result["b"])
		}
		if result["c"] != 4 {
			t.Errorf("expected c=4, got %v", result["c"])
		}
	})

	t.Run("merges nested maps", func(t *testing.T) {
		base := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{"c": 2, "d": 3},
		}
		override := map[string]interface{}{
			"b": map[string]interface{}{"d": 4, "e": 5},
		}

		result := utils.DeepMerge(base, override)

		nested := result["b"].(map[string]interface{})
		if nested["c"] != 2 {
			t.Errorf("expected c=2, got %v", nested["c"])
		}
		if nested["d"] != 4 {
			t.Errorf("expected d=4, got %v", nested["d"])
		}
		if nested["e"] != 5 {
			t.Errorf("expected e=5, got %v", nested["e"])
		}
	})
}
