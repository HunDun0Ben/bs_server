package api

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/HunDun0Ben/bs_server/app/internal/handler"
	"github.com/HunDun0Ben/bs_server/app/internal/repository"
	"github.com/HunDun0Ben/bs_server/app/internal/service/authsvc"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/iredis"

	_ "github.com/HunDun0Ben/bs_server/app/docs/swagger" // swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute(engine *gin.Engine) {
	engine.Use(otelgin.Middleware("bs_server"))

	// --- 1. Repository Layer ---
	userRepo := repository.NewUserRepository(&imongo.MongoCollection{Col: imongo.BizDataBase().Collection("user")})
	authRepo := repository.NewAuthRepository(iredis.GetRDB())
	butterflyRepo := repository.NewButterflyRepository(imongo.BizDataBase())

	// --- 2. Service Layer ---
	userService := usersvc.NewUserService(userRepo)
	authService := authsvc.NewAuthService(authRepo, userService)
	butterflyService := butterflysvc.NewButterflyService(butterflyRepo)

	// --- 3. Handler Layer ---
	loginHandler := handler.NewLoginHandler(userService, authService)
	userHandler := handler.NewUserHandler(userService, authService, butterflyService)
	lotteryHandler := handler.NewLotteryHandler()
	manageHandler := handler.NewManageHandler(butterflyService)
	commonHandler := handler.NewCommonHandler()

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
		public.POST("/login", loginHandler.Login)
		public.GET("/token/refresh", loginHandler.RefreshToken)
	}
	// 测试一些内容
	{
		otherTest := apiV1.Group("/test")
		otherTest.GET("/getAllProType", commonHandler.GetProType)
	}
	// jwt 认证路由组, 需要通过 JWT 认证
	auth := apiV1.Group("/")
	if conf.AppConfig.JWT.Enable {
		auth.Use(middleware.JWTAuth(authService))
	}
	auth.POST("/logout", loginHandler.Logout)

	// 彩票相关路由
	{
		lottery := auth.Group("/lottery")
		lottery.GET("/bigLottery/random", lotteryHandler.BigLotteryRandom)
	}
	{ // 管理路由
		manage := auth.Group("/manage")
		manage.GET("/initInsect", manageHandler.InitInsect)
		manage.GET("/initClassification", manageHandler.InitClassification)
	}
	{ // 用户路由
		user := auth.Group("/user")
		user.POST("/uploadImg", userHandler.UploadImg)
		user.GET("/getImgResult", userHandler.GetImgResult)
		user.GET("/insect", userHandler.InsectInfo)
		user.GET("/butterfly_type_info", userHandler.ButterflyInfo)
		user.GET("/mfa/setup/totp", userHandler.SetupTotp)
		user.POST("/mfa/verify/totp", userHandler.VerifyTotp)
	}
}
