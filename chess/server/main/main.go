package main

import (
	"chess/server/api"
	"chess/server/dao"
	"chess/server/tool"
)

func main() {
	tool.InitConfig()
	dao.InitMysql()
	dao.InitRedis()
	api.InitRouter()
}
