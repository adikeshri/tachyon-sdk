package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.*;

class AnalyticsAndOpsTest {

    @Test
    void analytics_top_scopedToCollection() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {
              "queries": [{"query":"wireless mouse","collection":"products","count":3,"zero_result_count":0,"last_result_count":12,"avg_latency_ms":1.8,"last_seen":1786625778150}],
              "tracked_queries": 1,
              "dropped_queries": 0
            }
            """))) {
            TachyonClient client = TestClients.forServer(server);

            AnalyticsQueriesResponse result = client.analytics.top(AnalyticsQueryParams.builder().collection("products").limit(10).build());

            assertEquals("wireless mouse", result.queries().get(0).query());
            var request = server.requests.get(0);
            assertEquals("/analytics/top", request.path());
            assertEquals("products", request.query().get("collection"));
            assertEquals("10", request.query().get("limit"));
        }
    }

    @Test
    void analytics_zeroResults_hitsRightPath() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"queries":[],"tracked_queries":0,"dropped_queries":0}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            client.analytics.zeroResults();

            assertEquals("/analytics/zero-results", server.requests.get(0).path());
        }
    }

    @Test
    void analytics_latency_returnsPercentiles() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"count":20,"mean_ms":1.9,"p50_ms":2.0,"p95_ms":4.0,"p99_ms":4.0,"max_ms":3.4,"total_searches":20,"uptime_seconds":61,"queries_per_second":0.33}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            AnalyticsLatencyResponse result = client.analytics.latency();

            assertEquals(4.0, result.p95Ms());
        }
    }

    @Test
    void health_returnsStatus() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.json(200,
            """
            {"ok":true,"version":"0.1.0","uptime_seconds":61,"num_collections":1}
            """))) {
            TachyonClient client = TestClients.forServer(server);

            HealthResponse health = client.health();

            assertTrue(health.ok());
        }
    }

    @Test
    void metrics_returnsRawTextNotJson() throws IOException {
        try (TestServer server = new TestServer(req -> TestServer.FakeResponse.text(200,
            "# HELP tachyon_uptime_seconds\ntachyon_uptime_seconds 61\n"))) {
            TachyonClient client = TestClients.forServer(server);

            String metrics = client.metrics();

            assertTrue(metrics.contains("tachyon_uptime_seconds 61"));
        }
    }
}
