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
// @Schemes   http
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token. Example: "Bearer {token}"
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
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/HunDun0Ben/bs_server/app/api"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	"github.com/HunDun0Ben/bs_server/app/pkg/logger"

	mcli "github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	if err := conf.InitConfig(); err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}
	// 初始化 OTel (Tracing + Logging)
	otel_clean, err := initOTel()
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}
	defer otel_clean()
	// 连接数据库
	loadDataBase()
	// 启动服务器
	startServer()
}

func initOTel() (func(), error) {
	if !conf.AppConfig.OTEL.Enable {
		// 如果未开启 OTel，依然使用基础的 OTelHandler (仅注入 ID 到 Stdout)
		slog.SetDefault(slog.New(logger.NewOTelHandler(slog.NewJSONHandler(os.Stdout, nil))))
		return func() {}, nil
	}

	ctx := context.Background()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(conf.AppConfig.OTEL.ServiceName),
		semconv.ServiceVersion(conf.AppConfig.OTEL.Version),
	)

	// --- 1. Tracing 配置 ---
	var traceOpts []otlptracegrpc.Option
	traceOpts = append(traceOpts, otlptracegrpc.WithEndpoint(conf.AppConfig.OTEL.Endpoint))
	if conf.AppConfig.OTEL.Insecure {
		traceOpts = append(traceOpts, otlptracegrpc.WithInsecure())
	}
	var tp *trace.TracerProvider
	traceExporter, err := otlptracegrpc.New(ctx, traceOpts...)
	if err != nil {
		// 严格要求开启 Otel
		if conf.AppConfig.OTEL.Strict {
			return nil, fmt.Errorf("failed to create trace exporter in strict mode: %w", err)
		}
		slog.Error("Failed to create trace exporter, continuing in non-strict mode", "error", err)
	} else {
		tp = trace.NewTracerProvider(
			trace.WithBatcher(traceExporter),
			trace.WithResource(res),
		)
		otel.SetTracerProvider(tp)
	}

	// --- 2. Logging 配置 (Direct push to collector) ---
	var logOpts []otlploggrpc.Option
	logOpts = append(logOpts, otlploggrpc.WithEndpoint(conf.AppConfig.OTEL.Endpoint))
	if conf.AppConfig.OTEL.Insecure {
		logOpts = append(logOpts, otlploggrpc.WithInsecure())
	}
	var lp *sdklog.LoggerProvider
	logExporter, err := otlploggrpc.New(ctx, logOpts...)
	if err != nil {
		if conf.AppConfig.OTEL.Strict {
			return nil, fmt.Errorf("failed to create log exporter in strict mode: %w", err)
		}
		slog.Error("Failed to create log exporter, continuing in non-strict mode", "error", err)
	} else {
		lp = sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
			sdklog.WithResource(res),
		)
		global.SetLoggerProvider(lp)
	}

	// --- 3. 整合 slog ---
	// 同时满足：1. 推送到 Collector; 2. 打印带 TraceID 的 JSON 到 Stdout
	stdoutHandler := slog.NewJSONHandler(os.Stdout, nil)
	var combinedHandler slog.Handler = stdoutHandler

	if lp != nil {
		otelslogHandler := otelslog.NewHandler(conf.AppConfig.OTEL.ServiceName)
		combinedHandler = logger.NewMultiHandler(
			stdoutHandler,
			otelslogHandler,
		)
	}

	slog.SetDefault(slog.New(logger.NewOTelHandler(combinedHandler)))

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if tp != nil {
			if err := tp.Shutdown(shutdownCtx); err != nil {
				slog.Error("Error shutting down tracer provider", "error", err)
			}
		}
		if lp != nil {
			if err := lp.Shutdown(shutdownCtx); err != nil {
				slog.Error("Error shutting down logger provider", "error", err)
			}
		}
	}, nil
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
