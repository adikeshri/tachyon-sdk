package io.github.adikeshri.tachyon;

import org.junit.jupiter.api.Test;

import static io.github.adikeshri.tachyon.TestSupport.*;
import static org.junit.jupiter.api.Assertions.*;

class OpsIT {

    @Test
    void health_doesNotRequireAnApiKey() {
        HealthResponse health = anonymousClient().health();
        assertTrue(health.ok());
        assertTrue(health.numCollections() >= 0);
    }

    @Test
    void metrics_exposesPrometheusText() {
        String metrics = adminClient().metrics();
        assertTrue(metrics.contains("tachyon_uptime_seconds"));
        assertTrue(metrics.contains("tachyon_collections"));
    }
}
