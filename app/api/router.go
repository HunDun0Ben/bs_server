package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	_ "github.com/HunDun0Ben/bs_server/app/docs/swagger" // swagger docs
	"github.com/HunDun0Ben/bs_server/app/internal/handler"
	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

func InitRoute(engine *gin.Engine) {
	engine.Use(otelgin.Middleware("bs_server"))

	// 自定义处理 404 请求
	engine.NoRoute(middleware.NoRouteHandler())

	// 配置 Swagger 路由
	// 访问 http://localhost:8080/swagger/index.html 即可查看文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建 API v1 路由组
	apiV1 := engine.Group("/api/v1")
	apiV1.Use(middleware.WebErrorHandler())

	// 公开路由
	{
		public := apiV1.Group("/")
		public.POST("/login", handler.Login)
		public.GET("/token/refresh", handler.RefreshToken)
	}
	// 测试一些内容
	{
		otherTest := apiV1.Group("/test")
		otherTest.GET("/getAllProType", handler.GetProType)
	}

	// jwt 认证路由组, 需要通过 JWT 认证
	auth := apiV1.Group("/")
	if conf.AppConfig.JWT.Enable {
		auth.Use(middleware.JWTAuth())
	}
	auth.POST("/logout", handler.Logout)

	{ // 管理路由
		manage := auth.Group("/manage")
		manage.GET("/initInsect", handler.InitInsect)
		manage.GET("/initClassification", handler.InitClassification)
	}
	{ // 用户路由
		user := auth.Group("/user")
		user.POST("/uploadImg", handler.UploadImg)
		user.GET("/getImgResult", handler.GetImgResult)
		user.GET("/insect", handler.InsectInfo)
		user.GET("/butterfly_type_info", handler.ButterflyInfo)
		user.GET("/mfa/setup/totp", handler.SetupTotp)
		user.GET("/mfa/verify/totp", handler.VerifyTotp)
	}
}
