package com.team10.mes.document.service;

import com.team10.mes.document.utils.DocumentIndexStore;
import com.team10.mes.document.utils.DocumentProperties;
import com.team10.mes.document.utils.DocumentText;
import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.*;
import java.util.stream.Stream;
import org.apache.pdfbox.Loader;
import org.apache.pdfbox.text.PDFTextStripper;
import org.springframework.stereotype.Service;

@Service
public class DocumentService {
  private final DocumentIndexStore store;
  private final Path root;

  public DocumentService(DocumentIndexStore store, DocumentProperties properties) {
    this.store = store;
    this.root = Paths.get(properties.getRoot()).normalize();
  }

  private Path dir(long history, long file) {
    return root.resolve(Long.toString(history)).resolve(Long.toString(file));
  }

  private Path latest(Path d, String ext) throws IOException {
    try (Stream<Path> s = Files.list(d)) {
      return s.filter(Files::isRegularFile)
          .filter(p -> p.getFileName().toString().toLowerCase().endsWith(ext))
          .max(
              Comparator.comparingLong(
                  p -> {
                    try {
                      return Files.getLastModifiedTime(p).toMillis();
                    } catch (IOException e) {
                      return 0;
                    }
                  }))
          .orElseThrow(() -> new FileNotFoundException("no " + ext + " file in " + d));
    }
  }

  public Map<String, Object> parse(long history, long file) throws IOException {
    Path pdf = latest(dir(history, file), ".pdf"),
        txt = pdf.resolveSibling(pdf.getFileName().toString().replaceFirst("(?i)\\.pdf$", ".txt"));
    String text = poppler(pdf, txt);
    if (text == null)
      try (var doc = Loader.loadPDF(pdf.toFile())) {
        text = DocumentText.clean(new PDFTextStripper().getText(doc));
      }
    Files.writeString(txt, text, StandardCharsets.UTF_8);
    return Map.of(
        "fileId",
        file,
        "textFilename",
        txt.getFileName().toString(),
        "textSize",
        (long) text.getBytes(StandardCharsets.UTF_8).length);
  }

  public Map<String, Object> index(long history, long file, long min, long max) throws IOException {
    Path txt = latest(dir(history, file), ".txt");
    DocumentText.Split split =
        DocumentText.splitWithParents(DocumentText.clean(Files.readString(txt)), min, max);
    Path parents = txt.getParent().resolve("parents");
    Files.createDirectories(parents);
    for (var parent : split.parents())
      Files.writeString(parents.resolve(parent.id() + ".txt"), parent.content());
    List<DocumentIndexStore.Child> children = new ArrayList<>();
    for (int i = 0; i < split.chunks().size(); i++)
      children.add(
          new DocumentIndexStore.Child(
              history, file, i + 1, split.parentIds().get(i), split.chunks().get(i).text()));
    store.index(children);
    return Map.of(
        "fileId",
        file,
        "chunkCount",
        (long) split.chunks().size(),
        "parentCount",
        (long) split.parents().size());
  }

  public Map<String, Object> search(long history, long file, String query, long topK) {
    return Map.of("parentIds", store.search(history, file, query, topK));
  }

  public Map<String, Object> delete(long history) {
    store.deleteHistory(history);
    return Map.of("success", true);
  }

  private String poppler(Path pdf, Path txt) {
    try {
      Process p =
          new ProcessBuilder(
                  "pdftotext", "-layout", "-enc", "UTF-8", pdf.toString(), txt.toString())
              .redirectErrorStream(true)
              .start();
      String output = new String(p.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
      if (p.waitFor() != 0) return null;
      return DocumentText.clean(Files.readString(txt));
    } catch (Exception ignored) {
      return null;
    }
  }
}
