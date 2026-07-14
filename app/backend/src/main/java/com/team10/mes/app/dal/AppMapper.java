package com.team10.mes.app.dal;

import java.time.LocalDateTime;
import java.util.List;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface AppMapper {
  record AppRow(
      Long id, String appName, Long userId, LocalDateTime createTime, LocalDateTime updateTime) {}

  int insert(MutableApp app);

  AppRow find(long id);

  List<AppRow> list(
      @Param("userId") Long userId,
      @Param("name") String name,
      @Param("limit") int limit,
      @Param("offset") long offset);

  long count(@Param("userId") Long userId, @Param("name") String name);

  int rename(@Param("id") long id, @Param("name") String name);

  int softDelete(long id);

  final class MutableApp {
    private Long id;
    private String name;
    private Long userId;

    public Long getId() {
      return id;
    }

    public void setId(Long v) {
      id = v;
    }

    public String getName() {
      return name;
    }

    public void setName(String v) {
      name = v;
    }

    public Long getUserId() {
      return userId;
    }

    public void setUserId(Long v) {
      userId = v;
    }
  }
}
