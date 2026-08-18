using System.Text.Json;
using System.Text.Json.Serialization;

namespace Tachyon.Sdk;

internal static class JsonOptions
{
    /// <summary>
    /// Shared serializer options: snake_case to match Tachyon's JSON API
    /// (verified to round-trip every DTO property name exactly, including
    /// digit-adjacent ones like P50Ms -&gt; p50_ms), nulls omitted so unset
    /// optional fields don't override the server's own defaults, and enums
    /// as their lowercase string form.
    /// </summary>
    public static readonly JsonSerializerOptions Default = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
        Converters = { new JsonStringEnumConverter(JsonNamingPolicy.SnakeCaseLower) },
    };
}
