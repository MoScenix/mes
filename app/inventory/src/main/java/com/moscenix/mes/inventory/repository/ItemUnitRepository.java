package com.moscenix.mes.inventory.repository;

import com.moscenix.mes.inventory.domain.ItemUnit;
import com.moscenix.mes.inventory.domain.QualityStatus;
import com.moscenix.mes.inventory.domain.StockStatus;
import jakarta.persistence.LockModeType;
import java.util.Collection;
import java.util.List;
import java.util.Optional;
import org.springframework.data.jpa.repository.EntityGraph;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface ItemUnitRepository extends JpaRepository<ItemUnit, Long>, JpaSpecificationExecutor<ItemUnit> {
    @EntityGraph(attributePaths = {"item", "engineeringOrder"})
    Optional<ItemUnit> findByIdAndDeletedAtIsNull(Long id);

    @EntityGraph(attributePaths = {"item", "engineeringOrder"})
    List<ItemUnit> findByIdInAndDeletedAtIsNull(Collection<Long> ids);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("select u from ItemUnit u join fetch u.item left join fetch u.engineeringOrder where u.id in :ids and u.deletedAt is null")
    List<ItemUnit> findForUpdateByIds(@Param("ids") Collection<Long> ids);

    long countByItemIdAndDeletedAtIsNull(Long itemId);

    long countByItemIdAndStockStatusAndDeletedAtIsNull(Long itemId, StockStatus stockStatus);

    long countByItemIdAndQualityStatusAndDeletedAtIsNull(Long itemId, QualityStatus qualityStatus);

    long countByItemIdAndStockStatusAndQualityStatusAndDeletedAtIsNull(
            Long itemId,
            StockStatus stockStatus,
            QualityStatus qualityStatus);

    long countByEngineeringOrderIdAndDeletedAtIsNull(Long engineeringOrderId);

    long countByEngineeringOrderIdAndQualityStatusAndDeletedAtIsNull(
            Long engineeringOrderId,
            QualityStatus qualityStatus);
}
