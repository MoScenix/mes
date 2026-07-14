package com.team10.mes.ai.controller;

import com.team10.mes.ai.service.AiService;
import com.team10.mes.ai.service.AiService.Identity;
import com.team10.mes.user.service.SessionIdentity;
import com.team10.mes.user.service.UnauthorizedException;
import jakarta.servlet.http.HttpSession;
import java.util.Map;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/history/ai")
public class AiController {
  private final AiService service;
  private final SessionIdentity sessionIdentity;

  public AiController(AiService service, SessionIdentity sessionIdentity) {
    this.service = service;
    this.sessionIdentity = sessionIdentity;
  }

  @PostMapping("/submit")
  public Response<Boolean> submit(@RequestBody SubmitRequest request, HttpSession session) {
    service.submit(request.historyId(), request.message(), identity(session));
    return Response.ok(true);
  }

  @PostMapping("/push")
  public Response<String> push(@RequestBody ControlRequest request, HttpSession session) {
    return Response.ok(service.push(request.historyId(), request.content(), identity(session)));
  }

  @PostMapping("/answer")
  public Response<Boolean> answer(@RequestBody ControlRequest request, HttpSession session) {
    service.answer(request.historyId(), request.answers(), identity(session));
    return Response.ok(true);
  }

  @PostMapping("/cancel")
  public Response<String> cancel(@RequestBody ControlRequest request, HttpSession session) {
    return Response.ok(service.cancel(request.historyId(), request.reason(), identity(session)));
  }

  @GetMapping("/state")
  public Response<AiService.AiState> state(@RequestParam long historyId, HttpSession session) {
    service.authorize(historyId, identity(session));
    return Response.ok(service.state(historyId));
  }

  @GetMapping("/events")
  public Response<AiService.EventPage> events(
      @RequestParam long historyId,
      @RequestParam(defaultValue = "0") String lastId,
      @RequestParam(defaultValue = "30000") long blockMs,
      @RequestParam(defaultValue = "50") int count,
      HttpSession session) {
    service.authorize(historyId, identity(session));
    return Response.ok(service.events(historyId, lastId, blockMs, count));
  }

  public record SubmitRequest(long historyId, String message) {}

  public record ControlRequest(
      long historyId, String content, String reason, Map<String, Answer> answers) {}

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
