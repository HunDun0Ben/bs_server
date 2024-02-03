package api

import (
	"demo/app/entities/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProType(cxt *gin.Context) {
	proTypeMap := make(map[string]map[string]int)
	proTypeMap["preProWayItems"] = constant.ProTypeMap
	proTypeMap["featureTypeItems"] = constant.FeatureTypeMap
	// proTypeMap["classifierTypeItems"] = nil
	cxt.JSON(http.StatusOK, proTypeMap)
}

func HelloWorld(cxt *gin.Context) {
	cxt.Writer.Write([]byte("Hello World."))
}

func Test(cxt *gin.Context) {
	cxt.Writer.Write([]byte("Hello World."))
}
