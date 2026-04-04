<?php

declare(strict_types=1);

namespace HelloMail\Exceptions;

/**
 * Represents a structured error response from the Hello Mail API.
 */
class ApiException extends HelloMailException
{
    public int $statusCode;

    public string $errorCode;

    public function __construct(int $statusCode, string $errorCode, string $message, ?\Throwable $previous = null)
    {
        $this->statusCode = $statusCode;
        $this->errorCode = $errorCode;

        parent::__construct($message, $statusCode, $previous);
    }
}
