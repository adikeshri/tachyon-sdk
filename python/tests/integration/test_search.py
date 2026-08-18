import pytest

from .support import admin_client, unique_name

pytestmark = pytest.mark.integration


@pytest.fixture(scope="module")
def dataset():
    client = admin_client()
    name = unique_name("search")
    client.collections.create(
        {
            "name": name,
            "fields": [
                {"name": "title", "type": "text"},
                {"name": "description", "type": "text"},
                {"name": "brand", "type": "keyword", "facet": True},
                {"name": "price", "type": "int", "filter": True, "sort": True},
                {"name": "in_stock", "type": "bool", "filter": True},
            ],
        }
    )
    coll = client.collection(name)
    coll.documents.index(
        [
            {
                "id": "1",
                "title": "Wireless Mouse",
                "description": "A great wireless mouse for everyday use",
                "brand": "Logitech",
                "price": 2999,
                "in_stock": True,
            },
            {
                "id": "2",
                "title": "Mechanical Keyboard",
                "description": "Clicky keys for typing enthusiasts",
                "brand": "Razer",
                "price": 8999,
                "in_stock": False,
            },
            {
                "id": "3",
                "title": "Wireless Keyboard",
                "description": "Silent wireless keyboard",
                "brand": "Logitech",
                "price": 5999,
                "in_stock": True,
            },
            {
                "id": "4",
                "title": "Gaming Mouse",
                "description": "Wired precision gaming mouse",
                "brand": "Razer",
                "price": 4999,
                "in_stock": True,
            },
            {
                "id": "5",
                "title": "USB Cable",
                "description": "A basic wired cable",
                "brand": "Anker",
                "price": 999,
                "in_stock": True,
            },
        ]
    )
    try:
        yield coll
    finally:
        client.collections.delete(name)


def test_matches_title_and_description_by_default_with_exact_count(dataset):
    results = dataset.search(q="wireless")
    assert results["found"] == 2
    assert results["found_is_exact"] is True
    assert sorted(h["document"]["id"] for h in results["hits"]) == ["1", "3"]


def test_query_by_restricts_matching_fields(dataset):
    only_description = dataset.search(q="Clicky", query_by="description")
    assert only_description["found"] == 1
    assert only_description["hits"][0]["document"]["id"] == "2"

    only_title = dataset.search(q="Clicky", query_by="title")
    assert only_title["found"] == 0


def test_query_by_accepts_a_list_and_joins_with_commas(dataset):
    results = dataset.search(q="Clicky", query_by=["title", "description"])
    assert results["found"] == 1


class TestFilters:
    def test_equality(self, dataset):
        results = dataset.search(filter="brand:=Logitech")
        assert results["found"] == 2

    def test_numeric_comparison(self, dataset):
        results = dataset.search(filter="price:<5000")
        assert results["found"] == 3  # 999, 2999, 4999

    def test_inclusive_range(self, dataset):
        results = dataset.search(filter="price:[1000..5000]")
        ids = sorted(h["document"]["id"] for h in results["hits"])
        assert ids == ["1", "4"]  # 2999 and 4999; 999 and 5999 fall outside

    def test_set_membership(self, dataset):
        results = dataset.search(filter="brand:=[Logitech,Razer]")
        assert results["found"] == 4

    def test_boolean_equality(self, dataset):
        results = dataset.search(filter="in_stock:=true")
        assert results["found"] == 4

    def test_and_or_with_grouping(self, dataset):
        results = dataset.search(filter="(brand:=Logitech || brand:=Razer) && price:<5000")
        ids = sorted(h["document"]["id"] for h in results["hits"])
        assert ids == ["1", "4"]

    def test_negation_only_matches_documents_that_have_the_field(self, dataset):
        results = dataset.search(filter="brand:!=Razer")
        ids = sorted(h["document"]["id"] for h in results["hits"])
        assert ids == ["1", "3", "5"]

    def test_negation_excludes_documents_missing_the_field_entirely(self, dataset):
        client = admin_client()
        name = unique_name("coll")
        client.collections.create(
            {
                "name": name,
                "fields": [{"name": "title", "type": "text"}, {"name": "brand", "type": "keyword", "filter": True}],
            }
        )
        try:
            scoped = client.collection(name)
            scoped.documents.index(
                [
                    {"id": "a", "title": "has brand", "brand": "Razer"},
                    {"id": "b", "title": "no brand at all"},
                ]
            )
            results = scoped.search(filter="brand:!=Razer")
            # 'b' has no brand field at all, so "not Razer" must not match it.
            assert results["found"] == 0
        finally:
            client.collections.delete(name)


class TestSorting:
    def test_ascending(self, dataset):
        results = dataset.search(q="", sort="price:asc", limit=10)
        prices = [h["document"]["price"] for h in results["hits"]]
        assert prices == sorted(prices)

    def test_descending(self, dataset):
        results = dataset.search(q="", sort="price:desc", limit=10)
        prices = [h["document"]["price"] for h in results["hits"]]
        assert prices == sorted(prices, reverse=True)


def test_pagination_moves_through_results_without_overlap(dataset):
    page1 = dataset.search(q="", sort="price:asc", limit=2, offset=0)
    page2 = dataset.search(q="", sort="price:asc", limit=2, offset=2)

    assert page1["found"] == 5
    assert page2["found"] == 5
    assert len(page1["hits"]) == 2
    assert len(page2["hits"]) == 2
    page1_ids = {h["document"]["id"] for h in page1["hits"]}
    page2_ids = {h["document"]["id"] for h in page2["hits"]}
    assert page1_ids.isdisjoint(page2_ids)


class TestPrefixMatching:
    def test_prefix_expands_the_final_token_by_default(self, dataset):
        results = dataset.search(q="wir")
        assert results["found"] >= 2

    def test_requires_full_token_when_prefix_disabled(self, dataset):
        results = dataset.search(q="wir", prefix=False)
        assert results["found"] == 0


class TestTypoTolerance:
    def test_corrects_a_typo_by_default(self, dataset):
        results = dataset.search(q="wirelss")
        assert results["found"] >= 1

    def test_finds_nothing_when_disabled(self, dataset):
        results = dataset.search(q="wirelss", typo_tolerance=False)
        assert results["found"] == 0


class TestMatchMode:
    def test_all_requires_every_token(self, dataset):
        results = dataset.search(q="wireless zzznonexistentterm")
        assert results["found"] == 0

    def test_any_requires_one_token(self, dataset):
        results = dataset.search(q="wireless zzznonexistentterm", match_mode="any")
        assert results["found"] >= 2


def test_phrase_queries_require_adjacency_within_a_field(dataset):
    results = dataset.search(q='"wireless mouse"')
    assert results["found"] == 1
    assert results["hits"][0]["document"]["id"] == "1"


def test_facets_count_every_matching_document_not_just_the_page(dataset):
    results = dataset.search(q="", facet="brand", limit=1)
    assert len(results["hits"]) == 1  # page size respected
    assert results["facets"]["brand"] == {"Logitech": 2, "Razer": 2, "Anker": 1}
