package conf

import "testing"

func TestNormalizeAIToolsConfigAppendsMissingDefaults(t *testing.T) {
	cfg := AITools{
		ToolGroups: map[string][]string{
			"common": {"ask_user", "SearchProjectFile"},
		},
		RoleGroups: map[string][]string{
			"admin": {"common"},
		},
	}

	normalizeAIToolsConfig(&cfg)

	if !contains(cfg.ToolGroups["common"], "search_users") {
		t.Fatalf("common tools = %v, want search_users appended", cfg.ToolGroups["common"])
	}
	if !contains(cfg.RoleGroups["admin"], "workorder") {
		t.Fatalf("admin role groups = %v, want default groups appended", cfg.RoleGroups["admin"])
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
