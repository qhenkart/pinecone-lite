package pinecone

import (
	"context"
	"encoding/json"
	"net/http"
)

// FetchByMetadataRequest describes the payload for POST /vectors/fetch_by_metadata.
// See: https://docs.pinecone.io/reference/api/2025-10/data-plane/fetch_by_metadata
type FetchByMetadataRequest struct {
	Namespace       string         `json:"namespace,omitempty"`
	Filter          map[string]any `json:"filter,omitempty"`
	Limit           int            `json:"limit,omitempty"`
	PaginationToken string         `json:"paginationToken,omitempty"`
}

// FetchedVector represents a vector returned by fetch_by_metadata.
type FetchedVector struct {
	ID       string         `json:"id"`
	Values   []float64      `json:"values,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Pagination holds pagination tokens for list-like responses.
type Pagination struct {
	Next string `json:"next"`
}

// FetchByMetadataResponse contains the response payload from fetch_by_metadata.
type FetchByMetadataResponse struct {
	Vectors    map[string]FetchedVector `json:"vectors"`
	Namespace  string                   `json:"namespace"`
	Usage      ReadUsage                `json:"usage"`
	Pagination Pagination               `json:"pagination"`
}

// FetchByMetadata retrieves vectors that match the provided metadata filter.
func (c *Client) FetchByMetadata(ctx context.Context, req *FetchByMetadataRequest) (*FetchByMetadataResponse, error) {
	resp, err := c.do(ctx, http.MethodPost, "/vectors/fetch_by_metadata", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, parseAPIError(resp)
	}

	var parsed FetchByMetadataResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}
