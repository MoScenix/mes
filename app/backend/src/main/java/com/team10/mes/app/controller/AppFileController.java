package com.team10.mes.app.controller;

import com.team10.mes.app.service.AppFileService;
import com.team10.mes.controller.ApiResponseAdvice.ApiResponse;
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
@RequestMapping("/app/file")
public class AppFileController {
  private final AppFileService service;
  private final SessionIdentity identity;

  public AppFileController(AppFileService service, SessionIdentity identity) {
    this.service = service;
    this.identity = identity;
  }

  @PostMapping(path = "/add", consumes = "multipart/form-data")
  public ApiResponse add(
      @RequestParam long appId, @RequestParam MultipartFile file, HttpSession session)
      throws IOException {
    Long userId = identity.userId(session);
    if (userId == null) throw new UnauthorizedException();
    boolean admin = "admin".equalsIgnoreCase(identity.role(session));
    return new ApiResponse(0, service.upload(appId, userId, admin, file), "success");
  }
}
