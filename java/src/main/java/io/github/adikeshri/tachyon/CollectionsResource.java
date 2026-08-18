package io.github.adikeshri.tachyon;

import com.fasterxml.jackson.core.type.TypeReference;

import java.util.List;

/** {@code /collections} — create, list, and remove collections. */
public final class CollectionsResource {
    private final HttpTransport transport;

    CollectionsResource(HttpTransport transport) {
        this.transport = transport;
    }

    /** {@code POST /collections}. Field types are immutable after creation. */
    public CollectionInfo create(CollectionSchema schema) {
        return transport.request("POST", "/collections", null, schema, CollectionInfo.class);
    }

    /** {@code GET /collections}. */
    public List<CollectionInfo> list() {
        return transport.request("GET", "/collections", null, null, new TypeReference<List<CollectionInfo>>() {
        });
    }

    /** {@code GET /collections/{name}}. */
    public CollectionInfo retrieve(String name) {
        return transport.request("GET", "/collections/" + UrlEncoding.pathSegment(name), null, null, CollectionInfo.class);
    }

    /** {@code DELETE /collections/{name}}. Removes the collection and all its data. */
    public void delete(String name) {
        transport.requestVoid("DELETE", "/collections/" + UrlEncoding.pathSegment(name));
    }
}
