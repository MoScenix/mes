package com.moscenix.mes.inventory.domain;

public enum StockStatus implements CodedEnum {
    UNKNOWN(0),
    IN_STOCK(1),
    RESERVED(2),
    OUT_STOCK(3);

    private final int code;

    StockStatus(int code) {
        this.code = code;
    }

    @Override
    public int getCode() {
        return code;
    }

    public static StockStatus fromCode(int code) {
        for (StockStatus status : values()) {
            if (status.code == code) {
                return status;
            }
        }
        return UNKNOWN;
    }
}
