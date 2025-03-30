package router

import (
	"github.com/HunDun0Ben/bs_server/app/api"
	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute(engine *gin.Engine) {
	// 公开路由
	public := engine.Group("/")
	public.POST("/login", api.Login)

	// 需要认证的路由
	authorized := engine.Group("/")
	authorized.Use(middleware.JWTAuth())

	// 测试路由
	test := authorized.Group("/test")
	test.GET("/getAllProType", api.GetProType)
	test.GET("/hello", api.HelloWorld)
	test.GET("/test", api.Test)

	// 管理路由
	manage := authorized.Group("/manage")
	manage.GET("/initImgDB")
	manage.GET("/initInsect", api.InitInsect)
	manage.GET("/initClassification", api.InitClassification)

	// 用户路由
	user := authorized.Group("/user")
	user.POST("/uploadImg", api.UploadImg)
	user.GET("/getImgResult", api.GetImgResult)
	user.GET("/insect", api.InsectInfo)
}
