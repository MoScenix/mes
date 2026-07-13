package rpc

import (
	"context"
	"sync"

	"github.com/MoScenix/mes/app/bff/biz/utils"
	"github.com/MoScenix/mes/app/bff/conf"
	"github.com/MoScenix/mes/common/clientsuit"
	"github.com/MoScenix/mes/common/rpcmeta"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/ai/aiservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/app/appservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/document/documentservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory/inventoryservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/user/userservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder/workorderservice"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/endpoint"
)

var (
	UserClient      userservice.Client
	AppClient       appservice.Client
	AiClient        aiservice.Client
	DocumentClient  documentservice.Client
	WorkOrderClient workorderservice.Client
	InventoryClient inventoryservice.Client
	once            sync.Once
	once1           sync.Once
	once2           sync.Once
	once3           sync.Once
	once4           sync.Once
	once5           sync.Once
)

func Init() {
	once.Do(initUserClient)
	once1.Do(initAppClient)
	once2.Do(initAiClient)
	once3.Do(initDocumentClient)
	once4.Do(initWorkOrderClient)
	once5.Do(initInventoryClient)
}
func initUserClient() {
	opts := newCommonClientOptions(true)
	var err error
	UserClient, err = userservice.NewClient(
		"user",
		opts...,
	)
	if err != nil {
		hlog.Fatal(err)
	}
}
func initAppClient() {
	opts := newCommonClientOptions(false)
	var err error
	AppClient, err = appservice.NewClient(
		"app",
		opts...,
	)
	if err != nil {
		hlog.Fatal(err)
	}
}
func initAiClient() {
	opts := newCommonClientOptions(true)
	var err error
	AiClient, err = aiservice.NewClient(
		"ai",
		opts...,
	)
	if err != nil {
		hlog.Fatal(err)
	}
}
func initDocumentClient() {
	opts := newCommonClientOptions(false)
	var err error
	DocumentClient, err = documentservice.NewClient(
		"document",
		opts...,
	)
	if err != nil {
		hlog.Fatal(err)
	}
}
func initWorkOrderClient() {
	opts := newCommonClientOptions(true)
	var err error
	WorkOrderClient, err = workorderservice.NewClient(
		"workorder",
		opts...,
	)
	if err != nil {
		hlog.Fatal(err)
	}
}
func initInventoryClient() {
	opts := newCommonClientOptions(true)
	var err error
	InventoryClient, err = inventoryservice.NewClient(
		"inventory",
		opts...,
	)
	if err != nil {
		hlog.Fatal(err)
	}
}

func newCommonClientOptions(enableGRPC bool) []client.Option {
	opts := clientsuit.CommonGrpcClientSuite{
		CurrentServiceName: conf.GetConf().Hertz.Service,
		RegistryAddr:       conf.GetConf().Consul.Address,
		EnableGRPC:         enableGRPC,
	}.Options()

	opts = append(opts, client.WithMiddleware(identityMiddleware))
	return opts
}

func identityMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req, resp interface{}) error {
		if rpcmeta.FromContext(ctx).OperatorID == "" {
			if userID, ok := utils.UserIDFromContext(ctx); ok {
				role, _ := ctx.Value(utils.UserRoleKey).(string)
				ctx = rpcmeta.WithOperator(ctx, userID, role)
			}
		}
		return next(ctx, req, resp)
	}
}
