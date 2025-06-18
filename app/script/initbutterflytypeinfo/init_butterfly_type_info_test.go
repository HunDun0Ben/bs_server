package initbutterflyinfo_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/xuri/excelize/v2"

	"github.com/HunDun0Ben/bs_server/app/entities/insect"
	"github.com/HunDun0Ben/bs_server/app/service/butterflytypesvc"
)

func main() {
	list, _ := loadTypeInfoFromCSV()
	slog.Info("list", slog.Any("list", list))
}
func init_butterfly_type_info() {

}

func Test_init_butterfly_type_info(t *testing.T) {
	list, _ := loadTypeInfoFromCSV()
	err := butterflytypesvc.NewButterflyService().InitAll(context.Background(), list)
	if err != nil {
		slog.Error("初始化蝴蝶信息失败", "err", err)
		return
	}
}

func loadTypeInfoFromCSV() ([]insect.Insect, error) {
	filepath := "./蝴蝶信息.xlsx"
	// headStr := [...]string{"中文名称", "英文名称", "拉丁学名","特征描述文本", "分布情况文本", "保护级别", "别名"}.
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("无法打开文件: %s, 错误: %v", filepath, err))
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("无法正确关闭资源文件.", "err", err)
		}
	}()
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		slog.Error("无法获取工作表的行数据.", "err", err)
		return nil, err
	}
	list := make([]insect.Insect, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		insect := insect.Insect{
			ChineseName:        row[0],
			LatinName:          row[1],
			EnglishName:        row[2],
			FeatureDescription: row[3],
		}
		list = append(list, insect)
	}
	return list, nil
}
