package resources

import (
	"context"
	"encoding/json"
	"fmt"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// WorkflowsResource provides access to workflow operations
type WorkflowsResource struct {
	client   *lunahttp.Client
	basePath string
}

// List retrieves all workflows
func (r *WorkflowsResource) List(ctx context.Context) (*WorkflowList, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   r.basePath,
	})
	if err != nil {
		return nil, err
	}

	var result WorkflowList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Trigger triggers a workflow execution
func (r *WorkflowsResource) Trigger(ctx context.Context, id string, params any) (*WorkflowRun, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "POST",
		Path:   fmt.Sprintf("%s/%s/trigger", r.basePath, id),
		Body:   params,
	})
	if err != nil {
		return nil, err
	}

	var result WorkflowRun
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// AutomationResource groups Automation service resources
type AutomationResource struct {
	Workflows *WorkflowsResource
}

// NewAutomationResource creates a new Automation resource
func NewAutomationResource(client *lunahttp.Client) *AutomationResource {
	return &AutomationResource{
		Workflows: &WorkflowsResource{
			client:   client,
			basePath: "/v1/automation/workflows",
		},
	}
}
