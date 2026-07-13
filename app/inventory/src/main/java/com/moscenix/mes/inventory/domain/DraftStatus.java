package com.moscenix.mes.inventory.domain;

public enum DraftStatus implements CodedEnum {
    UNKNOWN(0),
    DRAFT(1),
    SUBMITTED(2),
    DONE(3);

    private final int code;

    DraftStatus(int code) {
        this.code = code;
    }

    @Override
    public int getCode() {
        return code;
    }

    public static DraftStatus fromCode(int code) {
        for (DraftStatus status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return UNKNOWN;
    }
}
