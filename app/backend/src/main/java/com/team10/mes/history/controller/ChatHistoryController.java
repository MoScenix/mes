package com.team10.mes.history.controller;

import com.team10.mes.app.dal.AppMapper;
import com.team10.mes.app.service.AppService;
import com.team10.mes.history.model.HistoryMessage;
import com.team10.mes.history.service.HistoryMessageService;
import com.team10.mes.user.service.SessionIdentity;
import jakarta.servlet.http.HttpSession;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping("/chatHistory")
public class ChatHistoryController {
  private static final DateTimeFormatter TIME = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
  private final HistoryMessageService history;
  private final AppService apps;
  private final SessionIdentity identity;

  public ChatHistoryController(
      HistoryMessageService history, AppService apps, SessionIdentity identity) {
    this.history = history;
    this.apps = apps;
    this.identity = identity;
  }

  @GetMapping("/app/{appId}")
  public PageChatHistory app(
      @PathVariable long appId,
      @RequestParam(defaultValue = "20") int pageSize,
      @RequestParam(required = false) String lastCreateTime,
      HttpSession session) {
    requireOwnerOrAdmin(appId, session);
    LocalDateTime before = parse(lastCreateTime);
    var page = history.history(appId, pageSize, before, null);
    int size = Math.max(1, Math.min(pageSize, 100));
    return new PageChatHistory(
        page.messageList().stream().map(this::view).toList(),
        1,
        size,
        (page.total() + size - 1) / size,
        page.total(),
        true);
  }

  @PostMapping("/admin/list/page/vo")
  public PageChatHistory adminList(@RequestBody Query q, HttpSession session) {
    requireAdmin(session);
    var page =
        history.adminPage(
            q.id,
            q.message,
            q.messageType,
            q.appId,
            q.userId,
            parse(q.lastCreateTime),
            q.pageNum,
            q.pageSize);
    return new PageChatHistory(
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
    return new Result(history.deleteById(q.id));
  }

  private void requireOwnerOrAdmin(long appId, HttpSession session) {
    Long uid = identity.userId(session);
    if (uid == null) throw new ResponseStatusException(HttpStatus.UNAUTHORIZED);
    AppMapper.AppRow app = apps.get(appId);
    if (app == null) throw new ResponseStatusException(HttpStatus.NOT_FOUND);
    if (!uid.equals(app.userId()) && !"admin".equalsIgnoreCase(identity.role(session)))
      throw new ResponseStatusException(HttpStatus.FORBIDDEN);
  }

  private void requireAdmin(HttpSession session) {
    if (!"admin".equalsIgnoreCase(identity.role(session)))
      throw new ResponseStatusException(HttpStatus.FORBIDDEN);
  }

  private LocalDateTime parse(String value) {
    return value == null || value.isBlank() ? null : LocalDateTime.parse(value, TIME);
  }

  private ChatHistory view(HistoryMessage m) {
    return new ChatHistory(
        m.getId(),
        m.getContent(),
        m.getRole(),
        m.getAppId(),
        m.getUserId(),
        format(m.getCreateTime()),
        format(m.getUpdateTime()),
        0,
        Boolean.TRUE.equals(m.getIsFile()));
  }

  private String format(LocalDateTime value) {
    return value == null ? "" : value.format(TIME);
  }

  public record PageChatHistory(
      List<ChatHistory> records,
      long pageNumber,
      long pageSize,
      long totalPage,
      long totalRow,
      boolean optimizeCountQuery) {}

  public record ChatHistory(
      Long id,
      String message,
      String messageType,
      Long appId,
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
      Long appId,
      Long userId,
      String lastCreateTime) {}

  public record DeleteRequest(long id) {}

  public record Result(boolean success) {}
}
