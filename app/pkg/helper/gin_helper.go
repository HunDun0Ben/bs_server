package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/pkg/bserr"
)

func Success[T any](cxt *gin.Context, data T) {
	cxt.JSON(http.StatusOK, dto.NewBaseRes(200, nil, data))
}

func Failed(cxt *gin.Context, err *bserr.AppError) {
	cxt.Error(err)
}
