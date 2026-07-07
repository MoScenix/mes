package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MoScenix/mes/app/bff/conf"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	document "github.com/MoScenix/mes/rpc_gen/kitex_gen/document"
)

func deleteProjectFileData(ctx context.Context, projectID int64) error {
	if projectID <= 0 {
		return nil
	}
	if _, err := rpc.DocumentClient.DeleteProjectFileData(ctx, &document.DeleteProjectFileDataReq{
		ProjectId: projectID,
	}); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(conf.StaticRoot(), "document", fmt.Sprintf("%d", projectID))); err != nil {
		return err
	}
	return nil
}
