import type { Mailngine } from '../client.js';
import type {
  CreateDomainParams,
  CreateDomainResponse,
  DNSRecord,
  Domain,
  UpdateDomainParams,
} from '../types.js';

/**
 * Resource for managing sending domains.
 */
export class DomainsResource {
  constructor(private readonly client: Mailngine) {}

  /**
   * Register a new sending domain.
   *
   * POST /v1/domains
   */
  async create(params: CreateDomainParams): Promise<CreateDomainResponse> {
    return this.client.request<CreateDomainResponse>('/v1/domains', {
      method: 'POST',
      body: params,
    });
  }

  /**
   * List all domains for the current organization.
   *
   * GET /v1/domains
   */
  async list(): Promise<Domain[]> {
    return this.client.request<Domain[]>('/v1/domains');
  }

  /**
   * Retrieve a single domain by ID.
   *
   * GET /v1/domains/:id
   */
  async get(id: string): Promise<Domain> {
    return this.client.request<Domain>(`/v1/domains/${encodeURIComponent(id)}`);
  }

  /**
   * Update domain settings (tracking options).
   *
   * PATCH /v1/domains/:id
   */
  async update(id: string, params: UpdateDomainParams): Promise<Domain> {
    return this.client.request<Domain>(`/v1/domains/${encodeURIComponent(id)}`, {
      method: 'PATCH',
      body: params,
    });
  }

  /**
   * Delete a domain and its associated DNS records.
   *
   * DELETE /v1/domains/:id
   */
  async delete(id: string): Promise<void> {
    await this.client.request<unknown>(`/v1/domains/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    });
  }

  /**
   * Trigger DNS verification for a domain.
   *
   * POST /v1/domains/:id/verify
   */
  async verify(id: string): Promise<DNSRecord[]> {
    return this.client.request<DNSRecord[]>(
      `/v1/domains/${encodeURIComponent(id)}/verify`,
      { method: 'POST' },
    );
  }
}
