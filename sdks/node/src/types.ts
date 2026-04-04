// ---------------------------------------------------------------------------
// Common
// ---------------------------------------------------------------------------

/** Pagination options for list endpoints. */
export interface ListOptions {
  /** Page number (1-based). */
  page?: number;
  /** Number of items per page (max 100). */
  perPage?: number;
}

/** Envelope metadata returned alongside paginated lists. */
export interface ListMeta {
  page: number;
  per_page: number;
  total: number;
  has_more: boolean;
  next_cursor?: string;
}

/** Paginated list response. */
export interface ListResponse<T> {
  data: T[];
  meta: ListMeta;
}

// ---------------------------------------------------------------------------
// Emails
// ---------------------------------------------------------------------------

/** Parameters for sending an email. */
export interface SendEmailParams {
  from: string;
  to: string[];
  cc?: string[];
  bcc?: string[];
  reply_to?: string;
  subject: string;
  html?: string;
  text?: string;
  headers?: Record<string, string>;
  tags?: string[];
  template_id?: string;
  template_data?: Record<string, string>;
  idempotency_key?: string;
  scheduled_at?: string;
}

/** An email record returned by the API. */
export interface Email {
  id: string;
  from: string;
  to: string[];
  cc?: string[];
  bcc?: string[];
  subject: string;
  status: string;
  message_id?: string;
  scheduled_at?: string;
  sent_at?: string;
  delivered_at?: string;
  created_at: string;
}

// ---------------------------------------------------------------------------
// Domains
// ---------------------------------------------------------------------------

/** Parameters for creating a domain. */
export interface CreateDomainParams {
  name: string;
}

/** Parameters for updating a domain. */
export interface UpdateDomainParams {
  open_tracking?: boolean;
  click_tracking?: boolean;
}

/** A DNS record associated with a domain. */
export interface DNSRecord {
  id: string;
  domain_id: string;
  record_type: string;
  host: string;
  value: string;
  purpose: string;
  status: string;
  verified_at?: string;
  created_at: string;
}

/** A sending domain returned by the API. */
export interface Domain {
  id: string;
  org_id: string;
  name: string;
  status: string;
  region: string;
  dkim_selector: string;
  open_tracking: boolean;
  click_tracking: boolean;
  verified_at?: string;
  created_at: string;
  updated_at: string;
}

/** Response from creating a domain (includes DNS records to configure). */
export interface CreateDomainResponse {
  domain: Domain;
  dns_records: DNSRecord[];
}

// ---------------------------------------------------------------------------
// Webhooks
// ---------------------------------------------------------------------------

/** Parameters for creating a webhook. */
export interface CreateWebhookParams {
  url: string;
  events: string[];
}

/** Parameters for updating a webhook. */
export interface UpdateWebhookParams {
  url: string;
  events: string[];
  is_active: boolean;
}

/** A webhook endpoint returned by the API. */
export interface Webhook {
  id: string;
  url: string;
  events: string[];
  is_active: boolean;
  secret: string;
  created_at: string;
  updated_at: string;
}

/** A webhook delivery attempt. */
export interface WebhookDelivery {
  id: string;
  webhook_id: string;
  event_type: string;
  status_code: number;
  response_body: string;
  attempted_at: string;
}

// ---------------------------------------------------------------------------
// Templates
// ---------------------------------------------------------------------------

/** Parameters for creating a template. */
export interface CreateTemplateParams {
  name: string;
  subject: string;
  html_body: string;
  text_body?: string;
  variables?: string[];
}

/** Parameters for updating a template. */
export interface UpdateTemplateParams {
  name: string;
  subject: string;
  html_body: string;
  text_body?: string;
  variables?: string[];
}

/** Parameters for previewing a template. */
export interface PreviewTemplateParams {
  data: Record<string, string>;
}

/** A template returned by the API. */
export interface Template {
  id: string;
  org_id: string;
  name: string;
  subject: string;
  html_body: string;
  text_body?: string;
  variables?: string[];
  created_at: string;
  updated_at: string;
}

/** Preview rendering result. */
export interface TemplatePreview {
  subject: string;
  html_body: string;
  text_body: string;
}

// ---------------------------------------------------------------------------
// API Keys
// ---------------------------------------------------------------------------

/** Parameters for creating an API key. */
export interface CreateApiKeyParams {
  name: string;
  permission?: 'full' | 'send_only' | 'read_only';
  domain_id?: string;
  expires_at?: string;
}

/** An API key as returned at creation time (includes the full key). */
export interface ApiKeyCreateResponse {
  id: string;
  name: string;
  prefix: string;
  key: string;
  permission: string;
  expires_at?: string;
  created_at: string;
}

/** An API key in list responses (full key is never returned again). */
export interface ApiKey {
  id: string;
  name: string;
  prefix: string;
  permission: string;
  last_used_at?: string;
  expires_at?: string;
  created_at: string;
}
