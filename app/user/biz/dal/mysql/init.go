package mysql

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/MoScenix/mes/app/user/biz/model"
	"github.com/MoScenix/mes/app/user/conf"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, envOrDefault("MYSQL_USER", "root"), envOrDefault("MYSQL_PASSWORD", "YOUR_PASSWORD"), envOrDefault("MYSQL_HOST", "127.0.0.1"), envOrDefault("MYSQL_DATABASE", "YOU_DB"))
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	if err := ensureDefaultAdmin(); err != nil {
		panic(err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func ensureDefaultAdmin() error {
	ctx := context.Background()
	q := model.NewUserQuery(ctx, DB)
	root, err := q.GetUserByAccount("root")
	if err == nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte("rootroot"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		return q.UpdateUser(root.ID, model.User{
			Name:         "root",
			UserRole:     "admin",
			PasswordHash: string(hashed),
		})
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("rootroot"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = q.CreateUser(model.User{
		UserAccount:  "root",
		Name:         "root",
		UserRole:     "admin",
		PasswordHash: string(hashed),
	})
	return err
}
