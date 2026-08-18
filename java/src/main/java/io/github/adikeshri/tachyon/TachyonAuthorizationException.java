package io.github.adikeshri.tachyon;

/** 403 — a search key attempted a write. */
public class TachyonAuthorizationException extends TachyonException {
    public TachyonAuthorizationException(String message, String code, int statusCode) {
        super(message, code, statusCode);
    }
}
