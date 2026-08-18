package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class AnalyticsIT {

    @Test
    void recordsSearchCountsAndZeroResultsScopedByCollection() {
        withCollection(adminClient(), CollectionSchema.builder(uniqueName("coll")).fields(FieldSchema.builder("title", FieldType.TEXT).build()).build(), collection -> {
            collection.documents().index(new Document().set("id", "1").set("title", "gizmo"));

            collection.search(SearchParams.builder().q("gizmo").build());
            collection.search(SearchParams.builder().q("gizmo").build());
            collection.search(SearchParams.builder().q("zzz-totally-absent-term").build());

            AnalyticsQueriesResponse top = adminClient().analytics.top(AnalyticsQueryParams.builder().collection(collection.name()).build());
            AnalyticsQuery gizmo = top.queries().stream().filter(q -> q.query().equals("gizmo")).findFirst().orElse(null);
            assertNotNull(gizmo);
            assertEquals(2, gizmo.count());
            assertEquals(collection.name(), gizmo.collection());

            AnalyticsQueriesResponse zeroResults = adminClient().analytics.zeroResults(AnalyticsQueryParams.builder().collection(collection.name()).build());
            assertTrue(zeroResults.queries().stream().anyMatch(q -> q.query().equals("zzz-totally-absent-term")));
        });
    }

    @Test
    void top_respectsLimit() {
        AnalyticsQueriesResponse top = adminClient().analytics.top(AnalyticsQueryParams.builder().limit(1).build());
        assertTrue(top.queries().size() <= 1);
    }

    @Test
    void latency_reportsPercentilesAcrossAllSearchesSoFar() {
        AnalyticsLatencyResponse latency = adminClient().analytics.latency();
        assertTrue(latency.totalSearches() > 0);
        assertTrue(latency.p95Ms() >= latency.p50Ms());
        assertTrue(latency.p99Ms() >= latency.p95Ms());
    }
}
