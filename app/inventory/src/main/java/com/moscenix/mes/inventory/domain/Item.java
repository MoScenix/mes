package com.moscenix.mes.inventory.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;

@Entity
@Table(name = "items")
public class Item extends BaseEntity {
    @Column(name = "name", nullable = false, length = 100)
    private String name;

    @Column(name = "unit", nullable = false, length = 20)
    private String unit;

    @Column(name = "description", nullable = false, length = 255)
    private String description = "";

    @Column(name = "total_count", nullable = false)
    private Long totalCount = 0L;

    @Column(name = "in_stock_count", nullable = false)
    private Long inStockCount = 0L;

    @Column(name = "reserved_count", nullable = false)
    private Long reservedCount = 0L;

    @Column(name = "out_stock_count", nullable = false)
    private Long outStockCount = 0L;

    @Column(name = "pending_count", nullable = false)
    private Long pendingCount = 0L;

    @Column(name = "qualified_count", nullable = false)
    private Long qualifiedCount = 0L;

    @Column(name = "unqualified_count", nullable = false)
    private Long unqualifiedCount = 0L;

    @Column(name = "available_count", nullable = false)
    private Long availableCount = 0L;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getUnit() {
        return unit;
    }

    public void setUnit(String unit) {
        this.unit = unit;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Long getTotalCount() {
        return totalCount;
    }

    public void setTotalCount(Long totalCount) {
        this.totalCount = totalCount;
    }

    public Long getInStockCount() {
        return inStockCount;
    }

    public void setInStockCount(Long inStockCount) {
        this.inStockCount = inStockCount;
    }

    public Long getReservedCount() {
        return reservedCount;
    }

    public void setReservedCount(Long reservedCount) {
        this.reservedCount = reservedCount;
    }

    public Long getOutStockCount() {
        return outStockCount;
    }

    public void setOutStockCount(Long outStockCount) {
        this.outStockCount = outStockCount;
    }

    public Long getPendingCount() {
        return pendingCount;
    }

    public void setPendingCount(Long pendingCount) {
        this.pendingCount = pendingCount;
    }

    public Long getQualifiedCount() {
        return qualifiedCount;
    }

    public void setQualifiedCount(Long qualifiedCount) {
        this.qualifiedCount = qualifiedCount;
    }

    public Long getUnqualifiedCount() {
        return unqualifiedCount;
    }

    public void setUnqualifiedCount(Long unqualifiedCount) {
        this.unqualifiedCount = unqualifiedCount;
    }

    public Long getAvailableCount() {
        return availableCount;
    }

    public void setAvailableCount(Long availableCount) {
        this.availableCount = availableCount;
    }
}
