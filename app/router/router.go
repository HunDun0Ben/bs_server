package router

import (
	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/api"
	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/HunDun0Ben/bs_server/common/conf"
)

func InitRoute(engine *gin.Engine) {
	{
		// 公开路由
		public := engine.Group("/")
		public.POST("/login", api.Login)
		// 测试一些内容
		otherTest := engine.Group("/test")
		otherTest.GET("/getAllProType", api.GetProType)
	}

	// jwt 认证路由组
	auth := engine.Group("/")
	if conf.AppConfig.JWT.Enable {
		auth.Use(middleware.JWTAuth())
	}
	// 以下的路由都来自 auth 组, 故需要通过 JWT 认证
	{ // 管理路由
		manage := auth.Group("/manage")
		manage.GET("/initInsect", api.InitInsect)
		manage.GET("/initClassification", api.InitClassification)
	}
	{ // 用户路由
		user := auth.Group("/user")
		user.POST("/uploadImg", api.UploadImg)
		user.GET("/getImgResult", api.GetImgResult)
		user.GET("/insect", api.InsectInfo)
		user.GET("/butterfly_type_info", api.ButterflyInfo)
	}
}
