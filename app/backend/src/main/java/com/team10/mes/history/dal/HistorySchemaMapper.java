package com.team10.mes.history.dal;

import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

@Mapper
public interface HistorySchemaMapper {
  void createHistoryTable();

  int countColumn(@Param("table") String table, @Param("column") String column);

  int countIndex(@Param("table") String table, @Param("index") String index);

  void addMessagesHistoryIdColumn();

  void addMessagesUserIdColumn();

  void createMessagesHistoryTimeIndex();

  void createMessagesUserTimeIndex();
}
