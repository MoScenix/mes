package com.moscenix.mes.user.dto;

import java.util.ArrayList;
import java.util.List;

public class ListUserResponse {
    private List<GetUserResponse> userList = new ArrayList<>();
    private Long total;

    public List<GetUserResponse> getUserList() {
        return userList;
    }

    public void setUserList(List<GetUserResponse> userList) {
        this.userList = userList;
    }

    public Long getTotal() {
        return total;
    }

    public void setTotal(Long total) {
        this.total = total;
    }
}
