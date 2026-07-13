package com.moscenix.mes.user.dto;

public class RegisterResponse {
    private Long userId;
    private String userRole;

    public RegisterResponse() {
    }

    public RegisterResponse(Long userId, String userRole) {
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
