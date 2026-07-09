package infra

import (
	"sync"

	"github.com/MoScenix/mes/app/ai/conf"
	"github.com/MoScenix/mes/common/clientsuit"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/app/appservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/document/documentservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/inventory/inventoryservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/user/userservice"
	"github.com/MoScenix/mes/rpc_gen/kitex_gen/workorder/workorderservice"
	"github.com/cloudwego/kitex/client"
)

var (
	appClient           appservice.Client
	appClientOnce       sync.Once
	appClientErr        error
	documentClient      documentservice.Client
	documentClientOnce  sync.Once
	documentClientErr   error
	inventoryClient     inventoryservice.Client
	inventoryClientOnce sync.Once
	inventoryClientErr  error
	userClient          userservice.Client
	userClientOnce      sync.Once
	userClientErr       error
	workOrderClient     workorderservice.Client
	workOrderClientOnce sync.Once
	workOrderClientErr  error
)

func AppClient() (appservice.Client, error) {
	appClientOnce.Do(func() {
		appClient, appClientErr = appservice.NewClient("app", newCommonClientOptions(false)...)
	})
	return appClient, appClientErr
}

func DocumentClient() (documentservice.Client, error) {
	documentClientOnce.Do(func() {
		documentClient, documentClientErr = documentservice.NewClient("document", newCommonClientOptions(false)...)
	})
	return documentClient, documentClientErr
}

func InventoryClient() (inventoryservice.Client, error) {
	inventoryClientOnce.Do(func() {
		inventoryClient, inventoryClientErr = inventoryservice.NewClient("inventory", newCommonClientOptions(false)...)
	})
	return inventoryClient, inventoryClientErr
}

func UserClient() (userservice.Client, error) {
	userClientOnce.Do(func() {
		userClient, userClientErr = userservice.NewClient("user", newCommonClientOptions(false)...)
	})
	return userClient, userClientErr
}

func WorkOrderClient() (workorderservice.Client, error) {
	workOrderClientOnce.Do(func() {
		workOrderClient, workOrderClientErr = workorderservice.NewClient("workorder", newCommonClientOptions(false)...)
	})
	return workOrderClient, workOrderClientErr
}

func newCommonClientOptions(enableGRPC bool) []client.Option {
	return clientsuit.CommonGrpcClientSuite{
		CurrentServiceName: conf.GetConf().Kitex.Service,
		RegistryAddr:       conf.GetConf().Registry.RegistryAddress[0],
		EnableGRPC:         enableGRPC,
	}.Options()
}
