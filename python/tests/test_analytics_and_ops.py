from urllib.parse import parse_qs, urlparse

import responses

from .conftest import BASE_URL


@responses.activate
def test_analytics_top_scoped_to_collection(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/analytics/top",
        json={
            "queries": [
                {
                    "query": "wireless mouse",
                    "collection": "products",
                    "count": 3,
                    "zero_result_count": 0,
                    "last_result_count": 12,
                    "avg_latency_ms": 1.8,
                    "last_seen": 1786625778150,
                }
            ],
            "tracked_queries": 1,
            "dropped_queries": 0,
        },
        status=200,
    )

    result = client.analytics.top(collection="products", limit=10)

    assert result["queries"][0]["query"] == "wireless mouse"
    query = parse_qs(urlparse(responses.calls[0].request.url).query)
    assert query["collection"] == ["products"]
    assert query["limit"] == ["10"]


@responses.activate
def test_analytics_zero_results(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/analytics/zero-results",
        json={"queries": [], "tracked_queries": 0, "dropped_queries": 0},
        status=200,
    )

    client.analytics.zero_results()

    assert urlparse(responses.calls[0].request.url).path == "/analytics/zero-results"


@responses.activate
def test_analytics_latency(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/analytics/latency",
        json={
            "count": 20,
            "mean_ms": 1.9,
            "p50_ms": 2.0,
            "p95_ms": 4.0,
            "p99_ms": 4.0,
            "max_ms": 3.4,
            "total_searches": 20,
            "uptime_seconds": 61,
            "queries_per_second": 0.33,
        },
        status=200,
    )

    result = client.analytics.latency()

    assert result["p95_ms"] == 4.0


@responses.activate
def test_health(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/health",
        json={"ok": True, "version": "0.1.0", "uptime_seconds": 61, "num_collections": 1},
        status=200,
    )

    health = client.health()

    assert health["ok"] is True


@responses.activate
def test_metrics_returns_raw_prometheus_text_not_json(client):
    responses.add(
        responses.GET,
        f"{BASE_URL}/metrics",
        body="# HELP tachyon_uptime_seconds\ntachyon_uptime_seconds 61\n",
        status=200,
        content_type="text/plain",
    )

    metrics = client.metrics()

    assert "tachyon_uptime_seconds 61" in metrics
