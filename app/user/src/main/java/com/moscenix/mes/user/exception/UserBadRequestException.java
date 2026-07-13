package com.moscenix.mes.user.exception;

public class UserBadRequestException extends UserException {
    public UserBadRequestException(String message) {
        super(40000, message);
    }
}
