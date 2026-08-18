import json as json_lib
from typing import Any, Dict, Mapping, Optional, Tuple, cast

import requests

from .errors import TachyonConnectionError, TachyonTimeoutError, error_from_response

API_KEY_HEADER = "X-TACHYON-API-KEY"


class HttpClient:
    """Thin JSON-over-HTTP client shared by every resource in the SDK."""

    def __init__(
        self,
        url: str,
        *,
        api_key: Optional[str] = None,
        timeout: float = 15.0,
        headers: Optional[Mapping[str, str]] = None,
        session: Optional[requests.Session] = None,
    ) -> None:
        if not url:
            raise ValueError("Tachyon client requires a `url`.")
        self._base_url = url.rstrip("/")
        self._api_key = api_key
        self._timeout = timeout
        self._extra_headers = dict(headers or {})
        self._session = session or requests.Session()

    def request(
        self,
        method: str,
        path: str,
        *,
        query: Optional[Mapping[str, Any]] = None,
        body: Any = None,
    ) -> Any:
        status, text = self._send(method, path, query=query, body=body)
        if status == 204 or not text:
            return None

        try:
            payload = json_lib.loads(text)
        except ValueError:
            if status >= 400:
                raise error_from_response(status, text, "internal_error")
            raise TachyonConnectionError(f"Tachyon returned a non-JSON response: {_truncate(text)}")

        if status >= 400:
            raise error_from_response(status, _extract_message(payload, text), _extract_code(payload))

        return payload

    def request_text(self, method: str, path: str, *, query: Optional[Mapping[str, Any]] = None) -> str:
        status, text = self._send(method, path, query=query)
        if status >= 400:
            try:
                payload = json_lib.loads(text)
            except ValueError:
                raise error_from_response(status, text or f"HTTP {status}", "internal_error")
            raise error_from_response(status, _extract_message(payload, text), _extract_code(payload))
        return text

    def _send(
        self,
        method: str,
        path: str,
        *,
        query: Optional[Mapping[str, Any]] = None,
        body: Any = None,
    ) -> Tuple[int, str]:
        url = self._base_url + path
        headers: Dict[str, str] = {"Accept": "application/json", **self._extra_headers}
        if self._api_key:
            headers[API_KEY_HEADER] = self._api_key

        kwargs: Dict[str, Any] = {"headers": headers, "params": _clean_query(query), "timeout": self._timeout}
        if body is not None:
            headers["Content-Type"] = "application/json"
            kwargs["data"] = json_lib.dumps(body)

        try:
            response = self._session.request(method, url, **kwargs)
        except requests.Timeout as exc:
            raise TachyonTimeoutError(f"Request to {url} timed out after {self._timeout}s") from exc
        except requests.RequestException as exc:
            raise TachyonConnectionError(f"Failed to reach Tachyon at {url}: {exc}", exc) from exc

        return response.status_code, response.text


def _clean_query(query: Optional[Mapping[str, Any]]) -> Dict[str, Any]:
    """Drops unset params and lowercases booleans to match the JSON convention the API expects."""
    if not query:
        return {}
    cleaned: Dict[str, Any] = {}
    for key, value in query.items():
        if value is None:
            continue
        cleaned[key] = "true" if value is True else "false" if value is False else value
    return cleaned


def _extract_message(payload: Any, fallback: str) -> str:
    if isinstance(payload, dict):
        error = payload.get("error")
        if isinstance(error, dict) and isinstance(error.get("message"), str):
            return cast(str, error["message"])
    return fallback


def _extract_code(payload: Any) -> str:
    if isinstance(payload, dict):
        error = payload.get("error")
        if isinstance(error, dict) and isinstance(error.get("code"), str):
            return cast(str, error["code"])
    return "internal_error"


def _truncate(text: str, max_len: int = 200) -> str:
    return text if len(text) <= max_len else text[:max_len] + "…"
