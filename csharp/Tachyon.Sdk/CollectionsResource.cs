namespace Tachyon.Sdk;

/// <summary><c>/collections</c> — create, list, and remove collections.</summary>
public sealed class CollectionsResource
{
    private readonly HttpTransport _transport;

    internal CollectionsResource(HttpTransport transport) => _transport = transport;

    /// <summary><c>POST /collections</c>. Field types are immutable after creation.</summary>
    public Task<CollectionInfo> CreateAsync(CollectionSchema schema, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<CollectionInfo>(HttpMethod.Post, "/collections", body: schema, cancellationToken: cancellationToken);

    /// <summary><c>GET /collections</c>.</summary>
    public Task<IReadOnlyList<CollectionInfo>> ListAsync(CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<IReadOnlyList<CollectionInfo>>(HttpMethod.Get, "/collections", cancellationToken: cancellationToken);

    /// <summary><c>GET /collections/{name}</c>.</summary>
    public Task<CollectionInfo> RetrieveAsync(string name, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync<CollectionInfo>(HttpMethod.Get, $"/collections/{Uri.EscapeDataString(name)}", cancellationToken: cancellationToken);

    /// <summary><c>DELETE /collections/{name}</c>. Removes the collection and all its data.</summary>
    public Task DeleteAsync(string name, CancellationToken cancellationToken = default) =>
        _transport.RequestAsync(HttpMethod.Delete, $"/collections/{Uri.EscapeDataString(name)}", cancellationToken: cancellationToken);
}
