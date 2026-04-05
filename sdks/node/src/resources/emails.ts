import type { Mailngine } from '../client.js';
import type { Email, ListOptions, ListResponse, SendEmailParams } from '../types.js';

/**
 * Resource for sending and retrieving emails.
 *
 * @example
 * ```ts
 * const email = await client.emails.send({
 *   from: 'hello@example.com',
 *   to: ['user@example.com'],
 *   subject: 'Welcome!',
 *   html: '<h1>Hello</h1>',
 * });
 * ```
 */
export class EmailsResource {
  constructor(private readonly client: Mailngine) {}

  /**
   * Send an email.
   *
   * POST /v1/emails
   */
  async send(params: SendEmailParams): Promise<Email> {
    return this.client.request<Email>('/v1/emails', {
      method: 'POST',
      body: params,
    });
  }

  /**
   * Retrieve a single email by ID.
   *
   * GET /v1/emails/:id
   */
  async get(id: string): Promise<Email> {
    return this.client.request<Email>(`/v1/emails/${encodeURIComponent(id)}`);
  }

  /**
   * List emails with pagination.
   *
   * GET /v1/emails
   */
  async list(options?: ListOptions): Promise<ListResponse<Email>> {
    const params = new URLSearchParams();
    if (options?.page != null) params.set('page', String(options.page));
    if (options?.perPage != null) params.set('per_page', String(options.perPage));

    const qs = params.toString();
    const path = qs ? `/v1/emails?${qs}` : '/v1/emails';

    return this.client.request<ListResponse<Email>>(path, { envelope: 'list' });
  }
}
