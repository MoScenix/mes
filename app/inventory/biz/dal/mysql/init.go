package mysql

import (
	"fmt"
	"os"
	"strings"

	"github.com/MoScenix/mes/app/inventory/biz/model"
	"github.com/MoScenix/mes/app/inventory/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := conf.GetConf().MySQL.DSN
	if strings.Contains(dsn, "%s") {
		dsn = fmt.Sprintf(dsn, envOrDefault("MYSQL_USER", "root"), envOrDefault("MYSQL_PASSWORD", "YOUR_PASSWORD"), envOrDefault("MYSQL_HOST", "127.0.0.1"), envOrDefault("MYSQL_DATABASE", "YOU_DB"))
	}
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.Item{}, &model.Process{}, &model.ProcessItem{}, &model.EngineeringOrder{}, &model.ItemUnit{}, &model.InventoryFlow{}, &model.InventoryFlowItem{}, &model.InventoryFlowItemUnit{})
	if err != nil {
		panic(err)
	}
	dropInventoryFlowLegacyIndexes()
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func dropInventoryFlowLegacyIndexes() {
	legacyIndexes := []string{
		"idx_inventory_flow_filter",
		"idx_inventory_flow_from_time",
		"idx_inventory_flow_to_time",
		"idx_inventory_flow_from_status_time",
		"idx_inventory_flow_to_status_time",
		"idx_inventory_flow_from_deleted_time",
		"idx_inventory_flow_to_deleted_time",
	}
	for _, indexName := range legacyIndexes {
		if DB.Migrator().HasIndex(&model.InventoryFlow{}, indexName) {
			if err := DB.Migrator().DropIndex(&model.InventoryFlow{}, indexName); err != nil {
				panic(err)
			}
		}
	}
}
