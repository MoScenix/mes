package com.team10.mes.history.model;

import java.time.LocalDateTime;

public class HistoryMessage {
  private Long id;
  private Long appId;
  private Long userId;
  private String role;
  private String content;
  private LocalDateTime createTime;
  private LocalDateTime updateTime;
  private Boolean isFile;

  public Long getId() {
    return id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public Long getAppId() {
    return appId;
  }

  public void setAppId(Long appId) {
    this.appId = appId;
  }

  public Long getUserId() {
    return userId;
  }

  public void setUserId(Long userId) {
    this.userId = userId;
  }

  public String getRole() {
    return role;
  }

  public void setRole(String role) {
    this.role = role;
  }

  public String getContent() {
    return content;
  }

  public void setContent(String content) {
    this.content = content;
  }

  public LocalDateTime getCreateTime() {
    return createTime;
  }

  public void setCreateTime(LocalDateTime createTime) {
    this.createTime = createTime;
  }

  public LocalDateTime getUpdateTime() {
    return updateTime;
  }

  public void setUpdateTime(LocalDateTime updateTime) {
    this.updateTime = updateTime;
  }

  public Boolean getIsFile() {
    return isFile;
  }

  public void setIsFile(Boolean isFile) {
    this.isFile = isFile;
  }
}
