import pytest

from .support import admin_client, unique_name

pytestmark = pytest.mark.integration


@pytest.fixture(scope="module")
def dataset():
    client = admin_client()
    name = unique_name("suggest")
    client.collections.create({"name": name, "fields": [{"name": "title", "type": "text"}]})
    coll = client.collection(name)
    coll.documents.index(
        [
            {"id": "1", "title": "wireless mouse"},
            {"id": "2", "title": "wireless keyboard"},
            {"id": "3", "title": "wireless mouse"},
            {"id": "4", "title": "wired cable"},
        ]
    )
    try:
        yield coll
    finally:
        client.collections.delete(name)


def test_completes_prefix_with_live_document_counts(dataset):
    result = dataset.suggest("wir")
    texts = [s["text"] for s in result["suggestions"]]
    assert "wireless" in texts
    assert "wired" in texts

    wireless = next(s for s in result["suggestions"] if s["text"] == "wireless")
    assert wireless["count"] >= 2


def test_caps_suggestions_at_the_requested_limit(dataset):
    result = dataset.suggest("wir", limit=1)
    assert len(result["suggestions"]) <= 1
