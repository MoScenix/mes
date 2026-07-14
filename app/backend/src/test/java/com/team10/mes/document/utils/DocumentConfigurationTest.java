package com.team10.mes.document.utils;

import static org.junit.jupiter.api.Assertions.*;

import java.nio.file.*;
import org.junit.jupiter.api.Test;

class DocumentConfigurationTest {
  @Test
  void defaultsExistOnlyInApplicationConfiguration() throws Exception {
    String yaml = Files.readString(Path.of("src/main/resources/application.yml"));
    assertTrue(yaml.contains("root: ${DOCUMENT_ROOT:static/document}"));
    String service =
        Files.readString(
            Path.of("src/main/java/com/team10/mes/document/service/DocumentService.java"));
    String store =
        Files.readString(
            Path.of("src/main/java/com/team10/mes/document/utils/DocumentIndexStore.java"));
    assertFalse(service.contains("@Value"));
    assertFalse(store.contains("@Value"));
    assertFalse(service.contains("static/document"));
    assertFalse(store.contains("127.0.0.1"));
  }
}
