package com.team10.mes.ai.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.team10.mes.ai.control.AiControlWatcher;
import com.team10.mes.ai.controller.AiController.Answer;
import com.team10.mes.ai.graph.AiGraph;
import com.team10.mes.ai.node.coder.CoderNode;
import com.team10.mes.ai.node.designer.DesignerNode;
import com.team10.mes.ai.state.RedisAiStore;
import com.team10.mes.ai.task.ChatTask;
import com.team10.mes.ai.workpool.AiWorkPool;
import com.team10.mes.history.model.HistoryMessage;
import com.team10.mes.history.service.HistoryMessageService;
import com.team10.mes.history.service.HistorySessionService;
import com.team10.mes.inventory.service.InventoryService;
import com.team10.mes.user.service.UnauthorizedException;
import com.team10.mes.user.service.UserService;
import com.team10.mes.workorder.service.WorkOrderService;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicLong;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.messages.*;
import org.springframework.ai.support.ToolCallbacks;
import org.springframework.ai.tool.ToolCallback;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class AiService {
  private static final Set<String> ACTIVE =
      Set.of("queued", "running", "waiting_answer", "interrupted");
  private static final String AGENT = "Assistant";
  private static final String DEFAULT_PROMPT =
      """
        You are the assistant agent for an MES system. Help users with work orders, engineering orders,
        inventory flows, items and production. Never invent identifiers. If required information is
        missing, return only JSON in this shape: {"questions":[{"question":"...","options":[]}]}
        so the application can pause and ask the user. Otherwise answer normally.
        """;

  private final RedisAiStore store;
  private final HistoryMessageService histories;
  private final ChatClient chatClient;
  private final ObjectMapper json;
  private final HistorySessionService historySessions;
  private final WorkOrderService workOrders;
  private final InventoryService inventory;
  private final UserService users;
  private final AiWorkPool workPool;
  private final ConcurrentMap<Long, RuntimeTask> tasks = new ConcurrentHashMap<>();
  private final AiGraph graph = new AiGraph(new DesignerNode(), new CoderNode());
  private final int historyLimit;
  private final String systemPrompt;
  private final long controlBlockMs;
  private final int controlReadCount;

  public AiService(
      RedisAiStore store,
      HistoryMessageService histories,
      ChatClient.Builder builder,
      ObjectMapper json,
      HistorySessionService historySessions,
      WorkOrderService workOrders,
      InventoryService inventory,
      UserService users,
      AiWorkPool workPool,
      @Value("${mes.ai.history-limit:20}") int historyLimit,
      @Value("${mes.ai.system-prompt:}") String systemPrompt,
      @Value("${mes.ai.control.block-ms:1000}") long controlBlockMs,
      @Value("${mes.ai.control.read-count:10}") int controlReadCount) {
    this.store = store;
    this.histories = histories;
    this.chatClient = builder.build();
    this.json = json;
    this.historySessions = historySessions;
    this.workOrders = workOrders;
    this.inventory = inventory;
    this.users = users;
    this.workPool = workPool;
    this.historyLimit = Math.max(1, Math.min(100, historyLimit));
    this.systemPrompt =
        systemPrompt == null || systemPrompt.isBlank() ? DEFAULT_PROMPT : systemPrompt.trim();
    this.controlBlockMs = Math.max(1, controlBlockMs);
    this.controlReadCount = Math.max(1, controlReadCount);
  }

  public void submit(long historyId, String message, Identity identity) {
    requireHistory(historyId);
    authorize(historyId, identity);
    String text = trim(message);
    synchronized (lock(historyId)) {
      AiState current = store.state(historyId);
      if (current != null
          && ACTIVE.contains(current.status())
          && !"interrupted".equals(current.status())) {
        throw new IllegalStateException("AI task is already active");
      }
      if (!text.isEmpty()) histories.append(historyId, identity.userId(), "user", text, false);
      if (histories.history(historyId, 1, null, null).messageList().isEmpty())
        throw new IllegalArgumentException("message is required");
      historySessions.touch(historyId);
      store.resetEvents(historyId);
      start(historyId, "task accepted", identity, false, "");
    }
  }

  public String push(long historyId, String content, Identity identity) {
    requireHistory(historyId);
    authorize(historyId, identity);
    String text = trim(content);
    if (text.isEmpty()) throw new IllegalArgumentException("content is required");
    synchronized (lock(historyId)) {
      AiState state = store.state(historyId);
      if (state == null || !ACTIVE.contains(state.status()))
        throw new IllegalStateException("AI task is not active");
      return control(historyId, "push", text, null, null);
    }
  }

  public void answer(long historyId, Map<String, Answer> answers, Identity identity) {
    requireHistory(historyId);
    authorize(historyId, identity);
    Map<String, Answer> normalized =
        answers == null
            ? Map.of()
            : answers.entrySet().stream()
                .filter(e -> !trim(e.getKey()).isEmpty() && e.getValue() != null)
                .filter(
                    e ->
                        !trim(e.getValue().content()).isEmpty()
                            || (e.getValue().payload() != null
                                && !e.getValue().payload().isEmpty()))
                .collect(java.util.stream.Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));
    if (normalized.isEmpty()) throw new IllegalArgumentException("answers are required");
    synchronized (lock(historyId)) {
      AiState state = state(historyId);
      if (!state.exists() || state.pendingInterrupts().isEmpty())
        throw new IllegalStateException("pending interrupt target not found");
      QuestionCheckpoint checkpoint = store.checkpoint(historyId, QuestionCheckpoint.class);
      if (checkpoint == null || !Objects.equals(checkpoint.checkpointId(), state.checkpointId()))
        throw new IllegalStateException("AI checkpoint is missing or invalid");
      Set<String> targets = new HashSet<>();
      checkpoint.pendingInterrupts().forEach(i -> targets.add(i.id()));
      normalized
          .keySet()
          .forEach(
              id -> {
                if (!targets.contains(id))
                  throw new IllegalArgumentException(
                      "answer target does not match pending interrupts: " + id);
              });
      String answerText =
          normalized.values().stream()
              .map(a -> !trim(a.content()).isEmpty() ? trim(a.content()) : writeJson(a.payload()))
              .reduce((a, b) -> a + "\n" + b)
              .orElseThrow();
      histories.append(historyId, identity.userId(), "user", answerText, false);
      String target = normalized.keySet().iterator().next();
      control(historyId, "answer", answerText, target, Map.of("answers", normalized));
      publish(
          historyId,
          "answer",
          AGENT,
          answerText,
          target,
          null,
          null,
          Map.of("answers", normalized));
      historySessions.touch(historyId);
      start(historyId, "assistant resume accepted", identity, true, answerText);
    }
  }

  public String cancel(long historyId, String reason, Identity identity) {
    requireHistory(historyId);
    authorize(historyId, identity);
    String cause = trim(reason).isEmpty() ? "cancelled" : trim(reason);
    synchronized (lock(historyId)) {
      String controlId = control(historyId, "cancel", cause, null, null);
      if (tasks.containsKey(historyId)) return controlId;
      AiEvent cancelled = publish(historyId, "cancelled", AGENT, cause, null, null, null, null);
      cancelRuntime(historyId);
      saveState(historyId, "cancelled", AGENT, cancelled.id(), "", List.of(), cause, "", true);
      store.expireTerminal(historyId);
      return controlId;
    }
  }

  public AiState state(long historyId) {
    requireHistory(historyId);
    AiState current = store.state(historyId);
    return current == null
        ? new AiState(false, "", "", "", "", List.of(), "", "", false, 0)
        : current;
  }

  public EventPage events(long historyId, String lastId, long blockMs, int count) {
    requireHistory(historyId);
    List<AiEvent> result = store.events(historyId, lastId, blockMs, count);
    String next =
        result.isEmpty()
            ? (lastId == null || lastId.isBlank() ? "0" : lastId)
            : result.getLast().id();
    return new EventPage(result, next);
  }

  private void start(
      long historyId,
      String acceptedMessage,
      Identity identity,
      boolean resume,
      String resumeAnswer) {
    RuntimeTask previous = tasks.get(historyId);
    long generation = previous == null ? 1 : previous.generation.incrementAndGet();
    RuntimeTask runtime = new RuntimeTask(generation, identity, resumeAnswer);
    tasks.put(historyId, runtime);
    AiEvent accepted =
        publish(historyId, "accepted", AGENT, acceptedMessage, null, null, null, null);
    saveState(historyId, "queued", AGENT, accepted.id(), "", List.of(), acceptedMessage, "", false);
    runtime.task =
        new ChatTask(
            historyId,
            identity,
            resume,
            ignored -> run(historyId, runtime),
            (task, status) -> markLifecycle(historyId, status));
    if (resume) {
      QuestionCheckpoint checkpoint = store.checkpoint(historyId, QuestionCheckpoint.class);
      if (checkpoint != null && !checkpoint.pendingInterrupts().isEmpty())
        runtime.task.runtime().setControlCursor(checkpointControlCursor(checkpoint));
    }
    runtime.task.enqueue();
    workPool.submit(runtime.task);
  }

  private void run(long historyId, RuntimeTask runtime) {
    AiControlWatcher watcher =
        new AiControlWatcher(
            store,
            historyId,
            runtime.task.runtime().controlCursor(),
            controlBlockMs,
            controlReadCount,
            new AiControlWatcher.Handler() {
              @Override
              public void onPush(AiEvent event) {
                runtime.task.runtime().setControlCursor(event.id());
                if (!current(historyId, runtime)) return;
                String text = trim(event.content());
                if (text.isEmpty()) return;
                histories.append(historyId, runtime.identity.userId(), "user", text, false);
                publish(historyId, "push", AGENT, text, null, null, null, null);
                cancelRuntime(historyId);
                start(historyId, "push accepted", runtime.identity, false, "");
              }

              @Override
              public void onCancel(AiEvent event) {
                runtime.task.runtime().setControlCursor(event.id());
                if (!current(historyId, runtime)) return;
                runtime.task.cancel();
                AiEvent cancelled =
                    publish(
                        historyId,
                        "cancelled",
                        AGENT,
                        trim(event.content()).isEmpty() ? "cancelled" : event.content(),
                        null,
                        null,
                        null,
                        null);
                saveState(
                    historyId,
                    "cancelled",
                    AGENT,
                    cancelled.id(),
                    "",
                    List.of(),
                    event.content(),
                    runtime.task.runtime().buffer().value(),
                    true);
                store.expireTerminal(historyId);
              }

              @Override
              public void onAnswer(AiEvent event) {
                runtime.task.runtime().setControlCursor(event.id());
              }
            });
    Thread watcherThread = Thread.ofVirtual().name("ai-control-" + historyId).start(watcher);
    try {
      if (!current(historyId, runtime)) return;
      AiEvent started = publish(historyId, "agent_start", AGENT, "", null, null, null, null);
      saveState(historyId, "running", AGENT, started.id(), "", List.of(), "", "", false);
      List<Message> history = history(historyId);
      MesAiTools tools =
          new MesAiTools(
              workOrders, inventory, users, runtime.identity.userId(), runtime.identity.admin());
      DesignerNode.ModelRunner model =
          (resumeContext, chunks) -> {
            List<Message> promptHistory = new ArrayList<>(history);
            if (resumeContext != null && !resumeContext.isBlank())
              promptHistory.add(new UserMessage(resumeContext));
            StringBuilder output = new StringBuilder();
            for (String chunk :
                chatClient
                    .prompt()
                    .system(systemPrompt)
                    .messages(promptHistory)
                    .toolCallbacks(toolCallbacks(tools, runtime.identity.role(), historyId))
                    .stream()
                    .content()
                    .toIterable()) {
              if (!current(historyId, runtime)) throw new CancellationException();
              if (chunk != null && !chunk.isEmpty()) {
                output.append(chunk);
                runtime.task.runtime().buffer().append(chunk);
                chunks.accept(chunk);
              }
            }
            return output.toString();
          };
      java.util.function.Function<List<Question>, PendingInterrupt> interruptFactory =
          questions -> {
            String target = UUID.randomUUID().toString();
            Map<String, Object> payload =
                Map.of(
                    "questions",
                    questions,
                    "control_cursor",
                    runtime.task.runtime().controlCursor());
            AiEvent event =
                publish(
                    historyId,
                    "question",
                    AGENT,
                    questions.stream()
                        .map(Question::question)
                        .reduce((a, b) -> a + "\n" + b)
                        .orElse(""),
                    target,
                    null,
                    null,
                    payload);
            return new PendingInterrupt(target, AGENT, event.content(), writeJson(payload));
          };
      java.util.function.Function<PendingInterrupt, QuestionCheckpoint> checkpointFactory =
          pending -> {
            String id = UUID.randomUUID().toString();
            return new QuestionCheckpoint(
                id,
                List.of(pending),
                runtime.task.runtime().buffer().value(),
                store.state(historyId).lastEventId(),
                Instant.now().toEpochMilli());
          };
      AiGraph.CheckpointWriter checkpointWriter =
          checkpoint -> {
            store.saveCheckpoint(historyId, checkpoint);
            saveState(
                historyId,
                "waiting_answer",
                AGENT,
                checkpoint.lastEventId(),
                checkpoint.checkpointId(),
                checkpoint.pendingInterrupts(),
                "",
                checkpoint.modelOutput(),
                false);
          };
      QuestionCheckpoint persisted =
          runtime.task.needsResume() ? store.checkpoint(historyId, QuestionCheckpoint.class) : null;
      AiGraph.Result result =
          runtime.task.needsResume()
              ? graph.resume(
                  persisted,
                  runtime.resumeAnswer,
                  persisted == null ? "" : persisted.modelOutput(),
                  model,
                  this::parseQuestions,
                  chunk -> publish(historyId, "message", AGENT, chunk, null, null, null, null),
                  interruptFactory,
                  checkpointFactory,
                  checkpointWriter,
                  output ->
                      histories.append(
                          historyId, runtime.identity.userId(), "assistant", output, false))
              : graph.run(
                  model,
                  this::parseQuestions,
                  chunk -> publish(historyId, "message", AGENT, chunk, null, null, null, null),
                  interruptFactory,
                  checkpointFactory,
                  checkpointWriter,
                  output ->
                      histories.append(
                          historyId, runtime.identity.userId(), "assistant", output, false));
      if (result.status() == AiGraph.Status.INTERRUPTED) return;
      AiEvent done = publish(historyId, "done", AGENT, "", null, null, null, null);
      saveState(historyId, "done", AGENT, done.id(), "", List.of(), "", result.output(), false);
      store.deleteCheckpoint(historyId);
      store.expireTerminal(historyId);
    } catch (Throwable error) {
      if (!current(historyId, runtime)) return;
      String message =
          trim(error.getMessage()).isEmpty()
              ? error.getClass().getSimpleName()
              : trim(error.getMessage());
      AiEvent event = publish(historyId, "error", AGENT, truncate(message), null, null, null, null);
      saveState(historyId, "error", AGENT, event.id(), "", List.of(), message, "", false);
      store.expireTerminal(historyId);
    } finally {
      watcher.close();
      watcherThread.interrupt();
      tasks.computeIfPresent(historyId, (id, value) -> value == runtime ? null : value);
    }
  }

  private List<Message> history(long historyId) {
    List<HistoryMessage> rows =
        new ArrayList<>(histories.history(historyId, historyLimit, null, null).messageList());
    Collections.reverse(rows);
    return rows.stream()
        .map(
            row ->
                switch (empty(row.getRole()).toLowerCase(Locale.ROOT)) {
                  case "assistant", "ai" -> (Message) new AssistantMessage(historyContent(row));
                  case "system" -> new SystemMessage(historyContent(row));
                  default -> new UserMessage(historyContent(row));
                })
        .toList();
  }

  private String historyContent(HistoryMessage row) {
    if (!Boolean.TRUE.equals(row.getIsFile())) return empty(row.getContent());
    try {
      JsonNode meta = json.readTree(row.getContent());
      String filename = trim(meta.path("filename").asText());
      if (filename.isEmpty()) filename = "unnamed file";
      if (meta.path("isBig").asBoolean(false)) {
        return "User uploaded a large file. Search the history file when its contents are needed.\n"
            + "file_id: "
            + meta.path("fileId").asLong()
            + "\nfilename: "
            + filename
            + "\ncontent_type: "
            + meta.path("contentType").asText("");
      }
      String text = meta.path("text").asText("");
      return "User uploaded a file.\nfilename: "
          + filename
          + "\n\nFile content:\n"
          + (text.isBlank() ? "(empty or no text was extracted)" : text);
    } catch (Exception ignored) {
      return "User uploaded a file, but its metadata could not be parsed.";
    }
  }

  private AiEvent publish(
      long historyId,
      String type,
      String agent,
      String content,
      String targetId,
      String name,
      String status,
      Object payload) {
    String payloadJson = payload == null ? "" : writeJson(payload);
    AiEvent event =
        new AiEvent(
            "",
            String.valueOf(historyId),
            type,
            empty(agent),
            empty(content),
            empty(targetId),
            empty(name),
            empty(status),
            payloadJson,
            Instant.now().toEpochMilli(),
            "question".equals(type) ? parsePayloadQuestions(payloadJson) : List.of());
    String id = store.addEvent(historyId, event);
    return new AiEvent(
        id,
        event.historyId(),
        event.type(),
        event.agent(),
        event.content(),
        event.targetId(),
        event.name(),
        event.status(),
        event.payloadJson(),
        event.createdAt(),
        event.questions());
  }

  private String control(
      long historyId, String type, String content, String targetId, Object payload) {
    AiEvent event =
        new AiEvent(
            "",
            String.valueOf(historyId),
            type,
            AGENT,
            empty(content),
            empty(targetId),
            "",
            "",
            payload == null ? "" : writeJson(payload),
            Instant.now().toEpochMilli(),
            List.of());
    return store.addControl(historyId, event);
  }

  private void saveState(
      long historyId,
      String status,
      String agent,
      String lastEventId,
      String checkpointId,
      List<PendingInterrupt> pending,
      String message,
      String buffer,
      boolean cancelled) {
    store.saveState(
        historyId,
        new AiState(
            true,
            status,
            agent,
            empty(lastEventId),
            empty(checkpointId),
            pending,
            empty(message),
            empty(buffer),
            cancelled,
            Instant.now().toEpochMilli()));
  }

  private void markLifecycle(long historyId, String status) {
    AiState old = store.state(historyId);
    if (old != null)
      store.saveState(
          historyId,
          new AiState(
              true,
              status,
              old.agent(),
              old.lastEventId(),
              old.checkpointId(),
              old.pendingInterrupts(),
              old.message(),
              old.buffer(),
              old.isCancelled(),
              Instant.now().toEpochMilli()));
  }

  private Optional<List<Question>> parseQuestions(String text) {
    try {
      String source = text.replaceFirst("^```(?:json)?\\s*", "").replaceFirst("\\s*```$", "");
      JsonNode node = json.readTree(source).path("questions");
      if (!node.isArray() || node.isEmpty()) return Optional.empty();
      List<Question> out = new ArrayList<>();
      node.forEach(
          q -> {
            String question = trim(q.path("question").asText());
            if (!question.isEmpty()) {
              List<String> options = new ArrayList<>();
              q.path("options")
                  .forEach(
                      o -> {
                        if (!trim(o.asText()).isEmpty()) options.add(trim(o.asText()));
                      });
              out.add(new Question(question, options));
            }
          });
      return out.isEmpty() ? Optional.empty() : Optional.of(out);
    } catch (Exception ignored) {
      return Optional.empty();
    }
  }

  private List<Question> parsePayloadQuestions(String payload) {
    try {
      JsonNode n = json.readTree(payload).path("questions");
      List<Question> out = new ArrayList<>();
      n.forEach(
          q ->
              out.add(
                  new Question(
                      q.path("question").asText(),
                      json.convertValue(q.path("options"), new TypeReference<>() {}))));
      return out;
    } catch (Exception ignored) {
      return List.of();
    }
  }

  ToolCallback[] toolCallbacks(MesAiTools tools, String role, long historyId) {
    Set<String> common = Set.of("search_users");
    Set<String> work =
        Set.of(
            "list_work_orders",
            "mark_work_order_read",
            "create_work_order_draft",
            "update_work_order_draft");
    Set<String> engineering =
        Set.of(
            "create_engineering_order_draft",
            "update_engineering_order_draft",
            "list_engineering_orders",
            "get_engineering_order");
    Set<String> flows =
        Set.of("create_inventory_flow_draft", "list_inventory_flows", "get_inventory_flow");
    Set<String> flowRead = Set.of("list_inventory_flows", "get_inventory_flow");
    Set<String> items = Set.of("search_items", "get_item", "list_item_units");
    Set<String> warehouse = Set.of("list_pending_inventory_flows", "inventory_check");
    String normalized = empty(role).toLowerCase(Locale.ROOT);
    Set<String> allowed = new HashSet<>(common);
    switch (normalized) {
      case "admin", "administrator", "管理员" -> {
        allowed.addAll(work);
        allowed.addAll(engineering);
        allowed.addAll(flows);
        allowed.addAll(items);
        allowed.addAll(warehouse);
      }
      case "leader", "组长" -> {
        allowed.addAll(work);
        allowed.addAll(engineering);
        allowed.addAll(flows);
      }
      case "purchase", "采购专员", "sales", "销售" -> {
        allowed.addAll(work);
        allowed.addAll(flows);
      }
      case "process_engineer", "工艺工程师" -> {
        allowed.addAll(items);
        allowed.addAll(engineering);
      }
      case "warehouse_admin", "仓库管理员" -> {
        allowed.addAll(work);
        allowed.addAll(items);
        allowed.addAll(flowRead);
        allowed.addAll(warehouse);
      }
      default -> {
        allowed.addAll(work);
        allowed.addAll(flowRead);
      }
    }
    return Arrays.stream(ToolCallbacks.from(tools))
        .filter(t -> allowed.contains(t.getToolDefinition().name()))
        .map(
            delegate ->
                new ToolCallback() {
                  @Override
                  public org.springframework.ai.tool.definition.ToolDefinition getToolDefinition() {
                    return delegate.getToolDefinition();
                  }

                  @Override
                  public org.springframework.ai.tool.metadata.ToolMetadata getToolMetadata() {
                    return delegate.getToolMetadata();
                  }

                  @Override
                  public String call(String input) {
                    String target = UUID.randomUUID().toString();
                    String name = delegate.getToolDefinition().name();
                    publish(
                        historyId,
                        "tool_call",
                        AGENT,
                        "",
                        target,
                        name,
                        "running",
                        Map.of("arguments", input));
                    try {
                      String result = delegate.call(input);
                      publish(
                          historyId,
                          "tool_result",
                          AGENT,
                          truncate(empty(result)),
                          target,
                          name,
                          "success",
                          null);
                      return result;
                    } catch (RuntimeException error) {
                      publish(
                          historyId,
                          "tool_result",
                          AGENT,
                          truncate(empty(error.getMessage())),
                          target,
                          name,
                          "error",
                          null);
                      throw error;
                    }
                  }
                })
        .toArray(ToolCallback[]::new);
  }

  private String writeJson(Object value) {
    try {
      return json.writeValueAsString(value);
    } catch (JsonProcessingException e) {
      throw new IllegalStateException(e);
    }
  }

  private String checkpointControlCursor(QuestionCheckpoint checkpoint) {
    try {
      return json.readTree(checkpoint.pendingInterrupts().getFirst().payloadJson())
          .path("control_cursor")
          .asText("0");
    } catch (Exception ignored) {
      return "0";
    }
  }

  public void authorize(long historyId, Identity identity) {
    requireHistory(historyId);
    if (identity == null || identity.userId() <= 0) throw new UnauthorizedException();
    historySessions.authorize(historyId, identity.userId(), identity.role());
  }

  private Object lock(long historyId) {
    return tasks.computeIfAbsent(
        -historyId, ignored -> new RuntimeTask(0, new Identity(0, ""), ""));
  }

  private boolean current(long historyId, RuntimeTask task) {
    return tasks.get(historyId) == task
        && !task.cancelled
        && (task.task == null || !task.task.runtime().isCancelled());
  }

  private void cancelRuntime(long historyId) {
    RuntimeTask task = tasks.remove(historyId);
    if (task != null) {
      task.cancelled = true;
      if (task.task != null) task.task.cancel();
    }
  }

  private static void requireHistory(long historyId) {
    if (historyId <= 0) throw new IllegalArgumentException("historyId is required");
  }

  private static String trim(String value) {
    return value == null ? "" : value.trim();
  }

  private static String empty(String value) {
    return value == null ? "" : value;
  }

  private static String truncate(String value) {
    return value.length() <= 1200 ? value : value.substring(0, 1200);
  }

  private static final class RuntimeTask {
    final AtomicLong generation;
    final Identity identity;
    final String resumeAnswer;
    volatile ChatTask task;
    volatile boolean cancelled;

    RuntimeTask(long generation, Identity identity, String resumeAnswer) {
      this.generation = new AtomicLong(generation);
      this.identity = identity;
      this.resumeAnswer = resumeAnswer;
    }
  }

  public record Identity(long userId, String role) {
    public boolean admin() {
      return "admin".equalsIgnoreCase(role)
          || "administrator".equalsIgnoreCase(role)
          || "管理员".equals(role);
    }
  }

  public record Question(String question, List<String> options) {}

  public record PendingInterrupt(String id, String agent, String content, String payloadJson) {}

  public record QuestionCheckpoint(
      String checkpointId,
      List<PendingInterrupt> pendingInterrupts,
      String modelOutput,
      String lastEventId,
      long createdAt) {}

  public record AiState(
      boolean exists,
      String status,
      String agent,
      String lastEventId,
      String checkpointId,
      List<PendingInterrupt> pendingInterrupts,
      String message,
      String buffer,
      boolean isCancelled,
      long updatedAt) {}

  public record AiEvent(
      String id,
      String historyId,
      String type,
      String agent,
      String content,
      String targetId,
      String name,
      String status,
      String payloadJson,
      long createdAt,
      List<Question> questions) {}

  public record EventPage(List<AiEvent> events, String lastId) {}
}
