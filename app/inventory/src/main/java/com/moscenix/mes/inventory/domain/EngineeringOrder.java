package com.moscenix.mes.inventory.domain;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Convert;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;
import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "engineering_orders")
public class EngineeringOrder extends BaseEntity {
    @Column(name = "leader_user_id", nullable = false)
    private Long leaderUserId;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "process_id", nullable = false)
    private ProcessEntity process;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "item_id", nullable = false)
    private Item item;

    @Column(name = "name", nullable = false, length = 100)
    private String name = "";

    @Column(name = "expected_quantity", nullable = false)
    private Long expectedQuantity = 0L;

    @Column(name = "qualified_quantity", nullable = false)
    private Long qualifiedQuantity = 0L;

    @Column(name = "unqualified_quantity", nullable = false)
    private Long unqualifiedQuantity = 0L;

    @Column(name = "produced_quantity", nullable = false)
    private Long producedQuantity = 0L;

    @Convert(converter = DraftStatusConverter.class)
    @Column(name = "status", nullable = false)
    private DraftStatus status = DraftStatus.DRAFT;

    @Column(name = "description", nullable = false, length = 255)
    private String description = "";

    @OneToMany(mappedBy = "engineeringOrder", cascade = CascadeType.ALL)
    private List<ItemUnit> itemUnits = new ArrayList<>();

    public Long getLeaderUserId() {
        return leaderUserId;
    }

    public void setLeaderUserId(Long leaderUserId) {
        this.leaderUserId = leaderUserId;
    }

    public ProcessEntity getProcess() {
        return process;
    }

    public void setProcess(ProcessEntity process) {
        this.process = process;
    }

    public Item getItem() {
        return item;
    }

    public void setItem(Item item) {
        this.item = item;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Long getExpectedQuantity() {
        return expectedQuantity;
    }

    public void setExpectedQuantity(Long expectedQuantity) {
        this.expectedQuantity = expectedQuantity;
    }

    public Long getQualifiedQuantity() {
        return qualifiedQuantity;
    }

    public void setQualifiedQuantity(Long qualifiedQuantity) {
        this.qualifiedQuantity = qualifiedQuantity;
    }

    public Long getUnqualifiedQuantity() {
        return unqualifiedQuantity;
    }

    public void setUnqualifiedQuantity(Long unqualifiedQuantity) {
        this.unqualifiedQuantity = unqualifiedQuantity;
    }

    public Long getProducedQuantity() {
        return producedQuantity;
    }

    public void setProducedQuantity(Long producedQuantity) {
        this.producedQuantity = producedQuantity;
    }

    public DraftStatus getStatus() {
        return status;
    }

    public void setStatus(DraftStatus status) {
        this.status = status;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public List<ItemUnit> getItemUnits() {
        return itemUnits;
    }
}
