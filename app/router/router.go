package router

import (
	"demo/app/api"

	"github.com/gin-gonic/gin"
)

func InitRoute(engine *gin.Engine) {
	router := engine.Group("/test")
	{
		router.GET("/getAllProType", api.GetProType)
		router.GET("/hello", api.HelloWorld)
		router.GET("/test", api.Test)
	}

	manage := engine.Group("/manage")
	{
		manage.GET("/initImgDB")
		manage.GET("/initInsect", api.InitInsect)
		manage.GET("/initClassification", api.InitClassification)
	}

	user := engine.Group("/user")
	{
		user.POST("/uploadImg", api.UploadImg)
		user.GET("/getImgResult", api.GetImgResult)
		user.GET("/insect", api.InsectInfo)
	}
}
