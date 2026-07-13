package com.moscenix.mes.user.exception;

public class UserException extends RuntimeException {
    private final int code;

    public UserException(int code, String message) {
        super(message);
        this.code = code;
    }

    public int getCode() {
        return code;
    }
}
