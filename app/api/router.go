package api

import (
	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/handler"
	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

func InitRoute(engine *gin.Engine) {
	engine.Use(middleware.WebErrorHandler())
	{
		// 公开路由
		public := engine.Group("/")
		public.POST("/login", handler.Login)
		// 测试一些内容
		otherTest := engine.Group("/test")
		otherTest.GET("/getAllProType", handler.GetProType)
	}

	// jwt 认证路由组
	auth := engine.Group("/")
	if conf.AppConfig.JWT.Enable {
		auth.Use(middleware.JWTAuth)
	}
	// 以下的路由都来自 auth 组, 故需要通过 JWT 认证
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
	}
}
