package io.github.adikeshri.tachyon;

/**
 * Base exception for every non-2xx response from Tachyon's HTTP API. {@code
 * code} is the stable machine-readable string from the response body (see
 * https://github.com/adikeshri/tachyon/blob/main/docs/api.md#errors);
 * {@code statusCode} is the HTTP status code.
 */
public class TachyonException extends RuntimeException {
    private final String code;
    private final int statusCode;

    public TachyonException(String message, String code, int statusCode) {
        super(message);
        this.code = code;
        this.statusCode = statusCode;
    }

    public String getCode() {
        return code;
    }

    public int getStatusCode() {
        return statusCode;
    }

    @Override
    public String toString() {
        return getMessage() + " (code=" + code + ", status=" + statusCode + ")";
    }
}
