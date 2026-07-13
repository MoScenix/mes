package com.moscenix.mes.inventory.service;

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
import com.moscenix.mes.grpc.inventory.EngineeringOrderInfo;
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
import com.moscenix.mes.grpc.inventory.InventoryFlowInfo;
import com.moscenix.mes.grpc.inventory.InventoryFlowItemInfo;
import com.moscenix.mes.grpc.inventory.InventoryFlowItemReq;
import com.moscenix.mes.grpc.inventory.ItemInfo;
import com.moscenix.mes.grpc.inventory.ItemUnitInfo;
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
import com.moscenix.mes.grpc.inventory.ListScope;
import com.moscenix.mes.grpc.inventory.ProcessInfo;
import com.moscenix.mes.grpc.inventory.ProcessItemInfo;
import com.moscenix.mes.grpc.inventory.ProcessItemReq;
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
import com.moscenix.mes.inventory.domain.DraftStatus;
import com.moscenix.mes.inventory.domain.EngineeringOrder;
import com.moscenix.mes.inventory.domain.FlowStatus;
import com.moscenix.mes.inventory.domain.FlowType;
import com.moscenix.mes.inventory.domain.InventoryFlow;
import com.moscenix.mes.inventory.domain.InventoryFlowItem;
import com.moscenix.mes.inventory.domain.Item;
import com.moscenix.mes.inventory.domain.ItemUnit;
import com.moscenix.mes.inventory.domain.ProcessEntity;
import com.moscenix.mes.inventory.domain.ProcessItem;
import com.moscenix.mes.inventory.domain.QualityStatus;
import com.moscenix.mes.inventory.domain.StockStatus;
import com.moscenix.mes.inventory.repository.EngineeringOrderRepository;
import com.moscenix.mes.inventory.repository.InventoryFlowRepository;
import com.moscenix.mes.inventory.repository.ItemRepository;
import com.moscenix.mes.inventory.repository.ItemUnitRepository;
import com.moscenix.mes.inventory.repository.ProcessRepository;
import jakarta.persistence.criteria.JoinType;
import jakarta.persistence.criteria.Predicate;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.function.Function;
import java.util.stream.Collectors;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class InventoryService {
    private static final long DEFAULT_PAGE_SIZE = 10L;
    private static final long MAX_PAGE_SIZE = 100L;
    private static final DateTimeFormatter LIST_TIME_FORMATTER =
            DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    private final ItemRepository itemRepository;
    private final ProcessRepository processRepository;
    private final ItemUnitRepository itemUnitRepository;
    private final InventoryFlowRepository inventoryFlowRepository;
    private final EngineeringOrderRepository engineeringOrderRepository;

    public InventoryService(
            ItemRepository itemRepository,
            ProcessRepository processRepository,
            ItemUnitRepository itemUnitRepository,
            InventoryFlowRepository inventoryFlowRepository,
            EngineeringOrderRepository engineeringOrderRepository) {
        this.itemRepository = itemRepository;
        this.processRepository = processRepository;
        this.itemUnitRepository = itemUnitRepository;
        this.inventoryFlowRepository = inventoryFlowRepository;
        this.engineeringOrderRepository = engineeringOrderRepository;
    }

    @Transactional
    public AddItemResp addItem(AddItemReq request) {
        Item item = new Item();
        item.setName(request.getName());
        item.setUnit(request.getUnit());
        item.setDescription(valueOrEmpty(request.getDescription()));
        itemRepository.save(item);
        return AddItemResp.newBuilder().setId(item.getId()).build();
    }

    @Transactional
    public UpdateItemResp updateItem(UpdateItemReq request) {
        Item item = requireItem(request.getId());
        item.setName(request.getName());
        item.setUnit(request.getUnit());
        item.setDescription(valueOrEmpty(request.getDescription()));
        return UpdateItemResp.newBuilder().setSuccess(true).build();
    }

    @Transactional(readOnly = true)
    public GetItemResp getItem(GetItemReq request) {
        return GetItemResp.newBuilder().setItem(toItemInfo(requireItem(request.getId()))).build();
    }

    @Transactional(readOnly = true)
    public ListItemResp listItem(ListItemReq request) {
        long pageSize = normalizePageSize(request.getPageSize());
        Instant cursorUpdatedAt = parseOptionalTime(request.getCursorUpdatedAt(), "cursorUpdatedAt");
        Specification<Item> spec = (root, query, cb) -> {
            List<Predicate> predicates = activePredicates(root, cb);
            if (!request.getNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.get("name"), request.getNamePrefix() + "%"));
            }
            addCursor(predicates, root.get("updatedAt"), root.get("id"), cb, cursorUpdatedAt, request.getCursorId());
            return cb.and(predicates.toArray(new Predicate[0]));
        };
        List<Item> items = fetchPage(itemRepository::findAll, spec, pageSize);
        return buildListItemResp(items, pageSize);
    }

    @Transactional
    public CreateProcessDraftResp createProcessDraft(CreateProcessDraftReq request) {
        if (request.getOwnerUserId() <= 0) {
            throw new InventoryException("owner user id must be positive");
        }
        Item outputItem = requireItem(request.getItemId());
        List<ProcessItem> processItems = buildProcessItems(request.getItemsList());
        ProcessEntity process = new ProcessEntity();
        process.setItem(outputItem);
        process.setOwnerUserId(request.getOwnerUserId());
        process.setName(requireName(request.getName(), "process name is required"));
        process.setDescription(valueOrEmpty(request.getDescription()));
        process.setStatus(DraftStatus.DRAFT);
        processItems.forEach(process::addItem);
        processRepository.save(process);
        return CreateProcessDraftResp.newBuilder().setId(process.getId()).build();
    }

    @Transactional
    public UpdateProcessDraftResp updateProcessDraft(UpdateProcessDraftReq request) {
        ProcessEntity process = requireProcessForUpdate(request.getId());
        if (process.getStatus() != DraftStatus.DRAFT) {
            throw new InventoryException("only draft process can be updated");
        }
        if (request.getOwnerUserId() <= 0) {
            throw new InventoryException("owner user id must be positive");
        }
        process.setItem(requireItem(request.getItemId()));
        process.setOwnerUserId(request.getOwnerUserId());
        process.setName(requireName(request.getName(), "process name is required"));
        process.setDescription(valueOrEmpty(request.getDescription()));
        process.replaceItems(buildProcessItems(request.getItemsList()));
        return UpdateProcessDraftResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public DeleteProcessDraftResp deleteProcessDraft(DeleteProcessDraftReq request) {
        ProcessEntity process = requireProcessForUpdate(request.getId());
        if (process.getStatus() != DraftStatus.DRAFT) {
            throw new InventoryException("only draft process can be deleted");
        }
        process.softDelete();
        process.getItems().forEach(ProcessItem::softDelete);
        return DeleteProcessDraftResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public SubmitProcessResp submitProcess(SubmitProcessReq request) {
        ProcessEntity process = requireProcessForUpdate(request.getId());
        if (process.getItems().stream().filter(i -> i.getDeletedAt() == null).count() == 0) {
            throw new InventoryException("process items are required");
        }
        if (process.getStatus() != DraftStatus.DRAFT) {
            throw new InventoryException("only draft process can be submitted");
        }
        process.setStatus(DraftStatus.SUBMITTED);
        return SubmitProcessResp.newBuilder().setSuccess(true).build();
    }

    @Transactional(readOnly = true)
    public GetProcessResp getProcess(GetProcessReq request) {
        ProcessEntity process = processRepository.findWithItemsByIdAndDeletedAtIsNull(requirePositiveId(request.getId(), "process id"))
                .orElseThrow(() -> new InventoryException("process not found"));
        return GetProcessResp.newBuilder().setProcess(toProcessInfo(process, true)).build();
    }

    @Transactional(readOnly = true)
    public ListProcessResp listProcess(ListProcessReq request) {
        long pageSize = normalizePageSize(request.getPageSize());
        Instant sinceTime = parseSinceTime(request.getSinceTime(), request.getRecentSeconds());
        Instant cursorUpdatedAt = parseOptionalTime(request.getCursorUpdatedAt(), "cursorUpdatedAt");
        Specification<ProcessEntity> spec = (root, query, cb) -> {
            query.distinct(true);
            List<Predicate> predicates = activePredicates(root, cb);
            if (request.getOwnerUserId() > 0) {
                predicates.add(cb.equal(root.get("ownerUserId"), request.getOwnerUserId()));
            }
            if (request.getItemId() > 0) {
                predicates.add(cb.equal(root.get("item").get("id"), request.getItemId()));
            }
            DraftStatus status = DraftStatus.fromCode(request.getStatus().getNumber());
            if (status != DraftStatus.UNKNOWN) {
                predicates.add(cb.equal(root.get("status"), status));
            }
            if (!request.getNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.get("name"), request.getNamePrefix() + "%"));
            }
            if (!request.getItemNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.join("item", JoinType.INNER).get("name"), request.getItemNamePrefix() + "%"));
            }
            addSince(predicates, root.get("updatedAt"), cb, sinceTime);
            addCursor(predicates, root.get("updatedAt"), root.get("id"), cb, cursorUpdatedAt, request.getCursorId());
            return cb.and(predicates.toArray(new Predicate[0]));
        };
        List<ProcessEntity> processes = fetchPage(processRepository::findAll, spec, pageSize);
        ListProcessResp.Builder builder = ListProcessResp.newBuilder()
                .setTotal(Math.min(processes.size(), pageSize))
                .setHasMore(processes.size() > pageSize);
        trim(processes, pageSize).forEach(process -> builder.addProcessList(toProcessInfo(process, false)));
        setProcessCursor(builder, processes, pageSize);
        return builder.build();
    }

    @Transactional
    public AddItemUnitResp addItemUnit(AddItemUnitReq request) {
        if (!validStockStatus(request.getStockStatus().getNumber()) || !validQualityStatus(request.getQualityStatus().getNumber())) {
            throw new InventoryException("invalid item unit status");
        }
        Item item = requireItem(request.getItemId());
        ItemUnit unit = new ItemUnit();
        unit.setItem(item);
        unit.setStockStatus(StockStatus.OUT_STOCK);
        unit.setQualityStatus(QualityStatus.fromCode(request.getQualityStatus().getNumber()));
        unit.setDescription(valueOrEmpty(request.getDescription()));
        EngineeringOrder order = null;
        if (request.getEngineeringOrderId() > 0) {
            order = validateEngineeringOrderBinding(request.getEngineeringOrderId(), item.getId(), unit.getQualityStatus());
            unit.setEngineeringOrder(order);
        }
        itemUnitRepository.save(unit);
        recalculateItemCounts(item.getId());
        if (order != null) {
            recalculateEngineeringOrderProducedQuantity(order.getId());
        }
        return AddItemUnitResp.newBuilder().setId(unit.getId()).build();
    }

    @Transactional
    public UpdateItemUnitStatusResp updateItemUnitStatus(UpdateItemUnitStatusReq request) {
        if (!validStockStatus(request.getStockStatus().getNumber()) || !validQualityStatus(request.getQualityStatus().getNumber())) {
            throw new InventoryException("invalid item unit status");
        }
        ItemUnit unit = requireItemUnitForUpdate(request.getId());
        unit.setStockStatus(StockStatus.fromCode(request.getStockStatus().getNumber()));
        unit.setQualityStatus(QualityStatus.fromCode(request.getQualityStatus().getNumber()));
        recalculateItemCounts(unit.getItem().getId());
        if (unit.getEngineeringOrder() != null) {
            recalculateEngineeringOrderProducedQuantity(unit.getEngineeringOrder().getId());
        }
        return UpdateItemUnitStatusResp.newBuilder().setSuccess(true).build();
    }

    @Transactional(readOnly = true)
    public GetItemUnitResp getItemUnit(GetItemUnitReq request) {
        return GetItemUnitResp.newBuilder()
                .setItemUnit(toItemUnitInfo(requireItemUnit(request.getId())))
                .build();
    }

    @Transactional(readOnly = true)
    public ListItemUnitResp listItemUnit(ListItemUnitReq request) {
        long pageSize = normalizePageSize(request.getPageSize());
        Instant cursorUpdatedAt = parseOptionalTime(request.getCursorUpdatedAt(), "cursorUpdatedAt");
        Specification<ItemUnit> spec = (root, query, cb) -> {
            query.distinct(true);
            List<Predicate> predicates = activePredicates(root, cb);
            if (request.getItemId() > 0) {
                predicates.add(cb.equal(root.get("item").get("id"), request.getItemId()));
            }
            if (request.getEngineeringOrderId() > 0) {
                predicates.add(cb.equal(root.get("engineeringOrder").get("id"), request.getEngineeringOrderId()));
            }
            if (request.getInventoryFlowId() > 0) {
                predicates.add(cb.equal(root.join("inventoryFlows", JoinType.INNER).get("id"), request.getInventoryFlowId()));
            }
            StockStatus stockStatus = StockStatus.fromCode(request.getStockStatus().getNumber());
            if (stockStatus != StockStatus.UNKNOWN) {
                predicates.add(cb.equal(root.get("stockStatus"), stockStatus));
            }
            QualityStatus qualityStatus = QualityStatus.fromCode(request.getQualityStatus().getNumber());
            if (qualityStatus != QualityStatus.UNKNOWN) {
                predicates.add(cb.equal(root.get("qualityStatus"), qualityStatus));
            }
            if (!request.getItemNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.join("item", JoinType.INNER).get("name"), request.getItemNamePrefix() + "%"));
            }
            addCursor(predicates, root.get("updatedAt"), root.get("id"), cb, cursorUpdatedAt, request.getCursorId());
            return cb.and(predicates.toArray(new Predicate[0]));
        };
        List<ItemUnit> units = fetchPage(itemUnitRepository::findAll, spec, pageSize);
        ListItemUnitResp.Builder builder = ListItemUnitResp.newBuilder()
                .setTotal(Math.min(units.size(), pageSize))
                .setHasMore(units.size() > pageSize);
        trim(units, pageSize).forEach(unit -> builder.addItemUnitList(toItemUnitInfo(unit)));
        setItemUnitCursor(builder, units, pageSize);
        return builder.build();
    }

    @Transactional
    public CreateInventoryFlowResp createInventoryFlow(CreateInventoryFlowReq request) {
        FlowType flowType = requireFlowType(request.getFlowType().getNumber());
        InventoryFlow flow = new InventoryFlow();
        flow.setFromUserId(request.getFromUserId());
        flow.setToUserId(request.getToUserId());
        flow.setFlowType(flowType);
        flow.setFlowStatus(FlowStatus.DRAFT);
        flow.setName(requireName(request.getName(), "inventory flow name is required"));
        flow.setDescription(valueOrEmpty(request.getDescription()));
        buildFlowItems(request.getItemsList()).forEach(flow::addItem);
        flow.replaceItemUnits(findUnitsByIds(request.getItemUnitIdsList(), false));
        inventoryFlowRepository.save(flow);
        return CreateInventoryFlowResp.newBuilder().setId(flow.getId()).build();
    }

    @Transactional
    public UpdateInventoryFlowDraftResp updateInventoryFlowDraft(UpdateInventoryFlowDraftReq request) {
        InventoryFlow flow = requireFlowForUpdate(request.getId());
        if (flow.getFlowStatus() != FlowStatus.DRAFT) {
            throw new InventoryException("only draft inventory flow can be updated");
        }
        flow.setFromUserId(request.getFromUserId());
        flow.setToUserId(request.getToUserId());
        flow.setFlowType(requireFlowType(request.getFlowType().getNumber()));
        flow.setName(requireName(request.getName(), "inventory flow name is required"));
        flow.setDescription(valueOrEmpty(request.getDescription()));
        flow.replaceItems(buildFlowItems(request.getItemsList()));
        flow.replaceItemUnits(findUnitsByIds(request.getItemUnitIdsList(), false));
        return UpdateInventoryFlowDraftResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public DeleteInventoryFlowDraftResp deleteInventoryFlowDraft(DeleteInventoryFlowDraftReq request) {
        InventoryFlow flow = requireFlowForUpdate(request.getId());
        if (flow.getFlowStatus() != FlowStatus.DRAFT) {
            throw new InventoryException("only draft inventory flow can be deleted");
        }
        flow.getItemUnits().clear();
        flow.getItems().forEach(InventoryFlowItem::softDelete);
        flow.softDelete();
        return DeleteInventoryFlowDraftResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public SubmitInventoryFlowResp submitInventoryFlow(SubmitInventoryFlowReq request) {
        InventoryFlow flow = requireFlowForUpdate(request.getId());
        if (flow.getFlowStatus() != FlowStatus.DRAFT) {
            throw new InventoryException("only draft inventory flow can be submitted");
        }
        flow.setFlowStatus(FlowStatus.SUBMITTED);
        return SubmitInventoryFlowResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public AuditInventoryFlowResp auditInventoryFlow(AuditInventoryFlowReq request) {
        InventoryFlow flow = requireFlowForUpdate(request.getId());
        if (flow.getFlowStatus() != FlowStatus.SUBMITTED) {
            throw new InventoryException("only submitted inventory flow can be audited");
        }
        flow.setApprovedBy(request.getApprovedBy());
        flow.setApprovedAt(Instant.now());
        if (!request.getApproved()) {
            flow.setFlowStatus(FlowStatus.REJECTED);
        } else if (flow.getFlowType() == FlowType.IN || flow.getFlowType() == FlowType.OUT) {
            flow.setFlowStatus(FlowStatus.APPROVED);
        } else {
            throw new InventoryException("invalid flow type");
        }
        return AuditInventoryFlowResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public CompleteInventoryFlowResp completeInventoryFlow(CompleteInventoryFlowReq request) {
        InventoryFlow flow = requireFlowForUpdate(request.getId());
        if (flow.getFlowStatus() != FlowStatus.APPROVED) {
            throw new InventoryException("only approved inventory flow can be completed");
        }
        if (flow.getFlowType() != FlowType.IN && flow.getFlowType() != FlowType.OUT) {
            throw new InventoryException("invalid flow type");
        }
        List<InventoryFlowItem> details = activeFlowItems(flow);
        if (details.isEmpty()) {
            throw new InventoryException("inventory flow items are required");
        }
        List<ItemUnit> units = findUnitsByIds(request.getItemUnitIdsList(), true);
        Set<Long> existingUnitIds = flow.getItemUnits().stream().map(ItemUnit::getId).collect(Collectors.toSet());
        for (ItemUnit unit : units) {
            if (existingUnitIds.contains(unit.getId())) {
                throw new InventoryException("item unit " + unit.getId() + " has already been completed in this flow");
            }
        }
        Map<Long, Long> pendingByDetail = new LinkedHashMap<>();
        Set<Long> affectedItemIds = new LinkedHashSet<>();
        for (ItemUnit unit : units) {
            InventoryFlowItem detail = firstUnfinishedFlowItem(details, unit.getItem().getId(), pendingByDetail);
            if (detail == null) {
                throw new InventoryException("item unit " + unit.getId() + " has no unfinished matching flow item");
            }
            validateUnitForFlow(flow, unit);
            pendingByDetail.merge(detail.getId(), 1L, Long::sum);
            affectedItemIds.add(unit.getItem().getId());
        }
        flow.getItemUnits().addAll(units);
        for (Map.Entry<Long, Long> entry : pendingByDetail.entrySet()) {
            InventoryFlowItem detail = details.stream()
                    .filter(item -> item.getId().equals(entry.getKey()))
                    .findFirst()
                    .orElseThrow(() -> new InventoryException("inventory flow item " + entry.getKey() + " is invalid"));
            detail.setFinishedQuantity(detail.getFinishedQuantity() + entry.getValue());
        }
        for (ItemUnit unit : units) {
            if (flow.getFlowType() == FlowType.IN) {
                unit.setStockStatus(StockStatus.IN_STOCK);
                if (unit.getQualityStatus() == QualityStatus.UNKNOWN) {
                    unit.setQualityStatus(QualityStatus.PENDING);
                }
            } else {
                unit.setStockStatus(StockStatus.OUT_STOCK);
            }
        }
        affectedItemIds.forEach(this::recalculateItemCounts);
        return CompleteInventoryFlowResp.newBuilder().setSuccess(true).build();
    }

    @Transactional(readOnly = true)
    public GetInventoryFlowResp getInventoryFlow(GetInventoryFlowReq request) {
        InventoryFlow flow = inventoryFlowRepository.findWithDetailsByIdAndDeletedAtIsNull(requirePositiveId(request.getId(), "inventory flow id"))
                .orElseThrow(() -> new InventoryException("inventory flow not found"));
        return GetInventoryFlowResp.newBuilder().setInventoryFlow(toFlowInfo(flow)).build();
    }

    @Transactional(readOnly = true)
    public ListInventoryFlowResp listInventoryFlow(ListInventoryFlowReq request) {
        long pageSize = normalizePageSize(request.getPageSize());
        ListScope scope = request.getScope();
        boolean filterUser = scope != ListScope.LIST_SCOPE_ALL && scope != ListScope.LIST_SCOPE_AUDIT;
        if (filterUser && request.getUserId() <= 0) {
            throw new InventoryException("user id must be positive");
        }
        Instant sinceTime = parseSinceTime(request.getSinceTime(), request.getRecentSeconds());
        Instant cursorUpdatedAt = parseOptionalTime(request.getCursorUpdatedAt(), "cursorUpdatedAt");
        FlowStatus flowStatus = FlowStatus.fromCode(request.getFlowStatus().getNumber());
        if (scope == ListScope.LIST_SCOPE_AUDIT && flowStatus == FlowStatus.UNKNOWN) {
            flowStatus = FlowStatus.SUBMITTED;
        }
        FlowStatus finalFlowStatus = flowStatus;
        Specification<InventoryFlow> spec = (root, query, cb) -> {
            query.distinct(true);
            List<Predicate> predicates = activePredicates(root, cb);
            if (filterUser) {
                predicates.add(cb.equal(root.get(request.getIsTo() ? "toUserId" : "fromUserId"), request.getUserId()));
            }
            if (finalFlowStatus.isKnown()) {
                predicates.add(cb.equal(root.get("flowStatus"), finalFlowStatus));
            } else if (!filterUser) {
                predicates.add(cb.notEqual(root.get("flowStatus"), FlowStatus.DRAFT));
            }
            if (!request.getNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.get("name"), request.getNamePrefix() + "%"));
            }
            if (!request.getItemNamePrefix().isEmpty()) {
                String like = request.getItemNamePrefix() + "%";
                var detailItem = root.join("items", JoinType.LEFT).join("item", JoinType.LEFT);
                var unitItem = root.join("itemUnits", JoinType.LEFT).join("item", JoinType.LEFT);
                predicates.add(cb.or(cb.like(detailItem.get("name"), like), cb.like(unitItem.get("name"), like)));
            }
            if (request.getItemUnitId() > 0) {
                predicates.add(cb.equal(root.join("itemUnits", JoinType.INNER).get("id"), request.getItemUnitId()));
            }
            addSince(predicates, root.get("updatedAt"), cb, sinceTime);
            addCursor(predicates, root.get("updatedAt"), root.get("id"), cb, cursorUpdatedAt, request.getCursorId());
            return cb.and(predicates.toArray(new Predicate[0]));
        };
        List<InventoryFlow> flows = fetchPage(inventoryFlowRepository::findAll, spec, pageSize);
        ListInventoryFlowResp.Builder builder = ListInventoryFlowResp.newBuilder()
                .setTotal(Math.min(flows.size(), pageSize))
                .setHasMore(flows.size() > pageSize);
        trim(flows, pageSize).forEach(flow -> builder.addInventoryFlowList(toFlowInfo(flow)));
        setFlowCursor(builder, flows, pageSize);
        return builder.build();
    }

    @Transactional
    public CreateEngineeringOrderDraftResp createEngineeringOrderDraft(CreateEngineeringOrderDraftReq request) {
        if (request.getLeaderUserId() <= 0) {
            throw new InventoryException("leader user id must be positive");
        }
        validateEngineeringQuantities(request.getExpectedQuantity(), 0L);
        ProcessEntity process = getSubmittedProcessForEngineering(request.getProcessId(), request.getItemId());
        EngineeringOrder order = new EngineeringOrder();
        order.setLeaderUserId(request.getLeaderUserId());
        order.setProcess(process);
        order.setItem(process.getItem());
        order.setName(requireName(request.getName(), "engineering order name is required"));
        order.setExpectedQuantity(request.getExpectedQuantity());
        order.setStatus(DraftStatus.DRAFT);
        order.setDescription(valueOrEmpty(request.getDescription()));
        engineeringOrderRepository.save(order);
        return CreateEngineeringOrderDraftResp.newBuilder().setId(order.getId()).build();
    }

    @Transactional
    public UpdateEngineeringOrderDraftResp updateEngineeringOrderDraft(UpdateEngineeringOrderDraftReq request) {
        EngineeringOrder order = requireEngineeringOrderForUpdate(request.getId());
        if (order.getStatus() != DraftStatus.DRAFT) {
            throw new InventoryException("only draft engineering order can be updated");
        }
        if (request.getLeaderUserId() <= 0) {
            throw new InventoryException("leader user id must be positive");
        }
        validateEngineeringQuantities(request.getExpectedQuantity(), 0L);
        ProcessEntity process = getSubmittedProcessForEngineering(request.getProcessId(), request.getItemId());
        order.setLeaderUserId(request.getLeaderUserId());
        order.setProcess(process);
        order.setItem(process.getItem());
        order.setName(requireName(request.getName(), "engineering order name is required"));
        order.setExpectedQuantity(request.getExpectedQuantity());
        order.setDescription(valueOrEmpty(request.getDescription()));
        return UpdateEngineeringOrderDraftResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public DeleteEngineeringOrderDraftResp deleteEngineeringOrderDraft(DeleteEngineeringOrderDraftReq request) {
        EngineeringOrder order = requireEngineeringOrderForUpdate(request.getId());
        if (order.getStatus() != DraftStatus.DRAFT) {
            throw new InventoryException("only draft engineering order can be deleted");
        }
        order.softDelete();
        return DeleteEngineeringOrderDraftResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public SubmitEngineeringOrderResp submitEngineeringOrder(SubmitEngineeringOrderReq request) {
        EngineeringOrder order = requireEngineeringOrderForUpdate(request.getId());
        if (order.getStatus() != DraftStatus.DRAFT) {
            throw new InventoryException("only draft engineering order can be submitted");
        }
        order.setStatus(DraftStatus.SUBMITTED);
        return SubmitEngineeringOrderResp.newBuilder().setSuccess(true).build();
    }

    @Transactional
    public GetEngineeringOrderResp getEngineeringOrder(GetEngineeringOrderReq request) {
        Long id = requirePositiveId(request.getId(), "engineering order id");
        recalculateEngineeringOrderProducedQuantity(id);
        EngineeringOrder order = engineeringOrderRepository.findWithDetailsByIdAndDeletedAtIsNull(id)
                .orElseThrow(() -> new InventoryException("engineering order not found"));
        return GetEngineeringOrderResp.newBuilder().setEngineeringOrder(toEngineeringOrderInfo(order, false)).build();
    }

    @Transactional(readOnly = true)
    public ListEngineeringOrderResp listEngineeringOrder(ListEngineeringOrderReq request) {
        long pageSize = normalizePageSize(request.getPageSize());
        Instant sinceTime = parseSinceTime(request.getSinceTime(), request.getRecentSeconds());
        Instant cursorUpdatedAt = parseOptionalTime(request.getCursorUpdatedAt(), "cursorUpdatedAt");
        Specification<EngineeringOrder> spec = (root, query, cb) -> {
            query.distinct(true);
            List<Predicate> predicates = activePredicates(root, cb);
            if (request.getLeaderUserId() > 0) {
                predicates.add(cb.equal(root.get("leaderUserId"), request.getLeaderUserId()));
            }
            if (request.getItemId() > 0) {
                predicates.add(cb.equal(root.get("item").get("id"), request.getItemId()));
            }
            if (request.getProcessId() > 0) {
                predicates.add(cb.equal(root.get("process").get("id"), request.getProcessId()));
            }
            DraftStatus status = DraftStatus.fromCode(request.getStatus().getNumber());
            if (status != DraftStatus.UNKNOWN) {
                predicates.add(cb.equal(root.get("status"), status));
            }
            if (!request.getNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.get("name"), request.getNamePrefix() + "%"));
            }
            if (!request.getItemNamePrefix().isEmpty()) {
                predicates.add(cb.like(root.join("item", JoinType.INNER).get("name"), request.getItemNamePrefix() + "%"));
            }
            addSince(predicates, root.get("updatedAt"), cb, sinceTime);
            addCursor(predicates, root.get("updatedAt"), root.get("id"), cb, cursorUpdatedAt, request.getCursorId());
            return cb.and(predicates.toArray(new Predicate[0]));
        };
        List<EngineeringOrder> orders = fetchPage(engineeringOrderRepository::findAll, spec, pageSize);
        ListEngineeringOrderResp.Builder builder = ListEngineeringOrderResp.newBuilder()
                .setTotal(Math.min(orders.size(), pageSize))
                .setHasMore(orders.size() > pageSize);
        trim(orders, pageSize).forEach(order -> builder.addEngineeringOrderList(toEngineeringOrderInfo(order, false)));
        setEngineeringCursor(builder, orders, pageSize);
        return builder.build();
    }

    private List<ProcessItem> buildProcessItems(List<ProcessItemReq> requests) {
        if (requests.isEmpty()) {
            throw new InventoryException("process items are required");
        }
        Map<Long, Long> quantityByItem = new LinkedHashMap<>();
        for (ProcessItemReq request : requests) {
            Long itemId = requirePositiveId(request.getConsumeItemId(), "consume item id");
            if (request.getQuantity() <= 0) {
                throw new InventoryException("process item quantity must be positive");
            }
            quantityByItem.merge(itemId, request.getQuantity(), Long::sum);
        }
        return quantityByItem.entrySet().stream().map(entry -> {
            ProcessItem item = new ProcessItem();
            item.setConsumeItem(requireItem(entry.getKey()));
            item.setQuantity(entry.getValue());
            return item;
        }).toList();
    }

    private List<InventoryFlowItem> buildFlowItems(List<InventoryFlowItemReq> requests) {
        if (requests.isEmpty()) {
            throw new InventoryException("inventory flow items are required");
        }
        Map<Long, Long> quantityByItem = new LinkedHashMap<>();
        for (InventoryFlowItemReq request : requests) {
            Long itemId = requirePositiveId(request.getItemId(), "item id");
            if (request.getApplyQuantity() <= 0) {
                throw new InventoryException("apply quantity must be positive");
            }
            quantityByItem.merge(itemId, request.getApplyQuantity(), Long::sum);
        }
        return quantityByItem.entrySet().stream().map(entry -> {
            InventoryFlowItem item = new InventoryFlowItem();
            item.setItem(requireItem(entry.getKey()));
            item.setApplyQuantity(entry.getValue());
            item.setFinishedQuantity(0L);
            return item;
        }).toList();
    }

    private List<ItemUnit> findUnitsByIds(List<Long> ids, boolean forUpdate) {
        List<Long> uniqueIds = uniquePositiveIds(ids, "item unit id");
        if (uniqueIds.isEmpty()) {
            if (forUpdate) {
                throw new InventoryException("item units are required");
            }
            return List.of();
        }
        List<ItemUnit> units = forUpdate
                ? itemUnitRepository.findForUpdateByIds(uniqueIds)
                : itemUnitRepository.findByIdInAndDeletedAtIsNull(uniqueIds);
        if (units.size() != uniqueIds.size()) {
            throw new InventoryException("inventory flow contains unknown item unit");
        }
        Map<Long, ItemUnit> byId = units.stream().collect(Collectors.toMap(ItemUnit::getId, Function.identity()));
        return uniqueIds.stream().map(byId::get).toList();
    }

    private EngineeringOrder validateEngineeringOrderBinding(Long orderId, Long itemId, QualityStatus qualityStatus) {
        if (!validQualityStatus(qualityStatus.getCode())) {
            throw new InventoryException("invalid item unit quality status");
        }
        EngineeringOrder order = requireEngineeringOrderForUpdate(orderId);
        if (order.getStatus() != DraftStatus.SUBMITTED) {
            throw new InventoryException("item unit can only bind submitted engineering order");
        }
        if (!order.getItem().getId().equals(itemId)) {
            throw new InventoryException("item unit item does not match engineering order item");
        }
        recalculateEngineeringOrderProducedQuantity(orderId);
        if (order.getProducedQuantity() + 1 > order.getExpectedQuantity()) {
            throw new InventoryException("engineering order produced quantity exceeds expected quantity");
        }
        return order;
    }

    private ProcessEntity getSubmittedProcessForEngineering(Long processId, long requestItemId) {
        ProcessEntity process = requireProcessForUpdate(processId);
        if (process.getStatus() != DraftStatus.SUBMITTED) {
            throw new InventoryException("engineering order requires submitted process");
        }
        if (requestItemId > 0 && process.getItem().getId() != requestItemId) {
            throw new InventoryException("engineering order item does not match process item");
        }
        return process;
    }

    private void validateEngineeringQuantities(long expectedQuantity, long producedQuantity) {
        if (expectedQuantity <= 0) {
            throw new InventoryException("expected quantity must be positive");
        }
        if (expectedQuantity < producedQuantity) {
            throw new InventoryException("expected quantity cannot be less than produced quantity");
        }
    }

    private void validateUnitForFlow(InventoryFlow flow, ItemUnit unit) {
        if (flow.getFlowType() == FlowType.IN) {
            if (unit.getStockStatus() != StockStatus.OUT_STOCK) {
                throw new InventoryException("item unit " + unit.getId() + " is already in stock or reserved");
            }
            if (unit.getQualityStatus() != QualityStatus.QUALIFIED) {
                throw new InventoryException("item unit " + unit.getId() + " is not qualified");
            }
            return;
        }
        if (unit.getStockStatus() != StockStatus.IN_STOCK) {
            throw new InventoryException("item unit " + unit.getId() + " is not in stock");
        }
        if (unit.getQualityStatus() != QualityStatus.QUALIFIED) {
            throw new InventoryException("item unit " + unit.getId() + " is not qualified");
        }
    }

    private InventoryFlowItem firstUnfinishedFlowItem(
            List<InventoryFlowItem> details,
            Long itemId,
            Map<Long, Long> pendingByDetail) {
        for (InventoryFlowItem detail : details) {
            if (!detail.getItem().getId().equals(itemId)) {
                continue;
            }
            long pending = pendingByDetail.getOrDefault(detail.getId(), 0L);
            if (detail.getFinishedQuantity() + pending < detail.getApplyQuantity()) {
                return detail;
            }
        }
        return null;
    }

    private List<InventoryFlowItem> activeFlowItems(InventoryFlow flow) {
        return flow.getItems().stream()
                .filter(item -> item.getDeletedAt() == null)
                .sorted(Comparator.comparing(InventoryFlowItem::getId))
                .toList();
    }

    private void recalculateItemCounts(Long itemId) {
        Item item = requireItem(itemId);
        item.setTotalCount(itemUnitRepository.countByItemIdAndDeletedAtIsNull(itemId));
        item.setInStockCount(itemUnitRepository.countByItemIdAndStockStatusAndDeletedAtIsNull(itemId, StockStatus.IN_STOCK));
        item.setReservedCount(itemUnitRepository.countByItemIdAndStockStatusAndDeletedAtIsNull(itemId, StockStatus.RESERVED));
        item.setOutStockCount(itemUnitRepository.countByItemIdAndStockStatusAndDeletedAtIsNull(itemId, StockStatus.OUT_STOCK));
        item.setPendingCount(itemUnitRepository.countByItemIdAndQualityStatusAndDeletedAtIsNull(itemId, QualityStatus.PENDING));
        item.setQualifiedCount(itemUnitRepository.countByItemIdAndQualityStatusAndDeletedAtIsNull(itemId, QualityStatus.QUALIFIED));
        item.setUnqualifiedCount(itemUnitRepository.countByItemIdAndQualityStatusAndDeletedAtIsNull(itemId, QualityStatus.UNQUALIFIED));
        item.setAvailableCount(itemUnitRepository.countByItemIdAndStockStatusAndQualityStatusAndDeletedAtIsNull(
                itemId,
                StockStatus.IN_STOCK,
                QualityStatus.QUALIFIED));
    }

    private void recalculateEngineeringOrderProducedQuantity(Long orderId) {
        EngineeringOrder order = engineeringOrderRepository.findByIdAndDeletedAtIsNull(orderId)
                .orElseThrow(() -> new InventoryException("engineering order not found"));
        order.setProducedQuantity(itemUnitRepository.countByEngineeringOrderIdAndDeletedAtIsNull(orderId));
        order.setQualifiedQuantity(itemUnitRepository.countByEngineeringOrderIdAndQualityStatusAndDeletedAtIsNull(orderId, QualityStatus.QUALIFIED));
        order.setUnqualifiedQuantity(itemUnitRepository.countByEngineeringOrderIdAndQualityStatusAndDeletedAtIsNull(orderId, QualityStatus.UNQUALIFIED));
    }

    private Item requireItem(long id) {
        return itemRepository.findByIdAndDeletedAtIsNull(requirePositiveId(id, "item id"))
                .orElseThrow(() -> new InventoryException("item not found"));
    }

    private ProcessEntity requireProcessForUpdate(long id) {
        return processRepository.findForUpdate(requirePositiveId(id, "process id"))
                .orElseThrow(() -> new InventoryException("process not found"));
    }

    private ItemUnit requireItemUnit(long id) {
        return itemUnitRepository.findByIdAndDeletedAtIsNull(requirePositiveId(id, "item unit id"))
                .orElseThrow(() -> new InventoryException("item unit not found"));
    }

    private ItemUnit requireItemUnitForUpdate(long id) {
        List<ItemUnit> units = itemUnitRepository.findForUpdateByIds(List.of(requirePositiveId(id, "item unit id")));
        if (units.isEmpty()) {
            throw new InventoryException("item unit not found");
        }
        return units.get(0);
    }

    private InventoryFlow requireFlowForUpdate(long id) {
        return inventoryFlowRepository.findForUpdateWithDetails(requirePositiveId(id, "inventory flow id"))
                .orElseThrow(() -> new InventoryException("inventory flow not found"));
    }

    private EngineeringOrder requireEngineeringOrderForUpdate(long id) {
        return engineeringOrderRepository.findForUpdate(requirePositiveId(id, "engineering order id"))
                .orElseThrow(() -> new InventoryException("engineering order not found"));
    }

    private FlowType requireFlowType(int code) {
        FlowType flowType = FlowType.fromCode(code);
        if (flowType != FlowType.IN && flowType != FlowType.OUT) {
            throw new InventoryException("invalid flow type");
        }
        return flowType;
    }

    private boolean validStockStatus(int code) {
        StockStatus status = StockStatus.fromCode(code);
        return status == StockStatus.UNKNOWN || status == StockStatus.IN_STOCK || status == StockStatus.OUT_STOCK;
    }

    private boolean validQualityStatus(int code) {
        QualityStatus status = QualityStatus.fromCode(code);
        return status == QualityStatus.UNKNOWN || status == QualityStatus.PENDING
                || status == QualityStatus.QUALIFIED || status == QualityStatus.UNQUALIFIED;
    }

    private Long requirePositiveId(long id, String name) {
        if (id <= 0) {
            throw new InventoryException(name + " must be positive");
        }
        return id;
    }

    private List<Long> uniquePositiveIds(List<Long> ids, String name) {
        List<Long> result = new ArrayList<>();
        Set<Long> seen = new LinkedHashSet<>();
        for (Long id : ids) {
            Long value = requirePositiveId(id, name);
            if (seen.add(value)) {
                result.add(value);
            }
        }
        return result;
    }

    private String requireName(String name, String message) {
        String trimmed = name == null ? "" : name.trim();
        if (trimmed.isEmpty()) {
            throw new InventoryException(message);
        }
        return trimmed;
    }

    private String valueOrEmpty(String value) {
        return value == null ? "" : value;
    }

    private long normalizePageSize(long pageSize) {
        if (pageSize <= 0) {
            return DEFAULT_PAGE_SIZE;
        }
        return Math.min(pageSize, MAX_PAGE_SIZE);
    }

    private Instant parseSinceTime(String sinceTime, long recentSeconds) {
        if (sinceTime != null && !sinceTime.trim().isEmpty()) {
            return parseOptionalTime(sinceTime, "sinceTime");
        }
        return null;
    }

    private Instant parseOptionalTime(String value, String fieldName) {
        String trimmed = value == null ? "" : value.trim();
        if (trimmed.isEmpty()) {
            return null;
        }
        try {
            return LocalDateTime.parse(trimmed, LIST_TIME_FORMATTER)
                    .atZone(ZoneId.systemDefault())
                    .toInstant();
        } catch (DateTimeParseException ignored) {
            try {
                return OffsetDateTime.parse(trimmed, DateTimeFormatter.ISO_OFFSET_DATE_TIME).toInstant();
            } catch (DateTimeParseException ignoredAgain) {
                try {
                    return Instant.parse(trimmed);
                } catch (DateTimeParseException ignoredThird) {
                    throw new InventoryException(fieldName + " must use format yyyy-MM-dd HH:mm:ss");
                }
            }
        }
    }

    private String formatTime(Instant instant) {
        if (instant == null) {
            return "";
        }
        return DateTimeFormatter.ISO_OFFSET_DATE_TIME.format(instant.atOffset(ZoneOffset.UTC));
    }

    private <T extends com.moscenix.mes.inventory.domain.BaseEntity> List<T> trim(List<T> rows, long pageSize) {
        return rows.size() > pageSize ? rows.subList(0, Math.toIntExact(pageSize)) : rows;
    }

    private <T> List<T> fetchPage(
            Finder<T> finder,
            Specification<T> spec,
            long pageSize) {
        return finder.findAll(
                spec,
                PageRequest.of(
                        0,
                        Math.toIntExact(pageSize + 1),
                        Sort.by(Sort.Order.desc("updatedAt"), Sort.Order.desc("id"))))
                .getContent();
    }

    private <T> List<Predicate> activePredicates(jakarta.persistence.criteria.Root<T> root, jakarta.persistence.criteria.CriteriaBuilder cb) {
        List<Predicate> predicates = new ArrayList<>();
        predicates.add(cb.isNull(root.get("deletedAt")));
        return predicates;
    }

    private void addSince(
            List<Predicate> predicates,
            jakarta.persistence.criteria.Path<Instant> updatedAt,
            jakarta.persistence.criteria.CriteriaBuilder cb,
            Instant sinceTime) {
        if (sinceTime != null) {
            predicates.add(cb.greaterThan(updatedAt, sinceTime));
        }
    }

    private void addCursor(
            List<Predicate> predicates,
            jakarta.persistence.criteria.Path<Instant> updatedAt,
            jakarta.persistence.criteria.Path<Long> idPath,
            jakarta.persistence.criteria.CriteriaBuilder cb,
            Instant cursorUpdatedAt,
            long cursorId) {
        if (cursorUpdatedAt != null && cursorId > 0) {
            predicates.add(cb.or(
                    cb.lessThan(updatedAt, cursorUpdatedAt),
                    cb.and(cb.equal(updatedAt, cursorUpdatedAt), cb.lessThan(idPath, cursorId))));
        }
    }

    private ListItemResp buildListItemResp(List<Item> fetched, long pageSize) {
        List<Item> items = trim(fetched, pageSize);
        ListItemResp.Builder builder = ListItemResp.newBuilder()
                .setTotal(items.size())
                .setHasMore(fetched.size() > pageSize);
        items.forEach(item -> builder.addItemList(toItemInfo(item)));
        if (!items.isEmpty()) {
            Item last = items.get(items.size() - 1);
            builder.setNextCursorUpdatedAt(formatTime(last.getUpdatedAt())).setNextCursorId(last.getId());
        }
        return builder.build();
    }

    private void setProcessCursor(ListProcessResp.Builder builder, List<ProcessEntity> fetched, long pageSize) {
        List<ProcessEntity> rows = trim(fetched, pageSize);
        if (!rows.isEmpty()) {
            ProcessEntity last = rows.get(rows.size() - 1);
            builder.setNextCursorUpdatedAt(formatTime(last.getUpdatedAt())).setNextCursorId(last.getId());
        }
    }

    private void setItemUnitCursor(ListItemUnitResp.Builder builder, List<ItemUnit> fetched, long pageSize) {
        List<ItemUnit> rows = trim(fetched, pageSize);
        if (!rows.isEmpty()) {
            ItemUnit last = rows.get(rows.size() - 1);
            builder.setNextCursorUpdatedAt(formatTime(last.getUpdatedAt())).setNextCursorId(last.getId());
        }
    }

    private void setFlowCursor(ListInventoryFlowResp.Builder builder, List<InventoryFlow> fetched, long pageSize) {
        List<InventoryFlow> rows = trim(fetched, pageSize);
        if (!rows.isEmpty()) {
            InventoryFlow last = rows.get(rows.size() - 1);
            builder.setNextCursorUpdatedAt(formatTime(last.getUpdatedAt())).setNextCursorId(last.getId());
        }
    }

    private void setEngineeringCursor(ListEngineeringOrderResp.Builder builder, List<EngineeringOrder> fetched, long pageSize) {
        List<EngineeringOrder> rows = trim(fetched, pageSize);
        if (!rows.isEmpty()) {
            EngineeringOrder last = rows.get(rows.size() - 1);
            builder.setNextCursorUpdatedAt(formatTime(last.getUpdatedAt())).setNextCursorId(last.getId());
        }
    }

    private ItemInfo toItemInfo(Item item) {
        return ItemInfo.newBuilder()
                .setId(value(item.getId()))
                .setName(valueOrEmpty(item.getName()))
                .setUnit(valueOrEmpty(item.getUnit()))
                .setDescription(valueOrEmpty(item.getDescription()))
                .setTotalCount(value(item.getTotalCount()))
                .setInStockCount(value(item.getInStockCount()))
                .setReservedCount(value(item.getReservedCount()))
                .setOutStockCount(value(item.getOutStockCount()))
                .setPendingCount(value(item.getPendingCount()))
                .setQualifiedCount(value(item.getQualifiedCount()))
                .setUnqualifiedCount(value(item.getUnqualifiedCount()))
                .setAvailableCount(value(item.getAvailableCount()))
                .setCreateTime(formatTime(item.getCreatedAt()))
                .setUpdateTime(formatTime(item.getUpdatedAt()))
                .build();
    }

    private ItemUnitInfo toItemUnitInfo(ItemUnit unit) {
        ItemUnitInfo.Builder builder = ItemUnitInfo.newBuilder()
                .setId(value(unit.getId()))
                .setItemId(value(unit.getItem().getId()))
                .setStockStatus(com.moscenix.mes.grpc.inventory.StockStatus.forNumber(unit.getStockStatus().getCode()))
                .setQualityStatus(com.moscenix.mes.grpc.inventory.QualityStatus.forNumber(unit.getQualityStatus().getCode()))
                .setDescription(valueOrEmpty(unit.getDescription()))
                .setCreateTime(formatTime(unit.getCreatedAt()))
                .setUpdateTime(formatTime(unit.getUpdatedAt()));
        if (unit.getEngineeringOrder() != null) {
            builder.setEngineeringOrderId(value(unit.getEngineeringOrder().getId()));
        }
        return builder.build();
    }

    private ProcessItemInfo toProcessItemInfo(ProcessItem item) {
        return ProcessItemInfo.newBuilder()
                .setId(value(item.getId()))
                .setProcessId(value(item.getProcess().getId()))
                .setConsumeItemId(value(item.getConsumeItem().getId()))
                .setQuantity(value(item.getQuantity()))
                .setConsumeItem(toItemInfo(item.getConsumeItem()))
                .build();
    }

    private ProcessInfo toProcessInfo(ProcessEntity process, boolean withItems) {
        ProcessInfo.Builder builder = ProcessInfo.newBuilder()
                .setId(value(process.getId()))
                .setItemId(value(process.getItem().getId()))
                .setOwnerUserId(value(process.getOwnerUserId()))
                .setName(valueOrEmpty(process.getName()))
                .setDescription(valueOrEmpty(process.getDescription()))
                .setStatus(com.moscenix.mes.grpc.inventory.DraftStatus.forNumber(process.getStatus().getCode()))
                .setItem(toItemInfo(process.getItem()))
                .setCreateTime(formatTime(process.getCreatedAt()))
                .setUpdateTime(formatTime(process.getUpdatedAt()));
        if (withItems) {
            process.getItems().stream()
                    .filter(item -> item.getDeletedAt() == null)
                    .forEach(item -> builder.addItems(toProcessItemInfo(item)));
        }
        return builder.build();
    }

    private InventoryFlowItemInfo toFlowItemInfo(InventoryFlowItem item) {
        return InventoryFlowItemInfo.newBuilder()
                .setId(value(item.getId()))
                .setInventoryFlowId(value(item.getInventoryFlow().getId()))
                .setItemId(value(item.getItem().getId()))
                .setApplyQuantity(value(item.getApplyQuantity()))
                .setFinishedQuantity(value(item.getFinishedQuantity()))
                .setItem(toItemInfo(item.getItem()))
                .build();
    }

    private InventoryFlowInfo toFlowInfo(InventoryFlow flow) {
        InventoryFlowInfo.Builder builder = InventoryFlowInfo.newBuilder()
                .setId(value(flow.getId()))
                .setFromUserId(value(flow.getFromUserId()))
                .setToUserId(value(flow.getToUserId()))
                .setFlowType(com.moscenix.mes.grpc.inventory.FlowType.forNumber(flow.getFlowType().getCode()))
                .setFlowStatus(com.moscenix.mes.grpc.inventory.FlowStatus.forNumber(flow.getFlowStatus().getCode()))
                .setDescription(valueOrEmpty(flow.getDescription()))
                .setApprovedBy(value(flow.getApprovedBy()))
                .setApprovedAt(formatTime(flow.getApprovedAt()))
                .setCreateTime(formatTime(flow.getCreatedAt()))
                .setUpdateTime(formatTime(flow.getUpdatedAt()))
                .setName(valueOrEmpty(flow.getName()));
        activeFlowItems(flow).forEach(item -> builder.addItems(toFlowItemInfo(item)));
        flow.getItemUnits().stream()
                .filter(unit -> unit.getDeletedAt() == null)
                .forEach(unit -> builder.addItemUnits(toItemUnitInfo(unit)));
        return builder.build();
    }

    private EngineeringOrderInfo toEngineeringOrderInfo(EngineeringOrder order, boolean withUnits) {
        EngineeringOrderInfo.Builder builder = EngineeringOrderInfo.newBuilder()
                .setId(value(order.getId()))
                .setLeaderUserId(value(order.getLeaderUserId()))
                .setItemId(value(order.getItem().getId()))
                .setItem(toItemInfo(order.getItem()))
                .setExpectedQuantity(value(order.getExpectedQuantity()))
                .setQualifiedQuantity(value(order.getQualifiedQuantity()))
                .setProducedQuantity(value(order.getProducedQuantity()))
                .setDescription(valueOrEmpty(order.getDescription()))
                .setCreateTime(formatTime(order.getCreatedAt()))
                .setUpdateTime(formatTime(order.getUpdatedAt()))
                .setProcessId(value(order.getProcess().getId()))
                .setProcess(toProcessInfo(order.getProcess(), false))
                .setStatus(com.moscenix.mes.grpc.inventory.DraftStatus.forNumber(order.getStatus().getCode()))
                .setUnqualifiedQuantity(value(order.getUnqualifiedQuantity()))
                .setName(valueOrEmpty(order.getName()));
        if (withUnits) {
            order.getItemUnits().stream()
                    .filter(unit -> unit.getDeletedAt() == null)
                    .forEach(unit -> builder.addItemUnits(toItemUnitInfo(unit)));
        }
        return builder.build();
    }

    private long value(Long value) {
        return value == null ? 0L : value;
    }

    @FunctionalInterface
    private interface Finder<T> {
        org.springframework.data.domain.Page<T> findAll(Specification<T> spec, PageRequest pageRequest);
    }
}
