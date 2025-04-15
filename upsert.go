package pinecone

import (
	"context"
	"encoding/json"
	"net/http"
)

// Vector represents a single dense vector with optional metadata.
type Vector struct {
	ID       string    `json:"id"`
	Values   []float64 `json:"values"`
	Values32 []float32
	Metadata map[string]any `json:"metadata,omitempty"`
}

// UpsertRequest is the payload structure for upserting vectors.
type UpsertRequest struct {
	Vectors   []*Vector `json:"vectors"`
	Namespace string    `json:"namespace,omitempty"`
}

// UpsertResponse is the response from the Pinecone /vectors/upsert endpoint.
type UpsertResponse struct {
	UpsertedCount uint32 `json:"upsertedCount"`
}

// UpsertVectors inserts or updates one or more vectors into the specified namespace.
// Returns the number of vectors upserted or an error.
func (c *Client) UpsertVectors(ctx context.Context, vectors []*Vector, namespace string) (uint32, error) {
	payload := UpsertRequest{
		Vectors:   vectors,
		Namespace: namespace,
	}

	resp, err := c.do(ctx, http.MethodPost, "/vectors/upsert", payload)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return 0, parseAPIError(resp)
	}

	var parsed UpsertResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return 0, err
	}

	return parsed.UpsertedCount, nil
}
