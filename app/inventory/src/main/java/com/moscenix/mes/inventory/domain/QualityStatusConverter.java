package com.moscenix.mes.inventory.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class QualityStatusConverter implements AttributeConverter<QualityStatus, Integer> {
    @Override
    public Integer convertToDatabaseColumn(QualityStatus attribute) {
        return attribute == null ? QualityStatus.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public QualityStatus convertToEntityAttribute(Integer dbData) {
        return QualityStatus.fromCode(dbData == null ? 0 : dbData);
    }
}
