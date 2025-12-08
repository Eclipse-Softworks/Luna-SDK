package resources

import (
	"context"
	"encoding/json"
	"fmt"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// AiResource provides access to AI operations
type AiResource struct {
	client   *lunahttp.Client
	basePath string
}

// ChatCompletions generates chat completions
func (r *AiResource) ChatCompletions(ctx context.Context, params *CompletionRequest) (*CompletionResponse, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "POST",
		Path:   fmt.Sprintf("%s/chat/completions", r.basePath),
		Body:   params,
	})
	if err != nil {
		return nil, err
	}

	var result CompletionResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// NewAiResource creates a new AI resource
func NewAiResource(client *lunahttp.Client) *AiResource {
	return &AiResource{
		client:   client,
		basePath: "/v1/ai",
	}
}
