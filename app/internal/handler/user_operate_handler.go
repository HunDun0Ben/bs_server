package handler

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/internal/service/authsvc"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/bscxt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo/imongoutil"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

type UserHandler struct {
	userService      usersvc.UserService
	authService      authsvc.AuthService
	butterflyService butterflysvc.ButterflyService
}

func NewUserHandler(
	userService usersvc.UserService,
	authService authsvc.AuthService,
	butterflyService butterflysvc.ButterflyService,
) *UserHandler {
	return &UserHandler{
		userService:      userService,
		authService:      authService,
		butterflyService: butterflyService,
	}
}

// UploadImg godoc
// @Summary      上传图片
// @Description  上传图片文件到服务器
// @Tags         UserController
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "要上传的图片文件"
// @Success      200  {object}  dto.SwaggerResponse{data=map[string]interface{}} "成功响应，返回文件ID和文件名"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /user/uploadImg [post]
// @Security     BearerAuth
func (h *UserHandler) UploadImg(cxt *gin.Context) {
	file, header, err := cxt.Request.FormFile("file")
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的文件", nil, err))
		return
	}
	defer file.Close()
	slog.Info("Uploading file", "filename", header.Filename)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "读取文件失败", nil, err))
		return
	}
	// TODO: Move file storage logic to a dedicated service if needed.
	fileID, err := imongoutil.StoreFile(
		cxt,
		"updateImg",
		imongo.FileStoreData{FileName: header.Filename, Content: fileContent},
	)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "存储文件失败", nil, err))
		return
	}
	helper.Success(cxt, gin.H{
		"fileId":   fileID,
		"fileName": header.Filename,
	})
}

// GetImgResult godoc
// @Summary      获取图片处理结果
// @Description  根据图片ID获取图片处理的结果
// @Tags         UserController
// @Accept       json
// @Produce      json
// @Param        PreProWay  query  []int  false  "预处理方式数组"
// @Param        Feature    query  int    false  "特征类型"
// @Param        Classifier query  int    false  "分类器类型"
// @Param        ImgID      query  string true   "图片ID"
// @Success      200  {object}  dto.SwaggerResponse{data=map[string]interface{}} "成功响应，返回图片处理结果"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      404  {object}  dto.SwaggerResponse "图片未找到"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /user/getImgResult [get]
// @Security     BearerAuth
func (h *UserHandler) GetImgResult(cxt *gin.Context) {
	var req dto.GetImgResultReq
	if err := cxt.ShouldBindQuery(&req); err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的请求参数", nil, err))
		return
	}

	_, err := imongoutil.GetFile(cxt, "updateImg", req.ImgID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			cxt.Error(bsvo.NewAppError(http.StatusNotFound, "图片未找到", nil, err))
		} else {
			cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "获取图片结果失败", nil, err))
		}
		return
	}
	// TODO: 此处应添加实际的图片处理逻辑并返回结果。
	helper.Success(cxt, gin.H{
		"message": "图片结果尚未就绪。",
		"imgId":   req.ImgID,
	})
}

// InsectInfo godoc
// @Summary      获取昆虫信息
// @Description  获取昆虫的详细信息（待实现）
// @Tags         UserController
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse "成功响应"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /user/insect [get]
// @Security     BearerAuth
func (h *UserHandler) InsectInfo(cxt *gin.Context) {
}

// ButterflyInfo godoc
// @Summary      获取蝴蝶种类信息
// @Description  获取所有蝴蝶种类的详细信息列表
// @Tags         UserController
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.SwaggerResponse{data=[]insect.Insect} "成功响应，返回蝴蝶种类信息列表"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /user/butterfly_type_info [get]
// @Security     BearerAuth
func (h *UserHandler) ButterflyInfo(cxt *gin.Context) {
	insect_list, err := h.butterflyService.GetTypes(cxt.Request.Context())
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching butterfly info"})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   insect_list,
	})
}

// SetupTotp godoc
// @Summary      设置TOTP双因素认证
// @Description  为用户生成TOTP密钥和恢复码，返回包含二维码的配置信息
// @Tags         UserController
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  dto.SwaggerResponse{data=dto.TOTPSetupRes} "成功响应，返回TOTP配置信息"
// @Failure      401  {object}  dto.SwaggerResponse "未授权的访问"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /user/mfa/setup/totp [get]
// @Security     BearerAuth
func (h *UserHandler) SetupTotp(cxt *gin.Context) {
	// 从上下文获取用户名
	username := cxt.GetString(bscxt.ContextUsernameKey)
	if username == "" {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "未授权的访问", nil, nil))
		return
	}

	// 生成TOTP密钥
	key, err := h.authService.GenerateTOTPSecret(username)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成TOTP密钥失败", nil, err))
		return
	}

	// 生成恢复码
	recoveryCodes, err := h.authService.GenerateRecoveryCodes()
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成恢复码失败", nil, err))
		return
	}

	// 将TOTP密钥保存到数据库 (此时 mfaEnabled 为 false)
	err = h.authService.SaveMFASecret(cxt.Request.Context(), username, key.Secret())
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "保存MFA配置失败", nil, err))
		return
	}

	// 返回设置信息
	// Generate QR code image
	qrCode, err := key.Image(200, 200)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "生成QR码图片失败", nil, err))
		return
	}

	// Convert QR code image to base64
	var buf bytes.Buffer
	err = png.Encode(&buf, qrCode)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "编码QR码图片失败", nil, err))
		return
	}
	qrBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	helper.Success(cxt, dto.TOTPSetupRes{
		Secret:        key.Secret(),
		QRCode:        qrBase64,
		RecoveryCodes: recoveryCodes,
	})
}

// VerifyTotp godoc
// @Summary      验证TOTP
// @Description  验证用户提供的TOTP码，验证成功后激活MFA
// @Tags         UserController
// @Accept       json
// @Produce      json
// @Param        verify body dto.TOTPVerifyReq true "TOTP验证码"
// @Success      200  {object}  dto.SwaggerResponse{data=dto.TOTPVerifyRes} "验证成功响应"
// @Failure      400  {object}  dto.SwaggerResponse "请求参数错误"
// @Failure      401  {object}  dto.SwaggerResponse "未授权的访问"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /user/mfa/verify/totp [post]
// @Security     BearerAuth
func (h *UserHandler) VerifyTotp(cxt *gin.Context) {
	// 从上下文获取用户名
	username := cxt.GetString(bscxt.ContextUsernameKey)
	if username == "" {
		cxt.Error(bsvo.NewAppError(http.StatusUnauthorized, "未授权的访问", nil, nil))
		return
	}

	// 绑定并验证请求参数
	var req dto.TOTPVerifyReq
	if err := cxt.ShouldBindJSON(&req); err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的请求参数", nil, err))
		return
	}

	// 获取用户的TOTP secret
	secret, err := h.authService.GetUserMFASecret(cxt.Request.Context(), username)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "获取MFA配置失败", nil, err))
		return
	}

	// 验证TOTP码并激活MFA
	err = h.authService.VerifyAndActivateMFA(cxt.Request.Context(), username, secret, req.Code)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "TOTP验证失败", nil, err))
		return
	}

	helper.Success(cxt, dto.TOTPVerifyRes{
		Activated: true,
		Message:   "MFA验证成功并已激活",
	})
}
