package com.moscenix.mes.inventory.domain;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Convert;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.JoinTable;
import jakarta.persistence.ManyToMany;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;
import java.time.Instant;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;

@Entity
@Table(name = "inventory_flows")
public class InventoryFlow extends BaseEntity {
    @Column(name = "from_user_id", nullable = false)
    private Long fromUserId;

    @Column(name = "to_user_id", nullable = false)
    private Long toUserId;

    @Convert(converter = FlowTypeConverter.class)
    @Column(name = "flow_type", nullable = false)
    private FlowType flowType = FlowType.UNKNOWN;

    @Convert(converter = FlowStatusConverter.class)
    @Column(name = "flow_status", nullable = false)
    private FlowStatus flowStatus = FlowStatus.UNKNOWN;

    @Column(name = "name", nullable = false, length = 100)
    private String name = "";

    @Column(name = "description", nullable = false, length = 255)
    private String description = "";

    @Column(name = "approved_by", nullable = false)
    private Long approvedBy = 0L;

    @Column(name = "approved_at")
    private Instant approvedAt;

    @OneToMany(mappedBy = "inventoryFlow", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<InventoryFlowItem> items = new ArrayList<>();

    @ManyToMany(fetch = FetchType.LAZY)
    @JoinTable(
            name = "inventory_flow_item_units",
            joinColumns = @JoinColumn(name = "inventory_flow_id"),
            inverseJoinColumns = @JoinColumn(name = "item_unit_id"))
    private Set<ItemUnit> itemUnits = new LinkedHashSet<>();

    public Long getFromUserId() {
        return fromUserId;
    }

    public void setFromUserId(Long fromUserId) {
        this.fromUserId = fromUserId;
    }

    public Long getToUserId() {
        return toUserId;
    }

    public void setToUserId(Long toUserId) {
        this.toUserId = toUserId;
    }

    public FlowType getFlowType() {
        return flowType;
    }

    public void setFlowType(FlowType flowType) {
        this.flowType = flowType;
    }

    public FlowStatus getFlowStatus() {
        return flowStatus;
    }

    public void setFlowStatus(FlowStatus flowStatus) {
        this.flowStatus = flowStatus;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Long getApprovedBy() {
        return approvedBy;
    }

    public void setApprovedBy(Long approvedBy) {
        this.approvedBy = approvedBy;
    }

    public Instant getApprovedAt() {
        return approvedAt;
    }

    public void setApprovedAt(Instant approvedAt) {
        this.approvedAt = approvedAt;
    }

    public List<InventoryFlowItem> getItems() {
        return items;
    }

    public void replaceItems(List<InventoryFlowItem> newItems) {
        items.clear();
        newItems.forEach(this::addItem);
    }

    public void addItem(InventoryFlowItem item) {
        item.setInventoryFlow(this);
        items.add(item);
    }

    public Set<ItemUnit> getItemUnits() {
        return itemUnits;
    }

    public void replaceItemUnits(List<ItemUnit> units) {
        itemUnits.clear();
        itemUnits.addAll(units);
    }
}
