package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	r "github.com/HunDun0Ben/bs_server/app/router"
	"github.com/HunDun0Ben/bs_server/common/conf"

	mcli "github.com/HunDun0Ben/bs_server/common/data/imongo"
)

func main() {
	// 连接数据库
	loadDataBase()
	// 启动服务器
	startServer()
}

func loadDataBase() {
	mcli.Client()
}

func startServer() {
	router := gin.Default()
	r.InitRoute(router)

	serverAddr := fmt.Sprintf(":%d", conf.GlobalViper.GetInt("server.port"))
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Failed to start server: " + err.Error())
		} else {
			slog.Info("Server started.", "address", server.Addr)
		}
	}()

	// 优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞等待接收系统信号
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
	slog.Info("Server exited")
}
