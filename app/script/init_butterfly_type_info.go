package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/xuri/excelize/v2"

	"github.com/HunDun0Ben/bs_server/app/internal/model/insect"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
)

func main() {
	svc := butterflysvc.NewButterflyTypeSvc()
	count, err := svc.Count(context.Background())
	if err != nil {
		slog.Error("获取数据数量失败", "err", err)
		return
	}
	if count > 0 {
		slog.Info("已经初始化过数据")
		return
	}
	list, _ := loadTypeInfoFromCSV()
	if err := svc.InitAll(context.Background(), list); err != nil {
		slog.Error("初始化蝴蝶信息失败", "err", err)
		return
	}
}

func loadTypeInfoFromCSV() ([]insect.Insect, error) {
	filepath := "./initbutterflytypeinfo/蝴蝶信息.xlsx"
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
