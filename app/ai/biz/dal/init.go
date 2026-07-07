package dal

import (
	"github.com/MoScenix/mes/app/ai/biz/dal/redis"
)

func Init() {
	redis.Init()
	//mysql.Init()
}
