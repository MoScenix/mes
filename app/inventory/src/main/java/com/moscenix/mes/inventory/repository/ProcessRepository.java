package com.moscenix.mes.inventory.repository;

import com.moscenix.mes.inventory.domain.ProcessEntity;
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
public interface ProcessRepository extends JpaRepository<ProcessEntity, Long>, JpaSpecificationExecutor<ProcessEntity> {
    @EntityGraph(attributePaths = {"item", "items", "items.consumeItem"})
    @Query("select p from ProcessEntity p where p.id = :id and p.deletedAt is null")
    Optional<ProcessEntity> findWithItemsByIdAndDeletedAtIsNull(@Param("id") Long id);

    @EntityGraph(attributePaths = {"item"})
    Optional<ProcessEntity> findByIdAndDeletedAtIsNull(Long id);

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("select p from ProcessEntity p where p.id = :id and p.deletedAt is null")
    Optional<ProcessEntity> findForUpdate(@Param("id") Long id);
}
