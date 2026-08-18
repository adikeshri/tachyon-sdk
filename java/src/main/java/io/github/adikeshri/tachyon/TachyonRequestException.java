package io.github.adikeshri.tachyon;

/** 400 — invalid_schema, invalid_document, invalid_query, invalid_json. */
public class TachyonRequestException extends TachyonException {
    public TachyonRequestException(String message, String code, int statusCode) {
        super(message, code, statusCode);
    }
}
