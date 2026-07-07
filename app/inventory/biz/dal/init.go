package dal

import (
	"github.com/MoScenix/mes/app/inventory/biz/dal/mysql"
)

func Init() {
	mysql.Init()
}
