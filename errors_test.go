package pinecone

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("structured_message", func(t *testing.T) {
		body := `{"message": "something went wrong"}`
		resp := &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(strings.NewReader(body)),
		}

		err := parseAPIError(resp)
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Message != "something went wrong" {
			t.Errorf("unexpected message: %s", apiErr.Message)
		}
		if apiErr.StatusCode != 400 {
			t.Errorf("unexpected status code: %d", apiErr.StatusCode)
		}
	})

	t.Run("unstructured_message", func(t *testing.T) {
		body := `not json`
		resp := &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader(body)),
		}

		err := parseAPIError(resp)
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.Message != "not json" {
			t.Errorf("expected fallback message, got: %s", apiErr.Message)
		}
		if apiErr.StatusCode != 500 {
			t.Errorf("unexpected status code: %d", apiErr.StatusCode)
		}
	})
}
