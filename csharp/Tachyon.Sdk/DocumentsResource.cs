namespace Tachyon.Sdk;

/// <summary><c>/collections/{name}/documents</c> — index, fetch, and delete documents.</summary>
public sealed class DocumentsResource
{
    private readonly HttpTransport _transport;
    private readonly string _collectionName;

    internal DocumentsResource(HttpTransport transport, string collectionName)
    {
        _transport = transport;
        _collectionName = collectionName;
    }

    /// <summary>
    /// <c>POST /collections/{name}/documents</c>. Upserts one or more
    /// documents by id. Individual documents can fail without failing
    /// their neighbours — check <see cref="DocumentsIndexResponse.NumFailed"/>
    /// and <see cref="DocumentsIndexResponse.Results"/>.
    /// </summary>
    public Task<DocumentsIndexResponse> IndexAsync(params Document[] documents) =>
        IndexAsync(documents, CancellationToken.None);

    public Task<DocumentsIndexResponse> IndexAsync(IReadOnlyList<Document> documents, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<DocumentsIndexResponse>(
            HttpMethod.Post,
            $"/collections/{Uri.EscapeDataString(_collectionName)}/documents",
            body: documents,
            cancellationToken: cancellationToken);

    /// <summary><c>GET /collections/{name}/documents/{id}</c>.</summary>
    public Task<Document> RetrieveAsync(string id, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<Document>(
            HttpMethod.Get,
            $"/collections/{Uri.EscapeDataString(_collectionName)}/documents/{Uri.EscapeDataString(id)}",
            cancellationToken: cancellationToken);

    /// <summary><c>DELETE /collections/{name}/documents/{id}</c>.</summary>
    public Task DeleteAsync(string id, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync(
            HttpMethod.Delete,
            $"/collections/{Uri.EscapeDataString(_collectionName)}/documents/{Uri.EscapeDataString(id)}",
            cancellationToken: cancellationToken);
}
