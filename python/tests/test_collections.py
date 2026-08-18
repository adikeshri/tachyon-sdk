import json
from urllib.parse import urlparse

import responses

from .conftest import BASE_URL


@responses.activate
def test_create_collection(client):
    responses.add(
        responses.POST,
        f"{BASE_URL}/collections",
        json={
            "name": "products",
            "fields": [{"name": "title", "type": "text"}],
            "num_documents": 0,
            "num_segments": 0,
        },
        status=201,
    )

    info = client.collections.create({"name": "products", "fields": [{"name": "title", "type": "text"}]})

    assert info["name"] == "products"
    assert info["num_documents"] == 0
    assert len(responses.calls) == 1
    request = responses.calls[0].request
    assert request.method == "POST"
    assert request.headers["X-TACHYON-API-KEY"] == "admin-key"
    assert json.loads(request.body) == {"name": "products", "fields": [{"name": "title", "type": "text"}]}


@responses.activate
def test_list_collections(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections",
        json=[{"name": "products", "fields": [], "num_documents": 5, "num_segments": 1}],
        status=200,
    )

    collections = client.collections.list()

    assert len(collections) == 1
    assert collections[0]["name"] == "products"


@responses.activate
def test_retrieve_collection(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/products",
        json={"name": "products", "fields": [], "num_documents": 5, "num_segments": 1},
        status=200,
    )

    info = client.collections.retrieve("products")

    assert info["num_documents"] == 5


@responses.activate
def test_retrieve_url_encodes_collection_name(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/collections/my%20products",
        json={"name": "my products", "fields": [], "num_documents": 0, "num_segments": 0},
        status=200,
    )

    client.collections.retrieve("my products")

    request = responses.calls[0].request
    assert urlparse(request.url).path == "/collections/my%20products"


@responses.activate
def test_delete_collection_returns_none_on_204(client):
    responses.add(responses.DELETE, f"{BASE_URL}/collections/products", status=204)

    result = client.collections.delete("products")

    assert result is None
    request = responses.calls[0].request
    assert request.method == "DELETE"
