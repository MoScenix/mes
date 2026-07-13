package com.moscenix.mes.workorder.service;

import com.moscenix.mes.workorder.domain.WorkOrder;
import com.moscenix.mes.workorder.domain.WorkOrderReadStatus;
import com.moscenix.mes.workorder.domain.WorkOrderStatus;
import com.moscenix.mes.workorder.dto.CreateWorkOrderRequest;
import com.moscenix.mes.workorder.dto.CreateWorkOrderResponse;
import com.moscenix.mes.workorder.dto.GetWorkOrderResponse;
import com.moscenix.mes.workorder.dto.ListWorkOrderRequest;
import com.moscenix.mes.workorder.dto.ListWorkOrderResponse;
import com.moscenix.mes.workorder.dto.UpdateWorkOrderDraftRequest;
import com.moscenix.mes.workorder.dto.WorkOrderInfo;
import com.moscenix.mes.workorder.repository.WorkOrderRepository;
import jakarta.persistence.criteria.Predicate;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;
import java.util.ArrayList;
import java.util.List;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.data.jpa.domain.Specification;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class WorkOrderService {
    private static final DateTimeFormatter LIST_TIME_FORMATTER =
            DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    private final WorkOrderRepository workOrderRepository;

    public WorkOrderService(WorkOrderRepository workOrderRepository) {
        this.workOrderRepository = workOrderRepository;
    }

    @Transactional
    public CreateWorkOrderResponse create(CreateWorkOrderRequest request) {
        requireRequest(request);
        WorkOrder order = new WorkOrder();
        order.setFromUserId(request.getFromUserId());
        order.setToUserId(request.getToUserId());
        order.setName(requireName(request.getName()));
        order.setDescription(valueOrEmpty(request.getDescription()));
        order.setStatus(WorkOrderStatus.DRAFT);
        order.setReadStatus(WorkOrderReadStatus.UNREAD);

        WorkOrder saved = workOrderRepository.save(order);
        return new CreateWorkOrderResponse(saved.getId());
    }

    @Transactional
    public boolean updateDraft(Long id, UpdateWorkOrderDraftRequest request) {
        requireId(id);
        requireRequest(request);
        int affected = workOrderRepository.updateDraft(
                id,
                request.getFromUserId(),
                request.getToUserId(),
                requireName(request.getName()),
                valueOrEmpty(request.getDescription()),
                WorkOrderStatus.DRAFT,
                Instant.now());
        ensureDraftAffected(affected);
        return true;
    }

    @Transactional
    public boolean deleteDraft(Long id) {
        requireId(id);
        int affected = workOrderRepository.deleteDraft(id, WorkOrderStatus.DRAFT, Instant.now());
        ensureDraftAffected(affected);
        return true;
    }

    @Transactional
    public boolean submit(Long id) {
        requireId(id);
        int affected = workOrderRepository.submitDraft(
                id,
                WorkOrderStatus.DRAFT,
                WorkOrderStatus.SUBMITTED,
                Instant.now());
        ensureDraftAffected(affected);
        return true;
    }

    @Transactional(readOnly = true)
    public GetWorkOrderResponse get(Long id) {
        requireId(id);
        WorkOrder order = workOrderRepository.findByIdAndDeletedAtIsNull(id)
                .orElseThrow(() -> new WorkOrderNotFoundException(id));
        return new GetWorkOrderResponse(toInfo(order));
    }

    @Transactional(readOnly = true)
    public ListWorkOrderResponse list(ListWorkOrderRequest request) {
        requireRequest(request);
        long pageSize = normalizePageSize(request.getPageSize());
        WorkOrderStatus status = WorkOrderStatus.fromCode(request.getStatus());
        boolean isTo = Boolean.TRUE.equals(request.getIsTo());
        if (status == WorkOrderStatus.DRAFT) {
            isTo = false;
        }

        Instant sinceTime = parseOptionalListTime(request.getSinceTime(), "sinceTime");
        Instant cursorUpdatedAt = parseOptionalListTime(request.getCursorUpdatedAt(), "cursorUpdatedAt");
        long cursorId = request.getCursorId() == null ? 0L : request.getCursorId();

        Specification<WorkOrder> specification = buildListSpecification(
                request.getId(),
                isTo,
                Boolean.TRUE.equals(request.getIsUnread()),
                status,
                request.getNamePrefix(),
                sinceTime,
                cursorUpdatedAt,
                cursorId);

        List<WorkOrder> fetched = workOrderRepository.findAll(
                specification,
                PageRequest.of(
                        0,
                        Math.toIntExact(pageSize + 1),
                        Sort.by(Sort.Order.desc("updatedAt"), Sort.Order.desc("id"))))
                .getContent();

        boolean hasMore = fetched.size() > pageSize;
        List<WorkOrder> orders = hasMore ? fetched.subList(0, Math.toIntExact(pageSize)) : fetched;

        ListWorkOrderResponse response = new ListWorkOrderResponse();
        response.setTotal(orders.size());
        response.setHasMore(hasMore);
        response.setWorkOrderList(orders.stream().map(this::toInfo).toList());
        if (!orders.isEmpty()) {
            WorkOrder last = orders.get(orders.size() - 1);
            response.setNextCursorUpdatedAt(formatTime(last.getUpdatedAt()));
            response.setNextCursorId(last.getId());
        }
        return response;
    }

    @Transactional
    public boolean markRead(Long id) {
        requireId(id);
        int affected = workOrderRepository.markRead(id, WorkOrderReadStatus.READ, Instant.now());
        if (affected == 0) {
            throw new WorkOrderNotFoundException(id);
        }
        return true;
    }

    private Specification<WorkOrder> buildListSpecification(
            Long employeeId,
            boolean isTo,
            boolean onlyUnread,
            WorkOrderStatus status,
            String namePrefix,
            Instant sinceTime,
            Instant cursorUpdatedAt,
            long cursorId) {
        return (root, query, criteriaBuilder) -> {
            List<Predicate> predicates = new ArrayList<>();
            predicates.add(criteriaBuilder.isNull(root.get("deletedAt")));
            if (isTo) {
                predicates.add(criteriaBuilder.equal(root.<Long>get("toUserId"), employeeId));
            } else {
                predicates.add(criteriaBuilder.equal(root.<Long>get("fromUserId"), employeeId));
            }
            if (status.isKnown()) {
                predicates.add(criteriaBuilder.equal(root.<WorkOrderStatus>get("status"), status));
            } else if (isTo) {
                predicates.add(criteriaBuilder.notEqual(root.<WorkOrderStatus>get("status"), WorkOrderStatus.DRAFT));
            }
            if (onlyUnread) {
                predicates.add(criteriaBuilder.equal(root.<WorkOrderReadStatus>get("readStatus"), WorkOrderReadStatus.UNREAD));
            }
            if (namePrefix != null && !namePrefix.isEmpty()) {
                predicates.add(criteriaBuilder.like(root.<String>get("name"), namePrefix + "%"));
            }
            if (sinceTime != null) {
                predicates.add(criteriaBuilder.greaterThan(root.<Instant>get("updatedAt"), sinceTime));
            }
            if (cursorUpdatedAt != null && cursorId > 0) {
                predicates.add(criteriaBuilder.or(
                        criteriaBuilder.lessThan(root.<Instant>get("updatedAt"), cursorUpdatedAt),
                        criteriaBuilder.and(
                                criteriaBuilder.equal(root.<Instant>get("updatedAt"), cursorUpdatedAt),
                                criteriaBuilder.lessThan(root.<Long>get("id"), cursorId))));
            }
            return criteriaBuilder.and(predicates.toArray(new Predicate[0]));
        };
    }

    private WorkOrderInfo toInfo(WorkOrder order) {
        WorkOrderInfo info = new WorkOrderInfo();
        info.setId(order.getId());
        info.setFromUserId(order.getFromUserId());
        info.setToUserId(order.getToUserId());
        info.setName(order.getName());
        info.setDescription(order.getDescription());
        info.setStatus(order.getStatus());
        info.setReadStatus(order.getReadStatus());
        info.setCreateTime(formatTime(order.getCreatedAt()));
        info.setUpdateTime(formatTime(order.getUpdatedAt()));
        return info;
    }

    private String requireName(String name) {
        String trimmed = name == null ? "" : name.trim();
        if (trimmed.isEmpty()) {
            throw new WorkOrderBadRequestException("work order name is required");
        }
        return trimmed;
    }

    private void requireRequest(Object request) {
        if (request == null) {
            throw new WorkOrderBadRequestException("request body is required");
        }
    }

    private void requireId(Long id) {
        if (id == null) {
            throw new WorkOrderBadRequestException("id is required");
        }
    }

    private String valueOrEmpty(String value) {
        return value == null ? "" : value;
    }

    private long normalizePageSize(Long pageSize) {
        if (pageSize == null || pageSize <= 0) {
            return 10L;
        }
        return pageSize;
    }

    private void ensureDraftAffected(int affected) {
        if (affected == 0) {
            throw new DraftRequiredException();
        }
    }

    private Instant parseOptionalListTime(String value, String fieldName) {
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
                    throw new WorkOrderBadRequestException(fieldName + " must use format yyyy-MM-dd HH:mm:ss");
                }
            }
        }
    }

    private String formatTime(Instant instant) {
        return DateTimeFormatter.ISO_OFFSET_DATE_TIME.format(instant.atOffset(ZoneOffset.UTC));
    }
}
