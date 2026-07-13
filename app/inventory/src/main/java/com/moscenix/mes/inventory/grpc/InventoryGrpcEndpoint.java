package com.moscenix.mes.inventory.grpc;

import com.moscenix.mes.grpc.inventory.AddItemReq;
import com.moscenix.mes.grpc.inventory.AddItemResp;
import com.moscenix.mes.grpc.inventory.AddItemUnitReq;
import com.moscenix.mes.grpc.inventory.AddItemUnitResp;
import com.moscenix.mes.grpc.inventory.AuditInventoryFlowReq;
import com.moscenix.mes.grpc.inventory.AuditInventoryFlowResp;
import com.moscenix.mes.grpc.inventory.CompleteInventoryFlowReq;
import com.moscenix.mes.grpc.inventory.CompleteInventoryFlowResp;
import com.moscenix.mes.grpc.inventory.CreateEngineeringOrderDraftReq;
import com.moscenix.mes.grpc.inventory.CreateEngineeringOrderDraftResp;
import com.moscenix.mes.grpc.inventory.CreateInventoryFlowReq;
import com.moscenix.mes.grpc.inventory.CreateInventoryFlowResp;
import com.moscenix.mes.grpc.inventory.CreateProcessDraftReq;
import com.moscenix.mes.grpc.inventory.CreateProcessDraftResp;
import com.moscenix.mes.grpc.inventory.DeleteEngineeringOrderDraftReq;
import com.moscenix.mes.grpc.inventory.DeleteEngineeringOrderDraftResp;
import com.moscenix.mes.grpc.inventory.DeleteInventoryFlowDraftReq;
import com.moscenix.mes.grpc.inventory.DeleteInventoryFlowDraftResp;
import com.moscenix.mes.grpc.inventory.DeleteProcessDraftReq;
import com.moscenix.mes.grpc.inventory.DeleteProcessDraftResp;
import com.moscenix.mes.grpc.inventory.GetEngineeringOrderReq;
import com.moscenix.mes.grpc.inventory.GetEngineeringOrderResp;
import com.moscenix.mes.grpc.inventory.GetInventoryFlowReq;
import com.moscenix.mes.grpc.inventory.GetInventoryFlowResp;
import com.moscenix.mes.grpc.inventory.GetItemReq;
import com.moscenix.mes.grpc.inventory.GetItemResp;
import com.moscenix.mes.grpc.inventory.GetItemUnitReq;
import com.moscenix.mes.grpc.inventory.GetItemUnitResp;
import com.moscenix.mes.grpc.inventory.GetProcessReq;
import com.moscenix.mes.grpc.inventory.GetProcessResp;
import com.moscenix.mes.grpc.inventory.InventoryServiceGrpc;
import com.moscenix.mes.grpc.inventory.ListEngineeringOrderReq;
import com.moscenix.mes.grpc.inventory.ListEngineeringOrderResp;
import com.moscenix.mes.grpc.inventory.ListInventoryFlowReq;
import com.moscenix.mes.grpc.inventory.ListInventoryFlowResp;
import com.moscenix.mes.grpc.inventory.ListItemReq;
import com.moscenix.mes.grpc.inventory.ListItemResp;
import com.moscenix.mes.grpc.inventory.ListItemUnitReq;
import com.moscenix.mes.grpc.inventory.ListItemUnitResp;
import com.moscenix.mes.grpc.inventory.ListProcessReq;
import com.moscenix.mes.grpc.inventory.ListProcessResp;
import com.moscenix.mes.grpc.inventory.SubmitEngineeringOrderReq;
import com.moscenix.mes.grpc.inventory.SubmitEngineeringOrderResp;
import com.moscenix.mes.grpc.inventory.SubmitInventoryFlowReq;
import com.moscenix.mes.grpc.inventory.SubmitInventoryFlowResp;
import com.moscenix.mes.grpc.inventory.SubmitProcessReq;
import com.moscenix.mes.grpc.inventory.SubmitProcessResp;
import com.moscenix.mes.grpc.inventory.UpdateEngineeringOrderDraftReq;
import com.moscenix.mes.grpc.inventory.UpdateEngineeringOrderDraftResp;
import com.moscenix.mes.grpc.inventory.UpdateInventoryFlowDraftReq;
import com.moscenix.mes.grpc.inventory.UpdateInventoryFlowDraftResp;
import com.moscenix.mes.grpc.inventory.UpdateItemReq;
import com.moscenix.mes.grpc.inventory.UpdateItemResp;
import com.moscenix.mes.grpc.inventory.UpdateItemUnitStatusReq;
import com.moscenix.mes.grpc.inventory.UpdateItemUnitStatusResp;
import com.moscenix.mes.grpc.inventory.UpdateProcessDraftReq;
import com.moscenix.mes.grpc.inventory.UpdateProcessDraftResp;
import com.moscenix.mes.inventory.service.InventoryService;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import org.springframework.stereotype.Component;

