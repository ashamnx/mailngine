import type { HelloMail } from '../client.js';
import type { ApiKey, ApiKeyCreateResponse, CreateApiKeyParams } from '../types.js';

/**
 * Resource for managing API keys.
 */
export class ApiKeysResource {
  constructor(private readonly client: HelloMail) {}

  /**
   * Create a new API key.
   *
   * The full key value is only returned once in the response.
   *
   * POST /v1/api-keys
   */
  async create(params: CreateApiKeyParams): Promise<ApiKeyCreateResponse> {
    return this.client.request<ApiKeyCreateResponse>('/v1/api-keys', {
      method: 'POST',
      body: params,
    });
  }

  /**
   * List all active API keys for the current organization.
   *
   * GET /v1/api-keys
   */
  async list(): Promise<ApiKey[]> {
    return this.client.request<ApiKey[]>('/v1/api-keys');
  }

  /**
   * Revoke (delete) an API key.
   *
   * DELETE /v1/api-keys/:id
   */
  async revoke(id: string): Promise<void> {
    await this.client.request<unknown>(`/v1/api-keys/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    });
  }
}
