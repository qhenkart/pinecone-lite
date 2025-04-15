package pinecone

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("trailing_slash_trimmed", func(t *testing.T) {
		c := NewClient("https://example.com/", "key")
		if c.IndexURL != "https://example.com" {
			t.Errorf("expected trimmed URL, got %s", c.IndexURL)
		}
	})

	t.Run("sets_fields_correctly", func(t *testing.T) {
		c := NewClient("https://host", "my-key")
		if c.APIKey != "my-key" {
			t.Errorf("expected api key to be set")
		}
		if c.HTTPClient != http.DefaultClient {
			t.Errorf("expected default client")
		}
	})
}

func TestClientDo(t *testing.T) {
	t.Run("sends_correct_request", func(t *testing.T) {
		var gotReq *http.Request

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotReq = r
			w.WriteHeader(200)
		}))
		defer srv.Close()

		c := &Client{
			IndexURL:   srv.URL,
			APIKey:     "abc123",
			HTTPClient: srv.Client(),
		}

		ctx := context.Background()
		_, err := c.do(ctx, http.MethodPost, "/vectors/upsert", map[string]any{"key": "val"})
		if err != nil {
			t.Fatalf("do failed: %v", err)
		}

		if gotReq.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", gotReq.Method)
		}
		if gotReq.URL.Path != "/vectors/upsert" {
			t.Errorf("unexpected path: %s", gotReq.URL.Path)
		}
		if gotReq.Header.Get("Api-Key") != "abc123" {
			t.Errorf("missing Api-Key header")
		}
		if gotReq.Header.Get("X-Pinecone-API-Version") != "2025-04" {
			t.Errorf("missing version header")
		}
	})

	t.Run("handles_marshal_error", func(t *testing.T) {
		c := NewClient("http://localhost", "k")
		ctx := context.Background()
		_, err := c.do(ctx, http.MethodPost, "/foo", func() {})
		if err == nil {
			t.Fatal("expected marshal error")
		}
	})

	t.Run("handles_request_build_error", func(t *testing.T) {
		c := NewClient("%%%", "key")
		ctx := context.Background()
		_, err := c.do(ctx, http.MethodGet, "/bad", nil)
		if err == nil {
			t.Fatal("expected request build error")
		}
	})
}
