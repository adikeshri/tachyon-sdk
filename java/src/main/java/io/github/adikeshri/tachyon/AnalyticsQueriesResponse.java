package io.github.adikeshri.tachyon;

import java.util.List;

public record AnalyticsQueriesResponse(List<AnalyticsQuery> queries, int trackedQueries, int droppedQueries) {
}
