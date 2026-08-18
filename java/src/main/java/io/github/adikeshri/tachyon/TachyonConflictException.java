package io.github.adikeshri.tachyon;

/** 409 — collection_exists. */
public class TachyonConflictException extends TachyonException {
    public TachyonConflictException(String message, String code, int statusCode) {
        super(message, code, statusCode);
    }
}
