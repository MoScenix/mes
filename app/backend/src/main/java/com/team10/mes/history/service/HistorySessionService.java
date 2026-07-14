package com.team10.mes.history.service;

import com.team10.mes.history.dal.HistoryMapper;
import com.team10.mes.user.service.UnauthorizedException;
import java.util.List;
import java.util.Objects;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class HistorySessionService {
  private final HistoryMapper mapper;

  public HistorySessionService(HistoryMapper mapper) {
    this.mapper = mapper;
  }

  public HistoryMapper.HistoryRow get(long id) {
    return mapper.find(id);
  }

  public Page list(Long userId, String name, int page, int size) {
    int limit = Math.max(1, Math.min(size, 100));
    int pageNumber = Math.max(page, 1);
    long total = mapper.count(userId, name);
    return new Page(
        mapper.list(userId, name, limit, (long) (pageNumber - 1) * limit),
        pageNumber,
        limit,
        (total + limit - 1) / limit,
        total,
        true);
  }

  @Transactional
  public HistoryMapper.HistoryRow create(long userId, String initialPrompt) {
    if (userId <= 0) throw new UnauthorizedException();
    String prompt = initialPrompt == null ? "" : initialPrompt.trim();
    String normalized = prompt.replaceAll("\\s+", " ");
    String name =
        normalized.isEmpty() ? "新对话" : normalized.substring(0, Math.min(24, normalized.length()));
    HistoryMapper.MutableHistory history = new HistoryMapper.MutableHistory();
    history.setName(name);
    history.setUserId(userId);
    mapper.insert(history);
    return mapper.find(history.getId());
  }

  public boolean rename(long id, String name) {
    if (name == null || name.isBlank())
      throw new IllegalArgumentException("historyName must not be blank");
    return mapper.rename(id, name.trim()) > 0;
  }

  public boolean touch(long id) {
    return mapper.touch(id) > 0;
  }

  public boolean delete(long id) {
    return mapper.softDelete(id) > 0;
  }

  public void authorize(long historyId, long userId, String role) {
    if (historyId <= 0) throw new IllegalArgumentException("historyId is required");
    if (userId <= 0) throw new UnauthorizedException();
    HistoryMapper.HistoryRow history = mapper.find(historyId);
    if (history == null) throw new IllegalArgumentException("history not found");
    if (!"admin".equalsIgnoreCase(role) && !Objects.equals(history.userId(), userId))
      throw new IllegalStateException("forbidden: history owner or admin required");
  }

  public record Page(
      List<HistoryMapper.HistoryRow> records,
      int pageNumber,
      int pageSize,
      long totalPage,
      long totalRow,
      boolean optimizeCountQuery) {}
}
