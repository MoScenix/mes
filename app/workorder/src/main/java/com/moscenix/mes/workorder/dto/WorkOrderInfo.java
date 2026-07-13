package com.moscenix.mes.workorder.dto;

import com.moscenix.mes.workorder.domain.WorkOrderReadStatus;
import com.moscenix.mes.workorder.domain.WorkOrderStatus;

public class WorkOrderInfo {
    private Long id;
    private Long fromUserId;
    private Long toUserId;
    private String description;
    private WorkOrderStatus status;
    private String createTime;
    private String updateTime;
    private WorkOrderReadStatus readStatus;
    private String name;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getFromUserId() {
        return fromUserId;
    }

    public void setFromUserId(Long fromUserId) {
        this.fromUserId = fromUserId;
    }

    public Long getToUserId() {
        return toUserId;
    }

    public void setToUserId(Long toUserId) {
        this.toUserId = toUserId;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public WorkOrderStatus getStatus() {
        return status;
    }

    public void setStatus(WorkOrderStatus status) {
        this.status = status;
    }

    public String getCreateTime() {
        return createTime;
    }

    public void setCreateTime(String createTime) {
        this.createTime = createTime;
    }

    public String getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(String updateTime) {
        this.updateTime = updateTime;
    }

    public WorkOrderReadStatus getReadStatus() {
        return readStatus;
    }

    public void setReadStatus(WorkOrderReadStatus readStatus) {
        this.readStatus = readStatus;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
