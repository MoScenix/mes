package com.moscenix.mes.inventory.domain;

public enum FlowType implements CodedEnum {
    UNKNOWN(0),
    IN(1),
    OUT(2);

    private final int code;

    FlowType(int code) {
        this.code = code;
    }

    @Override
    public int getCode() {
        return code;
    }

    public static FlowType fromCode(int code) {
        for (FlowType type : values()) {
            if (type.code == code) {
                return type;
            }
        }
        return UNKNOWN;
    }
}
