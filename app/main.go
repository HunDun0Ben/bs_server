package main

import (
	r "demo/app/router"
	"demo/common/conf"
	mcli "demo/common/data/imongo"
	"net/http"

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
	router := gin.Default()
	r.InitRoute(router)

	s := &http.Server{
		Addr:    ":" + conf.GlobalViper.GetString("server.port"),
		Handler: router,
	}
	s.ListenAndServe()
}
