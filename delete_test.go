package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteVectorsByID(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != "/vectors/delete" {
				t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := &Client{
			IndexURL:   server.URL,
			APIKey:     "test-key",
			HTTPClient: server.Client(),
		}

		err := client.DeleteVectorsByID(context.Background(), []string{"vec-1", "vec-2"}, "ns")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
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

		err := client.DeleteVectorsByID(context.Background(), []string{"vec-1"}, "ns")
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

// DeleteAllRecordsInNamespace deletes all records in the specified namespace
// by removing the entire namespace from the index.
//
// This is a destructive operation: the namespace and all its associated data
// will be permanently deleted from the Pinecone index.

func TestDeleteAllRecordsInNamespace(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		var called bool
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			if r.Method != http.MethodDelete {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/namespaces/test-namespace" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if r.Header.Get("Api-Key") != "test-key" {
				t.Fatalf("missing or incorrect Api-Key header")
			}
			if r.Header.Get("X-Pinecone-API-Version") != "2025-04" {
				t.Fatalf("missing or incorrect API version header")
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer s.Close()

		client := &Client{
			IndexURL:   s.URL,
			APIKey:     "test-key",
			HTTPClient: s.Client(),
		}

		err := client.DeleteAllRecordsInNamespace(context.Background(), "test-namespace")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Fatal("handler was not called")
		}
	})

	t.Run("api_error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message":"bad request"}`))
		}))
		defer s.Close()

		client := &Client{
			IndexURL:   s.URL,
			APIKey:     "test-key",
			HTTPClient: s.Client(),
		}

		err := client.DeleteAllRecordsInNamespace(context.Background(), "test-namespace")
		if err == nil || !strings.Contains(err.Error(), "bad request") {
			t.Fatalf("expected API error, got: %v", err)
		}
	})
}

func TestDeleteVectorsByMetadata(t *testing.T) {
	t.Run("valid_delete_by_metadata", func(t *testing.T) {
		// Expectation: The endpoint returns 200 OK for a successful delete request.
		client := &Client{
			IndexURL: "https://example-index.svc.us-east1-gcp.io",
			APIKey:   "test-key",
			HTTPClient: &http.Client{
				Transport: roundTripFunc(func(req *http.Request) *http.Response {
					// Validate method and endpoint
					if req.Method != http.MethodPost {
						t.Errorf("expected POST, got %s", req.Method)
					}
					if req.URL.Path != "/vectors/delete" {
						t.Errorf("expected /vectors/delete, got %s", req.URL.Path)
					}
					// Validate body
					body, _ := io.ReadAll(req.Body)
					var reqData map[string]any
					json.Unmarshal(body, &reqData)
					if reqData["namespace"] != "example-namespace" {
						t.Errorf("expected namespace 'example-namespace', got %v", reqData["namespace"])
					}
					filter, ok := reqData["filter"].(map[string]any)
					if !ok || filter["genre"] != "documentary" {
						t.Errorf("expected filter genre=documentary, got %v", reqData["filter"])
					}

					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
						Header:     make(http.Header),
					}
				}),
			},
		}

		filter := map[string]any{"genre": "documentary"}
		err := client.DeleteVectorsByMetadata(context.Background(), "example-namespace", filter)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}
