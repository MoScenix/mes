package com.team10.mes.ai.service;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;

import com.team10.mes.document.service.DocumentService;
import com.team10.mes.inventory.service.InventoryService;
import com.team10.mes.user.service.UserService;
import com.team10.mes.workorder.service.WorkOrderService;
import java.util.Arrays;
import java.util.Set;
import java.util.stream.Collectors;
import org.junit.jupiter.api.Test;
import org.springframework.ai.support.ToolCallbacks;
import org.springframework.ai.tool.ToolCallback;

class MesAiToolsTest {
  @Test
  void springAiDiscoversEveryAnnotatedMesTool() {
    MesAiTools tools =
        new MesAiTools(
            mock(WorkOrderService.class),
            mock(InventoryService.class),
            mock(UserService.class),
            mock(DocumentService.class),
            42L,
            1L,
            true);

    ToolCallback[] callbacks = ToolCallbacks.from(tools);
    Set<String> names =
        Arrays.stream(callbacks)
            .map(callback -> callback.getToolDefinition().name())
            .collect(Collectors.toSet());

    assertEquals(19, callbacks.length);
    assertTrue(names.contains("list_work_orders"));
    assertTrue(names.contains("search_users"));
    assertTrue(names.contains("search_history_file"));
    assertTrue(names.contains("inventory_check"));
    assertTrue(names.contains("search_processes_by_item"));
  }
}
