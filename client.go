package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Client is a minimal REST client for the Pinecone vector database.
type Client struct {
	// IndexURL is the full index-specific Pinecone endpoint (e.g., https://example.svc.us-east1-gcp.pinecone.io)
	IndexURL string

	// APIKey is the Pinecone API key used for authentication.
	APIKey string

	// HTTPClient is the underlying HTTP client used for requests. Defaults to http.DefaultClient.
	HTTPClient *http.Client
}

// NewClient creates and returns a new Pinecone REST client.
func NewClient(indexURL, apiKey string) *Client {
	return &Client{
		IndexURL:   strings.TrimRight(indexURL, "/"),
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
}

// do sends an HTTP request to the Pinecone API with proper headers and optional JSON body.
// It returns the raw HTTP response or an error.
func (c *Client) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.IndexURL+path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.APIKey)
	req.Header.Set("X-Pinecone-API-Version", "2025-04")

	return c.HTTPClient.Do(req)
}
