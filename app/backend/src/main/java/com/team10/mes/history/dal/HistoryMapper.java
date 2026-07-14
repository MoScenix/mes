package com.team10.mes.history.dal;

import java.time.LocalDateTime;
import java.util.List;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface HistoryMapper {
  record HistoryRow(
      Long id,
      String historyName,
      Long userId,
      LocalDateTime createTime,
      LocalDateTime updateTime) {}

  int insert(MutableHistory history);

  HistoryRow find(long id);

  List<HistoryRow> list(
      @Param("userId") Long userId,
      @Param("name") String name,
      @Param("limit") int limit,
      @Param("offset") long offset);

  long count(@Param("userId") Long userId, @Param("name") String name);

  int rename(@Param("id") long id, @Param("name") String name);

  int touch(long id);

  int softDelete(long id);

  final class MutableHistory {
    private Long id;
    private String name;
    private Long userId;

    public Long getId() {
      return id;
    }

    public void setId(Long id) {
      this.id = id;
    }

    public String getName() {
      return name;
    }

    public void setName(String name) {
      this.name = name;
    }

    public Long getUserId() {
      return userId;
    }

    public void setUserId(Long userId) {
      this.userId = userId;
    }
  }
}
