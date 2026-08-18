package io.github.adikeshri.tachyon;

import java.util.List;

public record DocumentsIndexResponse(int numIndexed, int numFailed, List<DocumentIndexResult> results) {
}
