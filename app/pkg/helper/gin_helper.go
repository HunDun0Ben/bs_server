package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
)

func Success[T any](cxt *gin.Context, data T) {
	cxt.JSON(http.StatusOK, dto.NewBaseRes(200, nil, data))
}
