package com.team10.mes.app.service;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.team10.mes.app.dal.AppMapper;
import com.team10.mes.document.service.DocumentService;
import com.team10.mes.document.utils.DocumentProperties;
import com.team10.mes.history.model.HistoryMessage;
import com.team10.mes.history.service.HistoryMessageService;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.ArgumentCaptor;
import org.springframework.mock.web.MockMultipartFile;

class AppFileServiceTest {
  @TempDir Path root;

  private final AppService apps = mock(AppService.class);
  private final DocumentService documents = mock(DocumentService.class);
  private final HistoryMessageService histories = mock(HistoryMessageService.class);
  private final ObjectMapper json = new ObjectMapper();
  private AppFileProperties fileProperties;
  private AppFileService service;

  @BeforeEach
  void setUp() {
    DocumentProperties documentProperties = new DocumentProperties();
    documentProperties.setRoot(root.toString());
    fileProperties = new AppFileProperties();
    fileProperties.setBigThresholdBytes(5);
    fileProperties.setChunkMinSize(11);
    fileProperties.setChunkMaxSize(22);
    service =
        new AppFileService(apps, documents, histories, fileProperties, documentProperties, json);
    when(apps.get(7)).thenReturn(new AppMapper.AppRow(7L, "chat", 42L, null, null));
    HistoryMessage saved = new HistoryMessage();
    saved.setId(99L);
    when(histories.append(anyLong(), anyLong(), anyString(), anyString(), anyBoolean()))
        .thenReturn(saved);
  }

  @Test
  void savesSmallTxtAndAppendsOriginalFileMetadata() throws Exception {
    fileProperties.setBigThresholdBytes(100);
    var upload = new MockMultipartFile("file", "../note.TXT", "text/plain", "hello".getBytes());

    assertEquals("99", service.upload(7, 42, false, upload));

    verifyNoInteractions(documents);
    ArgumentCaptor<String> content = ArgumentCaptor.forClass(String.class);
    verify(histories).append(eq(7L), eq(42L), eq("user"), content.capture(), eq(true));
    JsonNode metadata = json.readTree(content.getValue());
    assertEquals("note.TXT", metadata.get("filename").asText());
    assertEquals("note.TXT", metadata.get("textFilename").asText());
    assertEquals(5, metadata.get("size").asLong());
    assertFalse(metadata.get("isBig").asBoolean());
    assertFalse(metadata.has("chunkCount"));
    Path saved =
        Files.walk(root)
            .filter(path -> path.getFileName().toString().equals("note.TXT"))
            .findFirst()
            .orElseThrow();
    assertEquals("hello", Files.readString(saved));
  }

  @Test
  void indexesBigTxtWithConfiguredChunkSizes() throws Exception {
    when(documents.index(eq(7L), anyLong(), eq(11L), eq(22L)))
        .thenReturn(Map.of("chunkCount", 3L, "parentCount", 2L));
    var upload = new MockMultipartFile("file", "large.txt", "text/plain", "123456".getBytes());

    service.upload(7, 42, false, upload);

    verify(documents).index(eq(7L), anyLong(), eq(11L), eq(22L));
    ArgumentCaptor<String> content = ArgumentCaptor.forClass(String.class);
    verify(histories).append(eq(7L), eq(42L), eq("user"), content.capture(), eq(true));
    JsonNode metadata = json.readTree(content.getValue());
    assertTrue(metadata.get("isBig").asBoolean());
    assertEquals(3, metadata.get("chunkCount").asLong());
    assertEquals(2, metadata.get("parentCount").asLong());
  }

  @Test
  void pdfUsesParsedTextMetadataBeforeBigFileDecision() throws Exception {
    when(documents.parse(eq(7L), anyLong()))
        .thenReturn(Map.of("textFilename", "manual.txt", "textSize", 6L));
    when(documents.index(eq(7L), anyLong(), eq(11L), eq(22L)))
        .thenReturn(Map.of("chunkCount", 1L, "parentCount", 1L));
    var upload = new MockMultipartFile("file", "manual.pdf", "application/pdf", "pdf".getBytes());

    service.upload(7, 42, false, upload);

    verify(documents).parse(eq(7L), anyLong());
    verify(documents).index(eq(7L), anyLong(), eq(11L), eq(22L));
  }

  @Test
  void rejectsNonOwnerBeforeWritingFile() {
    var upload =
        new MockMultipartFile(
            "file", "note.txt", "text/plain", "hello".getBytes(StandardCharsets.UTF_8));

    assertThrows(IllegalStateException.class, () -> service.upload(7, 8, false, upload));
    verifyNoInteractions(documents, histories);
    assertDoesNotThrow(() -> assertEquals(0, Files.list(root).count()));
  }

  @Test
  void acceptsAdminAndRejectsUnsupportedExtension() throws Exception {
    fileProperties.setBigThresholdBytes(100);
    var txt = new MockMultipartFile("file", "ok.txt", "text/plain", new byte[] {1});
    assertEquals("99", service.upload(7, 8, true, txt));

    var executable =
        new MockMultipartFile("file", "bad.exe", "application/octet-stream", new byte[] {1});
    assertThrows(IllegalArgumentException.class, () -> service.upload(7, 42, false, executable));
  }

  @Test
  void deletesRemoteIndexAndConfiguredProjectDirectory() throws Exception {
    Path nested = root.resolve("7/123/parents");
    Files.createDirectories(nested);
    Files.writeString(nested.resolve("1.txt"), "parent");

    service.deleteProjectFiles(7);

    verify(documents).delete(7);
    assertFalse(Files.exists(root.resolve("7")));
  }
}
