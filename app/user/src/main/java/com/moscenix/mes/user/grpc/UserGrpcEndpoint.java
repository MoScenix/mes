package com.moscenix.mes.user.grpc;

import com.moscenix.mes.grpc.user.AddUserReq;
import com.moscenix.mes.grpc.user.AddUserResp;
import com.moscenix.mes.grpc.user.DeleteUserReq;
import com.moscenix.mes.grpc.user.DeleteUserResp;
import com.moscenix.mes.grpc.user.GetUserReq;
import com.moscenix.mes.grpc.user.GetUserResp;
import com.moscenix.mes.grpc.user.ListUserReq;
import com.moscenix.mes.grpc.user.ListUserResp;
import com.moscenix.mes.grpc.user.LoginReq;
import com.moscenix.mes.grpc.user.LoginResp;
import com.moscenix.mes.grpc.user.RegisterReq;
import com.moscenix.mes.grpc.user.RegisterResp;
import com.moscenix.mes.grpc.user.UpdateReq;
import com.moscenix.mes.grpc.user.UpdateResp;
import com.moscenix.mes.grpc.user.UserServiceGrpc;
import com.moscenix.mes.user.dto.AddUserRequest;
import com.moscenix.mes.user.dto.GetUserResponse;
import com.moscenix.mes.user.dto.ListUserRequest;
import com.moscenix.mes.user.dto.LoginRequest;
import com.moscenix.mes.user.dto.RegisterRequest;
import com.moscenix.mes.user.dto.UpdateUserRequest;
import com.moscenix.mes.user.service.UserService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import org.springframework.stereotype.Component;

@Component
public class UserGrpcEndpoint extends UserServiceGrpc.UserServiceImplBase {
    private final UserService userService;

    public UserGrpcEndpoint(UserService userService) {
        this.userService = userService;
    }

    @Override
    public void login(LoginReq request, StreamObserver<LoginResp> responseObserver) {
        handle(responseObserver, () -> {
            LoginRequest dto = new LoginRequest();
            dto.setUserAccount(request.getUserAccount());
            dto.setUserPassword(request.getUserPassword());
            var result = userService.login(dto);
            return LoginResp.newBuilder()
                    .setUserId(result.getUserId().intValue())
                    .setUserRole(result.getUserRole())
                    .build();
        });
    }

    @Override
    public void register(RegisterReq request, StreamObserver<RegisterResp> responseObserver) {
        handle(responseObserver, () -> {
            RegisterRequest dto = new RegisterRequest();
            dto.setUserAccount(request.getUserAccount());
            dto.setUserPassword(request.getUserPassword());
            var result = userService.register(dto);
            return RegisterResp.newBuilder()
                    .setUserId(result.getUserId().intValue())
                    .setUserRole(result.getUserRole())
                    .build();
        });
    }

    @Override
    public void update(UpdateReq request, StreamObserver<UpdateResp> responseObserver) {
        handle(responseObserver, () -> {
            UpdateUserRequest dto = new UpdateUserRequest();
            dto.setId(request.getId());
            dto.setUserName(request.getUserName());
            dto.setUserAvatar(request.getUserAvatar());
            dto.setUserProfile(request.getUserProfile());
            dto.setUserRole(request.getUserRole());
            userService.update(dto);
            return UpdateResp.newBuilder().build();
        });
    }

    @Override
    public void getUser(GetUserReq request, StreamObserver<GetUserResp> responseObserver) {
        handle(responseObserver, () -> toProto(userService.getUser(request.getId())));
    }

    @Override
    public void deleteUser(DeleteUserReq request, StreamObserver<DeleteUserResp> responseObserver) {
        handle(responseObserver, () -> {
            userService.deleteUser(request.getUserId());
            return DeleteUserResp.newBuilder().build();
        });
    }

    @Override
    public void addUser(AddUserReq request, StreamObserver<AddUserResp> responseObserver) {
        handle(responseObserver, () -> {
            AddUserRequest dto = new AddUserRequest();
            dto.setUserAccount(request.getUserAccount());
            dto.setUserPassword(request.getUserPassword());
            dto.setUserName(request.getUserName());
            dto.setUserAvatar(request.getUserAvatar());
            dto.setUserProfile(request.getUserProfile());
            dto.setUserRole(request.getUserRole());
            userService.addUser(dto);
            return AddUserResp.newBuilder().build();
        });
    }

    @Override
    public void listUser(ListUserReq request, StreamObserver<ListUserResp> responseObserver) {
        handle(responseObserver, () -> {
            ListUserRequest dto = new ListUserRequest();
            dto.setPageNum(request.getPageNum());
            dto.setPageSize(request.getPageSize());
            dto.setUserName(request.getUserName());
            dto.setAccount(request.getAccount());
            var result = userService.listUser(dto);
            return ListUserResp.newBuilder()
                    .addAllUserList(result.getUserList().stream().map(this::toProto).toList())
                    .setTotal(result.getTotal())
                    .build();
        });
    }

    private GetUserResp toProto(GetUserResponse user) {
        return GetUserResp.newBuilder()
                .setId(value(user.getId()))
                .setUserAccount(value(user.getUserAccount()))
                .setUserPassword(value(user.getUserPassword()))
                .setUserName(value(user.getUserName()))
                .setUserAvatar(value(user.getUserAvatar()))
                .setUserProfile(value(user.getUserProfile()))
                .setUserRole(value(user.getUserRole()))
                .setEditTime(value(user.getEditTime()))
                .setCreateTime(value(user.getCreateTime()))
                .setUpdateTime(value(user.getUpdateTime()))
                .setIsDelete(Boolean.TRUE.equals(user.getIsDelete()))
                .build();
    }

    private String value(String value) {
        return value == null ? "" : value;
    }

    private long value(Long value) {
        return value == null ? 0L : value;
    }

    private <T> void handle(StreamObserver<T> responseObserver, GrpcCall<T> call) {
        try {
            responseObserver.onNext(call.invoke());
            responseObserver.onCompleted();
        } catch (RuntimeException e) {
            responseObserver.onError(Status.INVALID_ARGUMENT.withDescription(e.getMessage()).asRuntimeException());
        }
    }

    @FunctionalInterface
    private interface GrpcCall<T> {
        T invoke();
    }
}
