package com.moscenix.mes.inventory.repository;

import com.moscenix.mes.inventory.domain.InventoryFlowItem;
import java.util.List;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface InventoryFlowItemRepository extends JpaRepository<InventoryFlowItem, Long> {
    List<InventoryFlowItem> findByInventoryFlowIdAndDeletedAtIsNull(Long inventoryFlowId);
}
