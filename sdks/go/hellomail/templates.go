package hellomail

import "context"

// TemplatesResource provides access to the template management API.
type TemplatesResource struct {
	client *Client
}

// Create creates a new email template.
func (r *TemplatesResource) Create(ctx context.Context, params *CreateTemplateParams) (*Template, error) {
	var tmpl Template
	if err := r.client.do(ctx, "POST", "/v1/templates", params, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// Get retrieves a single template by its ID.
func (r *TemplatesResource) Get(ctx context.Context, id string) (*Template, error) {
	var tmpl Template
	if err := r.client.do(ctx, "GET", "/v1/templates/"+id, nil, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// List returns all templates for the authenticated organization.
func (r *TemplatesResource) List(ctx context.Context) ([]Template, error) {
	var templates []Template
	if err := r.client.do(ctx, "GET", "/v1/templates", nil, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

// Update modifies an existing template.
func (r *TemplatesResource) Update(ctx context.Context, id string, params *UpdateTemplateParams) (*Template, error) {
	var tmpl Template
	if err := r.client.do(ctx, "PATCH", "/v1/templates/"+id, params, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// Delete removes a template.
func (r *TemplatesResource) Delete(ctx context.Context, id string) error {
	return r.client.do(ctx, "DELETE", "/v1/templates/"+id, nil, nil)
}

// Preview renders a template with the given data and returns the result.
func (r *TemplatesResource) Preview(ctx context.Context, id string, params *PreviewTemplateParams) (*TemplatePreview, error) {
	var preview TemplatePreview
	if err := r.client.do(ctx, "POST", "/v1/templates/"+id+"/preview", params, &preview); err != nil {
		return nil, err
	}
	return &preview, nil
}
