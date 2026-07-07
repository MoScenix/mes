package dal

import (
	"github.com/MoScenix/mes/app/workorder/biz/dal/mysql"
)

func Init() {
	mysql.Init()
}
