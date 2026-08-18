/**
 * Base class for every error Tachyon's HTTP API returns. `code` is the
 * stable machine-readable string from the response body (see
 * https://github.com/adikeshri/tachyon/blob/main/docs/api.md#errors);
 * `status` is the HTTP status code.
 */
export class TachyonError extends Error {
  readonly code: string;
  readonly status: number;

  constructor(message: string, code: string, status: number) {
    super(message);
    this.name = 'TachyonError';
    this.code = code;
    this.status = status;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

/** 400 — invalid_schema, invalid_document, invalid_query, invalid_json. */
export class TachyonRequestError extends TachyonError {
  constructor(message: string, code: string, status: number) {
    super(message, code, status);
    this.name = 'TachyonRequestError';
  }
}

/** 401 — missing or wrong API key. */
export class TachyonAuthenticationError extends TachyonError {
  constructor(message: string, code: string, status: number) {
    super(message, code, status);
    this.name = 'TachyonAuthenticationError';
  }
}

/** 403 — a search key attempted a write. */
export class TachyonAuthorizationError extends TachyonError {
  constructor(message: string, code: string, status: number) {
    super(message, code, status);
    this.name = 'TachyonAuthorizationError';
  }
}

/** 404 — collection_not_found, document_not_found. */
export class TachyonNotFoundError extends TachyonError {
  constructor(message: string, code: string, status: number) {
    super(message, code, status);
    this.name = 'TachyonNotFoundError';
  }
}

/** 409 — collection_exists. */
export class TachyonConflictError extends TachyonError {
  constructor(message: string, code: string, status: number) {
    super(message, code, status);
    this.name = 'TachyonConflictError';
  }
}

/** 5xx — corrupt_data, io_error, internal_error. */
export class TachyonServerError extends TachyonError {
  constructor(message: string, code: string, status: number) {
    super(message, code, status);
    this.name = 'TachyonServerError';
  }
}

/** The request never reached the server, or the server never replied at all. */
export class TachyonConnectionError extends Error {
  readonly cause?: unknown;

  constructor(message: string, cause?: unknown) {
    super(message);
    this.name = 'TachyonConnectionError';
    this.cause = cause;
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

/** The request was aborted after exceeding its configured timeout. */
export class TachyonTimeoutError extends TachyonConnectionError {
  constructor(message: string) {
    super(message);
    this.name = 'TachyonTimeoutError';
  }
}

/** Maps an HTTP status + error code/message onto the right error subclass. */
export function errorFromResponse(status: number, message: string, code: string): TachyonError {
  switch (status) {
    case 400:
      return new TachyonRequestError(message, code, status);
    case 401:
      return new TachyonAuthenticationError(message, code, status);
    case 403:
      return new TachyonAuthorizationError(message, code, status);
    case 404:
      return new TachyonNotFoundError(message, code, status);
    case 409:
      return new TachyonConflictError(message, code, status);
    default:
      return status >= 500
        ? new TachyonServerError(message, code, status)
        : new TachyonError(message, code, status);
  }
}
