package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
)

func Success[T any](cxt *gin.Context, data T) {
	cxt.JSON(http.StatusOK, dto.NewBaseRes(200, nil, data))
}

func Failed(cxt *gin.Context, err *bsvo.AppError) {
	// 设置错误到 cxt 中, 交由 WebErrorHandler 处理
	cxt.Error(err)
}
