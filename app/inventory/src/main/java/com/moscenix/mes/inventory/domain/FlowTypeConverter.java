package com.moscenix.mes.inventory.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class FlowTypeConverter implements AttributeConverter<FlowType, Integer> {
    @Override
    public Integer convertToDatabaseColumn(FlowType attribute) {
        return attribute == null ? FlowType.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public FlowType convertToEntityAttribute(Integer dbData) {
        return FlowType.fromCode(dbData == null ? 0 : dbData);
    }
}
