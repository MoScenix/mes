package com.moscenix.mes.workorder.domain;

public enum WorkOrderStatus {
    UNKNOWN(0),
    DRAFT(1),
    SUBMITTED(2);

    private final int code;

    WorkOrderStatus(int code) {
        this.code = code;
    }

    public int getCode() {
        return code;
    }

    public static WorkOrderStatus fromJson(Object value) {
        if (value == null) {
            return UNKNOWN;
        }
        if (value instanceof Number number) {
            return fromCode(number.intValue());
        }
        String text = value.toString().trim();
        if (text.isEmpty()) {
            return UNKNOWN;
        }
        try {
            return fromCode(Integer.parseInt(text));
        } catch (NumberFormatException ignored) {
            return valueOf(text.toUpperCase());
        }
    }

    public static WorkOrderStatus fromCode(Integer code) {
        if (code == null) {
            return UNKNOWN;
        }
        for (WorkOrderStatus status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return UNKNOWN;
    }

    public boolean isKnown() {
        return this != UNKNOWN;
    }
}
