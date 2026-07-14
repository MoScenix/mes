package com.team10.mes.user.service;

import jakarta.servlet.http.HttpSession;
import java.util.Map;
import org.springframework.stereotype.Component;

@Component
public class SessionIdentity {
  public static final String USER_ID = "userId", USER_ROLE = "userRole", USER = "loginUser";

  public void login(HttpSession session, Map<String, Object> user) {
    session.setAttribute(USER_ID, user.get("id"));
    session.setAttribute(USER_ROLE, user.get("userRole"));
    session.setAttribute(USER, user);
  }

  public Long userId(HttpSession session) {
    Object value = session.getAttribute(USER_ID);
    return value instanceof Number n ? n.longValue() : null;
  }

  public String role(HttpSession session) {
    Object value = session.getAttribute(USER_ROLE);
    return value == null ? null : value.toString();
  }

  @SuppressWarnings("unchecked")
  public Map<String, Object> user(HttpSession session) {
    Object value = session.getAttribute(USER);
    return value instanceof Map<?, ?> ? (Map<String, Object>) value : null;
  }
}
