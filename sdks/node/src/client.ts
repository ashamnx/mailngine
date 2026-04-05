import { MailngineError } from './errors.js';
import { ApiKeysResource } from './resources/api-keys.js';
import { DomainsResource } from './resources/domains.js';
import { EmailsResource } from './resources/emails.js';
import { TemplatesResource } from './resources/templates.js';
import { WebhooksResource } from './resources/webhooks.js';
import type { ListMeta, ListResponse } from './types.js';

const VERSION = '0.1.0';
const DEFAULT_BASE_URL = 'https://api.mailngine.com';
const MAX_RETRIES = 3;
const RETRYABLE_STATUS_CODES = new Set([429, 500, 502, 503, 504]);

export interface MailngineOptions {
  /** Override the base URL (e.g. for self-hosted or staging). */
  baseURL?: string;
}

/** @internal Options for a single request. */
export interface RequestOptions {
  method?: string;
  body?: unknown;
  /**
   * When set to `'list'` the response is expected to have `data` + `meta` fields
   * and will be returned as a `ListResponse<T>` instead of unwrapping `data`.
   */
  envelope?: 'list';
}

/**
 * Mailngine API client.
 *
 * @example
 * ```ts
 * import { Mailngine } from 'mailngine';
 *
 * const client = new Mailngine('mn_live_...');
 *
 * const email = await client.emails.send({
 *   from: 'hello@example.com',
 *   to: ['user@example.com'],
 *   subject: 'Hello!',
 *   html: '<h1>Welcome</h1>',
 * });
 * ```
 */
export class Mailngine {
  private readonly apiKey: string;
  private readonly baseURL: string;

  readonly emails: EmailsResource;
  readonly domains: DomainsResource;
  readonly webhooks: WebhooksResource;
  readonly templates: TemplatesResource;
  readonly apiKeys: ApiKeysResource;

  constructor(apiKey: string, options?: MailngineOptions) {
    if (!apiKey) {
      throw new Error(
        'Missing API key. Pass it to the Mailngine constructor: new Mailngine("mn_live_...")',
      );
    }

    this.apiKey = apiKey;
    this.baseURL = (options?.baseURL ?? DEFAULT_BASE_URL).replace(/\/+$/, '');

    this.emails = new EmailsResource(this);
    this.domains = new DomainsResource(this);
    this.webhooks = new WebhooksResource(this);
    this.templates = new TemplatesResource(this);
    this.apiKeys = new ApiKeysResource(this);
  }

  /**
   * Internal method used by resource classes to make authenticated requests.
   *
   * Handles JSON serialization, the `{ data }` response envelope, retries
   * on 429 / 5xx, and error translation.
   *
   * @internal
   */
  async request<T>(path: string, options?: RequestOptions): Promise<T> {
    const url = `${this.baseURL}${path}`;
    const method = options?.method ?? 'GET';

    const headers: Record<string, string> = {
      Authorization: `Bearer ${this.apiKey}`,
      'Content-Type': 'application/json',
      'User-Agent': `mailngine-node/${VERSION}`,
    };

    const fetchOptions: RequestInit = { method, headers };

    if (options?.body !== undefined) {
      fetchOptions.body = JSON.stringify(options.body);
    }

    let lastError: MailngineError | undefined;

    for (let attempt = 0; attempt < MAX_RETRIES; attempt++) {
      if (attempt > 0) {
        await this.sleep(this.backoffMs(attempt));
      }

      const response = await fetch(url, fetchOptions);

      if (response.ok) {
        // Some endpoints (DELETE) may return no body.
        const text = await response.text();
        if (!text) {
          return undefined as T;
        }

        const json = JSON.parse(text) as { data?: T; meta?: ListMeta; error?: { code: string; message: string } };

        if (options?.envelope === 'list') {
          return { data: json.data, meta: json.meta } as T;
        }

        // Standard envelope: unwrap `data`.
        return json.data as T;
      }

      // Parse the error envelope.
      const errorBody = await this.parseErrorBody(response);
      lastError = new MailngineError(
        response.status,
        errorBody.code,
        errorBody.message,
      );

      // Only retry on retryable status codes and non-final attempts.
      if (!RETRYABLE_STATUS_CODES.has(response.status)) {
        throw lastError;
      }
    }

    // All retries exhausted.
    throw lastError ?? new MailngineError(500, 'unknown_error', 'Request failed after retries');
  }

  // ---------------------------------------------------------------------------
  // Helpers
  // ---------------------------------------------------------------------------

  private async parseErrorBody(
    response: Response,
  ): Promise<{ code: string; message: string }> {
    try {
      const json = (await response.json()) as {
        error?: { code?: string; message?: string };
      };
      return {
        code: json.error?.code ?? 'unknown_error',
        message: json.error?.message ?? response.statusText,
      };
    } catch {
      return { code: 'unknown_error', message: response.statusText };
    }
  }

  /** Exponential backoff: 500ms, 1s, 2s ... */
  private backoffMs(attempt: number): number {
    return Math.min(500 * Math.pow(2, attempt - 1), 5000);
  }

  private sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}
