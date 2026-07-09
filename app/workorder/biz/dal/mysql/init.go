package mysql

import (
	"fmt"
	"os"
	"strings"

	"github.com/MoScenix/mes/app/workorder/biz/model"
	"github.com/MoScenix/mes/app/workorder/conf"

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
	err = DB.AutoMigrate(&model.WorkOrder{})
	if err != nil {
		panic(err)
	}
	dropWorkOrderLegacyIndexes()
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func dropWorkOrderLegacyIndexes() {
	legacyIndexes := []string{
		"idx_from_user_created_at",
		"idx_to_user_created_at",
		"idx_from_user_read_status_created_at",
		"idx_to_user_read_status_created_at",
	}
	for _, indexName := range legacyIndexes {
		if DB.Migrator().HasIndex(&model.WorkOrder{}, indexName) {
			if err := DB.Migrator().DropIndex(&model.WorkOrder{}, indexName); err != nil {
				panic(err)
			}
		}
	}
}
