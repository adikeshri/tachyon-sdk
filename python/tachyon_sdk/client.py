from typing import Mapping, Optional, cast

import requests

from .analytics import AnalyticsApi
from .collection import Collection
from .collections import CollectionsApi
from .http import HttpClient
from .types import HealthResponse


class Tachyon:
    """
    Client for a single Tachyon server.

        client = Tachyon(url="http://localhost:8108", api_key="my-admin-key")
        client.collections.create({"name": "products", "fields": [{"name": "title", "type": "text"}]})
        client.collection("products").documents.index({"id": "1", "title": "Wireless Mouse"})
        results = client.collection("products").search(q="wireless mouse")
    """

    def __init__(
        self,
        url: Optional[str] = None,
        *,
        host: str = "localhost",
        port: int = 8108,
        protocol: str = "http",
        api_key: Optional[str] = None,
        timeout: float = 15.0,
        headers: Optional[Mapping[str, str]] = None,
        session: Optional[requests.Session] = None,
    ) -> None:
        """
        Args:
            url: Full base URL, e.g. `http://localhost:8108`. Takes precedence
                over host/port/protocol.
            host, port, protocol: Used to build the URL when `url` is omitted.
            api_key: Sent as `X-TACHYON-API-KEY`. Use an admin key for writes,
                a search key for read-only access.
            timeout: Per-request timeout in seconds. Default 15.
            headers: Extra headers merged into every request.
            session: Override the `requests.Session` (mainly for testing, or
                to share connection pooling with other clients).
        """
        resolved_url = url or f"{protocol}://{host}:{port}"
        self._http = HttpClient(resolved_url, api_key=api_key, timeout=timeout, headers=headers, session=session)
        self.collections = CollectionsApi(self._http)
        self.analytics = AnalyticsApi(self._http)

    def collection(self, name: str) -> Collection:
        """Get a handle scoped to one collection, for documents/search/suggest."""
        return Collection(self._http, name)

    def health(self) -> HealthResponse:
        """`GET /health`. Always reachable without an API key."""
        return cast(HealthResponse, self._http.request("GET", "/health"))

    def metrics(self) -> str:
        """`GET /metrics`. Prometheus exposition format, returned as plain text."""
        return self._http.request_text("GET", "/metrics")
