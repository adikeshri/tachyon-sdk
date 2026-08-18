import json
from urllib.parse import parse_qs, urlparse

import responses

from .conftest import BASE_URL


@responses.activate
def test_index_documents_reports_per_document_results(client):
    responses.add(
        responses.POST,
        f"{BASE_URL}/collections/products/documents",
        json={
            "num_indexed": 1,
            "num_failed": 1,
            "results": [
                {"success": True, "id": "1"},
                {"success": False, "code": "invalid_document", "error": "field `price`: expected an integer, got a string"},
            ],
        },
        status=200,
    )

    result = client.collection("products").documents.index(
        [{"id": "1", "title": "Wireless Mouse"}, {"id": "2", "price": "not a number"}]
    )

    assert result["num_indexed"] == 1
    assert result["num_failed"] == 1
    assert result["results"][1]["code"] == "invalid_document"
    request = responses.calls[0].request
    assert json.loads(request.body) == [
        {"id": "1", "title": "Wireless Mouse"},
        {"id": "2", "price": "not a number"},
    ]


@responses.activate
def test_retrieve_document(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products/documents/1",
        json={"id": "1", "title": "Wireless Mouse"},
        status=200,
    )

    doc = client.collection("products").documents.retrieve("1")

    assert doc["title"] == "Wireless Mouse"


@responses.activate
def test_delete_document(client):
    responses.add(responses.DELETE, f"{BASE_URL}/collections/products/documents/1", status=204)

    client.collection("products").documents.delete("1")

    request = responses.calls[0].request
    assert request.method == "DELETE"
    assert urlparse(request.url).path == "/collections/products/documents/1"


@responses.activate
def test_search_serializes_params_into_the_expected_query_string(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products/search",
        json={"found": 1, "found_is_exact": True, "search_time_ms": 1, "hits": []},
        status=200,
    )

    client.collection("products").search(
        "wireless mouse",
        query_by=["title", "description"],
        filter="brand:=Logitech && price:<5000",
        sort="_text_match:desc,price:asc",
        facet=["brand", "year"],
        limit=20,
        offset=40,
        prefix=False,
        typo_tolerance=True,
        match_mode="any",
    )

    request = responses.calls[0].request
    assert urlparse(request.url).path == "/collections/products/search"
    query = parse_qs(urlparse(request.url).query)
    assert query["q"] == ["wireless mouse"]
    assert query["query_by"] == ["title,description"]
    assert query["filter"] == ["brand:=Logitech && price:<5000"]
    assert query["sort"] == ["_text_match:desc,price:asc"]
    assert query["facet"] == ["brand,year"]
    assert query["limit"] == ["20"]
    assert query["offset"] == ["40"]
    assert query["prefix"] == ["false"]
    assert query["typo_tolerance"] == ["true"]
    assert query["match_mode"] == ["any"]


@responses.activate
def test_search_omits_unset_params(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products/search",
        json={"found": 0, "found_is_exact": True, "search_time_ms": 0, "hits": []},
        status=200,
    )

    client.collection("products").search()

    query = parse_qs(urlparse(responses.calls[0].request.url).query)
    assert "filter" not in query
    assert "limit" not in query


@responses.activate
def test_search_returns_hits_facets_and_found_is_exact(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products/search",
        json={
            "found": 1240,
            "found_is_exact": False,
            "search_time_ms": 12,
            "hits": [{"document": {"id": "1", "title": "Wireless Mouse"}, "text_match": 554.788}],
            "facets": {"brand": {"Logitech": 1240, "Razer": 830}},
        },
        status=200,
    )

    results = client.collection("products").search(q="wireless mouse")

    assert results["found"] == 1240
    assert results["found_is_exact"] is False
    assert results["hits"][0]["document"]["title"] == "Wireless Mouse"
    assert results["facets"]["brand"]["Logitech"] == 1240


@responses.activate
def test_suggest(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products/suggest",
        json={
            "suggestions": [
                {"text": "wireless", "count": 3, "typos": 0},
                {"text": "wired", "count": 2, "typos": 0},
            ],
            "search_time_ms": 0,
        },
        status=200,
    )

    result = client.collection("products").suggest("wir", limit=5)

    assert len(result["suggestions"]) == 2
    query = parse_qs(urlparse(responses.calls[0].request.url).query)
    assert query["q"] == ["wir"]
    assert query["limit"] == ["5"]
