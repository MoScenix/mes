package com.team10.mes.app.controller;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.*;

import com.team10.mes.app.dal.AppMapper;
import com.team10.mes.app.service.AppFileService;
import com.team10.mes.app.service.AppService;
import com.team10.mes.user.service.SessionIdentity;
import org.junit.jupiter.api.Test;
import org.springframework.mock.web.MockHttpSession;

class AppDeleteDocumentCleanupTest {
  private final AppService apps = mock(AppService.class);
  private final AppFileService files = mock(AppFileService.class);
  private final SessionIdentity identity = new SessionIdentity();
  private final AppController controller = new AppController(apps, files, identity);

  @Test
  void ownerDeleteCleansDocumentStateOnlyAfterSuccessfulSoftDelete() throws Exception {
    MockHttpSession session = session(42, "worker");
    when(apps.get(7)).thenReturn(new AppMapper.AppRow(7L, "chat", 42L, null, null));
    when(apps.delete(7)).thenReturn(true);

    assertTrue(controller.delete(new AppController.IdRequest(7), session).success());

    var order = inOrder(apps, files);
    order.verify(apps).delete(7);
    order.verify(files).deleteProjectFiles(7);
  }

  @Test
  void adminDeleteAlsoCleansAndFailedDeleteDoesNot() throws Exception {
    MockHttpSession admin = session(1, "admin");
    when(apps.delete(7)).thenReturn(true);
    when(apps.delete(8)).thenReturn(false);

    assertTrue(controller.adminDelete(new AppController.IdRequest(7), admin).success());
    controller.adminDelete(new AppController.IdRequest(8), admin);

    verify(files).deleteProjectFiles(7);
    verify(files, never()).deleteProjectFiles(8);
  }

  private static MockHttpSession session(long userId, String role) {
    MockHttpSession session = new MockHttpSession();
    session.setAttribute(SessionIdentity.USER_ID, userId);
    session.setAttribute(SessionIdentity.USER_ROLE, role);
    return session;
  }
}
