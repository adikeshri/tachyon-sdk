import pytest

from .support import collection

pytestmark = pytest.mark.integration


def test_records_search_counts_and_zero_results_scoped_by_collection(client):
    with collection(client, {"fields": [{"name": "title", "type": "text"}]}) as name:
        coll = client.collection(name)
        coll.documents.index({"id": "1", "title": "gizmo"})

        coll.search(q="gizmo")
        coll.search(q="gizmo")
        coll.search(q="zzz-totally-absent-term")

        top = client.analytics.top(collection=name)
        gizmo_query = next((q for q in top["queries"] if q["query"] == "gizmo"), None)
        assert gizmo_query is not None
        assert gizmo_query["count"] == 2
        assert gizmo_query["collection"] == name

        zero_results = client.analytics.zero_results(collection=name)
        assert any(q["query"] == "zzz-totally-absent-term" for q in zero_results["queries"])


def test_top_respects_limit(client):
    top = client.analytics.top(limit=1)
    assert len(top["queries"]) <= 1


def test_latency_reports_percentiles_across_all_searches_so_far(client):
    latency = client.analytics.latency()
    assert latency["total_searches"] > 0
    assert latency["p50_ms"] >= 0
    assert latency["p95_ms"] >= latency["p50_ms"]
    assert latency["p99_ms"] >= latency["p95_ms"]
