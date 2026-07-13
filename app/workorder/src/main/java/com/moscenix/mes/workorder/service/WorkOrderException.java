package com.moscenix.mes.workorder.service;

public class WorkOrderException extends RuntimeException {
    private final int code;

    public WorkOrderException(String message, int code) {
        super(message);
        this.code = code;
    }

    public int getCode() {
        return code;
    }
}
