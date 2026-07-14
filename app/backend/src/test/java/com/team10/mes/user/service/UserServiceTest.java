package com.team10.mes.user.service;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.argThat;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.inOrder;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.team10.mes.user.dal.UserCache;
import com.team10.mes.user.dal.UserMapper;
import java.util.Map;
import org.junit.jupiter.api.Test;
import org.mockito.InOrder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

class UserServiceTest {
  private final UserMapper mapper = org.mockito.Mockito.mock(UserMapper.class);
  private final UserCache cache = org.mockito.Mockito.mock(UserCache.class);
  private final UserService service = new UserService(mapper, cache);

  @Test
  void loginAcceptsExistingBcryptPasswordFromMapper() {
    when(mapper.findByAccount("root"))
        .thenReturn(
            Map.of(
                "id", 1L,
                "name", "Root",
                "userAccount", "root",
                "userRole", "admin",
                "passwordHash", new BCryptPasswordEncoder().encode("correct-password")));

    Map<String, Object> user =
        service.login(Map.of("userAccount", "root", "userPassword", "correct-password"));

    assertEquals(1L, user.get("id"));
    assertEquals("root", user.get("userAccount"));
    assertEquals("admin", user.get("userRole"));
  }

  @Test
  void loginRejectsWrongPassword() {
    when(mapper.findByAccount("root"))
        .thenReturn(
            Map.of(
                "id", 1L,
                "name", "Root",
                "userAccount", "root",
                "userRole", "admin",
                "passwordHash", new BCryptPasswordEncoder().encode("correct-password")));

    assertThrows(
        IllegalArgumentException.class,
        () -> service.login(Map.of("userAccount", "root", "userPassword", "wrong-password")));
  }

  @Test
  void getReturnsCachedUserWithoutQueryingMysql() {
    Map<String, Object> cached = Map.of("id", 1L, "userAccount", "root");
    when(cache.get(1L)).thenReturn(cached);

    assertEquals(cached, service.get(1L));

    verify(mapper, never()).findById(1L);
  }

  @Test
  void getQueriesMysqlAndBackfillsRedisOnCacheMiss() {
    Map<String, Object> row =
        Map.of("id", 1L, "name", "Root", "userAccount", "root", "userRole", "admin");
    when(cache.get(1L)).thenReturn(null);
    when(mapper.findById(1L)).thenReturn(row);

    Map<String, Object> result = service.get(1L);

    assertEquals("Root", result.get("userName"));
    verify(cache).put(1L, result);
  }

  @Test
  void updateWritesMysqlBeforeRefreshingRedis() {
    Map<String, Object> request =
        Map.of(
            "id", 1L,
            "name", "Updated",
            "avatar", "",
            "profile", "",
            "role", "admin");
    Map<String, Object> row =
        Map.of("id", 1L, "name", "Updated", "userAccount", "root", "userRole", "admin");
    when(mapper.update(request)).thenReturn(1);
    when(mapper.findById(1L)).thenReturn(row);

    service.update(request);

    InOrder order = inOrder(mapper, cache);
    order.verify(mapper).update(request);
    order.verify(mapper).findById(1L);
    order.verify(cache).put(eq(1L), argThat(user -> "Updated".equals(user.get("userName"))));
  }
}
