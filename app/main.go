package main

import (
	"demo/app/router"
	"demo/common/conf"
	mcli "demo/common/data/imongo"

	"github.com/gin-gonic/gin"
)

func main() {
	conf.InitConfig()
	loadDataBase()
	startServer()
}

func loadDataBase() {
	mcli.Init()
}

func startServer() {
	r := gin.Default()
	router.InitRoute(r)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
