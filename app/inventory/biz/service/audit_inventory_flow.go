package service

import (
	"context"
	inventory "github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory"
)

type AuditInventoryFlowService struct {
	ctx context.Context
} // NewAuditInventoryFlowService new AuditInventoryFlowService
func NewAuditInventoryFlowService(ctx context.Context) *AuditInventoryFlowService {
	return &AuditInventoryFlowService{ctx: ctx}
}

// Run create note info
func (s *AuditInventoryFlowService) Run(req *inventory.AuditInventoryFlowReq) (resp *inventory.AuditInventoryFlowResp, err error) {
	return runAuditInventoryFlow(s.ctx, req)
}
