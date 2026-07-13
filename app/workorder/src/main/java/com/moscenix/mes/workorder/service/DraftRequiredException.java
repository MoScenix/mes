package com.moscenix.mes.workorder.service;

public class DraftRequiredException extends WorkOrderException {
    public DraftRequiredException() {
        super("work order is not a draft or does not exist", 40900);
    }
}
