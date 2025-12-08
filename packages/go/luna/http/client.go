// Package http provides HTTP client functionality for the Luna SDK.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-softworks/luna-sdk-go/luna/auth"
	"github.com/eclipse-softworks/luna-sdk-go/luna/errors"
	"github.com/eclipse-softworks/luna-sdk-go/luna/telemetry"
)

// ClientConfig holds HTTP client configuration
type ClientConfig struct {
	BaseURL      string
	Timeout      int
	MaxRetries   int
	AuthProvider auth.Provider
	Logger       telemetry.Logger
}

// Client is the HTTP client for the Luna SDK
type Client struct {
	config     ClientConfig
	httpClient *http.Client
}

// NewClient creates a new HTTP client
func NewClient(config ClientConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Millisecond,
		},
	}
}

// RequestConfig holds request configuration
type RequestConfig struct {
	Method  string
	Path    string
	Query   url.Values
	Body    interface{}
	Timeout time.Duration
}

// Response holds response data
type Response struct {
	Data      json.RawMessage
	Status    int
	Headers   http.Header
	RequestID string
}

// Request makes an HTTP request with retry logic
func (c *Client) Request(ctx context.Context, config RequestConfig) (*Response, error) {
	reqURL := c.buildURL(config.Path, config.Query)
	requestID := c.generateRequestID()

	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err := c.executeRequest(ctx, reqURL, config, requestID)
		if err == nil {
			c.config.Logger.Info("HTTP request completed", map[string]interface{}{
				"request_id": requestID,
				"method":     config.Method,
				"path":       config.Path,
				"status":     resp.Status,
			})
			return resp, nil
		}

		lastErr = err

		// Check if retryable
		lunaErr, ok := err.(*errors.Error)
		if !ok {
			break
		}

		if !lunaErr.Retryable() || attempt >= c.config.MaxRetries {
			c.config.Logger.Error("HTTP request failed", map[string]interface{}{
				"request_id": requestID,
				"method":     config.Method,
				"path":       config.Path,
				"error":      lunaErr.Code,
				"attempt":    attempt,
			})
			break
		}

		c.config.Logger.Warn("HTTP request failed, retrying", map[string]interface{}{
			"request_id": requestID,
			"method":     config.Method,
			"path":       config.Path,
			"status":     lunaErr.Status,
			"attempt":    attempt,
		})

		// Get retry delay
		var retryAfter int
		if rateLimitErr, ok := err.(*errors.RateLimitError); ok {
			retryAfter = rateLimitErr.RetryAfter
		}

		c.waitForRetry(ctx, attempt, retryAfter)
	}

	return nil, lastErr
}

func (c *Client) executeRequest(ctx context.Context, reqURL string, config RequestConfig, requestID string) (*Response, error) {
	// Get auth headers
	authHeaders, err := c.config.AuthProvider.GetHeaders()
	if err != nil {
		return nil, &errors.NetworkError{Error: &errors.Error{
			Code:      errors.CodeNetworkConnection,
			Message:   "Failed to get auth headers",
			RequestID: requestID,
		}}
	}

	// Create request body
	var bodyReader io.Reader
	if config.Body != nil {
		bodyBytes, err := json.Marshal(config.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, config.Method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Request-Id", requestID)
	req.Header.Set("User-Agent", "luna-sdk-go/1.0.0")

	for key, value := range authHeaders {
		req.Header.Set(key, value)
	}

	c.config.Logger.Debug("Sending HTTP request", map[string]interface{}{
		"request_id": requestID,
		"method":     config.Method,
		"url":        reqURL,
	})

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &errors.NetworkError{Error: &errors.Error{
				Code:      errors.CodeNetworkTimeout,
				Message:   "Request timeout",
				RequestID: requestID,
			}}
		}
		return nil, &errors.NetworkError{Error: &errors.Error{
			Code:      errors.CodeNetworkConnection,
			Message:   "Connection error",
			RequestID: requestID,
		}}
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	serverRequestID := resp.Header.Get("X-Request-Id")
	if serverRequestID == "" {
		serverRequestID = requestID
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errBody struct {
			Code    string                 `json:"code"`
			Message string                 `json:"message"`
			Details map[string]interface{} `json:"details"`
		}
		json.Unmarshal(body, &errBody)

		retryAfter := 0
		if raHeader := resp.Header.Get("Retry-After"); raHeader != "" {
			retryAfter, _ = strconv.Atoi(raHeader)
		}

		return nil, errors.FromResponse(
			resp.StatusCode,
			errBody.Code,
			errBody.Message,
			serverRequestID,
			errBody.Details,
			retryAfter,
		)
	}

	return &Response{
		Data:      body,
		Status:    resp.StatusCode,
		Headers:   resp.Header,
		RequestID: serverRequestID,
	}, nil
}

func (c *Client) buildURL(path string, query url.Values) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u := c.config.BaseURL + path

	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	return u
}

func (c *Client) generateRequestID() string {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 36)
	random := strconv.FormatInt(rand.Int63(), 36)[:8]
	return fmt.Sprintf("req_%s%s", timestamp, random)
}

func (c *Client) waitForRetry(ctx context.Context, attempt int, retryAfter int) {
	var delay time.Duration

	if retryAfter > 0 {
		delay = time.Duration(retryAfter) * time.Second
	} else {
		baseDelay := 500 * time.Millisecond
		maxDelay := 30 * time.Second
		delay = time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))
		if delay > maxDelay {
			delay = maxDelay
		}
		// Add jitter
		jitter := time.Duration(float64(delay) * 0.1 * (rand.Float64()*2 - 1))
		delay += jitter
	}

	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}
