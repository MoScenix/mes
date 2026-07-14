package com.team10.mes.app.controller;

import static org.mockito.Mockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.multipart;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

import com.team10.mes.app.service.AppFileService;
import com.team10.mes.controller.ApiResponseAdvice;
import com.team10.mes.user.service.SessionIdentity;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.mock.web.MockHttpSession;
import org.springframework.mock.web.MockMultipartFile;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;

class AppFileContractTest {
  private final AppFileService files = mock(AppFileService.class);
  private final SessionIdentity identity = new SessionIdentity();
  private MockMvc mvc;

  @BeforeEach
  void setUp() {
    mvc =
        MockMvcBuilders.standaloneSetup(new AppFileController(files, identity))
            .setControllerAdvice(new ApiResponseAdvice())
            .build();
  }

  @Test
  void addFileKeepsLegacyMultipartAndBaseResponseStringContract() throws Exception {
    MockHttpSession session = new MockHttpSession();
    session.setAttribute(SessionIdentity.USER_ID, 42L);
    session.setAttribute(SessionIdentity.USER_ROLE, "worker");
    var upload = new MockMultipartFile("file", "note.txt", "text/plain", "hello".getBytes());
    when(files.upload(7, 42, false, upload)).thenReturn("99");

    mvc.perform(multipart("/app/file/add").file(upload).param("appId", "7").session(session))
        .andExpect(status().isOk())
        .andExpect(content().contentTypeCompatibleWith("application/json"))
        .andExpect(jsonPath("$.code").value(0))
        .andExpect(jsonPath("$.data").value("99"))
        .andExpect(jsonPath("$.message").value("success"));

    verify(files).upload(7, 42, false, upload);
  }

  @Test
  void missingSessionIsRejectedWithoutCallingService() throws Exception {
    var upload = new MockMultipartFile("file", "note.txt", "text/plain", "hello".getBytes());

    mvc.perform(multipart("/app/file/add").file(upload).param("appId", "7"))
        .andExpect(status().isOk())
        .andExpect(jsonPath("$.code").value(40100))
        .andExpect(jsonPath("$.data").doesNotExist());

    verifyNoInteractions(files);
  }
}
