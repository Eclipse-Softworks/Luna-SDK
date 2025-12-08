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
	Method      string
	Path        string
	Query       url.Values
	Body        interface{}
	BodyReader  io.Reader
	ContentType string
	Timeout     time.Duration
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
				"status":     resp.Status,
			})
			return resp, nil
		}

		lastErr = err

		// Check if retryable
		lunaErr, ok := err.(*errors.Error)
		// Note: The implementation of Error() method on the struct pointer vs embedding
		// might require type assertion check on specific error types or the base interface

		if !ok {
			// Try to unwrap or check if it implements the interface via other means,
			// or simply treat unknown errors as non-retryable for now unless they are network/timeout
			break
		}

		if !lunaErr.Retryable() || attempt >= c.config.MaxRetries {
			c.config.Logger.Error("HTTP request failed", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}

		c.config.Logger.Warn("Retrying request", map[string]interface{}{
			"attempt": attempt,
		})

		// Backoff Strategy
		var retryAfter int
		if rateLimitErr, ok := err.(*errors.RateLimitError); ok {
			retryAfter = rateLimitErr.RetryAfter
		}
		c.waitForRetry(ctx, attempt, retryAfter)
	}

	return nil, lastErr
}

func (c *Client) executeRequest(ctx context.Context, reqURL string, config RequestConfig, requestID string) (*Response, error) {
	// 1. Request Signing
	authHeaders, err := c.config.AuthProvider.GetHeaders()
	if err != nil {
		return nil, &errors.NetworkError{BaseError: &errors.Error{Code: errors.CodeNetworkConnection}}
	}

	var bodyReader io.Reader
	contentType := "application/json"

	if config.ContentType != "" {
		contentType = config.ContentType
	}

	if config.BodyReader != nil {
		bodyReader = config.BodyReader
	} else if config.Body != nil {
		bodyBytes, _ := json.Marshal(config.Body)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// 2. Timeout (via Context)
	req, err := http.NewRequestWithContext(ctx, config.Method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Request-Id", requestID)
	for k, v := range authHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &errors.NetworkError{BaseError: &errors.Error{Code: errors.CodeNetworkTimeout}}
		}
		return nil, &errors.NetworkError{BaseError: &errors.Error{Code: errors.CodeNetworkConnection}}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// 4. Error Normalization
	if serverRequestID := resp.Header.Get("X-Request-Id"); serverRequestID != "" {
		requestID = serverRequestID
	}

	if resp.StatusCode >= 400 {
		var errBody struct {
			Code    string                 `json:"code"`
			Message string                 `json:"message"`
			Details map[string]interface{} `json:"details"`
			Error   *struct {
				Code    string                 `json:"code"`
				Message string                 `json:"message"`
				Details map[string]interface{} `json:"details"`
			} `json:"error"`
		}
		json.Unmarshal(body, &errBody)

		if errBody.Error != nil {
			errBody.Code = errBody.Error.Code
			errBody.Message = errBody.Error.Message
			errBody.Details = errBody.Error.Details
		}

		retryAfter := 0
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			fmt.Sscanf(ra, "%d", &retryAfter)
		}

		return nil, errors.FromResponse(resp.StatusCode, errBody.Code, errBody.Message, requestID, errBody.Details, retryAfter)
	}

	return &Response{
		Data:      body,
		Status:    resp.StatusCode,
		Headers:   resp.Header,
		RequestID: requestID,
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
		// Exponential + Jitter
		delay = time.Duration(500*math.Pow(2, float64(attempt))) * time.Millisecond
		delay += time.Duration(rand.Int63n(int64(delay) / 10))
	}

	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}
