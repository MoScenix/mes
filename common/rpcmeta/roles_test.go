package rpcmeta

import "testing"

func TestNormalizeRole(t *testing.T) {
	cases := map[string]string{
		"组长":              RoleLeader,
		"leader":          RoleLeader,
		"采购专员":            RolePurchase,
		"purchase":        RolePurchase,
		"普通工人":            RoleWorker,
		"worker":          RoleWorker,
		"仓库管理员":           RoleWarehouseAdmin,
		"warehouse":       RoleWarehouseAdmin,
		"warehouse_admin": RoleWarehouseAdmin,
		"管理员":             RoleAdmin,
		"admin":           RoleAdmin,
		"":                RoleWorker,
	}
	for input, want := range cases {
		if got := NormalizeRole(input); got != want {
			t.Fatalf("NormalizeRole(%q)=%q, want %q", input, got, want)
		}
	}
}

func TestRoleCapabilities(t *testing.T) {
	if !CanCreateEngineeringOrder("组长") {
		t.Fatal("leader should be able to create engineering orders")
	}
	if CanCreateEngineeringOrder("采购专员") {
		t.Fatal("purchase should not be able to create engineering orders")
	}
	if !CanAuditInventoryFlow("warehouse") {
		t.Fatal("warehouse admin should be able to audit inventory flows")
	}
	if IsWarehouseAdmin("admin") {
		t.Fatal("admin should not normalize as warehouse admin")
	}
	if !CanAuditInventoryFlow("admin") {
		t.Fatal("admin should be able to audit inventory flows")
	}
	if CanAuditInventoryFlow("worker") {
		t.Fatal("worker should not be able to audit inventory flows")
	}
}
