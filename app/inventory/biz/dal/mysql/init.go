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
	dsn = withMultiStatements(dsn)
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
	runSQLFile()
}

func runSQLFile() {
	sqlBytes, err := os.ReadFile("biz/dal/mysql/mes_indexes.sql")
	if err != nil {
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	if _, err := sqlDB.Exec(string(sqlBytes)); err != nil {
		if !isIndexExistsError(err) {
			panic(err)
		}
	}
}

func isIndexExistsError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate key name")
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func withMultiStatements(dsn string) string {
	if strings.Contains(dsn, "multiStatements=") {
		return dsn
	}
	if strings.Contains(dsn, "?") {
		return dsn + "&multiStatements=true"
	}
	return dsn + "?multiStatements=true"
}
