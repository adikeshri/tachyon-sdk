namespace Tachyon.Sdk;

/// <summary>
/// Base exception for every non-2xx response from Tachyon's HTTP API.
/// <see cref="Code"/> is the stable machine-readable string from the
/// response body (see
/// https://github.com/adikeshri/tachyon/blob/main/docs/api.md#errors);
/// <see cref="StatusCode"/> is the HTTP status code.
/// </summary>
public class TachyonException : Exception
{
    public string Code { get; }
    public int StatusCode { get; }

    public TachyonException(string message, string code, int statusCode)
        : base(message)
    {
        Code = code;
        StatusCode = statusCode;
    }

    public override string ToString() => $"{Message} (code={Code}, status={StatusCode})";
}

/// <summary>400 — invalid_schema, invalid_document, invalid_query, invalid_json.</summary>
public sealed class TachyonRequestException : TachyonException
{
    public TachyonRequestException(string message, string code, int statusCode) : base(message, code, statusCode)
    {
    }
}

/// <summary>401 — missing or wrong API key.</summary>
public sealed class TachyonAuthenticationException : TachyonException
{
    public TachyonAuthenticationException(string message, string code, int statusCode) : base(message, code, statusCode)
    {
    }
}

/// <summary>403 — a search key attempted a write.</summary>
public sealed class TachyonAuthorizationException : TachyonException
{
    public TachyonAuthorizationException(string message, string code, int statusCode) : base(message, code, statusCode)
    {
    }
}

/// <summary>404 — collection_not_found, document_not_found.</summary>
public sealed class TachyonNotFoundException : TachyonException
{
    public TachyonNotFoundException(string message, string code, int statusCode) : base(message, code, statusCode)
    {
    }
}

/// <summary>409 — collection_exists.</summary>
public sealed class TachyonConflictException : TachyonException
{
    public TachyonConflictException(string message, string code, int statusCode) : base(message, code, statusCode)
    {
    }
}

/// <summary>5xx — corrupt_data, io_error, internal_error.</summary>
public sealed class TachyonServerException : TachyonException
{
    public TachyonServerException(string message, string code, int statusCode) : base(message, code, statusCode)
    {
    }
}

internal static class TachyonExceptionFactory
{
    public static TachyonException FromResponse(int status, string message, string code) => status switch
    {
        400 => new TachyonRequestException(message, code, status),
        401 => new TachyonAuthenticationException(message, code, status),
        403 => new TachyonAuthorizationException(message, code, status),
        404 => new TachyonNotFoundException(message, code, status),
        409 => new TachyonConflictException(message, code, status),
        >= 500 => new TachyonServerException(message, code, status),
        _ => new TachyonException(message, code, status),
    };
}

/// <summary>The request never reached the server, or the server never replied at all.</summary>
public class TachyonConnectionException : Exception
{
    public TachyonConnectionException(string message, Exception? innerException = null)
        : base(message, innerException)
    {
    }
}

/// <summary>The request was aborted after exceeding its configured timeout.</summary>
public sealed class TachyonTimeoutException : TachyonConnectionException
{
    public TachyonTimeoutException(string message, Exception? innerException = null)
        : base(message, innerException)
    {
    }
}
