package io.github.adikeshri.tachyon;

/** 401 — missing or wrong API key. */
public class TachyonAuthenticationException extends TachyonException {
    public TachyonAuthenticationException(String message, String code, int statusCode) {
        super(message, code, statusCode);
    }
}
