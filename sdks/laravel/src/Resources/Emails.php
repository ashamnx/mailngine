<?php

declare(strict_types=1);

namespace Mailngine\Resources;

use Mailngine\Mailngine;

/**
 * Manage emails through the Mailngine API.
 *
 * @see POST   /v1/emails       Send an email
 * @see GET    /v1/emails       List emails
 * @see GET    /v1/emails/{id}  Get an email
 */
class Emails
{
    public function __construct(private Mailngine $client)
    {
    }

    /**
     * Send an email.
     *
     * @param  array{
     *     from: string,
     *     to: string[],
     *     cc?: string[],
     *     bcc?: string[],
     *     reply_to?: string,
     *     subject: string,
     *     html?: string,
     *     text?: string,
     *     headers?: array<string, string>,
     *     tags?: string[],
     *     template_id?: string,
     *     template_data?: array<string, string>,
     *     idempotency_key?: string,
     *     scheduled_at?: string,
     * }  $params  Email parameters.
     * @return array  The created email object.
     */
    public function send(array $params): array
    {
        return $this->client->request('POST', '/v1/emails', $params);
    }

    /**
     * Get an email by ID.
     *
     * @param  string  $id  Email UUID.
     * @return array  The email object.
     */
    public function get(string $id): array
    {
        return $this->client->request('GET', '/v1/emails/' . urlencode($id));
    }

    /**
     * List emails with pagination.
     *
     * @param  array{
     *     page?: int,
     *     per_page?: int,
     * }  $options  Pagination options.
     * @return array  List of email objects (with meta if paginated).
     */
    public function list(array $options = []): array
    {
        return $this->client->request('GET', '/v1/emails', null, $options);
    }
}
