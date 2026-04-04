<?php

declare(strict_types=1);

namespace HelloMail\Resources;

use HelloMail\HelloMail;

/**
 * Manage sending domains through the Hello Mail API.
 *
 * @see POST   /v1/domains              Create a domain
 * @see GET    /v1/domains              List domains
 * @see GET    /v1/domains/{id}         Get a domain
 * @see PATCH  /v1/domains/{id}         Update a domain
 * @see DELETE /v1/domains/{id}         Delete a domain
 * @see POST   /v1/domains/{id}/verify  Trigger DNS verification
 */
class Domains
{
    public function __construct(private HelloMail $client)
    {
    }

    /**
     * Register a new sending domain.
     *
     * @param  string  $name  The domain name (e.g. "example.com").
     * @return array  The created domain with its required DNS records.
     */
    public function create(string $name): array
    {
        return $this->client->request('POST', '/v1/domains', ['name' => $name]);
    }

    /**
     * List all domains for the organization.
     *
     * @return array  List of domain objects.
     */
    public function list(): array
    {
        return $this->client->request('GET', '/v1/domains');
    }

    /**
     * Get a domain by ID.
     *
     * @param  string  $id  Domain UUID.
     * @return array  The domain object.
     */
    public function get(string $id): array
    {
        return $this->client->request('GET', '/v1/domains/' . urlencode($id));
    }

    /**
     * Update domain settings.
     *
     * @param  string  $id      Domain UUID.
     * @param  array{
     *     open_tracking?: bool,
     *     click_tracking?: bool,
     * }  $params  Fields to update.
     * @return array  The updated domain object.
     */
    public function update(string $id, array $params): array
    {
        return $this->client->request('PATCH', '/v1/domains/' . urlencode($id), $params);
    }

    /**
     * Delete a domain.
     *
     * @param  string  $id  Domain UUID.
     * @return array  Confirmation message.
     */
    public function delete(string $id): array
    {
        return $this->client->request('DELETE', '/v1/domains/' . urlencode($id));
    }

    /**
     * Trigger DNS verification for a domain.
     *
     * @param  string  $id  Domain UUID.
     * @return array  The DNS record verification results.
     */
    public function verify(string $id): array
    {
        return $this->client->request('POST', '/v1/domains/' . urlencode($id) . '/verify');
    }
}
