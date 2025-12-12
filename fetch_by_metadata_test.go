package pinecone

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchByMetadata(t *testing.T) {
	t.Run("successful_fetch", func(t *testing.T) {
		var receivedBody map[string]any

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/vectors/fetch_by_metadata" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.Header.Get("X-Pinecone-API-Version") != "2025-10" {
				t.Fatalf("unexpected version header: %s", r.Header.Get("X-Pinecone-API-Version"))
			}

			bodyBytes, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(bodyBytes, &receivedBody)

			response := FetchByMetadataResponse{
				Vectors: map[string]FetchedVector{
					"id-1": {
						ID:     "id-1",
						Values: []float32{0.1, 0.2},
						Metadata: map[string]any{
							"rating": 4,
						},
					},
					"id-2": {
						ID:     "id-2",
						Values: []float32{0.3, 0.4},
						Metadata: map[string]any{
							"rating": 1,
						},
					},
				},
				Namespace: "example-namespace",
				Usage: ReadUsage{
					ReadUnits: 1,
				},
				Pagination: Pagination{
					Next: "next-token",
				},
			}

			data, _ := json.Marshal(response)

			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}))
		defer server.Close()

		client := &Client{
			IndexURL:   server.URL,
			APIKey:     "test-key",
			HTTPClient: server.Client(),
		}

		req := &FetchByMetadataRequest{
			Namespace: "example-namespace",
			Filter: map[string]any{
				"rating": map[string]any{
					"$lt": 5,
				},
			},
			Limit:           2,
			PaginationToken: "token-123",
		}

		resp, err := client.FetchByMetadata(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if receivedBody["namespace"] != "example-namespace" {
			t.Fatalf("unexpected namespace in request: %v", receivedBody["namespace"])
		}

		filter, ok := receivedBody["filter"].(map[string]any)
		if !ok {
			t.Fatalf("expected filter map, got %T", receivedBody["filter"])
		}
		rating, ok := filter["rating"].(map[string]any)
		if !ok || rating["$lt"] != float64(5) {
			t.Fatalf("unexpected filter: %v", receivedBody["filter"])
		}

		if resp.Namespace != "example-namespace" {
			t.Fatalf("unexpected namespace: %s", resp.Namespace)
		}
		if len(resp.Vectors) != 2 {
			t.Fatalf("expected 2 vectors, got %d", len(resp.Vectors))
		}
		if resp.Pagination.Next != "next-token" {
			t.Fatalf("unexpected pagination token: %s", resp.Pagination.Next)
		}
	})

	t.Run("api_error_response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message":"invalid request"}`))
		}))
		defer server.Close()

		client := &Client{
			IndexURL:   server.URL,
			APIKey:     "test-key",
			HTTPClient: server.Client(),
		}

		_, err := client.FetchByMetadata(context.Background(), &FetchByMetadataRequest{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Message != "invalid request" {
			t.Fatalf("unexpected message: %s", apiErr.Message)
		}
	})
}

