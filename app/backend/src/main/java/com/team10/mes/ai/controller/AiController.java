package com.team10.mes.ai.controller;

import com.team10.mes.ai.service.AiService;
import com.team10.mes.ai.service.AiService.Identity;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import java.util.Map;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/app/ai")
public class AiController {
  private final AiService service;
  private final SessionIdentity sessionIdentity;

  public AiController(AiService service, SessionIdentity sessionIdentity) {
    this.service = service;
    this.sessionIdentity = sessionIdentity;
  }

  @PostMapping("/submit")
  public Response<Boolean> submit(@RequestBody SubmitRequest request, HttpSession session) {
    service.submit(request.appId(), request.message(), identity(session));
    return Response.ok(true);
  }

  @PostMapping("/push")
  public Response<String> push(@RequestBody ControlRequest request, HttpSession session) {
    return Response.ok(service.push(request.appId(), request.content(), identity(session)));
  }

  @PostMapping("/answer")
  public Response<Boolean> answer(@RequestBody ControlRequest request, HttpSession session) {
    service.answer(request.appId(), request.answers(), identity(session));
    return Response.ok(true);
  }

  @PostMapping("/cancel")
  public Response<String> cancel(@RequestBody ControlRequest request, HttpSession session) {
    return Response.ok(service.cancel(request.appId(), request.reason(), identity(session)));
  }

  @GetMapping("/state")
  public Response<AiService.AiState> state(@RequestParam long appId, HttpSession session) {
    service.authorize(appId, identity(session));
    return Response.ok(service.state(appId));
  }

  @GetMapping("/events")
  public Response<AiService.EventPage> events(
      @RequestParam long appId,
      @RequestParam(defaultValue = "0") String lastId,
      @RequestParam(defaultValue = "30000") long blockMs,
      @RequestParam(defaultValue = "50") int count,
      HttpSession session) {
    service.authorize(appId, identity(session));
    return Response.ok(service.events(appId, lastId, blockMs, count));
  }

  public record SubmitRequest(long appId, String message) {}

  public record ControlRequest(
      long appId, String content, String reason, Map<String, Answer> answers) {}

  public record Answer(String content, Map<String, Object> payload) {}

  private Identity identity(HttpSession session) {
    Long id = sessionIdentity.userId(session);
    if (id == null) throw new UnauthorizedException();
    return new Identity(id, sessionIdentity.role(session));
  }

  public record Response<T>(int code, T data, String message) {
    public static <T> Response<T> ok(T data) {
      return new Response<>(0, data, "success");
    }
  }

  @ExceptionHandler({IllegalArgumentException.class, IllegalStateException.class})
  public Response<Object> requestError(RuntimeException error) {
    return new Response<>(1, null, error.getMessage());
  }
}
