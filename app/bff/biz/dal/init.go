package dal

import (
	"github.com/MoScenix/mes/app/bff/biz/dal/redis"
)

func Init() {
	redis.Init()
	//mysql.Init()
}
