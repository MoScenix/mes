package com.team10.mes.user.service;

import com.team10.mes.user.dal.UserMapper;
import java.util.HashMap;
import java.util.Map;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Component;

@Component
public class UserBootstrapInitializer implements ApplicationRunner {
  private static final String ROOT_ACCOUNT = "root";
  private static final String ROOT_PASSWORD = "rootroot";

  private final UserMapper mapper;
  private final BCryptPasswordEncoder encoder = new BCryptPasswordEncoder();

  public UserBootstrapInitializer(UserMapper mapper) {
    this.mapper = mapper;
  }

  @Override
  public void run(ApplicationArguments args) {
    if (mapper.findByAccount(ROOT_ACCOUNT) != null) return;
    Map<String, Object> root = rootUser();
    if (mapper.countByAccount(ROOT_ACCOUNT) > 0) {
      mapper.restoreByAccount(root);
      return;
    }
    mapper.insert(root);
  }

  private Map<String, Object> rootUser() {
    Map<String, Object> root = new HashMap<>();
    root.put("name", "root");
    root.put("passwordHash", encoder.encode(ROOT_PASSWORD));
    root.put("account", ROOT_ACCOUNT);
    root.put("avatar", "");
    root.put("profile", "");
    root.put("role", "admin");
    return root;
  }
}
