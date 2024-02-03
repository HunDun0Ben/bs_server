package main

import (
	"demo/app/router"
	mcli "demo/common/data/mongodb"

	"github.com/gin-gonic/gin"
)

func main() {
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
