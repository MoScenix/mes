package com.moscenix.mes.user.exception;

public class UserNotFoundException extends UserException {
    public UserNotFoundException(Long id) {
        super(40400, "user not found: " + id);
    }
}
