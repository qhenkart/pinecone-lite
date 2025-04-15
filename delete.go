package pinecone

import (
	"context"
	"net/http"
)

// DeleteVectorsByID deletes the specified vector IDs from the given namespace.
// It returns an error if the request fails.
func (c *Client) DeleteVectorsByID(ctx context.Context, ids []string, namespace string) error {
	payload := map[string]any{
		"ids":       ids,
		"namespace": namespace,
	}

	resp, err := c.do(ctx, http.MethodPost, "/vectors/delete", payload)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return parseAPIError(resp)
	}
	return nil
}

// DeleteAllRecordsInNamespace deletes all records in the specified namespace
// by removing the entire namespace from the index.
//
// This is a destructive operation: the namespace and all its associated data
// will be permanently deleted from the Pinecone index.
func (c *Client) DeleteAllRecordsInNamespace(ctx context.Context, namespace string) error {
	resp, err := c.do(ctx, http.MethodDelete, "/namespaces/"+namespace, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return parseAPIError(resp)
	}
	return nil
}
