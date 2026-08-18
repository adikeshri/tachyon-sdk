from typing import Optional, cast

from .http import HttpClient
from .types import AnalyticsLatencyResponse, AnalyticsQueriesResponse


class AnalyticsApi:
    """`/analytics/*` — recorded automatically from search traffic, in memory only."""

    def __init__(self, http: HttpClient) -> None:
        self._http = http

    def top(self, *, collection: Optional[str] = None, limit: Optional[int] = None) -> AnalyticsQueriesResponse:
        """`GET /analytics/top`."""
        return cast(
            AnalyticsQueriesResponse,
            self._http.request("GET", "/analytics/top", query={"collection": collection, "limit": limit}),
        )

    def zero_results(self, *, collection: Optional[str] = None, limit: Optional[int] = None) -> AnalyticsQueriesResponse:
        """`GET /analytics/zero-results`. Ranks by how often a query came back empty."""
        return cast(
            AnalyticsQueriesResponse,
            self._http.request("GET", "/analytics/zero-results", query={"collection": collection, "limit": limit}),
        )

    def latency(self) -> AnalyticsLatencyResponse:
        """`GET /analytics/latency`."""
        return cast(AnalyticsLatencyResponse, self._http.request("GET", "/analytics/latency"))
