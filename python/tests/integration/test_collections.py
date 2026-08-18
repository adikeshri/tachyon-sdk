import pytest

from .support import unique_name
from tachyon_sdk import TachyonConflictError, TachyonError, TachyonNotFoundError, TachyonRequestError

pytestmark = pytest.mark.integration


def test_create_fills_in_field_and_collection_defaults(client):
    name = unique_name("coll")
    try:
        created = client.collections.create(
            {
                "name": name,
                "fields": [
                    {"name": "title", "type": "text"},
                    {"name": "brand", "type": "keyword", "facet": True},
                    {"name": "price", "type": "int", "filter": True, "sort": True},
                ],
            }
        )
        assert created["name"] == name
        assert created["num_documents"] == 0
        assert created["num_segments"] == 0
        assert len(created["fields"]) == 3

        price = next(f for f in created["fields"] if f["name"] == "price")
        assert price["filter"] is True
        assert price["sort"] is True
        assert price["optional"] is True
        assert price["index"] is True

        brand = next(f for f in created["fields"] if f["name"] == "brand")
        assert brand["facet"] is True
    finally:
        try:
            client.collections.delete(name)
        except Exception:
            pass


def test_create_rejects_duplicate_name(client):
    name = unique_name("coll")
    client.collections.create({"name": name, "fields": [{"name": "title", "type": "text"}]})
    try:
        with pytest.raises(TachyonConflictError) as excinfo:
            client.collections.create({"name": name, "fields": [{"name": "title", "type": "text"}]})
        assert excinfo.value.code == "collection_exists"
        assert excinfo.value.status == 409
    finally:
        client.collections.delete(name)


def test_list_includes_a_newly_created_collection(client):
    name = unique_name("coll")
    client.collections.create({"name": name, "fields": [{"name": "title", "type": "text"}]})
    try:
        names = [c["name"] for c in client.collections.list()]
        assert name in names
    finally:
        client.collections.delete(name)


def test_retrieve_returns_the_matching_collection(client):
    name = unique_name("coll")
    client.collections.create({"name": name, "fields": [{"name": "title", "type": "text"}]})
    try:
        retrieved = client.collections.retrieve(name)
        assert retrieved["name"] == name
    finally:
        client.collections.delete(name)


def test_retrieve_unknown_collection_404s(client):
    with pytest.raises(TachyonNotFoundError) as excinfo:
        client.collections.retrieve(unique_name("missing"))
    assert excinfo.value.code == "collection_not_found"
    assert excinfo.value.status == 404


def test_delete_removes_it_from_retrieve_and_list(client):
    name = unique_name("coll")
    client.collections.create({"name": name, "fields": [{"name": "title", "type": "text"}]})

    client.collections.delete(name)

    with pytest.raises(TachyonNotFoundError):
        client.collections.retrieve(name)
    names = [c["name"] for c in client.collections.list()]
    assert name not in names


def test_create_rejects_a_schema_with_no_text_field(client):
    with pytest.raises(TachyonRequestError) as excinfo:
        client.collections.create({"name": unique_name("coll"), "fields": [{"name": "price", "type": "int"}]})
    assert excinfo.value.code == "invalid_schema"
    assert excinfo.value.status == 400


def test_create_rejects_unrecognized_field_type_at_json_layer(client):
    # An unknown `type` enum value fails JSON extraction before the request reaches
    # the handler, so the server never wraps it in the documented {error:{code,...}}
    # shape -- it comes back as plain text with a 422. The SDK still surfaces
    # status + message, just via the base TachyonError rather than a subclass.
    with pytest.raises(TachyonError) as excinfo:
        client.collections.create({"name": unique_name("coll"), "fields": [{"name": "x", "type": "not-a-real-type"}]})
    assert excinfo.value.status == 422


def test_collection_level_typo_tolerance_and_default_sorting_field(client):
    name = unique_name("coll")
    try:
        created = client.collections.create(
            {
                "name": name,
                "fields": [
                    {"name": "title", "type": "text"},
                    {"name": "popularity", "type": "int", "sort": True},
                ],
                "typo_tolerance": {"enabled": False},
                "default_sorting_field": "popularity",
            }
        )
        assert created["typo_tolerance"]["enabled"] is False
        assert created["default_sorting_field"] == "popularity"
    finally:
        client.collections.delete(name)
