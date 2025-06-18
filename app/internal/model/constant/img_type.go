package constant

// 图片类型.
const (
	NormalSampleImg int = 1 << (iota) // 未特殊处理过的样本图片
	ProSampleImg                      // 经过处理过的样本图片
	UserImg                           // 用户上传图片
)
