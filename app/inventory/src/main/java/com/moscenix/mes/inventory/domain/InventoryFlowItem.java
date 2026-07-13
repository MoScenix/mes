package com.moscenix.mes.inventory.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;

@Entity
@Table(name = "inventory_flow_items")
public class InventoryFlowItem extends BaseEntity {
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "inventory_flow_id", nullable = false)
    private InventoryFlow inventoryFlow;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "item_id", nullable = false)
    private Item item;

    @Column(name = "apply_quantity", nullable = false)
    private Long applyQuantity = 0L;

    @Column(name = "finished_quantity", nullable = false)
    private Long finishedQuantity = 0L;

    public InventoryFlow getInventoryFlow() {
        return inventoryFlow;
    }

    public void setInventoryFlow(InventoryFlow inventoryFlow) {
        this.inventoryFlow = inventoryFlow;
    }

    public Item getItem() {
        return item;
    }

    public void setItem(Item item) {
        this.item = item;
    }

    public Long getApplyQuantity() {
        return applyQuantity;
    }

    public void setApplyQuantity(Long applyQuantity) {
        this.applyQuantity = applyQuantity;
    }

    public Long getFinishedQuantity() {
        return finishedQuantity;
    }

    public void setFinishedQuantity(Long finishedQuantity) {
        this.finishedQuantity = finishedQuantity;
    }
}
