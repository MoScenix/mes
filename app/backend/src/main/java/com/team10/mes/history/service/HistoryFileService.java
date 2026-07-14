package com.team10.mes.history.service;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.team10.mes.document.service.DocumentService;
import com.team10.mes.document.utils.DocumentProperties;
import com.team10.mes.history.model.HistoryMessage;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Locale;
import java.util.Map;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Stream;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

@Service
public class HistoryFileService {
  private static final AtomicLong LAST_FILE_ID =
      new AtomicLong(System.currentTimeMillis() * 1_000_000L);

  private final HistorySessionService historySessions;
  private final DocumentService documents;
  private final HistoryMessageService histories;
  private final HistoryFileProperties fileProperties;
  private final ObjectMapper json;
  private final Path documentRoot;

  public HistoryFileService(
      HistorySessionService historySessions,
      DocumentService documents,
      HistoryMessageService histories,
      HistoryFileProperties fileProperties,
      DocumentProperties documentProperties,
      ObjectMapper json) {
    this.historySessions = historySessions;
    this.documents = documents;
    this.histories = histories;
    this.fileProperties = fileProperties;
    this.documentRoot = Paths.get(documentProperties.getRoot()).normalize();
    this.json = json;
  }

  public String upload(long historyId, long userId, String role, MultipartFile upload)
      throws IOException {
    requireOwner(historyId, userId, role);
    if (upload == null || upload.isEmpty()) {
      throw new IllegalArgumentException("missing upload file");
    }

    String filename = safeFilename(upload.getOriginalFilename());
    if (filename.isEmpty()) throw new IllegalArgumentException("missing upload filename");
    String extension = extension(filename);
    if (!extension.equals(".pdf") && !extension.equals(".txt")) {
      throw new IllegalArgumentException("unsupported file type: " + extension);
    }

    long fileId = nextFileId();
    Path directory = documentRoot.resolve(Long.toString(historyId)).resolve(Long.toString(fileId));
    Files.createDirectories(directory);
    Path savedFile = directory.resolve(filename).normalize();
    if (!savedFile.getParent().equals(directory.normalize())) {
      throw new IllegalArgumentException("invalid upload filename");
    }
    upload.transferTo(savedFile);

    long size = Files.size(savedFile);
    String textFilename = filename;
    long textSize = size;
    if (extension.equals(".pdf")) {
      Map<String, Object> parsed = documents.parse(historyId, fileId);
      textFilename = stringValue(parsed, "textFilename");
      textSize = longValue(parsed, "textSize");
    }

    boolean big =
        fileProperties.getBigThresholdBytes() > 0
            && textSize > fileProperties.getBigThresholdBytes();
    long chunkCount = 0;
    long parentCount = 0;
    if (big) {
      Map<String, Object> indexed =
          documents.index(
              historyId,
              fileId,
              fileProperties.getChunkMinSize(),
              fileProperties.getChunkMaxSize());
      chunkCount = longValue(indexed, "chunkCount");
      parentCount = longValue(indexed, "parentCount");
    }

    FileMessageContent content =
        new FileMessageContent(
            fileId,
            filename,
            upload.getContentType() == null ? "" : upload.getContentType(),
            size,
            textFilename,
            textSize,
            big,
            chunkCount,
            parentCount);
    HistoryMessage message =
        histories.append(historyId, userId, "user", writeContent(content), true);
    historySessions.touch(historyId);
    return Long.toString(message.getId());
  }

  public void deleteHistoryFiles(long historyId) throws IOException {
    documents.delete(historyId);
    Path projectDirectory = documentRoot.resolve(Long.toString(historyId)).normalize();
    if (!projectDirectory.startsWith(documentRoot) || !Files.exists(projectDirectory)) return;
    try (Stream<Path> paths = Files.walk(projectDirectory)) {
      for (Path path : paths.sorted((left, right) -> right.compareTo(left)).toList()) {
        Files.deleteIfExists(path);
      }
    }
  }

  private void requireOwner(long historyId, long userId, String role) {
    if (userId <= 0) throw new IllegalStateException("unauthorized");
    historySessions.authorize(historyId, userId, role);
  }

  private String writeContent(FileMessageContent content) {
    try {
      return json.writeValueAsString(content);
    } catch (JsonProcessingException e) {
      throw new IllegalStateException("serialize file message failed", e);
    }
  }

  static String safeFilename(String original) {
    if (original == null) return "";
    String cleaned = original.trim().replace("\0", "").replace('\\', '/');
    int slash = cleaned.lastIndexOf('/');
    return slash < 0 ? cleaned : cleaned.substring(slash + 1);
  }

  private static String extension(String filename) {
    int dot = filename.lastIndexOf('.');
    return dot < 0 ? "" : filename.substring(dot).toLowerCase(Locale.ROOT);
  }

  private static long nextFileId() {
    long now = System.currentTimeMillis() * 1_000_000L;
    return LAST_FILE_ID.updateAndGet(previous -> Math.max(now, previous + 1));
  }

  private static long longValue(Map<String, Object> values, String key) {
    Object value = values.get(key);
    if (value instanceof Number number) return number.longValue();
    throw new IllegalStateException("document response missing " + key);
  }

  private static String stringValue(Map<String, Object> values, String key) {
    Object value = values.get(key);
    if (value instanceof String text && !text.isBlank()) return text;
    throw new IllegalStateException("document response missing " + key);
  }

  public record FileMessageContent(
      long fileId,
      String filename,
      String contentType,
      long size,
      String textFilename,
      long textSize,
      boolean isBig,
      @JsonInclude(JsonInclude.Include.NON_DEFAULT) long chunkCount,
      @JsonInclude(JsonInclude.Include.NON_DEFAULT) long parentCount) {}
}
