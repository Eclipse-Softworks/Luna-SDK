package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	lunahttp "github.com/eclipse-softworks/luna-sdk-go/luna/http"
)

// ResidencesResource provides access to residence operations
type ResidencesResource struct {
	client   *lunahttp.Client
	basePath string
}

// List searches for residences
func (r *ResidencesResource) List(ctx context.Context, params *ResidenceSearch) (*ResidenceList, error) {
	query := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			query.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Cursor != "" {
			query.Set("cursor", params.Cursor)
		}
		if params.Query != "" {
			query.Set("query", params.Query)
		}
		if params.NSFAS != nil {
			query.Set("nsfas", strconv.FormatBool(*params.NSFAS))
		}
		if params.MinPrice > 0 {
			query.Set("min_price", fmt.Sprintf("%f", params.MinPrice))
		}
		if params.MaxPrice > 0 {
			query.Set("max_price", fmt.Sprintf("%f", params.MaxPrice))
		}
		if params.Gender != "" {
			query.Set("gender", params.Gender)
		}
		if params.CampusID != "" {
			query.Set("campus_id", params.CampusID)
		}
		if params.Radius > 0 {
			query.Set("radius", fmt.Sprintf("%f", params.Radius))
		}
		if params.MinRating > 0 {
			query.Set("min_rating", fmt.Sprintf("%f", params.MinRating))
		}
	}

	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   r.basePath,
		Query:  query,
	})
	if err != nil {
		return nil, err
	}

	var result ResidenceList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Get retrieves a residence by ID
func (r *ResidencesResource) Get(ctx context.Context, id string) (*Residence, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   fmt.Sprintf("%s/%s", r.basePath, id),
	})
	if err != nil {
		return nil, err
	}

	var result Residence
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// Iterate returns a paginator for iterating over residences
func (r *ResidencesResource) Iterate(ctx context.Context, params *ResidenceSearch) *Paginator[Residence] {
	return NewPaginator(ctx, func(ctx context.Context, cursor string) (*ListResponse[Residence], error) {
		p := params
		if p == nil {
			p = &ResidenceSearch{}
		}
		// Create a shallow copy to modify cursor without affecting original params if re-used (though pointer is passed)
		// For robustness, we should copy.
		newParams := *p
		newParams.Cursor = cursor
		return r.List(ctx, &newParams)
	})
}

// CampusesResource provides access to campus operations
type CampusesResource struct {
	client   *lunahttp.Client
	basePath string
}

// List retrieves all campuses
func (r *CampusesResource) List(ctx context.Context) (*CampusList, error) {
	resp, err := r.client.Request(ctx, lunahttp.RequestConfig{
		Method: "GET",
		Path:   r.basePath,
	})
	if err != nil {
		return nil, err
	}

	var result CampusList
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ResMateResource groups ResMate service resources
type ResMateResource struct {
	Residences *ResidencesResource
	Campuses   *CampusesResource
}

// NewResMateResource creates a new ResMate resource
func NewResMateResource(client *lunahttp.Client) *ResMateResource {
	return &ResMateResource{
		Residences: &ResidencesResource{
			client:   client,
			basePath: "/v1/resmate/residences",
		},
		Campuses: &CampusesResource{
			client:   client,
			basePath: "/v1/resmate/campuses",
		},
	}
}
