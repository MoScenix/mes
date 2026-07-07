package dal

import (
	"github.com/MoScenix/mes/app/document/biz/dal/mysql"
	"github.com/MoScenix/mes/app/document/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
