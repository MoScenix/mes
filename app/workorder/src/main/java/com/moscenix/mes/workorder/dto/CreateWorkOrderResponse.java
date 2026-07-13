package com.moscenix.mes.workorder.dto;

public class CreateWorkOrderResponse {
    private Long id;

    public CreateWorkOrderResponse(Long id) {
        this.id = id;
    }

    public Long getId() {
        return id;
    }
}
