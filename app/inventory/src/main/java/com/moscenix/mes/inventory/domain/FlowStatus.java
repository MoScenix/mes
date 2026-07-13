package com.moscenix.mes.inventory.domain;

public enum FlowStatus implements CodedEnum {
    UNKNOWN(0),
    DRAFT(1),
    SUBMITTED(2),
    APPROVED(3),
    REJECTED(4);

    private final int code;

    FlowStatus(int code) {
        this.code = code;
    }

    @Override
    public int getCode() {
        return code;
    }

    public boolean isKnown() {
        return this != UNKNOWN;
    }

    public static FlowStatus fromCode(int code) {
        for (FlowStatus status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return UNKNOWN;
    }
}
