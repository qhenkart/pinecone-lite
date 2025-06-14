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

// DeleteVectorsByMetadata deletes all vectors in the specified namespace that match the provided metadata filter.
//
// The filter parameter must use Pinecone's filter expression syntax. If no comparison operator is specified for a field, $eq is assumed.
// See: https://docs.pinecone.io/guides/index-data/indexing-overview#metadata-filter-expressions
//
// Example:
//
//	filter := map[string]any{"genre": "documentary"}  // defaults to $eq
//	err := client.DeleteVectorsByMetadata(ctx, "example-namespace", filter)
//	if err != nil {
//	    // handle error
//	}
func (c *Client) DeleteVectorsByMetadata(ctx context.Context, namespace string, filter map[string]any) error {
	body := map[string]any{
		"namespace": namespace,
		"filter":    filter,
	}

	resp, err := c.do(ctx, http.MethodPost, "/vectors/delete", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return parseAPIError(resp)
	}
	return nil
}
