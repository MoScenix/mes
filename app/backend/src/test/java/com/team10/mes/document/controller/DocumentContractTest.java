package com.team10.mes.document.controller;

import static org.junit.jupiter.api.Assertions.*;

import java.lang.reflect.Method;
import java.util.*;
import org.junit.jupiter.api.Test;
import org.springframework.web.bind.annotation.*;

class DocumentContractTest {
  @Test
  void exposesExactlyTheFourOriginalPostOperations() {
    RequestMapping root = DocumentController.class.getAnnotation(RequestMapping.class);
    assertArrayEquals(new String[] {"/mes/document"}, root.value());
    Set<String> paths = new TreeSet<>();
    for (Method m : DocumentController.class.getDeclaredMethods()) {
      PostMapping p = m.getAnnotation(PostMapping.class);
      if (p != null) paths.addAll(List.of(p.value()));
    }
    assertEquals(
        Set.of(
            "/parse-pdf-to-text", "/index-text-file", "/search-file", "/delete-project-file-data"),
        paths);
  }
}
