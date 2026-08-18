package io.github.adikeshri.tachyon;

public record AnalyticsQuery(
    String query,
    String collection,
    int count,
    int zeroResultCount,
    int lastResultCount,
    double avgLatencyMs,
    long lastSeen
) {
}
