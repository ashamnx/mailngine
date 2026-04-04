package hellomail

import "context"

// WebhooksResource provides access to the webhook management API.
type WebhooksResource struct {
	client *Client
}

// Create registers a new webhook endpoint.
func (r *WebhooksResource) Create(ctx context.Context, params *CreateWebhookParams) (*Webhook, error) {
	var webhook Webhook
	if err := r.client.do(ctx, "POST", "/v1/webhooks", params, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// Get retrieves a single webhook by its ID.
func (r *WebhooksResource) Get(ctx context.Context, id string) (*Webhook, error) {
	var webhook Webhook
	if err := r.client.do(ctx, "GET", "/v1/webhooks/"+id, nil, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// List returns all webhooks for the authenticated organization.
func (r *WebhooksResource) List(ctx context.Context) ([]Webhook, error) {
	var webhooks []Webhook
	if err := r.client.do(ctx, "GET", "/v1/webhooks", nil, &webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

// Update modifies an existing webhook.
func (r *WebhooksResource) Update(ctx context.Context, id string, params *UpdateWebhookParams) (*Webhook, error) {
	var webhook Webhook
	if err := r.client.do(ctx, "PATCH", "/v1/webhooks/"+id, params, &webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}

// Delete removes a webhook.
func (r *WebhooksResource) Delete(ctx context.Context, id string) error {
	return r.client.do(ctx, "DELETE", "/v1/webhooks/"+id, nil, nil)
}
