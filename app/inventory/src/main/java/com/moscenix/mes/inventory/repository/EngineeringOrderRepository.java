package com.moscenix.mes.inventory.repository;

import com.moscenix.mes.inventory.domain.EngineeringOrder;
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
public interface EngineeringOrderRepository extends JpaRepository<EngineeringOrder, Long>, JpaSpecificationExecutor<EngineeringOrder> {
    @EntityGraph(attributePaths = {"item", "process", "process.item", "itemUnits"})
    @Query("select e from EngineeringOrder e where e.id = :id and e.deletedAt is null")
    Optional<EngineeringOrder> findWithDetailsByIdAndDeletedAtIsNull(@Param("id") Long id);

    @EntityGraph(attributePaths = {"item", "process", "process.item"})
    Optional<EngineeringOrder> findByIdAndDeletedAtIsNull(Long id);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @EntityGraph(attributePaths = {"item", "process"})
    @Query("select e from EngineeringOrder e where e.id = :id and e.deletedAt is null")
    Optional<EngineeringOrder> findForUpdate(@Param("id") Long id);
}
