package com.moscenix.mes.workorder.repository;

import com.moscenix.mes.workorder.domain.WorkOrder;
import com.moscenix.mes.workorder.domain.WorkOrderReadStatus;
import com.moscenix.mes.workorder.domain.WorkOrderStatus;
import java.time.Instant;
import java.util.Optional;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface WorkOrderRepository extends JpaRepository<WorkOrder, Long>, JpaSpecificationExecutor<WorkOrder> {
    Optional<WorkOrder> findByIdAndDeletedAtIsNull(Long id);

    @Modifying
    @Query("""
            update WorkOrder w
               set w.fromUserId = :fromUserId,
                   w.toUserId = :toUserId,
                   w.name = :name,
                   w.description = :description,
                   w.updatedAt = :updatedAt
             where w.id = :id
               and w.status = :draft
               and w.deletedAt is null
            """)
    int updateDraft(
            @Param("id") Long id,
            @Param("fromUserId") Long fromUserId,
            @Param("toUserId") Long toUserId,
            @Param("name") String name,
            @Param("description") String description,
            @Param("draft") WorkOrderStatus draft,
            @Param("updatedAt") Instant updatedAt);

    @Modifying
    @Query("""
            update WorkOrder w
               set w.deletedAt = :deletedAt,
                   w.updatedAt = :deletedAt
             where w.id = :id
               and w.status = :draft
               and w.deletedAt is null
            """)
    int deleteDraft(
            @Param("id") Long id,
            @Param("draft") WorkOrderStatus draft,
            @Param("deletedAt") Instant deletedAt);

    @Modifying
    @Query("""
            update WorkOrder w
               set w.status = :submitted,
                   w.updatedAt = :updatedAt
             where w.id = :id
               and w.status = :draft
               and w.deletedAt is null
            """)
    int submitDraft(
            @Param("id") Long id,
            @Param("draft") WorkOrderStatus draft,
            @Param("submitted") WorkOrderStatus submitted,
            @Param("updatedAt") Instant updatedAt);

    @Modifying
    @Query("""
            update WorkOrder w
               set w.readStatus = :read,
                   w.updatedAt = :updatedAt
             where w.id = :id
               and w.deletedAt is null
            """)
    int markRead(
            @Param("id") Long id,
            @Param("read") WorkOrderReadStatus read,
            @Param("updatedAt") Instant updatedAt);
}
