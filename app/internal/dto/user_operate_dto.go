package dto

type GetImgResultReq struct {
	PreProWay  []int
	Feature    int
	Classifier int
	ImgID      string
}

type TOTPSetupRes struct {
	Secret        string   `json:"secret"`         // TOTP 密钥
	QRCode        string   `json:"qr_code"`        // QR码(base64编码)
	RecoveryCodes []string `json:"recovery_codes"` // 恢复码
}

type TOTPVerifyReq struct {
	Code string `json:"code" binding:"required"` // TOTP验证码
}

type TOTPVerifyRes struct {
	Activated bool   `json:"activated"` // MFA是否已激活
	Message   string `json:"message"`   // 响应消息
}
