package io.github.adikeshri.tachyon;

import java.util.List;

public record SuggestResponse(List<Suggestion> suggestions, int searchTimeMs) {
}
