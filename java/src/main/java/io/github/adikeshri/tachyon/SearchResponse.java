package io.github.adikeshri.tachyon;

import java.util.List;
import java.util.Map;

public record SearchResponse(
    int found,
    /**
     * False once block-max WAND pruning has skipped part of a term's
     * postings for this query, at which point found (and facet counts)
     * become a lower bound rather than an exact count.
     */
    boolean foundIsExact,
    int searchTimeMs,
    List<SearchHit> hits,
    Map<String, Map<String, Integer>> facets
) {
}
