package com.moscenix.mes.workorder.service;

public class WorkOrderBadRequestException extends WorkOrderException {
    public WorkOrderBadRequestException(String message) {
        super(message, 40000);
    }
}
