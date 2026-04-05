/**
 * Error thrown when the Mailngine API returns a non-2xx response.
 */
export class MailngineError extends Error {
  /** HTTP status code from the API. */
  public readonly statusCode: number;

  /** Machine-readable error code from the API (e.g. "bad_request", "not_found"). */
  public readonly code: string;

  constructor(statusCode: number, code: string, message: string) {
    super(message);
    this.name = 'MailngineError';
    this.statusCode = statusCode;
    this.code = code;

    // Maintain proper prototype chain for instanceof checks.
    Object.setPrototypeOf(this, new.target.prototype);
  }
}
