package io.github.adikeshri.tachyon;

/** 404 — collection_not_found, document_not_found. */
public class TachyonNotFoundException extends TachyonException {
    public TachyonNotFoundException(String message, String code, int statusCode) {
        super(message, code, statusCode);
    }
}
