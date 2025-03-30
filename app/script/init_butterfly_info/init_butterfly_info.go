package main

import (
	"context"
	"log/slog"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/HunDun0Ben/bs_server/app/entities/insect"
	mcli "github.com/HunDun0Ben/bs_server/common/data/imongo"
)

func main() {
	col := mcli.Database("bs_server_db").Collection("butterfly_info")
	// list := loadInfoFromCSV()
	// col.InsertMany(context.Background(), list)
	var insect insect.Insect
	col.FindOne(context.Background(), bson.D{}).Decode(&insect)
	slog.Info("", "insect", insect)
}

func loadInfoFromCSV() []any {
	filepath := "./蝴蝶信息.xlsx"
	// headStr := [...]string{"中文名称", "英文名称", "拉丁学名",
	// 	"特征描述文本", "分布情况文本", "保护级别", "别名"}
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("", "err", err)
		}
	}()

	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		panic(err)
	}
	list := make([]interface{}, 0)
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
	slog.Info("list", "list", list)
	return list
}
