package io.github.adikeshri.tachyon;

/** 5xx — corrupt_data, io_error, internal_error. */
public class TachyonServerException extends TachyonException {
    public TachyonServerException(String message, String code, int statusCode) {
        super(message, code, statusCode);
    }
}
