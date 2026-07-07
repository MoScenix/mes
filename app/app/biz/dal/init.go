package dal

import (
	"github.com/MoScenix/mes/app/app/biz/dal/mysql"
	"github.com/MoScenix/mes/app/app/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
