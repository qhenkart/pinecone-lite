package pinecone

import (
	"context"
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
