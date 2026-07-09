package graph

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func Buildaicode(ctx context.Context) (r compose.Runnable[any, any], err error) {
	g := compose.NewGraph[any, any]()
	_ = g.AddLambdaNode(assistantNode, compose.InvokableLambda(runAssistantNode))
	_ = g.AddEdge(compose.START, assistantNode)
	_ = g.AddEdge(assistantNode, compose.END)
	opts := []compose.GraphCompileOption{compose.WithGraphName("aicode")}
	if store, ok := newGraphCheckpointStore(ctx); ok {
		opts = append(opts, compose.WithCheckPointStore(store))
	}
	r, err = g.Compile(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return r, err
}