@Component
public class InventoryGrpcEndpoint extends InventoryServiceGrpc.InventoryServiceImplBase {
    private final InventoryService inventoryService;

    public InventoryGrpcEndpoint(InventoryService inventoryService) {
        this.inventoryService = inventoryService;
    }

    @Override
    public void addItem(AddItemReq request, StreamObserver<AddItemResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.addItem(request));
    }

    @Override
    public void updateItem(UpdateItemReq request, StreamObserver<UpdateItemResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.updateItem(request));
    }

    @Override
    public void getItem(GetItemReq request, StreamObserver<GetItemResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.getItem(request));
    }

    @Override
    public void listItem(ListItemReq request, StreamObserver<ListItemResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.listItem(request));
    }

    @Override
    public void createProcessDraft(CreateProcessDraftReq request, StreamObserver<CreateProcessDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.createProcessDraft(request));
    }

    @Override
    public void updateProcessDraft(UpdateProcessDraftReq request, StreamObserver<UpdateProcessDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.updateProcessDraft(request));
    }

    @Override
    public void deleteProcessDraft(DeleteProcessDraftReq request, StreamObserver<DeleteProcessDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.deleteProcessDraft(request));
    }

    @Override
    public void submitProcess(SubmitProcessReq request, StreamObserver<SubmitProcessResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.submitProcess(request));
    }

    @Override
    public void getProcess(GetProcessReq request, StreamObserver<GetProcessResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.getProcess(request));
    }

    @Override
    public void listProcess(ListProcessReq request, StreamObserver<ListProcessResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.listProcess(request));
    }

    @Override
    public void addItemUnit(AddItemUnitReq request, StreamObserver<AddItemUnitResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.addItemUnit(request));
    }

    @Override
    public void updateItemUnitStatus(UpdateItemUnitStatusReq request, StreamObserver<UpdateItemUnitStatusResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.updateItemUnitStatus(request));
    }

    @Override
    public void getItemUnit(GetItemUnitReq request, StreamObserver<GetItemUnitResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.getItemUnit(request));
    }

    @Override
    public void listItemUnit(ListItemUnitReq request, StreamObserver<ListItemUnitResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.listItemUnit(request));
    }

    @Override
    public void createInventoryFlow(CreateInventoryFlowReq request, StreamObserver<CreateInventoryFlowResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.createInventoryFlow(request));
    }

    @Override
    public void updateInventoryFlowDraft(UpdateInventoryFlowDraftReq request, StreamObserver<UpdateInventoryFlowDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.updateInventoryFlowDraft(request));
    }

    @Override
    public void deleteInventoryFlowDraft(DeleteInventoryFlowDraftReq request, StreamObserver<DeleteInventoryFlowDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.deleteInventoryFlowDraft(request));
    }

    @Override
    public void submitInventoryFlow(SubmitInventoryFlowReq request, StreamObserver<SubmitInventoryFlowResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.submitInventoryFlow(request));
    }

    @Override
    public void completeInventoryFlow(CompleteInventoryFlowReq request, StreamObserver<CompleteInventoryFlowResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.completeInventoryFlow(request));
    }

    @Override
    public void auditInventoryFlow(AuditInventoryFlowReq request, StreamObserver<AuditInventoryFlowResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.auditInventoryFlow(request));
    }

    @Override
    public void getInventoryFlow(GetInventoryFlowReq request, StreamObserver<GetInventoryFlowResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.getInventoryFlow(request));
    }

    @Override
    public void listInventoryFlow(ListInventoryFlowReq request, StreamObserver<ListInventoryFlowResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.listInventoryFlow(request));
    }

    @Override
    public void createEngineeringOrderDraft(CreateEngineeringOrderDraftReq request, StreamObserver<CreateEngineeringOrderDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.createEngineeringOrderDraft(request));
    }

    @Override
    public void updateEngineeringOrderDraft(UpdateEngineeringOrderDraftReq request, StreamObserver<UpdateEngineeringOrderDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.updateEngineeringOrderDraft(request));
    }

    @Override
    public void deleteEngineeringOrderDraft(DeleteEngineeringOrderDraftReq request, StreamObserver<DeleteEngineeringOrderDraftResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.deleteEngineeringOrderDraft(request));
    }

    @Override
    public void submitEngineeringOrder(SubmitEngineeringOrderReq request, StreamObserver<SubmitEngineeringOrderResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.submitEngineeringOrder(request));
    }

    @Override
    public void getEngineeringOrder(GetEngineeringOrderReq request, StreamObserver<GetEngineeringOrderResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.getEngineeringOrder(request));
    }

    @Override
    public void listEngineeringOrder(ListEngineeringOrderReq request, StreamObserver<ListEngineeringOrderResp> responseObserver) {
        handle(responseObserver, () -> inventoryService.listEngineeringOrder(request));
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
