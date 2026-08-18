package io.github.adikeshri.tachyon;

import java.util.List;

/**
 * A collection schema as returned by the server, with live counters. Java
 * records can't extend {@link CollectionSchema}, so this flattens the same
 * fields directly rather than composing it.
 */
public record CollectionInfo(
    String name,
    List<FieldSchema> fields,
    TypoToleranceConfig typoTolerance,
    String defaultSortingField,
    int numDocuments,
    int numSegments
) {
}
