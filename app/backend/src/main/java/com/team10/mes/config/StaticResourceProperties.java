package com.team10.mes.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "mes.static")
public class StaticResourceProperties {
  private String root;

  public String getRoot() {
    return root;
  }

  public void setRoot(String root) {
    this.root = root;
  }
}
