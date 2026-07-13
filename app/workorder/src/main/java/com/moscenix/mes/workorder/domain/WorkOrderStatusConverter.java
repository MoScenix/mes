package com.moscenix.mes.workorder.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter(autoApply = false)
public class WorkOrderStatusConverter implements AttributeConverter<WorkOrderStatus, Integer> {
    @Override
    public Integer convertToDatabaseColumn(WorkOrderStatus attribute) {
        return attribute == null ? WorkOrderStatus.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public WorkOrderStatus convertToEntityAttribute(Integer dbData) {
        return WorkOrderStatus.fromCode(dbData);
    }
}
