package rpcmeta

import "strings"

const (
	RoleLeader          = "leader"
	RolePurchase        = "purchase"
	RoleWorker          = "worker"
	RoleProcessEngineer = "process_engineer"
	RoleWarehouseAdmin  = "warehouse_admin"
	RoleSales           = "sales"
	RoleAdmin           = "admin"
)

var defaultRoleAliases = map[string]string{
	"组长":               RoleLeader,
	"leader":           RoleLeader,
	"采购专员":             RolePurchase,
	"purchase":         RolePurchase,
	"普通工人":             RoleWorker,
	"worker":           RoleWorker,
	"工艺工程师":            RoleProcessEngineer,
	"process_engineer": RoleProcessEngineer,
	"仓库管理员":            RoleWarehouseAdmin,
	"warehouse":        RoleWarehouseAdmin,
	"warehouse_admin":  RoleWarehouseAdmin,
	"销售":               RoleSales,
	"sales":            RoleSales,
	"管理员":              RoleAdmin,
	"admin":            RoleAdmin,
}

func NormalizeRole(role string, aliases ...map[string]string) string {
	trimmed := strings.TrimSpace(role)
	lower := strings.ToLower(trimmed)
	for _, aliasMap := range aliases {
		if aliasMap == nil {
			continue
		}
		if normalized := aliasMap[trimmed]; normalized != "" {
			return NormalizeRole(normalized)
		}
		if normalized := aliasMap[lower]; normalized != "" {
			return NormalizeRole(normalized)
		}
	}
	if normalized := defaultRoleAliases[trimmed]; normalized != "" {
		return normalized
	}
	if normalized := defaultRoleAliases[lower]; normalized != "" {
		return normalized
	}
	return RoleWorker
}

func IsAdmin(role string) bool {
	return NormalizeRole(role) == RoleAdmin
}

func IsWarehouseAdmin(role string) bool {
	return NormalizeRole(role) == RoleWarehouseAdmin
}

func CanCreateEngineeringOrder(role string) bool {
	normalized := NormalizeRole(role)
	return normalized == RoleLeader || normalized == RoleAdmin
}

func CanCreateInventoryFlow(role string) bool {
	normalized := NormalizeRole(role)
	return normalized == RoleLeader || normalized == RolePurchase || normalized == RoleAdmin
}

func CanAuditInventoryFlow(role string) bool {
	return IsWarehouseAdmin(role) || IsAdmin(role)
}
