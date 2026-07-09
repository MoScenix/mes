package graph

import (
	"context"

	assistantnode "github.com/MoScenix/mes/app/ai/node/designer"
)

func runAssistantNode(ctx context.Context, input map[string]any) (output map[string]any, err error) {
	return assistantnode.Run(ctx)
}
