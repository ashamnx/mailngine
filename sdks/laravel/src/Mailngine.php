<?php

declare(strict_types=1);

namespace Mailngine;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\ConnectException;
use GuzzleHttp\Exception\RequestException;
use Mailngine\Exceptions\ApiException;
use Mailngine\Exceptions\MailngineException;
use Mailngine\Resources\ApiKeys;
use Mailngine\Resources\Domains;
use Mailngine\Resources\Emails;
use Mailngine\Resources\Templates;
use Mailngine\Resources\Webhooks;

/**
 * Mailngine API client.
 *
 * @see https://mailngine.com/docs
 */
class Mailngine
{
    private const VERSION = '0.1.0';

    private const DEFAULT_BASE_URL = 'https://api.mailngine.com';

    private const MAX_RETRIES = 3;

    private const RETRY_DELAY_MS = 500;

    private string $apiKey;

    private string $baseUrl;

    private Client $http;

    /**
     * @param  string  $apiKey   Your Mailngine API key.
     * @param  array{
     *     base_url?: string,
     *     timeout?: int,
     *     http_client?: Client,
     * }  $options  Optional configuration.
     */
    public function __construct(string $apiKey, array $options = [])
    {
        if ($apiKey === '') {
            throw new MailngineException('API key is required');
        }

        $this->apiKey = $apiKey;
        $this->baseUrl = rtrim($options['base_url'] ?? self::DEFAULT_BASE_URL, '/');

        $this->http = $options['http_client'] ?? new Client([
            'base_uri' => $this->baseUrl,
            'timeout' => $options['timeout'] ?? 30,
            'http_errors' => false,
        ]);
    }

    public function emails(): Emails
    {
        return new Emails($this);
    }

    public function domains(): Domains
    {
        return new Domains($this);
    }

    public function webhooks(): Webhooks
    {
        return new Webhooks($this);
    }

    public function templates(): Templates
    {
        return new Templates($this);
    }

    public function apiKeys(): ApiKeys
    {
        return new ApiKeys($this);
    }

    /**
     * Send an HTTP request to the Mailngine API.
     *
     * Automatically retries on 429 (rate limit) and 5xx (server error) responses
     * up to MAX_RETRIES attempts with exponential backoff.
     *
     * @param  string      $method  HTTP method (GET, POST, PATCH, DELETE).
     * @param  string      $path    API path (e.g. "/v1/emails").
     * @param  array|null  $body    Request body (JSON-encoded for POST/PATCH/PUT).
     * @param  array       $query   Query parameters for GET requests.
     * @return array                Decoded response data from the "data" envelope.
     *
     * @throws ApiException           On API error responses.
     * @throws MailngineException     On connection or unexpected errors.
     */
    public function request(string $method, string $path, ?array $body = null, array $query = []): array
    {
        $url = $this->baseUrl . $path;

        $requestOptions = [
            'headers' => [
                'Authorization' => 'Bearer ' . $this->apiKey,
                'Content-Type' => 'application/json',
                'Accept' => 'application/json',
                'User-Agent' => 'mailngine-laravel/' . self::VERSION,
            ],
        ];

        if ($body !== null && in_array(strtoupper($method), ['POST', 'PATCH', 'PUT'], true)) {
            $requestOptions['json'] = $body;
        }

        if (! empty($query)) {
            $requestOptions['query'] = $query;
        }

        $attempt = 0;
        $lastException = null;

        while ($attempt < self::MAX_RETRIES) {
            $attempt++;

            try {
                $response = $this->http->request($method, $url, $requestOptions);
            } catch (ConnectException $e) {
                $lastException = new MailngineException(
                    'Failed to connect to Mailngine API: ' . $e->getMessage(),
                    0,
                    $e
                );

                if ($attempt < self::MAX_RETRIES) {
                    $this->sleep($attempt);
                    continue;
                }

                throw $lastException;
            } catch (\Throwable $e) {
                throw new MailngineException(
                    'Unexpected error communicating with Mailngine API: ' . $e->getMessage(),
                    0,
                    $e
                );
            }

            $statusCode = $response->getStatusCode();

            // Retry on 429 or 5xx
            if ($statusCode === 429 || $statusCode >= 500) {
                if ($attempt < self::MAX_RETRIES) {
                    $this->sleep($attempt);
                    continue;
                }
            }

            $responseBody = (string) $response->getBody();
            $decoded = json_decode($responseBody, true);

            if (json_last_error() !== JSON_ERROR_NONE) {
                throw new MailngineException(
                    'Invalid JSON response from Mailngine API'
                );
            }

            // Handle error responses
            if ($statusCode >= 400) {
                $errorCode = $decoded['error']['code'] ?? 'unknown_error';
                $errorMessage = $decoded['error']['message'] ?? 'An unknown error occurred';

                throw new ApiException($statusCode, $errorCode, $errorMessage);
            }

            // Unwrap the data envelope
            return $decoded['data'] ?? $decoded;
        }

        // Should not reach here, but handle gracefully
        throw $lastException ?? new MailngineException('Request failed after ' . self::MAX_RETRIES . ' attempts');
    }

    /**
     * Sleep with exponential backoff between retry attempts.
     */
    private function sleep(int $attempt): void
    {
        $delayMs = self::RETRY_DELAY_MS * (2 ** ($attempt - 1));
        usleep($delayMs * 1000);
    }
}
