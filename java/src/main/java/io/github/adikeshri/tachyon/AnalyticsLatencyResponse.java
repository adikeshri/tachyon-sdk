package io.github.adikeshri.tachyon;

public record AnalyticsLatencyResponse(
    int count,
    double meanMs,
    double p50Ms,
    double p95Ms,
    double p99Ms,
    double maxMs,
    int totalSearches,
    long uptimeSeconds,
    double queriesPerSecond
) {
}
