package com.team10.mes.history.controller;

import com.team10.mes.history.dal.HistoryMapper;
import com.team10.mes.history.model.HistoryMessage;
import com.team10.mes.history.service.HistoryFileService;
import com.team10.mes.history.service.HistoryMessageService;
import com.team10.mes.history.service.HistorySessionService;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import java.io.IOException;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/history")
public class HistoryController {
  private static final DateTimeFormatter TIME = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
  private final HistoryMessageService messages;
  private final HistorySessionService histories;
  private final HistoryFileService files;
  private final SessionIdentity identity;

  public HistoryController(
      HistoryMessageService messages,
      HistorySessionService histories,
      HistoryFileService files,
      SessionIdentity identity) {
    this.messages = messages;
    this.histories = histories;
    this.files = files;
    this.identity = identity;
  }

  @PostMapping("/add")
  public Long create(@RequestBody CreateRequest request, HttpSession session) {
    return histories.create(userId(session), request.initPrompt()).id();
  }

  @PostMapping("/my/list/page/vo")
  public HistorySessionService.Page listMine(
      @RequestBody ListRequest request, HttpSession session) {
    return histories.list(
        userId(session), request.historyName(), request.pageNum(), request.pageSize());
  }

  @PostMapping("/delete")
  public Result delete(@RequestBody DeleteRequest request, HttpSession session) throws IOException {
    histories.authorize(request.id(), userId(session), identity.role(session));
    boolean deleted = histories.delete(request.id());
    if (deleted) files.deleteHistoryFiles(request.id());
    return new Result(deleted);
  }

  @GetMapping("/{historyId}/messages")
  public PageHistoryMessage messages(
      @PathVariable long historyId,
      @RequestParam(defaultValue = "20") int pageSize,
      @RequestParam(required = false) String lastCreateTime,
      HttpSession session) {
    histories.authorize(historyId, userId(session), identity.role(session));
    LocalDateTime before = parse(lastCreateTime);
    var page = messages.history(historyId, pageSize, before, null);
    int size = Math.max(1, Math.min(pageSize, 100));
    return new PageHistoryMessage(
        page.messageList().stream().map(this::view).toList(),
        1,
        size,
        (page.total() + size - 1) / size,
        page.total(),
        true);
  }

  @PostMapping("/admin/list/page/vo")
  public PageHistoryMessage adminList(@RequestBody Query q, HttpSession session) {
    requireAdmin(session);
    var page =
        messages.adminPage(
            q.id,
            q.message,
            q.messageType,
            q.historyId,
            q.userId,
            parse(q.lastCreateTime),
            q.pageNum,
            q.pageSize);
    return new PageHistoryMessage(
        page.records().stream().map(this::view).toList(),
        page.pageNumber(),
        page.pageSize(),
        page.totalPage(),
        page.totalRow(),
        page.optimizeCountQuery());
  }

  @PostMapping("/admin/delete")
  public Result adminDelete(@RequestBody DeleteRequest q, HttpSession session) {
    requireAdmin(session);
    return new Result(messages.deleteById(q.id));
  }

  private long userId(HttpSession session) {
    Long id = identity.userId(session);
    if (id == null) throw new UnauthorizedException();
    return id;
  }

  private void requireAdmin(HttpSession session) {
    if (!"admin".equalsIgnoreCase(identity.role(session))) throw new UnauthorizedException();
  }

  private LocalDateTime parse(String value) {
    return value == null || value.isBlank() ? null : LocalDateTime.parse(value, TIME);
  }

  private HistoryMessageView view(HistoryMessage m) {
    return new HistoryMessageView(
        m.getId(),
        m.getContent(),
        m.getRole(),
        m.getHistoryId(),
        m.getUserId(),
        format(m.getCreateTime()),
        format(m.getUpdateTime()),
        0,
        Boolean.TRUE.equals(m.getIsFile()));
  }

  private String format(LocalDateTime value) {
    return value == null ? "" : value.format(TIME);
  }

  public record CreateRequest(String initPrompt) {}

  public record ListRequest(int pageNum, int pageSize, String historyName) {}

  public record PageHistoryMessage(
      List<HistoryMessageView> records,
      long pageNumber,
      long pageSize,
      long totalPage,
      long totalRow,
      boolean optimizeCountQuery) {}

  public record HistoryMessageView(
      Long id,
      String message,
      String messageType,
      Long historyId,
      Long userId,
      String createTime,
      String updateTime,
      long isDelete,
      boolean isFile) {}

  public record Query(
      int pageNum,
      int pageSize,
      String sortField,
      String sortOrder,
      Long id,
      String message,
      String messageType,
      Long historyId,
      Long userId,
      String lastCreateTime) {}

  public record DeleteRequest(long id) {}

  public record Result(boolean success) {}

  public record HistoryVO(
      Long id,
      String historyName,
      Long userId,
      LocalDateTime createTime,
      LocalDateTime updateTime) {
    public static HistoryVO from(HistoryMapper.HistoryRow row) {
      return new HistoryVO(
          row.id(), row.historyName(), row.userId(), row.createTime(), row.updateTime());
    }
  }
}
