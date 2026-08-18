package io.github.adikeshri.tachyon;

/** The request never reached the server, or the server never replied at all. */
public class TachyonConnectionException extends RuntimeException {
    public TachyonConnectionException(String message) {
        super(message);
    }

    public TachyonConnectionException(String message, Throwable cause) {
        super(message, cause);
    }
}
