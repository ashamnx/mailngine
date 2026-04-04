/**
 * Error thrown when the Hello Mail API returns a non-2xx response.
 */
export class HelloMailError extends Error {
  /** HTTP status code from the API. */
  public readonly statusCode: number;

  /** Machine-readable error code from the API (e.g. "bad_request", "not_found"). */
  public readonly code: string;

  constructor(statusCode: number, code: string, message: string) {
    super(message);
    this.name = 'HelloMailError';
    this.statusCode = statusCode;
    this.code = code;

    // Maintain proper prototype chain for instanceof checks.
    Object.setPrototypeOf(this, new.target.prototype);
  }
}
