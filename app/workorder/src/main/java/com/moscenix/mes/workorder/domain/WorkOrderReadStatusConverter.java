package com.moscenix.mes.workorder.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter(autoApply = false)
public class WorkOrderReadStatusConverter implements AttributeConverter<WorkOrderReadStatus, Integer> {
    @Override
    public Integer convertToDatabaseColumn(WorkOrderReadStatus attribute) {
        return attribute == null ? WorkOrderReadStatus.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public WorkOrderReadStatus convertToEntityAttribute(Integer dbData) {
        return WorkOrderReadStatus.fromCode(dbData);
    }
}
