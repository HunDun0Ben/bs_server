package constant

// 分类器类型,
// iota 是一个预声明的标识符,
// 每遇到一个 const 关键字声明块，iota 的计数值都会重置为 0.
const (
	SVM   int = 1 << (iota) // 一开始 itoa 为 0, 即 SVM = 1 << 0, 即 SVM = 0
	DTREE                   // DTREE 1 << 2 = 2, next variable 1 << 3 = 4
)
