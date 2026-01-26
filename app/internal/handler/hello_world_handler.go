package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/model/constant"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

// GetProType godoc
// @Summary      获取产品类型信息
// @Description  获取预处理方式和特征类型等产品配置信息
// @Tags         TestController
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse{data=map[string]map[string]int} "成功响应，返回产品类型映射信息"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /test/getAllProType [get]
func GetProType(cxt *gin.Context) {
	proTypeMap := make(map[string]map[string]int)
	proTypeMap["preProWayItems"] = constant.ProTypeMap
	proTypeMap["featureTypeItems"] = constant.FeatureTypeMap
	// proTypeMap["classifierTypeItems"] = nil
	helper.Success(cxt, proTypeMap)
}
