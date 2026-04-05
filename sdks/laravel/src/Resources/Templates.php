<?php

declare(strict_types=1);

namespace Mailngine\Resources;

use Mailngine\Mailngine;

/**
 * Manage email templates through the Mailngine API.
 *
 * @see POST   /v1/templates               Create a template
 * @see GET    /v1/templates               List templates
 * @see GET    /v1/templates/{id}          Get a template
 * @see PATCH  /v1/templates/{id}          Update a template
 * @see DELETE /v1/templates/{id}          Delete a template
 * @see POST   /v1/templates/{id}/preview  Preview a rendered template
 */
class Templates
{
    public function __construct(private Mailngine $client)
    {
    }

    /**
     * Create a new email template.
     *
     * @param  array{
     *     name: string,
     *     subject: string,
     *     html_body: string,
     *     text_body?: string,
     *     variables?: string[],
     * }  $params  Template content.
     * @return array  The created template object.
     */
    public function create(array $params): array
    {
        return $this->client->request('POST', '/v1/templates', $params);
    }

    /**
     * List all templates for the organization.
     *
     * @return array  List of template objects.
     */
    public function list(): array
    {
        return $this->client->request('GET', '/v1/templates');
    }

    /**
     * Get a template by ID.
     *
     * @param  string  $id  Template UUID.
     * @return array  The template object.
     */
    public function get(string $id): array
    {
        return $this->client->request('GET', '/v1/templates/' . urlencode($id));
    }

    /**
     * Update a template.
     *
     * @param  string  $id      Template UUID.
     * @param  array{
     *     name: string,
     *     subject: string,
     *     html_body: string,
     *     text_body?: string,
     *     variables?: string[],
     * }  $params  Updated template content.
     * @return array  The updated template object.
     */
    public function update(string $id, array $params): array
    {
        return $this->client->request('PATCH', '/v1/templates/' . urlencode($id), $params);
    }

    /**
     * Delete a template.
     *
     * @param  string  $id  Template UUID.
     * @return array  Confirmation message.
     */
    public function delete(string $id): array
    {
        return $this->client->request('DELETE', '/v1/templates/' . urlencode($id));
    }

    /**
     * Preview a rendered template with sample data.
     *
     * @param  string                   $id    Template UUID.
     * @param  array<string, string>    $data  Variable values to substitute.
     * @return array{subject: string, html_body: string, text_body: string}
     */
    public function preview(string $id, array $data): array
    {
        return $this->client->request('POST', '/v1/templates/' . urlencode($id) . '/preview', [
            'data' => $data,
        ]);
    }
}
