package com.team10.mes.document.utils;

import java.util.*;

public final class DocumentText {
  private DocumentText() {}

  public record Chunk(int start, int end, String text) {}

  public record Parent(long id, String content) {}

  public record Split(List<Chunk> chunks, List<List<Long>> parentIds, List<Parent> parents) {}

  public static String clean(String text) {
    text = text.replace("\r\n", "\n").replace('\r', '\n');
    StringBuilder b = new StringBuilder(text.length());
    text.codePoints()
        .filter(c -> c == '\n' || c == '\t' || !Character.isISOControl(c))
        .forEach(b::appendCodePoint);
    String[] lines = b.toString().split("\n", -1);
    for (int i = 0; i < lines.length; i++) lines[i] = lines[i].stripTrailing();
    return String.join("\n", lines).replaceAll("\n{3,}", "\n\n").trim();
  }

  public static Split splitWithParents(String text, long min, long max) {
    List<Chunk> chunks = split(text, (int) min, (int) max);
    List<List<Long>> ids = new ArrayList<>();
    for (int i = 0; i < chunks.size(); i++) ids.add(new ArrayList<>());
    List<Parent> parents = new ArrayList<>();
    long id = 1;
    for (int start = 0; start < chunks.size(); start += 3) {
      int end = Math.min(chunks.size(), start + 5);
      if (start >= end) break;
      List<String> parts = new ArrayList<>();
      for (int i = start; i < end; i++) {
        parts.add(chunks.get(i).text().trim());
        ids.get(i).add(id);
      }
      parents.add(new Parent(id++, String.join("\n\n", parts)));
      if (end == chunks.size()) break;
    }
    return new Split(chunks, ids.stream().map(List::copyOf).toList(), parents);
  }

  public static List<Chunk> split(String text, int min, int max) {
    if (max <= 0) max = 400;
    if (min <= 0) min = 200;
    if (min > max) min = max;
    int[] runes = text.codePoints().toArray();
    List<Chunk> out = new ArrayList<>();
    for (int cursor = 0; cursor < runes.length; ) {
      int end = nextEnd(runes, cursor, runes.length, min, max, 0);
      if (end <= cursor) end = Math.min(cursor + max, runes.length);
      out.add(new Chunk(cursor, end, new String(runes, cursor, end - cursor)));
      cursor = end;
    }
    return out;
  }

  private static int nextEnd(int[] r, int start, int limit, int min, int max, int level) {
    if (limit - start <= max) return limit;
    if (level >= 3) return Math.min(start + max, limit);
    int current = start, unit = start;
    while (unit < limit) {
      int scan = Math.min(limit, start + max + 1), end = scanEnd(r, unit, scan, level);
      if (end <= unit) end = unit + 1;
      if (end - start > max) {
        if (current - start >= min) return current;
        return nextEnd(r, start, end, min, max, level + 1);
      }
      current = end;
      unit = end;
    }
    return current > start ? current : Math.min(start + max, limit);
  }

  private static int scanEnd(int[] r, int start, int limit, int level) {
    for (int i = start; i < limit; i++)
      if (boundary(r[i], level)) {
        int end = i + 1;
        if (level == 0) while (end < limit && r[end] == '\n') end++;
        return end;
      }
    return limit;
  }

  private static boolean boundary(int c, int level) {
    return switch (level) {
      case 0 -> c == '\n';
      case 1 -> "。！？!?".indexOf(c) >= 0;
      case 2 -> "，,；;、：:".indexOf(c) >= 0;
      default -> false;
    };
  }

  public static List<Long> fuse(List<List<List<Long>>> resultSets, int topK) {
    Map<Long, Double> scores = new HashMap<>();
    for (List<List<Long>> results : resultSets)
      for (int rank = 0; rank < results.size(); rank++)
        for (long id : results.get(rank))
          if (id > 0) scores.merge(id, 1d / (60d + rank + 1), Double::sum);
    return scores.entrySet().stream()
        .sorted(
            Comparator.<Map.Entry<Long, Double>>comparingDouble(Map.Entry::getValue)
                .reversed()
                .thenComparingLong(Map.Entry::getKey))
        .limit(topK > 0 ? topK : 5)
        .map(Map.Entry::getKey)
        .toList();
  }
}
