package com.team10.mes.workorder.service;

import com.team10.mes.workorder.dal.WorkOrderMapper;
import com.team10.mes.workorder.dal.WorkOrderMapper.WorkOrderRow;
import java.time.*;
import java.time.format.*;
import java.util.List;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class WorkOrderService {
  private static final DateTimeFormatter LEGACY_TIME =
      DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
  private final WorkOrderMapper mapper;

  public WorkOrderService(WorkOrderMapper mapper) {
    this.mapper = mapper;
  }

  public record DraftRequest(Long toUserId, String description, String name) {}

  public record UpdateRequest(Long id, Long toUserId, String description, String name) {}

  public record ListRequest(
      Long pageNum,
      Long pageSize,
      Long id,
      Boolean isTo,
      Boolean isUnread,
      String sinceTime,
      Long recentSeconds,
      String cursorUpdatedAt,
      Long cursorId,
      String namePrefix,
      Integer status,
      Integer scope) {}

  public record View(
      Long id,
      Long fromUserId,
      Long toUserId,
      String description,
      Integer status,
      String createTime,
      String updateTime,
      Integer readStatus,
      String name) {}

  public record Page(
      List<View> records,
      long pageNumber,
      long pageSize,
      long totalPage,
      long totalRow,
      boolean optimizeCountQuery,
      boolean hasMore,
      String nextCursorUpdatedAt,
      Long nextCursorId) {}

  @Transactional
  public long create(long userId, DraftRequest req) {
    requireUser(userId);
    require(req != null, "request body is required");
    WorkOrderRow row = new WorkOrderRow();
    row.fromUserId = userId;
    row.toUserId = positive(req.toUserId(), "toUserId");
    row.name = name(req.name());
    row.description = empty(req.description());
    mapper.insert(row);
    return row.id;
  }

  @Transactional
  public void update(long userId, boolean admin, UpdateRequest req) {
    requireUser(userId);
    require(req != null, "request body is required");
    WorkOrderRow current = getPermitted(req.id(), userId, admin, true, false);
    WorkOrderRow row = new WorkOrderRow();
    row.id = current.id;
    row.fromUserId = current.fromUserId;
    row.toUserId = positive(req.toUserId(), "toUserId");
    row.name = name(req.name());
    row.description = empty(req.description());
    require(mapper.updateDraft(row) > 0, "work order is not a draft or does not exist");
  }

  @Transactional
  public void delete(long id, long userId, boolean admin) {
    getPermitted(id, userId, admin, true, false);
    require(mapper.deleteDraft(id) > 0, "work order is not a draft or does not exist");
  }

  @Transactional
  public void submit(long id, long userId, boolean admin) {
    getPermitted(id, userId, admin, true, false);
    require(mapper.submit(id) > 0, "work order is not a draft or does not exist");
  }

  @Transactional
  public void read(long id, long userId, boolean admin) {
    getPermitted(id, userId, admin, false, true);
    require(mapper.markRead(id) > 0, "work order not found");
  }

  public View get(long id, long userId, boolean admin) {
    return view(getPermitted(id, userId, admin, false, false));
  }

  public Page list(ListRequest req, long currentUserId, boolean admin) {
    requireUser(currentUserId);
    require(req != null, "request body is required");
    long requested = req.id() == null ? 0 : req.id();
    long userId = admin && requested > 0 ? requested : currentUserId;
    require(admin || requested == 0 || requested == currentUserId, "forbidden: no permission");
    int size =
        (int) Math.min(200, req.pageSize() == null || req.pageSize() <= 0 ? 10 : req.pageSize());
    long page = req.pageNum() == null || req.pageNum() <= 0 ? 1 : req.pageNum();
    boolean to = Boolean.TRUE.equals(req.isTo());
    Integer status = req.status();
    if (status != null && status == 1) to = false;
    LocalDateTime since = parse(req.sinceTime());
    if (since == null && req.recentSeconds() != null && req.recentSeconds() > 0)
      since = LocalDateTime.now().minusSeconds(req.recentSeconds());
    LocalDateTime cursor = parse(req.cursorUpdatedAt());
    long offset = cursor == null ? (page - 1) * size : 0;
    List<WorkOrderRow> rows =
        mapper.list(
            userId,
            to,
            Boolean.TRUE.equals(req.isUnread()),
            status,
            req.namePrefix(),
            since,
            cursor,
            req.cursorId(),
            offset,
            size + 1);
    boolean more = rows.size() > size;
    if (more) rows = rows.subList(0, size);
    long total =
        mapper.count(
            userId, to, Boolean.TRUE.equals(req.isUnread()), status, req.namePrefix(), since);
    WorkOrderRow last = rows.isEmpty() ? null : rows.get(rows.size() - 1);
    return new Page(
        rows.stream().map(this::view).toList(),
        page,
        size,
        (total + size - 1) / size,
        total,
        true,
        more,
        last == null ? "" : format(last.updatedAt),
        last == null ? 0L : last.id);
  }

  private WorkOrderRow getPermitted(
      Long id, long userId, boolean admin, boolean senderOnly, boolean receiverOnly) {
    requireUser(userId);
    require(id != null && id > 0, "id is required");
    WorkOrderRow row = mapper.find(id);
    require(row != null, "work order not found");
    boolean allowed =
        admin
            || (senderOnly
                ? row.fromUserId == userId
                : receiverOnly
                    ? row.toUserId == userId
                    : (row.status == 1
                        ? row.fromUserId == userId
                        : row.fromUserId == userId || row.toUserId == userId));
    require(allowed, "forbidden: no permission");
    return row;
  }

  private View view(WorkOrderRow r) {
    return new View(
        r.id,
        r.fromUserId,
        r.toUserId,
        r.description,
        r.status,
        format(r.createdAt),
        format(r.updatedAt),
        r.readStatus,
        r.name);
  }

  private String format(LocalDateTime t) {
    return t == null
        ? ""
        : t.atZone(ZoneId.systemDefault())
            .withZoneSameInstant(ZoneOffset.UTC)
            .format(DateTimeFormatter.ISO_OFFSET_DATE_TIME);
  }

  private LocalDateTime parse(String s) {
    if (s == null || s.isBlank()) return null;
    try {
      return LocalDateTime.parse(s, LEGACY_TIME);
    } catch (DateTimeParseException e) {
      try {
        return OffsetDateTime.parse(s).atZoneSameInstant(ZoneId.systemDefault()).toLocalDateTime();
      } catch (DateTimeParseException x) {
        throw new IllegalArgumentException("time must use format yyyy-MM-dd HH:mm:ss");
      }
    }
  }

  private long positive(Long v, String field) {
    require(v != null && v > 0, field + " is required");
    return v;
  }

  private String name(String v) {
    String n = v == null ? "" : v.trim();
    require(!n.isEmpty(), "work order name is required");
    return n;
  }

  private String empty(String v) {
    return v == null ? "" : v;
  }

  private void requireUser(long id) {
    require(id > 0, "unauthorized: user id is required");
  }

  private void require(boolean ok, String message) {
    if (!ok) throw new IllegalArgumentException(message);
  }
}
