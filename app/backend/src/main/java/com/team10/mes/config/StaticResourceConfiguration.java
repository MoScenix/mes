package com.team10.mes.config;

import java.nio.file.Paths;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import org.springframework.web.servlet.resource.PathResourceResolver;

@Configuration
public class StaticResourceConfiguration implements WebMvcConfigurer {
  private final StaticResourceProperties properties;

  public StaticResourceConfiguration(StaticResourceProperties properties) {
    this.properties = properties;
  }

  @Override
  public void addResourceHandlers(ResourceHandlerRegistry registry) {
    String location =
        Paths.get(properties.getRoot()).toAbsolutePath().normalize().toUri().toString();
    if (!location.endsWith("/")) location += "/";
    registry
        .addResourceHandler("/static/**")
        .addResourceLocations(location)
        .resourceChain(true)
        .addResolver(new PathResourceResolver());
  }
}
