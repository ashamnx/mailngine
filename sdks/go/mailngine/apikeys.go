package mailngine

import "context"

// APIKeysResource provides access to the API key management API.
type APIKeysResource struct {
	client *Client
}

// Create generates a new API key. The full key value is only available in the
// returned APIKey; it cannot be retrieved again after creation.
func (r *APIKeysResource) Create(ctx context.Context, params *CreateAPIKeyParams) (*APIKey, error) {
	var key APIKey
	if err := r.client.do(ctx, "POST", "/v1/api-keys", params, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

// List returns all active API keys for the authenticated organization.
// The full key value is not included in list responses.
func (r *APIKeysResource) List(ctx context.Context) ([]APIKey, error) {
	var keys []APIKey
	if err := r.client.do(ctx, "GET", "/v1/api-keys", nil, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// Revoke deactivates an API key, preventing it from being used for authentication.
func (r *APIKeysResource) Revoke(ctx context.Context, id string) error {
	return r.client.do(ctx, "DELETE", "/v1/api-keys/"+id, nil, nil)
}
