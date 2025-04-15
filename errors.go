package pinecone

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents a structured error returned by Pinecone's API.
type APIError struct {
	StatusCode int
	Message    string
	Body       []byte
}

// Error returns the string representation of the API error.
func (e *APIError) Error() string {
	return fmt.Sprintf("pinecone: %s (status %d)", e.Message, e.StatusCode)
}

// parseAPIError parses a non-2xx HTTP response into an APIError.
func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var parsed struct {
		Message string `json:"message"`
	}
	_ = json.Unmarshal(body, &parsed)

	msg := parsed.Message
	if msg == "" {
		msg = string(body)
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    msg,
		Body:       body,
	}
}
