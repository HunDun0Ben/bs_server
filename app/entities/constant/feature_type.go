package constant

var ProTypeMap map[string]int

const (
	GARING int = 1 << (iota)
	GAUSSIAN_BLUR
	EQUALIZE_HIST
	GRAPH_CUT
)

func init() {
	ProTypeMap = make(map[string]int)
	ProTypeMap["灰度化"] = GARING
	ProTypeMap["高斯模糊"] = GAUSSIAN_BLUR
	ProTypeMap["直方图均衡化"] = EQUALIZE_HIST
	ProTypeMap["图像分割"] = GRAPH_CUT
}
