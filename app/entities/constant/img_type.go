package constant

// 图片类型
const (
	// 未特殊处理过的样本图片
	NORMAL_SAMPLE_IMG int = 1 << (iota)
	// 经过处理过的样本图片
	PRO_SAMPLE_IMG
	// 用户上传图片
	USER_IMG
)
