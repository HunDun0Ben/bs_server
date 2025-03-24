package main

import (
	"context"
	"demo/app/entities/insect"
	mcli "demo/common/data/imongo"
	"log"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	mcli.Init()
	col := mcli.Client.Database("bs_server_db").Collection("butterfly_info")
	// list := loadInfoFromCSV()
	// col.InsertMany(context.Background(), list)

	var insect insect.Insect
	col.FindOne(context.Background(), bson.D{}).Decode(&insect)
	log.Println(insect)
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
			log.Print(err)
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
	log.Println(list)
	return list
}
