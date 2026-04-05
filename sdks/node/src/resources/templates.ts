import type { Mailngine } from '../client.js';
import type {
  CreateTemplateParams,
  PreviewTemplateParams,
  Template,
  TemplatePreview,
  UpdateTemplateParams,
} from '../types.js';

/**
 * Resource for managing email templates.
 */
export class TemplatesResource {
  constructor(private readonly client: Mailngine) {}

  /**
   * Create a new template.
   *
   * POST /v1/templates
   */
  async create(params: CreateTemplateParams): Promise<Template> {
    return this.client.request<Template>('/v1/templates', {
      method: 'POST',
      body: params,
    });
  }

  /**
   * List all templates for the current organization.
   *
   * GET /v1/templates
   */
  async list(): Promise<Template[]> {
    return this.client.request<Template[]>('/v1/templates');
  }

  /**
   * Retrieve a single template by ID.
   *
   * GET /v1/templates/:id
   */
  async get(id: string): Promise<Template> {
    return this.client.request<Template>(`/v1/templates/${encodeURIComponent(id)}`);
  }

  /**
   * Update a template.
   *
   * PATCH /v1/templates/:id
   */
  async update(id: string, params: UpdateTemplateParams): Promise<Template> {
    return this.client.request<Template>(`/v1/templates/${encodeURIComponent(id)}`, {
      method: 'PATCH',
      body: params,
    });
  }

  /**
   * Delete a template.
   *
   * DELETE /v1/templates/:id
   */
  async delete(id: string): Promise<void> {
    await this.client.request<unknown>(`/v1/templates/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    });
  }

  /**
   * Preview a template with sample data.
   *
   * POST /v1/templates/:id/preview
   */
  async preview(id: string, params: PreviewTemplateParams): Promise<TemplatePreview> {
    return this.client.request<TemplatePreview>(
      `/v1/templates/${encodeURIComponent(id)}/preview`,
      { method: 'POST', body: params },
    );
  }
}
