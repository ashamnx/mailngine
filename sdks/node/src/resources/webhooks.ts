import type { Mailngine } from '../client.js';
import type {
  CreateWebhookParams,
  ListOptions,
  UpdateWebhookParams,
  Webhook,
  WebhookDelivery,
} from '../types.js';

/**
 * Resource for managing webhook endpoints.
 */
export class WebhooksResource {
  constructor(private readonly client: Mailngine) {}

  /**
   * Create a new webhook endpoint.
   *
   * POST /v1/webhooks
   */
  async create(params: CreateWebhookParams): Promise<Webhook> {
    return this.client.request<Webhook>('/v1/webhooks', {
      method: 'POST',
      body: params,
    });
  }

  /**
   * List all webhooks for the current organization.
   *
   * GET /v1/webhooks
   */
  async list(): Promise<Webhook[]> {
    return this.client.request<Webhook[]>('/v1/webhooks');
  }

  /**
   * Retrieve a single webhook by ID.
   *
   * GET /v1/webhooks/:id
   */
  async get(id: string): Promise<Webhook> {
    return this.client.request<Webhook>(`/v1/webhooks/${encodeURIComponent(id)}`);
  }

  /**
   * Update a webhook endpoint.
   *
   * PATCH /v1/webhooks/:id
   */
  async update(id: string, params: UpdateWebhookParams): Promise<Webhook> {
    return this.client.request<Webhook>(`/v1/webhooks/${encodeURIComponent(id)}`, {
      method: 'PATCH',
      body: params,
    });
  }

  /**
   * Delete a webhook endpoint.
   *
   * DELETE /v1/webhooks/:id
   */
  async delete(id: string): Promise<void> {
    await this.client.request<unknown>(`/v1/webhooks/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    });
  }

  /**
   * List delivery attempts for a webhook.
   *
   * GET /v1/webhooks/:id/deliveries
   */
  async listDeliveries(id: string, options?: ListOptions): Promise<WebhookDelivery[]> {
    const params = new URLSearchParams();
    if (options?.page != null) params.set('page', String(options.page));
    if (options?.perPage != null) params.set('per_page', String(options.perPage));

    const qs = params.toString();
    const path = qs
      ? `/v1/webhooks/${encodeURIComponent(id)}/deliveries?${qs}`
      : `/v1/webhooks/${encodeURIComponent(id)}/deliveries`;

    return this.client.request<WebhookDelivery[]>(path);
  }
}
