package com.moscenix.mes.inventory.domain;

import jakarta.persistence.AttributeConverter;
import jakarta.persistence.Converter;

@Converter
public class DraftStatusConverter implements AttributeConverter<DraftStatus, Integer> {
    @Override
    public Integer convertToDatabaseColumn(DraftStatus attribute) {
        return attribute == null ? DraftStatus.UNKNOWN.getCode() : attribute.getCode();
    }

    @Override
    public DraftStatus convertToEntityAttribute(Integer dbData) {
        return DraftStatus.fromCode(dbData == null ? 0 : dbData);
    }
}
