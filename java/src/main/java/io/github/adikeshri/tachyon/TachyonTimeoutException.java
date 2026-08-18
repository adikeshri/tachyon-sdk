package io.github.adikeshri.tachyon;

/** The request was aborted after exceeding its configured timeout. */
public class TachyonTimeoutException extends TachyonConnectionException {
    public TachyonTimeoutException(String message) {
        super(message);
    }

    public TachyonTimeoutException(String message, Throwable cause) {
        super(message, cause);
    }
}
