"""
Exception hierarchy for Tachyon's HTTP API. `code` is the stable
machine-readable string from the response body (see
https://github.com/adikeshri/tachyon/blob/main/docs/api.md#errors);
`status` is the HTTP status code.
"""

from typing import Dict, Optional, Type


class TachyonError(Exception):
    def __init__(self, message: str, code: str, status: int) -> None:
        super().__init__(message)
        self.message = message
        self.code = code
        self.status = status

    def __str__(self) -> str:
        return f"{self.message} (code={self.code}, status={self.status})"


class TachyonRequestError(TachyonError):
    """400 — invalid_schema, invalid_document, invalid_query, invalid_json."""


class TachyonAuthenticationError(TachyonError):
    """401 — missing or wrong API key."""


class TachyonAuthorizationError(TachyonError):
    """403 — a search key attempted a write."""


class TachyonNotFoundError(TachyonError):
    """404 — collection_not_found, document_not_found."""


class TachyonConflictError(TachyonError):
    """409 — collection_exists."""


class TachyonServerError(TachyonError):
    """5xx — corrupt_data, io_error, internal_error."""


class TachyonConnectionError(Exception):
    """The request never reached the server, or the server never replied at all."""

    def __init__(self, message: str, cause: Optional[BaseException] = None) -> None:
        super().__init__(message)
        self.cause = cause


class TachyonTimeoutError(TachyonConnectionError):
    """The request was aborted after exceeding its configured timeout."""


_STATUS_TO_ERROR: Dict[int, Type[TachyonError]] = {
    400: TachyonRequestError,
    401: TachyonAuthenticationError,
    403: TachyonAuthorizationError,
    404: TachyonNotFoundError,
    409: TachyonConflictError,
}


def error_from_response(status: int, message: str, code: str) -> TachyonError:
    """Maps an HTTP status + error code/message onto the right error subclass."""
    error_cls = _STATUS_TO_ERROR.get(status)
    if error_cls is None:
        error_cls = TachyonServerError if status >= 500 else TachyonError
    return error_cls(message, code, status)
