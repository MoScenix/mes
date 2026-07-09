package service

import "context"

type AuditInventoryFlowService struct {
	ctx context.Context
} // NewAuditInventoryFlowService new AuditInventoryFlowService
func NewAuditInventoryFlowService(ctx context.Context) *AuditInventoryFlowService {
	return &AuditInventoryFlowService{ctx: ctx}
}

// Run create note info
