package com.moscenix.mes.workorder.grpc;

import com.moscenix.mes.grpc.workorder.CreateWorkOrderReq;
import com.moscenix.mes.grpc.workorder.CreateWorkOrderResp;
import com.moscenix.mes.grpc.workorder.DeleteWorkOrderDraftReq;
import com.moscenix.mes.grpc.workorder.DeleteWorkOrderDraftResp;
import com.moscenix.mes.grpc.workorder.GetWorkOrderReq;
import com.moscenix.mes.grpc.workorder.GetWorkOrderResp;
import com.moscenix.mes.grpc.workorder.ListWorkOrderReq;
import com.moscenix.mes.grpc.workorder.ListWorkOrderResp;
import com.moscenix.mes.grpc.workorder.MarkWorkOrderReadReq;
import com.moscenix.mes.grpc.workorder.MarkWorkOrderReadResp;
import com.moscenix.mes.grpc.workorder.SubmitWorkOrderReq;
import com.moscenix.mes.grpc.workorder.SubmitWorkOrderResp;
import com.moscenix.mes.grpc.workorder.UpdateWorkOrderDraftReq;
import com.moscenix.mes.grpc.workorder.UpdateWorkOrderDraftResp;
import com.moscenix.mes.grpc.workorder.WorkOrderReadStatus;
import com.moscenix.mes.grpc.workorder.WorkOrderServiceGrpc;
import com.moscenix.mes.grpc.workorder.WorkOrderStatus;
import com.moscenix.mes.workorder.dto.CreateWorkOrderRequest;
import com.moscenix.mes.workorder.dto.ListWorkOrderRequest;
import com.moscenix.mes.workorder.dto.UpdateWorkOrderDraftRequest;
import com.moscenix.mes.workorder.dto.WorkOrderInfo;
import com.moscenix.mes.workorder.service.WorkOrderService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import org.springframework.stereotype.Component;

@Component
public class WorkOrderGrpcEndpoint extends WorkOrderServiceGrpc.WorkOrderServiceImplBase {
    private final WorkOrderService workOrderService;

    public WorkOrderGrpcEndpoint(WorkOrderService workOrderService) {
        this.workOrderService = workOrderService;
    }

    @Override
    public void createWorkOrder(CreateWorkOrderReq request, StreamObserver<CreateWorkOrderResp> responseObserver) {
        handle(responseObserver, () -> {
            CreateWorkOrderRequest dto = new CreateWorkOrderRequest();
            dto.setFromUserId(request.getFromUserId());
            dto.setToUserId(request.getToUserId());
            dto.setDescription(request.getDescription());
            dto.setName(request.getName());
            var result = workOrderService.create(dto);
            return CreateWorkOrderResp.newBuilder().setId(result.getId()).build();
        });
    }

    @Override
    public void updateWorkOrderDraft(UpdateWorkOrderDraftReq request, StreamObserver<UpdateWorkOrderDraftResp> responseObserver) {
        handle(responseObserver, () -> {
            UpdateWorkOrderDraftRequest dto = new UpdateWorkOrderDraftRequest();
            dto.setFromUserId(request.getFromUserId());
            dto.setToUserId(request.getToUserId());
            dto.setDescription(request.getDescription());
            dto.setName(request.getName());
            boolean success = workOrderService.updateDraft(request.getId(), dto);
            return UpdateWorkOrderDraftResp.newBuilder().setSuccess(success).build();
        });
    }

    @Override
    public void deleteWorkOrderDraft(DeleteWorkOrderDraftReq request, StreamObserver<DeleteWorkOrderDraftResp> responseObserver) {
        handle(responseObserver, () -> DeleteWorkOrderDraftResp.newBuilder()
                .setSuccess(workOrderService.deleteDraft(request.getId()))
                .build());
    }

    @Override
    public void submitWorkOrder(SubmitWorkOrderReq request, StreamObserver<SubmitWorkOrderResp> responseObserver) {
        handle(responseObserver, () -> SubmitWorkOrderResp.newBuilder()
                .setSuccess(workOrderService.submit(request.getId()))
                .build());
    }

    @Override
    public void getWorkOrder(GetWorkOrderReq request, StreamObserver<GetWorkOrderResp> responseObserver) {
        handle(responseObserver, () -> GetWorkOrderResp.newBuilder()
                .setWorkOrder(toProto(workOrderService.get(request.getId()).getWorkOrder()))
                .build());
    }

    @Override
    public void listWorkOrder(ListWorkOrderReq request, StreamObserver<ListWorkOrderResp> responseObserver) {
        handle(responseObserver, () -> {
            ListWorkOrderRequest dto = new ListWorkOrderRequest();
            dto.setPageNum(request.getPageNum());
            dto.setPageSize(request.getPageSize());
            dto.setId(request.getId());
            dto.setIsTo(request.getIsTo());
            dto.setIsUnread(request.getIsUnread());
            dto.setSinceTime(request.getSinceTime());
            dto.setRecentSeconds(request.getRecentSeconds());
            dto.setCursorUpdatedAt(request.getCursorUpdatedAt());
            dto.setCursorId(request.getCursorId());
            dto.setNamePrefix(request.getNamePrefix());
            dto.setStatus(request.getStatus().getNumber());
            var result = workOrderService.list(dto);
            return ListWorkOrderResp.newBuilder()
                    .addAllWorkOrderList(result.getWorkOrderList().stream().map(this::toProto).toList())
                    .setTotal(result.getTotal())
                    .setHasMore(result.isHasMore())
                    .setNextCursorUpdatedAt(value(result.getNextCursorUpdatedAt()))
                    .setNextCursorId(result.getNextCursorId() == null ? 0L : result.getNextCursorId())
                    .build();
        });
    }

    @Override
    public void markWorkOrderRead(MarkWorkOrderReadReq request, StreamObserver<MarkWorkOrderReadResp> responseObserver) {
        handle(responseObserver, () -> MarkWorkOrderReadResp.newBuilder()
                .setSuccess(workOrderService.markRead(request.getId()))
                .build());
    }

    private com.moscenix.mes.grpc.workorder.WorkOrderInfo toProto(WorkOrderInfo order) {
        return com.moscenix.mes.grpc.workorder.WorkOrderInfo.newBuilder()
                .setId(order.getId() == null ? 0L : order.getId())
                .setFromUserId(order.getFromUserId() == null ? 0L : order.getFromUserId())
                .setToUserId(order.getToUserId() == null ? 0L : order.getToUserId())
                .setDescription(value(order.getDescription()))
                .setStatus(WorkOrderStatus.forNumber(order.getStatus() == null ? 0 : order.getStatus().getCode()))
                .setCreateTime(value(order.getCreateTime()))
                .setUpdateTime(value(order.getUpdateTime()))
                .setReadStatus(WorkOrderReadStatus.forNumber(order.getReadStatus() == null ? 0 : order.getReadStatus().getCode()))
                .setName(value(order.getName()))
                .build();
    }

    private String value(String value) {
        return value == null ? "" : value;
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
