package io.github.adikeshri.tachyon;

final class TachyonExceptionFactory {
    private TachyonExceptionFactory() {
    }

    static TachyonException fromResponse(int status, String message, String code) {
        return switch (status) {
            case 400 -> new TachyonRequestException(message, code, status);
            case 401 -> new TachyonAuthenticationException(message, code, status);
            case 403 -> new TachyonAuthorizationException(message, code, status);
            case 404 -> new TachyonNotFoundException(message, code, status);
            case 409 -> new TachyonConflictException(message, code, status);
            default -> status >= 500
                ? new TachyonServerException(message, code, status)
                : new TachyonException(message, code, status);
        };
    }
}
