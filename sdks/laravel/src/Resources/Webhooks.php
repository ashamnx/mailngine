<?php

declare(strict_types=1);

namespace HelloMail\Resources;

use HelloMail\HelloMail;

/**
 * Manage webhooks through the Hello Mail API.
 *
 * @see POST   /v1/webhooks                    Create a webhook
 * @see GET    /v1/webhooks                    List webhooks
 * @see GET    /v1/webhooks/{id}               Get a webhook
 * @see PATCH  /v1/webhooks/{id}               Update a webhook
 * @see DELETE /v1/webhooks/{id}               Delete a webhook
 * @see GET    /v1/webhooks/{id}/deliveries    List deliveries
 */
class Webhooks
{
    public function __construct(private HelloMail $client)
    {
    }

    /**
     * Create a new webhook endpoint.
     *
     * @param  array{
     *     url: string,
     *     events: string[],
     * }  $params  Webhook configuration.
     * @return array  The created webhook object.
     */
    public function create(array $params): array
    {
        return $this->client->request('POST', '/v1/webhooks', $params);
    }

    /**
     * List all webhooks for the organization.
     *
     * @return array  List of webhook objects.
     */
    public function list(): array
    {
        return $this->client->request('GET', '/v1/webhooks');
    }

    /**
     * Get a webhook by ID.
     *
     * @param  string  $id  Webhook UUID.
     * @return array  The webhook object.
     */
    public function get(string $id): array
    {
        return $this->client->request('GET', '/v1/webhooks/' . urlencode($id));
    }

    /**
     * Update a webhook.
     *
     * @param  string  $id      Webhook UUID.
     * @param  array{
     *     url: string,
     *     events: string[],
     *     is_active: bool,
     * }  $params  Updated webhook configuration.
     * @return array  The updated webhook object.
     */
    public function update(string $id, array $params): array
    {
        return $this->client->request('PATCH', '/v1/webhooks/' . urlencode($id), $params);
    }

    /**
     * Delete a webhook.
     *
     * @param  string  $id  Webhook UUID.
     * @return array  Confirmation message.
     */
    public function delete(string $id): array
    {
        return $this->client->request('DELETE', '/v1/webhooks/' . urlencode($id));
    }

    /**
     * List delivery attempts for a webhook.
     *
     * @param  string  $id       Webhook UUID.
     * @param  array{
     *     page?: int,
     *     per_page?: int,
     * }  $options  Pagination options.
     * @return array  List of delivery objects.
     */
    public function listDeliveries(string $id, array $options = []): array
    {
        return $this->client->request('GET', '/v1/webhooks/' . urlencode($id) . '/deliveries', null, $options);
    }
}
