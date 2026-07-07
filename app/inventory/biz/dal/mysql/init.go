package mysql

import (
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
	DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.Item{}, &model.ItemUnit{}, &model.InventoryFlow{}, &model.InventoryFlowItem{})
	if err != nil {
		panic(err)
	}
}
