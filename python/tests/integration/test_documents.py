import pytest

from .support import collection
from tachyon_sdk import TachyonNotFoundError

pytestmark = pytest.mark.integration


def test_index_single_document_not_wrapped_in_a_list(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        docs = client.collection(name).documents
        result = docs.index({"id": "1", "title": "Hello World"})
        assert result["num_indexed"] == 1
        assert result["num_failed"] == 0
        assert result["results"] == [{"success": True, "id": "1"}]


def test_retrieve_document_by_id(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        docs = client.collection(name).documents
        docs.index({"id": "1", "title": "Hello World"})

        doc = docs.retrieve("1")
        assert doc["id"] == "1"
        assert doc["title"] == "Hello World"


def test_index_upserts_by_id(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        docs = client.collection(name).documents
        docs.index({"id": "1", "title": "First title"})
        docs.index({"id": "1", "title": "Second title"})

        doc = docs.retrieve("1")
        assert doc["title"] == "Second title"


def test_delete_document_then_retrieve_404s(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        docs = client.collection(name).documents
        docs.index({"id": "1", "title": "Hello World"})

        docs.delete("1")

        with pytest.raises(TachyonNotFoundError) as excinfo:
            docs.retrieve("1")
        assert excinfo.value.code == "document_not_found"
        assert excinfo.value.status == 404


def test_retrieve_and_delete_unknown_id_404s(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        docs = client.collection(name).documents
        with pytest.raises(TachyonNotFoundError):
            docs.retrieve("never-existed")
        with pytest.raises(TachyonNotFoundError):
            docs.delete("never-existed")


def test_batch_index_reports_per_document_success_and_failure(client):
    schema = {"fields": [{"name": "title", "type": "text"}, {"name": "price", "type": "int"}]}
    with collection(client, schema) as name:
        docs = client.collection(name).documents
        result = docs.index(
            [
                {"id": "1", "price": 100},
                {"id": "2", "price": "not-a-number"},
                {"id": "3", "price": 300},
            ]
        )

        assert result["num_indexed"] == 2
        assert result["num_failed"] == 1
        assert result["results"][0]["success"] is True
        assert result["results"][0]["id"] == "1"
        assert result["results"][1]["success"] is False
        assert result["results"][1]["code"] == "invalid_document"
        assert result["results"][2]["success"] is True
        assert result["results"][2]["id"] == "3"

        with pytest.raises(TachyonNotFoundError):
            docs.retrieve("2")


def test_undeclared_fields_are_stored_but_not_indexed(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        coll = client.collection(name)
        coll.documents.index({"id": "1", "title": "Hello", "undeclared_field": "surprise"})

        doc = coll.documents.retrieve("1")
        assert doc["undeclared_field"] == "surprise"

        results = coll.search(q="surprise")
        assert len(results["hits"]) == 0
