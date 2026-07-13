package com.moscenix.mes.inventory.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Convert;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToMany;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import java.util.LinkedHashSet;
import java.util.Set;

@Entity
@Table(name = "item_units")
public class ItemUnit extends BaseEntity {
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "item_id", nullable = false)
    private Item item;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "engineering_order_id")
    private EngineeringOrder engineeringOrder;

    @Convert(converter = StockStatusConverter.class)
    @Column(name = "stock_status", nullable = false)
    private StockStatus stockStatus = StockStatus.UNKNOWN;

    @Convert(converter = QualityStatusConverter.class)
    @Column(name = "quality_status", nullable = false)
    private QualityStatus qualityStatus = QualityStatus.UNKNOWN;

    @Column(name = "description", nullable = false, length = 255)
    private String description = "";

    @ManyToMany(mappedBy = "itemUnits")
    private Set<InventoryFlow> inventoryFlows = new LinkedHashSet<>();

    public Item getItem() {
        return item;
    }

    public void setItem(Item item) {
        this.item = item;
    }

    public EngineeringOrder getEngineeringOrder() {
        return engineeringOrder;
    }

    public void setEngineeringOrder(EngineeringOrder engineeringOrder) {
        this.engineeringOrder = engineeringOrder;
    }

    public StockStatus getStockStatus() {
        return stockStatus;
    }

    public void setStockStatus(StockStatus stockStatus) {
        this.stockStatus = stockStatus;
    }

    public QualityStatus getQualityStatus() {
        return qualityStatus;
    }

    public void setQualityStatus(QualityStatus qualityStatus) {
        this.qualityStatus = qualityStatus;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }
}
