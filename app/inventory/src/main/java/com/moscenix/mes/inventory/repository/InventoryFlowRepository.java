package com.moscenix.mes.inventory.repository;

import com.moscenix.mes.inventory.domain.InventoryFlow;
import jakarta.persistence.LockModeType;
import java.util.Optional;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface InventoryFlowRepository extends JpaRepository<InventoryFlow, Long>, JpaSpecificationExecutor<InventoryFlow> {
    @EntityGraph(attributePaths = {"items", "items.item", "itemUnits", "itemUnits.item", "itemUnits.engineeringOrder"})
    @Query("select f from InventoryFlow f where f.id = :id and f.deletedAt is null")
    Optional<InventoryFlow> findWithDetailsByIdAndDeletedAtIsNull(@Param("id") Long id);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @EntityGraph(attributePaths = {"items", "items.item", "itemUnits", "itemUnits.item"})
    @Query("select f from InventoryFlow f where f.id = :id and f.deletedAt is null")
    Optional<InventoryFlow> findForUpdateWithDetails(@Param("id") Long id);
}
