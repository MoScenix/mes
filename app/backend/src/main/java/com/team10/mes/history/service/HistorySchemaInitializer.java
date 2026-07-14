package com.team10.mes.history.service;

import com.team10.mes.history.dal.HistorySchemaMapper;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.stereotype.Component;

@Component
public class HistorySchemaInitializer implements ApplicationRunner {
  private final HistorySchemaMapper mapper;

  public HistorySchemaInitializer(HistorySchemaMapper mapper) {
    this.mapper = mapper;
  }

  @Override
  public void run(ApplicationArguments args) {
    mapper.createHistoryTable();
    addColumnIfMissing("history_id", mapper::addMessagesHistoryIdColumn);
    addColumnIfMissing("user_id", mapper::addMessagesUserIdColumn);
    addIndexIfMissing("idx_messages_history_time", mapper::createMessagesHistoryTimeIndex);
    addIndexIfMissing("idx_messages_user_time", mapper::createMessagesUserTimeIndex);
  }

  private void addColumnIfMissing(String column, Runnable ddl) {
    if (mapper.countColumn("messages", column) == 0) ddl.run();
  }

  private void addIndexIfMissing(String index, Runnable ddl) {
    if (mapper.countIndex("messages", index) == 0) ddl.run();
  }
}
