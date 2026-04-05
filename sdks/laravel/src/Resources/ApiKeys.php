<?php

declare(strict_types=1);

namespace Mailngine\Resources;

use Mailngine\Mailngine;

/**
 * Manage API keys through the Mailngine API.
 *
 * @see POST   /v1/api-keys       Create an API key
 * @see GET    /v1/api-keys       List API keys
 * @see DELETE /v1/api-keys/{id}  Revoke an API key
 */
class ApiKeys
{
    public function __construct(private Mailngine $client)
    {
    }

    /**
     * Create a new API key.
     *
     * The full key value is only returned once in the create response.
     * Store it securely -- it cannot be retrieved again.
     *
     * @param  array{
     *     name: string,
     *     permission?: string,
     *     domain_id?: string,
     *     expires_at?: string,
     * }  $params  API key configuration.
     * @return array  The created API key object (includes the full key).
     */
    public function create(array $params): array
    {
        return $this->client->request('POST', '/v1/api-keys', $params);
    }

    /**
     * List all active API keys for the organization.
     *
     * @return array  List of API key objects (without full key values).
     */
    public function list(): array
    {
        return $this->client->request('GET', '/v1/api-keys');
    }

    /**
     * Revoke an API key.
     *
     * @param  string  $id  API key UUID.
     * @return array  Confirmation message.
     */
    public function revoke(string $id): array
    {
        return $this->client->request('DELETE', '/v1/api-keys/' . urlencode($id));
    }
}
