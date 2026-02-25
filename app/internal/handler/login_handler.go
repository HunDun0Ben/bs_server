package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/internal/service/authsvc"
	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/bscxt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

type LoginHandler struct {
	userService usersvc.UserService
	authService authsvc.AuthService
}

func NewLoginHandler(userService usersvc.UserService, authService authsvc.AuthService) *LoginHandler {
	return &LoginHandler{
		userService: userService,
		authService: authService,
	}
}

// Login godoc
// @Summary      用户登录
// @Description  用户使用用户名和密码进行登录，成功后返回 JWT Token。
// @Tags         LoginController
// @Accept       json
// @Produce      json
// @Param        login body dto.LoginRequest true "登录凭据"
// @Success      200  {object}  dto.SwaggerResponse{data=dto.LoginResponse} "成功响应，Access Token 在响应体中，Refresh Token 在 HttpOnly Cookie 中"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      401  {object}  dto.SwaggerResponse "用户名或密码错误"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /login [post]
func (h *LoginHandler) Login(cxt *gin.Context) {
	var req dto.LoginRequest

	if err := cxt.Bind(&req); err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的请求参数", nil, err))
		return
	}
	// 查看用户信息
	user, err := h.userService.FindByLogin(cxt, req.Username, req.Password)
	if err != nil || user == nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "用户名或密码错误", nil, err))
		return
	}
	accessTokenStr, refreshTokenStr, claims, err := bsjwt.GenerateTokenPair(*user)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成token失败", nil, err))
		return
	}
	// 存储 refresh token 到 redis 中.
	err = h.authService.StoreRefreshToken(cxt, claims.ID, user.Username, time.Until(claims.ExpiresAt.Time))
	if err != nil {
		slog.Error("存储Token失败")
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}
	maxAge := int(time.Until(claims.ExpiresAt.Time).Seconds())
	cxt.SetCookie("refreshToken", refreshTokenStr, maxAge, "/api/token/refresh", "", true, true)
	helper.Success(cxt, dto.LoginResponse{
		AccessToken: accessTokenStr,
	})
}

// RefreshToken godoc
// @Summary      刷新 Access Token
// @Description  使用 Refresh Token 获取一个新的 Access Token
// @Tags         LoginController
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse{data=dto.RefreshTokenResponse} "成功响应，返回新的 Access Token，Refresh Token 在 HttpOnly Cookie 中"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      401  {object}  dto.SwaggerResponse "Refresh Token 无效或已过期"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /token/refresh [get]
func (h *LoginHandler) RefreshToken(cxt *gin.Context) {
	refreshTokenStr, err := cxt.Cookie("refreshToken")
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "Refresh Token 不存在", nil, err))
		return
	}

	claims, err := bsjwt.ParseToken(refreshTokenStr)
	if err != nil || claims.Subject != bsjwt.RefreshToken {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "Refresh Token 无效", nil, err))
		return
	}
	// 查找对应 jti 的 refresh token 是否存在
	storedUsername, err := h.authService.IsRefreshTokenValid(cxt, claims.ID)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "Refresh Token 已失效或不存在", nil, err))
		return
	}

	authHeader := cxt.GetHeader(bsjwt.AuthHeaderName)
	claims, err = bsjwt.ParseTokenByHeader(authHeader)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "", nil, err))
		cxt.Abort()
		return
	}

	if storedUsername != claims.Username {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "Refresh Token 与用户不匹配", nil, err))
		return
	}

	// 查找用户信息
	user, err := h.userService.FindByUsername(cxt, storedUsername)
	if err != nil || user == nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "用户名或密码错误", nil, err))
		return
	}

	// 生成对应 access and refresh token
	accessTokenStr, refreshTokenStr, claims, err := bsjwt.GenerateTokenPair(*user)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成token失败", nil, err))
		return
	}
	err = h.authService.StoreRefreshToken(cxt, claims.ID, user.Username, time.Until(claims.ExpiresAt.Time))
	if err != nil {
		slog.Error("存储Token失败")
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}

	cxt.SetCookie("refreshToken", refreshTokenStr, int(time.Until(claims.ExpiresAt.Time).Seconds()), "/api/token/refresh", "", true, true)
	helper.Success(cxt, dto.LoginResponse{
		AccessToken: accessTokenStr,
	})
}

// Logout godoc
// @Summary      用户登出
// @Description  用户登出，将当前 Access Token 加入黑名单使其失效
// @Tags         LoginController
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse "成功响应"
// @Failure      401  {object}  dto.SwaggerResponse "认证失败"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /logout [post]
func (h *LoginHandler) Logout(cxt *gin.Context) {
	jti := cxt.GetString(bscxt.ContextJTIKey)
	ExpiresAt := cxt.GetTime(bscxt.ExpiresAtKey)

	remainingTime := time.Until(ExpiresAt)
	if remainingTime <= 0 {
		helper.Success(cxt, gin.H{"message": "登出成功，Token已自然过期"})
		return
	}

	_ = h.authService.InvalidateRefreshToken(cxt, jti)
	err := h.authService.InvalidateAccessToken(cxt, jti, remainingTime)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登出操作失败", nil, err))
		return
	}

	helper.Success(cxt, gin.H{"message": "登出成功"})
}
