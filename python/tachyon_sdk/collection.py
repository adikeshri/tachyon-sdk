from typing import List, Optional, Union, cast
from urllib.parse import quote

from .documents import DocumentsApi
from .http import HttpClient
from .types import MatchMode, SearchResponse, SuggestResponse


def _join(value: Optional[Union[str, List[str]]]) -> Optional[str]:
    if value is None:
        return None
    return ",".join(value) if isinstance(value, list) else value


class Collection:
    """
    A handle scoped to one collection. Get one via `client.collection(name)`;
    it does not verify the collection exists until you make a request.
    """

    def __init__(self, http: HttpClient, name: str) -> None:
        self.name = name
        self.documents = DocumentsApi(http, name)
        self._http = http

    def search(
        self,
        q: Optional[str] = None,
        *,
        query_by: Optional[Union[str, List[str]]] = None,
        filter: Optional[str] = None,
        sort: Optional[str] = None,
        facet: Optional[Union[str, List[str]]] = None,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        prefix: Optional[bool] = None,
        typo_tolerance: Optional[bool] = None,
        match_mode: Optional[MatchMode] = None,
    ) -> SearchResponse:
        """`GET /collections/{name}/search`."""
        query = {
            "q": q,
            "query_by": _join(query_by),
            "filter": filter,
            "sort": sort,
            "facet": _join(facet),
            "limit": limit,
            "offset": offset,
            "prefix": prefix,
            "typo_tolerance": typo_tolerance,
            "match_mode": match_mode,
        }
        return cast(
            SearchResponse,
            self._http.request("GET", f"/collections/{quote(self.name, safe='')}/search", query=query),
        )

    def suggest(
        self,
        q: str,
        *,
        query_by: Optional[Union[str, List[str]]] = None,
        limit: Optional[int] = None,
        typo_tolerance: Optional[bool] = None,
    ) -> SuggestResponse:
        """`GET /collections/{name}/suggest`."""
        query = {
            "q": q,
            "query_by": _join(query_by),
            "limit": limit,
            "typo_tolerance": typo_tolerance,
        }
        return cast(
            SuggestResponse,
            self._http.request("GET", f"/collections/{quote(self.name, safe='')}/suggest", query=query),
        )
