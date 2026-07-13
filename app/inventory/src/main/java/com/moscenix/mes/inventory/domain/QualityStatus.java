package com.moscenix.mes.inventory.domain;

public enum QualityStatus implements CodedEnum {
    UNKNOWN(0),
    PENDING(1),
    QUALIFIED(2),
    UNQUALIFIED(3);

    private final int code;

    QualityStatus(int code) {
        this.code = code;
    }

    @Override
    public int getCode() {
        return code;
    }

    public static QualityStatus fromCode(int code) {
        for (QualityStatus status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return UNKNOWN;
    }
}
