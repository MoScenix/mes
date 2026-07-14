package com.team10.mes.document.utils;

import static org.junit.jupiter.api.Assertions.*;

import java.util.*;
import org.junit.jupiter.api.Test;

class DocumentTextTest {
  @Test
  void defaultsAndRuneOffsetsMatchOriginalSplitter() {
    String text = "甲".repeat(199) + "。" + "乙".repeat(250);
    var chunks = DocumentText.split(text, 0, 0);
    assertEquals(2, chunks.size());
    assertEquals(200, chunks.getFirst().end());
    assertEquals(
        200, chunks.getFirst().text().codePointCount(0, chunks.getFirst().text().length()));
    assertEquals(450, chunks.getLast().end());
  }

  @Test
  void paragraphHasPriorityAndMinGreaterThanMaxIsClamped() {
    String text = "a".repeat(8) + "\n" + "b".repeat(12);
    var chunks = DocumentText.split(text, 20, 10);
    assertEquals(List.of(10, 20, 21), chunks.stream().map(DocumentText.Chunk::end).toList());
  }

  @Test
  void parentsOverlapFiveByThreeAndChildrenKeepEveryParent() {
    StringBuilder text = new StringBuilder();
    for (int i = 0; i < 8; i++) text.append((char) ('a' + i)).append(".".repeat(9));
    var split = DocumentText.splitWithParents(text.toString(), 10, 10);
    assertEquals(8, split.chunks().size());
    assertEquals(2, split.parents().size());
    assertEquals(List.of(1L, 2L), split.parentIds().get(3));
    assertEquals(List.of(1L, 2L), split.parentIds().get(4));
    assertEquals(List.of(2L), split.parentIds().get(7));
  }

  @Test
  void rrfDeduplicatesOrdersByScoreThenParentIdAndDefaultsTopK() {
    var es = List.of(List.of(2L, 1L), List.of(3L));
    var mv = List.of(List.of(3L), List.of(2L));
    assertEquals(List.of(2L, 3L, 1L), DocumentText.fuse(List.of(es, mv), 5));
    assertEquals(3, DocumentText.fuse(List.of(es, mv), 0).size());
  }

  @Test
  void cleanMatchesControlAndWhitespaceRules() {
    assertEquals("a\n\n b", DocumentText.clean(" a\r\n\u0000\n \n b  "));
  }
}
