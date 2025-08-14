// @title           蝴蝶识别系统 API
// @version         1.0
// @description     这是一个使用 Go, Gin, Gocv 和 MongoDB 构建的蝴蝶识别系统的服务端 API 文档。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	"github.com/HunDun0Ben/bs_server/app/api"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	mcli "github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
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
	api.InitRoute(router)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.AppConfig.Server.Port),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Failed to start server: " + err.Error())
		}
		slog.Info("Server started.", "address", server.Addr)
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
