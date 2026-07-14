package com.team10.mes.user.dal;

import java.util.List;
import java.util.Map;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface UserMapper {
  Map<String, Object> findById(Long id);

  Map<String, Object> findByAccount(String account);

  int countByAccount(String account);

  int insert(Map<String, Object> user);

  int update(Map<String, Object> user);

  int softDelete(Long id);

  List<Map<String, Object>> page(
      @Param("name") String name,
      @Param("account") String account,
      @Param("offset") int offset,
      @Param("size") int size);

  long pageCount(@Param("name") String name, @Param("account") String account);
}
