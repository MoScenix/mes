package com.moscenix.mes.inventory.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class FlowStatusConverter implements AttributeConverter<FlowStatus, Integer> {
    @Override
    public Integer convertToDatabaseColumn(FlowStatus attribute) {
        return attribute == null ? FlowStatus.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public FlowStatus convertToEntityAttribute(Integer dbData) {
        return FlowStatus.fromCode(dbData == null ? 0 : dbData);
    }
}
