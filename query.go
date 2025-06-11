package pinecone

import (
	"context"
	"encoding/json"
	"net/http"
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
