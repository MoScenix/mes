package tools

import (
	"context"
	"fmt"

	lutils "github.com/MoScenix/mes/app/ai/utils"
	"github.com/MoScenix/mes/common/filestore/project"
)

func projectStoreFromContext(ctx context.Context) (project.Store, error) {
	store, ok := ctx.Value(lutils.ProjectFileStore).(project.Store)
	if !ok || store == nil {
		return nil, fmt.Errorf("ProjectFileStore is missing in context")
	}
	return store, nil
}
