package com.moscenix.mes.workorder.dto;

public class GetWorkOrderResponse {
    private WorkOrderInfo workOrder;

    public GetWorkOrderResponse(WorkOrderInfo workOrder) {
        this.workOrder = workOrder;
    }

    public WorkOrderInfo getWorkOrder() {
        return workOrder;
    }
}
