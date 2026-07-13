package com.moscenix.mes.user.exception;

public class UserConflictException extends UserException {
    public UserConflictException(String message) {
        super(40900, message);
    }
}
