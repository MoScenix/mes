package com.team10.mes.ai.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.team10.mes.ai.control.AiControlWatcher;
import com.team10.mes.ai.controller.AiController.Answer;
import com.team10.mes.ai.state.RedisAiStore;
import com.team10.mes.ai.task.ChatTask;
import com.team10.mes.ai.workpool.AiWorkPool;
import com.team10.mes.document.service.DocumentService;
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
import org.springframework.ai.tool.definition.DefaultToolDefinition;
import org.springframework.ai.tool.definition.ToolDefinition;
import org.springframework.ai.tool.metadata.DefaultToolMetadata;
import org.springframework.ai.tool.metadata.ToolMetadata;
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
        missing, call the ask_user tool and wait for the user's answer. Never print JSON questions
        or tool-call payloads as assistant text. Prefer providing concrete options in ask_user when
        there are clear choices; use an empty options array for open-ended questions. Otherwise answer
        normally.

        When the user asks about an uploaded large file, use search_history_file with the file_id
        shown in the conversation history. Do not claim that you cannot read or search uploaded
        files before trying the tool.

        The authenticated operator is always the current user. When creating work orders,
        engineering orders, or inventory flow drafts, do not ask who the initiator/from user/leader is
        unless the user explicitly asks to create it on behalf of another user and the tool schema
        supports that field. For inventory flow drafts, the from user is supplied by the backend from
        the current session; only ask for the receiver, item, quantity, type, or description when
        those are missing.
        """;

  private final RedisAiStore store;
  private final HistoryMessageService histories;
  private final ChatClient chatClient;
  private final ObjectMapper json;
  private final HistorySessionService historySessions;
  private final WorkOrderService workOrders;
  private final InventoryService inventory;
  private final UserService users;
  private final DocumentService documents;
  private final AiWorkPool workPool;
  private final ConcurrentMap<Long, RuntimeTask> tasks = new ConcurrentHashMap<>();
  private final ConcurrentMap<String, PendingAsk> pendingAsks = new ConcurrentHashMap<>();
  private final int historyLimit;
  private final String systemPrompt;
  private final long controlBlockMs;
  private final int controlReadCount;
  private final long askUserTimeoutMs;

  public AiService(
      RedisAiStore store,
      HistoryMessageService histories,
      ChatClient.Builder builder,
      ObjectMapper json,
      HistorySessionService historySessions,
      WorkOrderService workOrders,
      InventoryService inventory,
      UserService users,
      DocumentService documents,
      AiWorkPool workPool,
      @Value("${mes.ai.history-limit:20}") int historyLimit,
      @Value("${mes.ai.system-prompt:}") String systemPrompt,
      @Value("${mes.ai.control.block-ms:1000}") long controlBlockMs,
      @Value("${mes.ai.control.read-count:10}") int controlReadCount,
      @Value("${mes.ai.ask-user.timeout-ms:1800000}") long askUserTimeoutMs) {
    this.store = store;
    this.histories = histories;
    this.chatClient = builder.build();
    this.json = json;
    this.historySessions = historySessions;
    this.workOrders = workOrders;
    this.inventory = inventory;
    this.users = users;
    this.documents = documents;
    this.workPool = workPool;
    this.historyLimit = Math.max(1, Math.min(100, historyLimit));
    this.systemPrompt =
        systemPrompt == null || systemPrompt.isBlank() ? DEFAULT_PROMPT : systemPrompt.trim();
    this.controlBlockMs = Math.max(1, controlBlockMs);
    this.controlReadCount = Math.max(1, controlReadCount);
    this.askUserTimeoutMs = Math.max(60_000, askUserTimeoutMs);
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
      String target = normalized.keySet().iterator().next();
      PendingAsk pending = pendingAsks.get(target);
      if (pending == null || pending.historyId() != historyId)
        throw new IllegalStateException("pending ask_user target not found");
      String answerText =
          normalized.values().stream()
              .map(a -> !trim(a.content()).isEmpty() ? trim(a.content()) : writeJson(a.payload()))
              .reduce((a, b) -> a + "\n" + b)
              .orElseThrow();
      histories.append(historyId, identity.userId(), "user", answerText, false);
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
      completePendingAsk(historyId, target, answerText);
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
                if (!current(historyId, runtime)) return;
                completePendingAsk(historyId, event.targetId(), event.content());
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
              workOrders,
              inventory,
              users,
              documents,
              historyId,
              runtime.identity.userId(),
              runtime.identity.admin());
      StringBuilder output = new StringBuilder();
      for (String chunk :
          chatClient
              .prompt()
              .system(systemPrompt)
              .messages(history)
              .toolCallbacks(toolCallbacks(tools, runtime.identity.role(), historyId))
              .stream()
              .content()
              .toIterable()) {
        if (!current(historyId, runtime)) throw new CancellationException();
        if (chunk != null && !chunk.isEmpty()) {
          output.append(chunk);
          runtime.task.runtime().buffer().append(chunk);
          publish(historyId, "message", AGENT, chunk, null, null, null, null);
        }
      }
      histories.append(historyId, runtime.identity.userId(), "assistant", output.toString(), false);
      AiEvent done = publish(historyId, "done", AGENT, "", null, null, null, null);
      saveState(historyId, "done", AGENT, done.id(), "", List.of(), "", output.toString(), false);
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
    Set<String> common = Set.of("search_users", "search_history_file");
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
    List<ToolCallback> callbacks =
        Arrays.stream(ToolCallbacks.from(tools))
            .filter(t -> allowed.contains(t.getToolDefinition().name()))
            .map(delegate -> wrapToolCallback(delegate, historyId))
            .collect(java.util.stream.Collectors.toCollection(ArrayList::new));
    callbacks.add(wrapToolCallback(askUserCallback(historyId), historyId));
    return callbacks.toArray(ToolCallback[]::new);
  }

  private ToolCallback wrapToolCallback(ToolCallback delegate, long historyId) {
    return new ToolCallback() {
      @Override
      public ToolDefinition getToolDefinition() {
        return delegate.getToolDefinition();
      }

      @Override
      public ToolMetadata getToolMetadata() {
        return delegate.getToolMetadata();
      }

      @Override
      public String call(String input) {
        String target = UUID.randomUUID().toString();
        String name = delegate.getToolDefinition().name();
        if ("ask_user".equals(name)) return callAskUser(historyId, input, target);
        publish(
            historyId, "tool_call", AGENT, "", target, name, "running", Map.of("arguments", input));
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
    };
  }

  private ToolCallback askUserCallback(long historyId) {
    ToolDefinition definition =
        DefaultToolDefinition.builder()
            .name("ask_user")
            .description(
                "Ask the current user for missing information. Prefer concrete options when there are clear choices; use an empty options array for open-ended questions. The UI also provides a free-text input for other answers. Use this instead of printing JSON or guessing. The tool blocks until the user answers.")
            .inputSchema(
                """
                {
                  "type": "object",
                  "properties": {
                    "questions": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "question": {"type": "string"},
                          "options": {"type": "array", "items": {"type": "string"}}
                        },
                        "required": ["question", "options"]
                      }
                    },
                    "question": {"type": "string"},
                    "options": {"type": "array", "items": {"type": "string"}}
                  }
                }
                """)
            .build();
    ToolMetadata metadata = DefaultToolMetadata.builder().returnDirect(false).build();
    return new ToolCallback() {
      @Override
      public ToolDefinition getToolDefinition() {
        return definition;
      }

      @Override
      public ToolMetadata getToolMetadata() {
        return metadata;
      }

      @Override
      public String call(String input) {
        throw new IllegalStateException("ask_user must be called through AiService wrapper");
      }
    };
  }

  private String callAskUser(long historyId, String input, String target) {
    List<Question> questions = parseAskUserQuestions(input);
    if (questions.isEmpty()) throw new IllegalArgumentException("ask_user question is required");
    CompletableFuture<String> answer = new CompletableFuture<>();
    PendingAsk pending = new PendingAsk(historyId, target, answer, Instant.now().toEpochMilli());
    pendingAsks.put(target, pending);
    Map<String, Object> payload = Map.of("arguments", input, "questions", questions);
    AiEvent toolCall =
        publish(
            historyId,
            "tool_call",
            AGENT,
            questions.stream().map(Question::question).reduce((a, b) -> a + "\n" + b).orElse(""),
            target,
            "ask_user",
            "waiting",
            payload);
    saveState(
        historyId,
        "waiting_answer",
        AGENT,
        toolCall.id(),
        "",
        List.of(new PendingInterrupt(target, AGENT, toolCall.content(), writeJson(payload))),
        "",
        "",
        false);
    try {
      String value = answer.get(askUserTimeoutMs, TimeUnit.MILLISECONDS);
      saveState(historyId, "running", AGENT, toolCall.id(), "", List.of(), "", "", false);
      return value;
    } catch (TimeoutException error) {
      pendingAsks.remove(target, pending);
      publish(
          historyId,
          "tool_result",
          AGENT,
          "ask_user timed out waiting for answer",
          target,
          "ask_user",
          "error",
          null);
      throw new IllegalStateException("ask_user timed out waiting for answer", error);
    } catch (InterruptedException error) {
      Thread.currentThread().interrupt();
      pendingAsks.remove(target, pending);
      throw new CancellationException("ask_user interrupted");
    } catch (ExecutionException error) {
      pendingAsks.remove(target, pending);
      Throwable cause = error.getCause() == null ? error : error.getCause();
      if (cause instanceof RuntimeException runtime) throw runtime;
      throw new IllegalStateException(cause);
    }
  }

  private List<Question> parseAskUserQuestions(String input) {
    try {
      JsonNode root = json.readTree(empty(input));
      JsonNode source = root.has("in") && root.path("in").isObject() ? root.path("in") : root;
      List<Question> questions = new ArrayList<>();
      JsonNode array = source.path("questions");
      if (array.isArray()) {
        array.forEach(
            item -> {
              String text = trim(item.path("question").asText());
              if (!text.isEmpty())
                questions.add(new Question(text, parseOptions(item.path("options"))));
            });
      }
      String single = trim(source.path("question").asText());
      if (questions.isEmpty() && !single.isEmpty())
        questions.add(new Question(single, parseOptions(source.path("options"))));
      return questions;
    } catch (Exception error) {
      throw new IllegalArgumentException("invalid ask_user input", error);
    }
  }

  private List<String> parseOptions(JsonNode node) {
    List<String> options = new ArrayList<>();
    if (node != null && node.isArray()) {
      node.forEach(
          option -> {
            String value = trim(option.asText());
            if (!value.isEmpty()) options.add(value);
          });
    }
    return options;
  }

  private void completePendingAsk(long historyId, String targetId, String answer) {
    String target = trim(targetId);
    if (target.isEmpty()) return;
    PendingAsk pending = pendingAsks.get(target);
    if (pending == null || pending.historyId() != historyId) return;
    if (pendingAsks.remove(target, pending)) pending.answer().complete(empty(answer));
  }

  private String writeJson(Object value) {
    try {
      return json.writeValueAsString(value);
    } catch (JsonProcessingException e) {
      throw new IllegalStateException(e);
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
    pendingAsks
        .entrySet()
        .removeIf(
            entry -> {
              PendingAsk ask = entry.getValue();
              if (ask.historyId() != historyId) return false;
              ask.answer().completeExceptionally(new CancellationException("AI task cancelled"));
              return true;
            });
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
    if (value.length() <= 1200) return value;
    return value.substring(0, 1200) + "... (+" + (value.length() - 1200) + " chars)";
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

  private record PendingAsk(
      long historyId, String targetId, CompletableFuture<String> answer, long createdAt) {}

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
