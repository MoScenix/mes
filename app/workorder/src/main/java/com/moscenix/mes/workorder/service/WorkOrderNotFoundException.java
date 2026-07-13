package com.moscenix.mes.workorder.service;

public class WorkOrderNotFoundException extends WorkOrderException {
    public WorkOrderNotFoundException(Long id) {
        super("work order not found: " + id, 40400);
    }
}
