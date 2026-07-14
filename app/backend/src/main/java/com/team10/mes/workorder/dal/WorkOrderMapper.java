package com.team10.mes.workorder.dal;

import java.time.LocalDateTime;
import java.util.List;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface WorkOrderMapper {
  class WorkOrderRow {
    public Long id;
    public Long fromUserId;
    public Long toUserId;
    public String name;
    public String description;
    public Integer status;
    public Integer readStatus;
    public LocalDateTime createdAt;
    public LocalDateTime updatedAt;
  }

  int insert(WorkOrderRow row);

  WorkOrderRow find(long id);

  int updateDraft(WorkOrderRow row);

  int deleteDraft(long id);

  int submit(long id);

  int markRead(long id);

  List<WorkOrderRow> list(
      @Param("userId") long userId,
      @Param("isTo") boolean isTo,
      @Param("unread") boolean unread,
      @Param("status") Integer status,
      @Param("namePrefix") String namePrefix,
      @Param("since") LocalDateTime since,
      @Param("cursorTime") LocalDateTime cursorTime,
      @Param("cursorId") Long cursorId,
      @Param("offset") long offset,
      @Param("limit") int limit);

  long count(
      @Param("userId") long userId,
      @Param("isTo") boolean isTo,
      @Param("unread") boolean unread,
      @Param("status") Integer status,
      @Param("namePrefix") String namePrefix,
      @Param("since") LocalDateTime since);
}
