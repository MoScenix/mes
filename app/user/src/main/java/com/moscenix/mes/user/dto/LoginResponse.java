package com.moscenix.mes.user.dto;

public class LoginResponse {
    private Long userId;
    private String userRole;

    public LoginResponse() {
    }

    public LoginResponse(Long userId, String userRole) {
        this.userId = userId;
        this.userRole = userRole;
    }

    public Long getUserId() {
        return userId;
    }

    public void setUserId(Long userId) {
        this.userId = userId;
    }

    public String getUserRole() {
        return userRole;
    }

    public void setUserRole(String userRole) {
        this.userRole = userRole;
    }
}
