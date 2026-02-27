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
	userService            usersvc.UserService
	authService            authsvc.AuthService
	mfaVerificationService *authsvc.MFAVerificationService
}

func NewLoginHandler(userService usersvc.UserService, authService authsvc.AuthService, mfaVerificationService *authsvc.MFAVerificationService) *LoginHandler {
	return &LoginHandler{
		userService:            userService,
		authService:            authService,
		mfaVerificationService: mfaVerificationService,
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

	clientIP := cxt.ClientIP()
	mfaRequired, requiredTypes := h.userService.IsHighRisk(cxt, user, clientIP)

	accessTokenStr, refreshTokenStr, claims, err := bsjwt.GenerateTokenPair(*user, mfaRequired, requiredTypes)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成token失败", nil, err))
		return
	}

	// 非完整权限操作. 保证以后当非完整权限功能范围扩大时候, 知道副作用在哪里进行执行.
	// 更新登录信息可能只是其中之一.
	if !mfaRequired {
		err = h.userService.UpdateLoginInfo(cxt, user.ID, clientIP)
		if err != nil {
			slog.Error("更新用户登录信息失败", "error", err)
			cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登录失败", nil, err))
			return
		}
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

	response := dto.LoginResponse{
		AccessToken: accessTokenStr,
	}
	if mfaRequired {
		response.MFARequired = true
		response.RequiredTypes = requiredTypes
	}
	helper.Success(cxt, response)
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

	refreshClaims, err := bsjwt.ParseToken(refreshTokenStr)
	if err != nil || refreshClaims.Subject != bsjwt.RefreshToken {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "Refresh Token 无效", nil, err))
		return
	}
	// 查找对应 jti 的 refresh token 是否存在
	storedUsername, err := h.authService.IsRefreshTokenValid(cxt, refreshClaims.ID)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "Refresh Token 已失效或不存在", nil, err))
		return
	}

	authHeader := cxt.GetHeader(bsjwt.AuthHeaderName)
	accessClaims, err := bsjwt.ParseTokenByHeader(authHeader)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "", nil, err))
		cxt.Abort()
		return
	}

	if storedUsername != accessClaims.Username {
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
	// Refresh Token 过程保持原有的 MFA 状态
	accessTokenStr, refreshTokenStr, newClaims, err := bsjwt.GenerateTokenPair(*user, refreshClaims.MFAPending, refreshClaims.RequiredTypes)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成token失败", nil, err))
		return
	}
	err = h.authService.StoreRefreshToken(cxt, newClaims.ID, user.Username, time.Until(newClaims.ExpiresAt.Time))
	if err != nil {
		slog.Error("存储Token失败")
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}

	cxt.SetCookie("refreshToken", refreshTokenStr, int(time.Until(newClaims.ExpiresAt.Time).Seconds()), "/api/token/refresh", "", true, true)
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

// VerifyMFA godoc
// @Summary      二次验证 MFA
// @Description  用户提交 MFA 验证码进行二次验证
// @Tags         LoginController
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        mfa body dto.MFAVerifyRequest true "MFA 验证请求"
// @Success      200  {object}  dto.SwaggerResponse{data=dto.LoginResponse} "成功响应"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      401  {object}  dto.SwaggerResponse "MFA 验证失败"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /login/mfa-verify [post]
func (h *LoginHandler) VerifyMFA(cxt *gin.Context) {
	var req dto.MFAVerifyRequest
	if err := cxt.Bind(&req); err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的请求参数", nil, err))
		return
	}

	claims, exists := cxt.Get(bscxt.ContextClaimsKey)
	if !exists {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "未提供认证信息", nil, nil))
		return
	}
	customClaims, ok := claims.(*bsjwt.CustomClaims)
	if !ok {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "认证信息类型错误", nil, nil))
		return
	}

	if !customClaims.MFAPending {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "当前不需要 MFA 验证", nil, nil))
		return
	}

	user, err := h.userService.FindByUsername(cxt, customClaims.Username)
	if err != nil || user == nil {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "用户不存在", nil, err))
		return
	}

	secret, enabled, err := h.userService.GetMFAInfo(cxt, user.Username)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "获取 MFA 信息失败", nil, err))
		return
	}

	if !enabled {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "用户未开启 MFA", nil, nil))
		return
	}

	// 默认使用第一个要求的验证方式
	providerType := "totp"
	if len(customClaims.RequiredTypes) > 0 {
		providerType = customClaims.RequiredTypes[0]
	}

	valid, err := h.mfaVerificationService.VerifyCode(cxt, providerType, secret, req.Code)
	if err != nil || !valid {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "MFA 验证失败", nil, err))
		return
	}

	// 验证成功，签发新的全权限 Token
	// 非完整权限操作. 保证以后当非完整权限功能范围扩大时候, 知道副作用在哪里进行执行.
	// 更新登录信息可能只是其中之一.
	err = h.userService.UpdateLoginInfo(cxt, user.ID, cxt.ClientIP())
	if err != nil {
		slog.Error("更新用户登录信息失败", "error", err)
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}

	accessTokenStr, refreshTokenStr, newClaims, err := bsjwt.GenerateTokenPair(*user, false, nil)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成token失败", nil, err))
		return
	}

	// 作废旧的 Access Token
	remainingTime := time.Until(customClaims.ExpiresAt.Time)
	if remainingTime > 0 {
		_ = h.authService.InvalidateAccessToken(cxt, customClaims.ID, remainingTime)
	}

	// 从 cookie 获取并作废旧的 Refresh Token
	if oldRefreshTokenStr, err := cxt.Cookie("refreshToken"); err == nil {
		if oldRefreshClaims, err := bsjwt.ParseToken(oldRefreshTokenStr); err == nil {
			_ = h.authService.InvalidateRefreshToken(cxt, oldRefreshClaims.ID)
		}
	}

	// 存储新的 refresh token 到 redis
	err = h.authService.StoreRefreshToken(cxt, newClaims.ID, user.Username, time.Until(newClaims.ExpiresAt.Time))
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "登录失败", nil, err))
		return
	}

	maxAge := int(time.Until(newClaims.ExpiresAt.Time).Seconds())
	cxt.SetCookie("refreshToken", refreshTokenStr, maxAge, "/api/token/refresh", "", true, true)

	helper.Success(cxt, dto.LoginResponse{
		AccessToken: accessTokenStr,
	})
}
