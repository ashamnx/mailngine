package mailngine

import (
	"context"
	"fmt"
)

// EmailsResource provides access to the email sending and retrieval API.
type EmailsResource struct {
	client *Client
}

// Send sends an email with the given parameters.
func (r *EmailsResource) Send(ctx context.Context, params *SendEmailParams) (*Email, error) {
	var email Email
	if err := r.client.do(ctx, "POST", "/v1/emails", params, &email); err != nil {
		return nil, err
	}
	return &email, nil
}

// Get retrieves a single email by its ID.
func (r *EmailsResource) Get(ctx context.Context, id string) (*Email, error) {
	var email Email
	if err := r.client.do(ctx, "GET", "/v1/emails/"+id, nil, &email); err != nil {
		return nil, err
	}
	return &email, nil
}

// List returns a paginated list of emails.
func (r *EmailsResource) List(ctx context.Context, opts *ListOptions) (*ListResponse[Email], error) {
	path := "/v1/emails"
	if opts != nil {
		path += buildListQuery(opts)
	}

	var emails []Email
	meta, err := r.client.doWithMeta(ctx, "GET", path, nil, &emails)
	if err != nil {
		return nil, err
	}

	resp := &ListResponse[Email]{Data: emails}
	if meta != nil {
		resp.Meta = *meta
	}
	return resp, nil
}

// buildListQuery encodes ListOptions into URL query parameters.
func buildListQuery(opts *ListOptions) string {
	if opts == nil {
		return ""
	}

	q := "?"
	sep := ""
	if opts.Page > 0 {
		q += fmt.Sprintf("%spage=%d", sep, opts.Page)
		sep = "&"
	}
	if opts.PerPage > 0 {
		q += fmt.Sprintf("%sper_page=%d", sep, opts.PerPage)
	}

	if q == "?" {
		return ""
	}
	return q
}
