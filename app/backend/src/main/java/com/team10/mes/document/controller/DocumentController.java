package com.team10.mes.document.controller;

import com.team10.mes.document.service.DocumentService;
import java.io.IOException;
import java.util.Map;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/mes/document")
public class DocumentController {
  private final DocumentService service;

  public DocumentController(DocumentService s) {
    service = s;
  }

  @PostMapping("/parse-pdf-to-text")
  public Map<String, Object> parse(@RequestBody Map<String, Object> r) throws IOException {
    return service.parse(num(r, "projectId"), num(r, "fileId"));
  }

  @PostMapping("/index-text-file")
  public Map<String, Object> index(@RequestBody Map<String, Object> r) throws IOException {
    return service.index(
        num(r, "projectId"), num(r, "fileId"), num(r, "minSize"), num(r, "maxSize"));
  }

  @PostMapping("/search-file")
  public Map<String, Object> search(@RequestBody Map<String, Object> r) {
    return service.search(
        num(r, "projectId"),
        num(r, "fileId"),
        String.valueOf(r.getOrDefault("query", "")),
        num(r, "topK"));
  }

  @PostMapping("/delete-project-file-data")
  public Map<String, Object> delete(@RequestBody Map<String, Object> r) throws IOException {
    return service.delete(num(r, "projectId"));
  }

  private long num(Map<String, Object> r, String k) {
    Object v = r.get(k);
    return v == null ? 0 : Long.parseLong(v.toString());
  }
}
