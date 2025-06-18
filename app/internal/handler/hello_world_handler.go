package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/model/constant"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

func GetProType(cxt *gin.Context) {
	proTypeMap := make(map[string]map[string]int)
	proTypeMap["preProWayItems"] = constant.ProTypeMap
	proTypeMap["featureTypeItems"] = constant.FeatureTypeMap
	// proTypeMap["classifierTypeItems"] = nil
	helper.Success(cxt, proTypeMap)
}
