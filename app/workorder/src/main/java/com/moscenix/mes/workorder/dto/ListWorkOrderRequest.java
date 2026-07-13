package com.moscenix.mes.workorder.dto;

public class ListWorkOrderRequest {
    private Long pageNum;
    private Long pageSize;
    private Long id;
    private Boolean isTo;
    private Boolean isUnread;
    private String sinceTime;
    private Long recentSeconds;
    private String cursorUpdatedAt;
    private Long cursorId;
    private String namePrefix;
    private Integer status;

    public Long getPageNum() {
        return pageNum;
    }

    public void setPageNum(Long pageNum) {
        this.pageNum = pageNum;
    }

    public Long getPageSize() {
        return pageSize;
    }

    public void setPageSize(Long pageSize) {
        this.pageSize = pageSize;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Boolean getIsTo() {
        return isTo;
    }

    public void setIsTo(Boolean isTo) {
        this.isTo = isTo;
    }

    public Boolean getIsUnread() {
        return isUnread;
    }

    public void setIsUnread(Boolean isUnread) {
        this.isUnread = isUnread;
    }

    public String getSinceTime() {
        return sinceTime;
    }

    public void setSinceTime(String sinceTime) {
        this.sinceTime = sinceTime;
    }

    public Long getRecentSeconds() {
        return recentSeconds;
    }

    public void setRecentSeconds(Long recentSeconds) {
        this.recentSeconds = recentSeconds;
    }

    public String getCursorUpdatedAt() {
        return cursorUpdatedAt;
    }

    public void setCursorUpdatedAt(String cursorUpdatedAt) {
        this.cursorUpdatedAt = cursorUpdatedAt;
    }

    public Long getCursorId() {
        return cursorId;
    }

    public void setCursorId(Long cursorId) {
        this.cursorId = cursorId;
    }

    public String getNamePrefix() {
        return namePrefix;
    }

    public void setNamePrefix(String namePrefix) {
        this.namePrefix = namePrefix;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }
}
