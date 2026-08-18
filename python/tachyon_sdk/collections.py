from typing import List, cast
from urllib.parse import quote

from .http import HttpClient
from .types import CollectionInfo, CollectionSchema


class CollectionsApi:
    """`/collections` — create, list, and remove collections."""

    def __init__(self, http: HttpClient) -> None:
        self._http = http

    def create(self, schema: CollectionSchema) -> CollectionInfo:
        """`POST /collections`. Field types are immutable after creation."""
        return cast(CollectionInfo, self._http.request("POST", "/collections", body=schema))

    def list(self) -> List[CollectionInfo]:
        """`GET /collections`."""
        return cast(List[CollectionInfo], self._http.request("GET", "/collections"))

    def retrieve(self, name: str) -> CollectionInfo:
        """`GET /collections/{name}`."""
        return cast(CollectionInfo, self._http.request("GET", f"/collections/{quote(name, safe='')}"))

    def delete(self, name: str) -> None:
        """`DELETE /collections/{name}`. Removes the collection and all its data."""
        self._http.request("DELETE", f"/collections/{quote(name, safe='')}")
