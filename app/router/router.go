package router

import (
	"demo/app/api"

	"github.com/gin-gonic/gin"
)

func InitRoute(engine *gin.Engine) {
	test := engine.Group("/test")
	{
		test.GET("/getAllProType", api.GetProType)
		test.GET("/hello", api.HelloWorld)
		test.GET("/test")
	}

	manage := engine.Group("/manage")
	{
		manage.GET("/initImgDB")
		manage.GET("/initInsect")
		manage.GET("/initClassification")
	}

	user := engine.Group("/user")
	{
		user.POST("/uploadImg")
		user.GET("/getImgResult")
		user.GET("/insect")
	}
}
