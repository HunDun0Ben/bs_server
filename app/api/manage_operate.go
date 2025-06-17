package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func InitImgDB(cxt *gin.Context) {
}

func InitInsect(cxt *gin.Context) {
	filepath := "/home/workspace/data/leedsbutterfly/butterfly_type_info.xlsx"
	headstr := [...]string{
		"分类器id", "中文名称", "英文名称", "拉丁学名",
		"特征描述文本", "分布情况文本", "保护级别文本",
	}
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		panic(err)
	}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		a := strconv.Itoa(i)
		for _, colCell := range row {
			a += colCell + "\t"
		}
		log.Println(a)
	}
	log.Println(headstr)
}

func InitClassification(cxt *gin.Context) {
}
