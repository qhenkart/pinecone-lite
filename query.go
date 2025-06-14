package pinecone

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// MatchResult represents a result match returned from a query.
type MatchResult struct {
	ID       string         `json:"id"`
	Score    float64        `json:"score"`
	Values   []float32      `json:"values,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// QueryByVectorRequest represents a request to query similar vectors.
type QueryByVectorRequest struct {
	Vector          []float64
	TopK            int
	Namespace       string
	Filter          map[string]any
	IncludeValues   bool
	IncludeMetadata bool
}

// QueryByVectorResponse represents the response from a vector query.
type QueryByVectorResponse struct {
	Matches   []MatchResult `json:"matches"`
	Namespace string        `json:"namespace"`
	Usage     ReadUsage     `json:"usage"`
}

// ReadUsage tracks query usage
type ReadUsage struct {
	ReadUnits uint32 `json:"readUnits"`
}

// QueryByVector performs a similarity search using a dense vector.
func (c *Client) QueryByVector(ctx context.Context, req *QueryByVectorRequest) (*QueryByVectorResponse, error) {
	body := map[string]any{
		"vector":          req.Vector,
		"topK":            req.TopK,
		"namespace":       req.Namespace,
		"includeValues":   req.IncludeValues,
		"includeMetadata": req.IncludeMetadata,
	}

	if req.Filter != nil {
		body["filter"] = req.Filter
	}

	resp, err := c.do(ctx, http.MethodPost, "/query", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, parseAPIError(resp)
	}

	var parsed QueryByVectorResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

// ListVectorIDs retrieves vector IDs from a namespace, with optional prefix, limit, and pagination.
//
// Parameters:
//
//	namespace - namespace string to list IDs from
//	prefix - optional string to filter IDs by prefix (pass "" for no filter)
//	limit - max IDs per page (default 100)
//	paginationToken - token for next page (pass "" for first page)
//
// Returns:
//
//	ids - slice of vector IDs
//	nextPaginationToken - token for next page; empty string if no further pages
//	error - error if request fails or API returns error
//
// Example usage:
//
//	// Retrieve up to 100 vector IDs from the "production" namespace
//	ids, nextToken, err := client.ListVectors(ctx, "production", "", 100, "")
//	if err != nil {
//	    // handle error
//	}
//
//	// Process the vector IDs in ids
//
//	// If nextToken is not empty, retrieve the next page:
//	// moreIDs, _, err := client.ListVectors(ctx, "production", "", 100, nextToken)
func (c *Client) ListVectorIDs(ctx context.Context, namespace, prefix string, limit int, paginationToken string) ([]string, string, error) {
	params := map[string]string{
		"namespace": namespace,
	}
	if prefix != "" {
		params["prefix"] = prefix
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if paginationToken != "" {
		params["paginationToken"] = paginationToken
	}

	resp, err := c.do(ctx, http.MethodGet, "/vectors/list", params)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, "", parseAPIError(resp)
	}

	var result struct {
		Vectors []struct {
			ID string `json:"id"`
		} `json:"vectors"`
		Pagination struct {
			Next string `json:"next"`
		} `json:"pagination"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", err
	}

	ids := make([]string, len(result.Vectors))
	for i, v := range result.Vectors {
		ids[i] = v.ID
	}

	return ids, result.Pagination.Next, nil
}
