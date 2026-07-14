package com.team10.mes.history.controller;

import com.team10.mes.controller.ApiResponseAdvice.ApiResponse;
import com.team10.mes.history.service.HistoryFileService;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import java.io.IOException;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/history/file")
public class HistoryFileController {
  private final HistoryFileService service;
  private final SessionIdentity identity;

  public HistoryFileController(HistoryFileService service, SessionIdentity identity) {
    this.service = service;
    this.identity = identity;
  }

  @PostMapping(path = "/add", consumes = "multipart/form-data")
  public ApiResponse add(
      @RequestParam long historyId, @RequestParam MultipartFile file, HttpSession session)
      throws IOException {
    Long userId = identity.userId(session);
    if (userId == null) throw new UnauthorizedException();
    return new ApiResponse(
        0, service.upload(historyId, userId, identity.role(session), file), "success");
  }
}
