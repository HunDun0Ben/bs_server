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

// Login godoc
// @Summary      用户登录
// @Description  用户使用用户名和密码进行登录，成功后返回 JWT Token。
// @Tags         公开路由
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "登录凭据"
// @Success      200  {object}  dto.SwaggerResponse{data=dto.LoginResponse} "成功响应，返回 JWT Token"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      401  {object}  dto.SwaggerResponse "用户名或密码错误"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /login [post]
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
	helper.Success(cxt, dto.LoginResponse{Token: token})
}
