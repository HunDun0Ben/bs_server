package handler

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func InitImgDB(cxt *gin.Context) {
}

// InitInsect godoc
// @Summary      初始化昆虫信息
// @Description  从服务器的预定路径读取 Excel 文件，并将蝴蝶物种信息初始化到数据库中。这是一个管理接口。
// @Tags         ManagerController
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse{data=string} "成功响应，返回操作成功的消息"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /manage/initInsect [get]
// @Security     BearerAuth
func InitInsect(cxt *gin.Context) {
	filepath := "/home/workspace/data/leedsbutterfly/butterfly_type_info.xlsx"
	headstr := [...]string{
		"分类器id", "中文名称", "英文名称", "拉丁学名",
		"特征描述文本", "分布情况文本", "保护级别文本",
	}
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		slog.Error("open file error", "err", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("close file error", "err", err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		slog.Error("get rows error", "err", err)
		return
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		a := strconv.Itoa(i)
		for _, colCell := range row {
			a += colCell + "\t"
		}
		slog.Info(a)
	}
	slog.Info("excel head", "head", headstr)
}

// InitClassification godoc
// @Summary      初始化分类器
// @Description  执行分类器的初始化或训练任务。这是一个管理接口。
// @Tags         ManagerController
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse{data=string} "成功响应，返回操作成功的消息"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /manage/initClassification [get]
// @Security     BearerAuth
func InitClassification(cxt *gin.Context) {
}
