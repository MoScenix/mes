package com.moscenix.mes.user.mapper;

import com.moscenix.mes.user.dto.GetUserResponse;
import com.moscenix.mes.user.entity.UserEntity;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import org.springframework.stereotype.Component;

@Component
public class UserMapper {
    private static final DateTimeFormatter FORMATTER = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    public GetUserResponse toGetUserResponse(UserEntity user, boolean includePassword) {
        GetUserResponse response = new GetUserResponse();
        response.setId(user.getId());
        response.setUserAccount(user.getUserAccount());
        response.setUserPassword(includePassword ? user.getPasswordHash() : "");
        response.setUserName(user.getName());
        response.setUserAvatar(user.getUserAvatar());
        response.setUserProfile(user.getUserProfile());
        response.setUserRole(user.getUserRole());
        response.setEditTime("");
        response.setCreateTime(format(user.getCreatedAt()));
        response.setUpdateTime(format(user.getUpdatedAt()));
        response.setIsDelete(user.getDeletedAt() != null);
        return response;
    }

    private String format(LocalDateTime time) {
        if (time == null) {
            return "";
        }
        return time.format(FORMATTER);
    }
}
