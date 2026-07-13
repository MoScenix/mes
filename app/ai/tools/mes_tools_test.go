package tools

import (
	"context"
	"testing"

	"github.com/MoScenix/mes/app/ai/conf"
	"github.com/MoScenix/mes/common/rpcmeta"
)

func TestToolNamesForRoleSkipsFinalEffectTools(t *testing.T) {
	cfg := conf.AITools{
		ToolGroups: map[string][]string{
			"workorder":         {"list_work_orders", "submit_work_order"},
			"engineering_order": {"create_engineering_order_draft", "update_engineering_order_draft", "list_engineering_orders"},
			"warehouse_admin":   {"list_pending_inventory_flows", "audit_inventory_flow"},
		},
		RoleGroups: map[string][]string{
			rpcmeta.RoleAdmin: {"workorder", "engineering_order", "warehouse_admin"},
		},
	}

	names := toolNamesForRole(rpcmeta.RoleAdmin, cfg)
	for _, forbidden := range []string{"submit_work_order", "audit_inventory_flow"} {
		for _, name := range names {
			if name == forbidden {
				t.Fatalf("toolNamesForRole injected forbidden tool %q in %v", forbidden, names)
			}
		}
	}

	for _, required := range []string{"create_engineering_order_draft", "update_engineering_order_draft"} {
		found := false
		for _, name := range names {
			if name == required {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("toolNamesForRole did not inject draft tool %q in %v", required, names)
		}
	}
}

func TestEffectiveUserIDForcesCurrentOperatorForNonAdmin(t *testing.T) {
	ctx := rpcmeta.WithOperator(context.Background(), 42, rpcmeta.RoleLeader)

	userID, err := effectiveUserID(ctx, 99)
	if err != nil {
		t.Fatal(err)
	}
	if userID != 42 {
		t.Fatalf("effectiveUserID() = %d, want current operator 42", userID)
	}

	optionalUserID, err := effectiveOptionalUserID(ctx, 99)
	if err != nil {
		t.Fatal(err)
	}
	if optionalUserID != 42 {
		t.Fatalf("effectiveOptionalUserID() = %d, want current operator 42", optionalUserID)
	}
}

func TestEffectiveUserIDAllowsAdminOverride(t *testing.T) {
	ctx := rpcmeta.WithOperator(context.Background(), 42, rpcmeta.RoleAdmin)

	userID, err := effectiveUserID(ctx, 99)
	if err != nil {
		t.Fatal(err)
	}
	if userID != 99 {
		t.Fatalf("effectiveUserID() = %d, want admin override 99", userID)
	}

	optionalUserID, err := effectiveOptionalUserID(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	if optionalUserID != 0 {
		t.Fatalf("effectiveOptionalUserID() = %d, want admin all-users filter 0", optionalUserID)
	}
}
