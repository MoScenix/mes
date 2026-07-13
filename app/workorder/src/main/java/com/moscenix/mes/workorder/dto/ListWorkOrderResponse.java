package com.moscenix.mes.workorder.dto;

import java.util.ArrayList;
import java.util.List;

public class ListWorkOrderResponse {
    private List<WorkOrderInfo> workOrderList = new ArrayList<>();
    private long total;
    private boolean hasMore;
    private String nextCursorUpdatedAt;
    private Long nextCursorId;

    public List<WorkOrderInfo> getWorkOrderList() {
        return workOrderList;
    }

    public void setWorkOrderList(List<WorkOrderInfo> workOrderList) {
        this.workOrderList = workOrderList;
    }

    public long getTotal() {
        return total;
    }

    public void setTotal(long total) {
        this.total = total;
    }

    public boolean isHasMore() {
        return hasMore;
    }

    public void setHasMore(boolean hasMore) {
        this.hasMore = hasMore;
    }

    public String getNextCursorUpdatedAt() {
        return nextCursorUpdatedAt;
    }

    public void setNextCursorUpdatedAt(String nextCursorUpdatedAt) {
        this.nextCursorUpdatedAt = nextCursorUpdatedAt;
    }

    public Long getNextCursorId() {
        return nextCursorId;
    }

    public void setNextCursorId(Long nextCursorId) {
        this.nextCursorId = nextCursorId;
    }
}
