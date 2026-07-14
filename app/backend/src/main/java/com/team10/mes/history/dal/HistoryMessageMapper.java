package com.team10.mes.history.dal;

import com.team10.mes.history.model.HistoryMessage;
import java.time.LocalDateTime;
import java.util.List;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface HistoryMessageMapper {
  HistoryMessage findById(@Param("historyId") long historyId, @Param("id") long id);

  int insert(HistoryMessage message);

  long count(@Param("historyId") long historyId);

  List<HistoryMessage> page(
      @Param("historyId") long historyId,
      @Param("before") LocalDateTime before,
      @Param("beforeId") Long beforeId,
      @Param("limit") int limit);

  int delete(@Param("historyId") long historyId, @Param("id") long id);

  int deleteById(long id);

  List<HistoryMessage> adminPage(
      @Param("id") Long id,
      @Param("message") String message,
      @Param("messageType") String messageType,
      @Param("historyId") Long historyId,
      @Param("userId") Long userId,
      @Param("lastCreateTime") LocalDateTime lastCreateTime,
      @Param("offset") long offset,
      @Param("limit") int limit);

  long adminCount(
      @Param("id") Long id,
      @Param("message") String message,
      @Param("messageType") String messageType,
      @Param("historyId") Long historyId,
      @Param("userId") Long userId,
      @Param("lastCreateTime") LocalDateTime lastCreateTime);
}
