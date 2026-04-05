package mailngine

import "context"

// DomainsResource provides access to the domain management API.
type DomainsResource struct {
	client *Client
}

// Create registers a new sending domain and returns it along with the
// DNS records that must be configured for verification.
func (r *DomainsResource) Create(ctx context.Context, name string) (*CreateDomainResponse, error) {
	body := map[string]string{"name": name}
	var resp CreateDomainResponse
	if err := r.client.do(ctx, "POST", "/v1/domains", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single domain by its ID.
func (r *DomainsResource) Get(ctx context.Context, id string) (*Domain, error) {
	var domain Domain
	if err := r.client.do(ctx, "GET", "/v1/domains/"+id, nil, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

// List returns all domains for the authenticated organization.
func (r *DomainsResource) List(ctx context.Context) ([]Domain, error) {
	var domains []Domain
	if err := r.client.do(ctx, "GET", "/v1/domains", nil, &domains); err != nil {
		return nil, err
	}
	return domains, nil
}

// Verify triggers DNS verification for a domain and returns the updated
// DNS records with their verification status.
func (r *DomainsResource) Verify(ctx context.Context, id string) ([]DNSRecord, error) {
	var records []DNSRecord
	if err := r.client.do(ctx, "POST", "/v1/domains/"+id+"/verify", nil, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Update modifies domain settings such as open and click tracking.
func (r *DomainsResource) Update(ctx context.Context, id string, params *UpdateDomainParams) (*Domain, error) {
	var domain Domain
	if err := r.client.do(ctx, "PATCH", "/v1/domains/"+id, params, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

// Delete removes a domain and its associated DNS records.
func (r *DomainsResource) Delete(ctx context.Context, id string) error {
	return r.client.do(ctx, "DELETE", "/v1/domains/"+id, nil, nil)
}
