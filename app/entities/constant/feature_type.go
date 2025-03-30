package constant

var ProTypeMap map[string]int

const (
	Garing int = 1 << (iota)
	GaussianBlur
	EqualizeHist
	GraphCut
)

func init() {
	ProTypeMap = make(map[string]int)
	ProTypeMap["灰度化"] = Garing
	ProTypeMap["高斯模糊"] = GaussianBlur
	ProTypeMap["直方图均衡化"] = EqualizeHist
	ProTypeMap["图像分割"] = GraphCut
}
