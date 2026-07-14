package com.team10.mes.history.service;

import com.team10.mes.history.dal.HistoryMessageMapper;
import com.team10.mes.history.dto.MessagePage;
import com.team10.mes.history.model.HistoryMessage;
import java.time.LocalDateTime;
import java.util.List;
import org.springframework.stereotype.Service;

@Service
public class HistoryMessageService {
  private final HistoryMessageMapper mapper;

  public HistoryMessageService(HistoryMessageMapper mapper) {
    this.mapper = mapper;
  }

  public HistoryMessage append(
      long appId, Long userId, String role, String content, Boolean isFile) {
    if (content == null || content.isBlank())
      throw new IllegalArgumentException("content must not be blank");
    HistoryMessage m = new HistoryMessage();
    m.setAppId(appId);
    m.setUserId(userId);
    m.setRole(role);
    m.setContent(content);
    m.setIsFile(Boolean.TRUE.equals(isFile));
    mapper.insert(m);
    return m;
  }

  public MessagePage history(long appId, int size, LocalDateTime before, Long beforeId) {
    int limit = Math.max(1, Math.min(size, 100));
    var rows = mapper.page(appId, before, beforeId, limit + 1);
    boolean more = rows.size() > limit;
    if (more) rows = rows.subList(0, limit);
    return new MessagePage(rows, mapper.count(appId), more);
  }

  public boolean delete(long appId, long id) {
    return mapper.delete(appId, id) > 0;
  }

  public boolean deleteById(long id) {
    return mapper.deleteById(id) > 0;
  }

  public AdminPage adminPage(
      Long id,
      String message,
      String messageType,
      Long appId,
      Long userId,
      LocalDateTime lastCreateTime,
      int pageNum,
      int pageSize) {
    int number = Math.max(1, pageNum), size = Math.max(1, Math.min(pageSize, 100));
    long total = mapper.adminCount(id, message, messageType, appId, userId, lastCreateTime);
    return new AdminPage(
        mapper.adminPage(
            id,
            message,
            messageType,
            appId,
            userId,
            lastCreateTime,
            (long) (number - 1) * size,
            size),
        number,
        size,
        (total + size - 1) / size,
        total,
        true);
  }

  public record AdminPage(
      List<HistoryMessage> records,
      long pageNumber,
      long pageSize,
      long totalPage,
      long totalRow,
      boolean optimizeCountQuery) {}
}
