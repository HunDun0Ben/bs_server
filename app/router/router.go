package router

import (
	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/api"
	"github.com/HunDun0Ben/bs_server/app/middleware"
)

func InitRoute(engine *gin.Engine) {
	// Global middleware - applies to all routes.

	// 公开路由
	public := engine.Group("/")
	public.POST("/login", api.Login)

	// 需要认证的路由
	web := engine.Group("/")
	web.Use(middleware.JWTAuth(),
		middleware.WebErrorHandler())

	// 测试路由
	test := web.Group("/test")
	// test.GET("/getAllProType", api.GetProType)
	test.GET("/hello", api.HelloWorld)
	test.GET("/test", api.Test)

	abb := engine.Group("/test")
	abb.GET("/getAllProType", api.GetProType)

	// 管理路由
	manage := web.Group("/manage")
	manage.GET("/initImgDB")
	manage.GET("/initInsect", api.InitInsect)
	manage.GET("/initClassification", api.InitClassification)

	// 用户路由
	user := web.Group("/user")
	user.POST("/uploadImg", api.UploadImg)
	user.GET("/getImgResult", api.GetImgResult)
	user.GET("/insect", api.InsectInfo)
}
