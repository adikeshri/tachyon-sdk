from typing import List, Union, cast
from urllib.parse import quote

from .http import HttpClient
from .types import DocumentsIndexResponse, TachyonDocument


class DocumentsApi:
    """`/collections/{name}/documents` — index, fetch, and delete documents."""

    def __init__(self, http: HttpClient, collection_name: str) -> None:
        self._http = http
        self._collection_name = collection_name

    def index(self, documents: Union[TachyonDocument, List[TachyonDocument]]) -> DocumentsIndexResponse:
        """
        `POST /collections/{name}/documents`. Accepts one document or a list;
        documents are upserted by `id`. Individual documents can fail without
        failing their neighbours — check `num_failed` and `results`.
        """
        return cast(
            DocumentsIndexResponse,
            self._http.request(
                "POST",
                f"/collections/{quote(self._collection_name, safe='')}/documents",
                body=documents,
            ),
        )

    def retrieve(self, document_id: str) -> TachyonDocument:
        """`GET /collections/{name}/documents/{id}`."""
        return cast(
            TachyonDocument,
            self._http.request(
                "GET",
                f"/collections/{quote(self._collection_name, safe='')}/documents/{quote(document_id, safe='')}",
            ),
        )

    def delete(self, document_id: str) -> None:
        """`DELETE /collections/{name}/documents/{id}`."""
        self._http.request(
            "DELETE",
            f"/collections/{quote(self._collection_name, safe='')}/documents/{quote(document_id, safe='')}",
        )
