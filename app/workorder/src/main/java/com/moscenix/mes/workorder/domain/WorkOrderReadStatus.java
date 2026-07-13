package com.moscenix.mes.workorder.domain;

public enum WorkOrderReadStatus {
    UNKNOWN(0),
    UNREAD(1),
    READ(2);

    private final int code;

    WorkOrderReadStatus(int code) {
        this.code = code;
    }

    public int getCode() {
        return code;
    }

    public static WorkOrderReadStatus fromJson(Object value) {
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

    public static WorkOrderReadStatus fromCode(Integer code) {
        if (code == null) {
            return UNKNOWN;
        }
        for (WorkOrderReadStatus status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return UNKNOWN;
    }
}
