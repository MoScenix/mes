package com.team10.mes.controller;

import com.team10.mes.user.service.UnauthorizedException;
import java.io.IOException;
import java.nio.file.NoSuchFileException;
import java.util.Arrays;
import java.util.Map;
import java.util.NoSuchElementException;
import org.springframework.core.MethodParameter;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.MediaType;
import org.springframework.http.converter.HttpMessageConverter;
import org.springframework.http.server.ServerHttpRequest;
import org.springframework.http.server.ServerHttpResponse;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.servlet.mvc.method.annotation.ResponseBodyAdvice;

@RestControllerAdvice(basePackages = "com.team10.mes")
public class ApiResponseAdvice implements ResponseBodyAdvice<Object> {
  public record ApiResponse(int code, Object data, String message) {}

  @Override
  public boolean supports(
      MethodParameter returnType, Class<? extends HttpMessageConverter<?>> converterType) {
    return !ApiResponse.class.isAssignableFrom(returnType.getParameterType());
  }

  @Override
  public Object beforeBodyWrite(
      Object body,
      MethodParameter returnType,
      MediaType contentType,
      Class<? extends HttpMessageConverter<?>> converterType,
      ServerHttpRequest request,
      ServerHttpResponse response) {
    if (isResponseEnvelope(body)) return body;
    return new ApiResponse(0, body, "ok");
  }

  private boolean isResponseEnvelope(Object body) {
    if (body instanceof ApiResponse) return true;
    if (body instanceof Map<?, ?> map)
      return map.containsKey("code") && map.containsKey("data") && map.containsKey("message");
    if (body == null || !body.getClass().isRecord()) return false;
    var names =
        Arrays.stream(body.getClass().getRecordComponents())
            .map(component -> component.getName())
            .toList();
    return names.contains("code") && names.contains("data") && names.contains("message");
  }

  @ExceptionHandler(IllegalArgumentException.class)
  public ApiResponse badRequest(IllegalArgumentException error) {
    return new ApiResponse(40000, null, error.getMessage());
  }

  @ExceptionHandler(NoSuchElementException.class)
  public ApiResponse notFound(NoSuchElementException error) {
    return new ApiResponse(40400, null, error.getMessage());
  }

  @ExceptionHandler(NoSuchFileException.class)
  public ApiResponse fileNotFound(NoSuchFileException error) {
    return new ApiResponse(40400, null, "file not found");
  }

  @ExceptionHandler(IOException.class)
  public ApiResponse ioError(IOException error) {
    return new ApiResponse(40000, null, error.getMessage());
  }

  @ExceptionHandler(DataIntegrityViolationException.class)
  public ApiResponse dataConflict(DataIntegrityViolationException error) {
    return new ApiResponse(40900, null, "data constraint violation");
  }

  @ExceptionHandler(IllegalStateException.class)
  public ApiResponse conflict(IllegalStateException error) {
    return new ApiResponse(40900, null, error.getMessage());
  }

  @ExceptionHandler(UnauthorizedException.class)
  public ApiResponse unauthorized(UnauthorizedException error) {
    return new ApiResponse(40100, null, error.getMessage());
  }
}
