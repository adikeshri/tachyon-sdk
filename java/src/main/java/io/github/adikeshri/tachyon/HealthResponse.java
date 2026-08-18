package io.github.adikeshri.tachyon;

public record HealthResponse(boolean ok, String version, long uptimeSeconds, int numCollections) {
}
