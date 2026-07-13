package com.moscenix.mes.inventory.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class StockStatusConverter implements AttributeConverter<StockStatus, Integer> {
    @Override
    public Integer convertToDatabaseColumn(StockStatus attribute) {
        return attribute == null ? StockStatus.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public StockStatus convertToEntityAttribute(Integer dbData) {
        return StockStatus.fromCode(dbData == null ? 0 : dbData);
    }
}
