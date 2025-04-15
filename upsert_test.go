package pinecone

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpsertVectors(t *testing.T) {
	t.Run("empty_input", func(t *testing.T) {
		client := &Client{
			IndexURL:   "https://example.com",
			APIKey:     "test-key",
			HTTPClient: http.DefaultClient,
		}

		_, err := client.UpsertVectors(context.Background(), []*Vector{}, "test-ns")
		if err == nil {
			t.Fatal("expected error on empty input, got nil")
		}
	})

	t.Run("http_failure", func(t *testing.T) {
		client := &Client{
			IndexURL:   "http://invalid host",
			APIKey:     "test-key",
			HTTPClient: http.DefaultClient,
		}

		_, err := client.UpsertVectors(context.Background(), []*Vector{
			{ID: "vec1", Values: []float64{0.1}},
		}, "ns")
		if err == nil {
			t.Fatal("expected error on HTTP request failure")
		}
	})

	t.Run("invalid_json_response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`invalid`))
		}))
		defer ts.Close()

		client := &Client{
			IndexURL:   ts.URL,
			APIKey:     "key",
			HTTPClient: ts.Client(),
		}

		_, err := client.UpsertVectors(context.Background(), []*Vector{
			{ID: "v1", Values: []float64{0.1}},
		}, "ns")

		if err == nil || !strings.Contains(err.Error(), "invalid character") {
			t.Fatalf("expected JSON parse error, got: %v", err)
		}
	})
}
