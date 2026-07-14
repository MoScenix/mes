package com.team10.mes.ai.service;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.team10.mes.ai.controller.AiController.Answer;
import com.team10.mes.ai.state.RedisAiStore;
import com.team10.mes.ai.workpool.AiWorkPool;
import com.team10.mes.document.service.DocumentService;
import com.team10.mes.history.service.HistoryMessageService;
import com.team10.mes.history.service.HistorySessionService;
import com.team10.mes.inventory.service.InventoryService;
import com.team10.mes.user.service.UserService;
import com.team10.mes.workorder.service.WorkOrderService;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.tool.ToolCallback;

class AiServiceToolContextTest {
  @Test
  void toolCallbackKeepsAppContextWhenSpringAiExecutesOnAnotherThread() throws Exception {
    RedisAiStore store = mock(RedisAiStore.class);
    UserService users = mock(UserService.class);
    AiWorkPool pool = new AiWorkPool(1, 1, 4, 1, 0, 1);
    when(store.addEvent(eq(42L), any())).thenReturn("1-0", "2-0");
    when(users.page(any())).thenReturn(Map.of("records", java.util.List.of(), "total", 0));
    AiService service = service(store, users, pool);
    MesAiTools tools =
        new MesAiTools(
            mock(WorkOrderService.class),
            mock(InventoryService.class),
            users,
            mock(DocumentService.class),
            42L,
            7L,
            false);
    ToolCallback callback =
        java.util.Arrays.stream(service.toolCallbacks(tools, "worker", 42L))
            .filter(tool -> "search_users".equals(tool.getToolDefinition().name()))
            .findFirst()
            .orElseThrow();

    String result;
    try {
      result =
          CompletableFuture.supplyAsync(
                  () ->
                      callback.call(
                          "{\"in\":{\"id\":0,\"name\":\"\",\"account\":\"\",\"role\":\"\",\"page_size\":10}}"))
              .get(5, TimeUnit.SECONDS);
    } finally {
      pool.close();
    }

    assertEquals(
        new ObjectMapper().readTree("{\"total\":0,\"records\":[]}"),
        new ObjectMapper().readTree(result));
    ArgumentCaptor<AiService.AiEvent> events = ArgumentCaptor.forClass(AiService.AiEvent.class);
    verify(store, org.mockito.Mockito.times(2)).addEvent(eq(42L), events.capture());
    assertEquals("tool_call", events.getAllValues().get(0).type());
    assertEquals("tool_result", events.getAllValues().get(1).type());
    assertEquals("success", events.getAllValues().get(1).status());
  }

  @Test
  void askUserToolBlocksUntilAnswerEndpointCompletesPendingChannel() throws Exception {
    RedisAiStore store = mock(RedisAiStore.class);
    UserService users = mock(UserService.class);
    AiWorkPool pool = new AiWorkPool(1, 1, 4, 1, 0, 1);
    AtomicInteger ids = new AtomicInteger();
    when(store.addEvent(eq(42L), any())).thenAnswer(invocation -> ids.incrementAndGet() + "-0");
    when(store.addControl(eq(42L), any())).thenReturn("answer-0");
    AiService service = service(store, users, pool);
    MesAiTools tools =
        new MesAiTools(
            mock(WorkOrderService.class),
            mock(InventoryService.class),
            users,
            mock(DocumentService.class),
            42L,
            7L,
            false);
    ToolCallback askUser =
        java.util.Arrays.stream(service.toolCallbacks(tools, "worker", 42L))
            .filter(tool -> "ask_user".equals(tool.getToolDefinition().name()))
            .findFirst()
            .orElseThrow();

    CompletableFuture<String> result =
        CompletableFuture.supplyAsync(
            () -> askUser.call("{\"question\":\"流转单的发起人是谁？\",\"options\":[\"root\"]}"));
    ArgumentCaptor<AiService.AiEvent> events = ArgumentCaptor.forClass(AiService.AiEvent.class);
    verify(store, org.mockito.Mockito.timeout(3000).atLeast(1)).addEvent(eq(42L), events.capture());
    String target =
        events.getAllValues().stream()
            .filter(event -> "tool_call".equals(event.type()) && "ask_user".equals(event.name()))
            .findFirst()
            .orElseThrow()
            .targetId();

    service.answer(
        42L,
        Map.of(target, new Answer("purchase", Map.of())),
        new AiService.Identity(7L, "worker"));

    try {
      assertEquals("purchase", result.get(5, TimeUnit.SECONDS));
    } finally {
      pool.close();
    }
  }

  private AiService service(RedisAiStore store, UserService users, AiWorkPool pool) {
    ChatClient.Builder builder = mock(ChatClient.Builder.class);
    when(builder.build()).thenReturn(mock(ChatClient.class));
    return new AiService(
        store,
        mock(HistoryMessageService.class),
        builder,
        new ObjectMapper(),
        mock(HistorySessionService.class),
        mock(WorkOrderService.class),
        mock(InventoryService.class),
        users,
        mock(DocumentService.class),
        pool,
        20,
        "test",
        1000,
        10,
        1_800_000);
  }
}
