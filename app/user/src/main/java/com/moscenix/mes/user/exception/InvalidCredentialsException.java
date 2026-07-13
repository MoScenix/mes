package com.moscenix.mes.user.exception;

public class InvalidCredentialsException extends UserException {
    public InvalidCredentialsException() {
        super(40100, "invalid user account or password");
    }
}
