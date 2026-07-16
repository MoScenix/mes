package com.team10.mes.inventory.controller;

import com.team10.mes.inventory.service.DashboardService;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import java.util.Map;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/mes/dashboard")
public class DashboardController {
  private final DashboardService service;
  private final SessionIdentity identity;

  public DashboardController(DashboardService service, SessionIdentity identity) {
    this.service = service;
    this.identity = identity;
  }

  @GetMapping("/overview")
  public Map<String, Object> overview(HttpSession session) {
    Long userId = identity.userId(session);
    if (userId == null || userId <= 0) throw new UnauthorizedException();
    return service.overview();
  }
}
