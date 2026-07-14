package com.team10.mes.app.controller;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

import com.team10.mes.app.service.StaticResourceProperties;
import com.team10.mes.controller.HealthController;
import java.nio.file.Files;
import java.nio.file.Path;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.test.web.servlet.MockMvc;

@WebMvcTest(
    controllers = HealthController.class,
    properties = "mes.static.root=target/static-resource-test")
@Import({StaticResourceConfiguration.class, StaticResourceProperties.class})
class StaticResourceConfigurationTest {
  private static final Path ROOT = Path.of("target/static-resource-test");

  @Autowired MockMvc mvc;

  @BeforeAll
  static void createStaticFile() throws Exception {
    Files.createDirectories(ROOT.resolve("document"));
    Files.writeString(ROOT.resolve("document/file.txt"), "mapped");
    Files.writeString(Path.of("target/outside-static.txt"), "secret");
  }

  @AfterAll
  static void cleanStaticFiles() throws Exception {
    Files.deleteIfExists(ROOT.resolve("document/file.txt"));
    Files.deleteIfExists(ROOT.resolve("document"));
    Files.deleteIfExists(ROOT);
    Files.deleteIfExists(Path.of("target/outside-static.txt"));
  }

  @Test
  void mapsConfiguredStaticRootAndRejectsTraversal() throws Exception {
    mvc.perform(get("/static/document/file.txt"))
        .andExpect(status().isOk())
        .andExpect(content().string("mapped"));
    mvc.perform(get("/static/../outside-static.txt")).andExpect(status().is4xxClientError());
  }
}
