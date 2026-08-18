package io.github.adikeshri.tachyon;

import java.util.List;

/** {@code /collections/{name}/documents} — index, fetch, and delete documents. */
public final class DocumentsResource {
    private final HttpTransport transport;
    private final String collectionName;

    DocumentsResource(HttpTransport transport, String collectionName) {
        this.transport = transport;
        this.collectionName = collectionName;
    }

    /**
     * {@code POST /collections/{name}/documents}. Upserts one or more
     * documents by id. Individual documents can fail without failing their
     * neighbours — check {@link DocumentsIndexResponse#numFailed} and
     * {@link DocumentsIndexResponse#results}.
     */
    public DocumentsIndexResponse index(Document... documents) {
        return index(List.of(documents));
    }

    public DocumentsIndexResponse index(List<Document> documents) {
        String path = "/collections/" + UrlEncoding.pathSegment(collectionName) + "/documents";
        return transport.request("POST", path, null, documents, DocumentsIndexResponse.class);
    }

    /** {@code GET /collections/{name}/documents/{id}}. */
    public Document retrieve(String id) {
        String path = "/collections/" + UrlEncoding.pathSegment(collectionName) + "/documents/" + UrlEncoding.pathSegment(id);
        return transport.request("GET", path, null, null, Document.class);
    }

    /** {@code DELETE /collections/{name}/documents/{id}}. */
    public void delete(String id) {
        String path = "/collections/" + UrlEncoding.pathSegment(collectionName) + "/documents/" + UrlEncoding.pathSegment(id);
        transport.requestVoid("DELETE", path);
    }
}
