package com.team10.mes.app.service;

import com.team10.mes.app.dal.AppMapper;
import java.util.List;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class AppService {
  private final AppMapper mapper;

  public AppService(AppMapper mapper) {
    this.mapper = mapper;
  }

  public AppMapper.AppRow get(long id) {
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
  public AppMapper.AppRow create(long userId, String initialPrompt) {
    if (userId <= 0) throw new IllegalArgumentException("userId must be positive");
    String prompt = initialPrompt == null ? "" : initialPrompt.trim();
    String name =
        prompt.isEmpty() ? "New conversation" : prompt.substring(0, Math.min(12, prompt.length()));
    AppMapper.MutableApp app = new AppMapper.MutableApp();
    app.setName(name);
    app.setUserId(userId);
    mapper.insert(app);
    return mapper.find(app.getId());
  }

  public boolean rename(long id, String name) {
    if (name == null || name.isBlank())
      throw new IllegalArgumentException("appName must not be blank");
    return mapper.rename(id, name.trim()) > 0;
  }

  public boolean delete(long id) {
    return mapper.softDelete(id) > 0;
  }

  public record Page(
      List<AppMapper.AppRow> records,
      int pageNumber,
      int pageSize,
      long totalPage,
      long totalRow,
      boolean optimizeCountQuery) {}
}
