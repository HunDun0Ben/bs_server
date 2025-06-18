package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/bserr"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

func Login(cxt *gin.Context) {
	var req dto.LoginRequest
	if err := cxt.ShouldBindJSON(&req); err != nil {
		cxt.Error(bserr.NewAppError(http.StatusBadRequest, "无效的请求参数", nil, err))
		return
	}

	user, err := usersvc.NewUserService().FindByLogin(cxt, req.Username, req.Password)
	if err != nil {
		cxt.Error(bserr.NewAppError(http.StatusUnauthorized, "用户名或密码错误", nil, err))
		return
	}
	if user == nil {
		cxt.Error(bserr.NewAppError(http.StatusUnauthorized, "用户不存在", nil, err))
		return
	}

	token, err := bsjwt.GenerateToken(*user)
	if err != nil {
		cxt.Error(bserr.NewAppError(http.StatusInternalServerError, "生成token失败", nil, err))
		return
	}
	helper.Success(cxt, gin.H{"token": token})
}
