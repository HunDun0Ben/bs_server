package gexcel

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

func WriteData(records interface{}) error {
	xlsx := excelize.NewFile()
	sheetName := "sheet1"
	index, err := xlsx.NewSheet("sheet1")
	if err != nil {
		slog.Error("create new sheet error", "err", err)
		return err
	}
	xlsx.SetActiveSheet(index)
	sType := reflect.TypeOf(records)
	if sType.Kind() != reflect.Slice {
		return errors.New("数据非 slice 类型")
	}
	sValue := reflect.ValueOf(records)
	if sValue.Len() < 1 {
		return errors.New("数据长度为0. 无效数据长度")
	}
	setTitleRow(xlsx, sheetName, sValue.Index(0).Interface())
	setRow(xlsx, sheetName, records)
	// 根据指定路径保存文件
	if err := xlsx.SaveAs("Book1.xlsx"); err != nil {
		slog.Error("save file error", "err", err)
	}
	return nil
}

var basicColIndex = 65 // "A"

func setTitleRow(xlsx *excelize.File, sheetName string, obj interface{}) {
	elemType := reflect.TypeOf(obj)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("xlsx")
		if strings.TrimSpace(tag) == "" {
			continue
		}
		name := tag
		position := fmt.Sprintf("%s%d", string(rune(basicColIndex+i)), 1)
		err := xlsx.SetCellValue(sheetName, position, name)
		if err != nil {
			slog.Error("set cell value error", "err", err)
		}
	}
}

func setRow(xlsx *excelize.File, sheetName string, objs interface{}) {
	// sliceType := reflect.TypeOf(objs)
	slice := reflect.ValueOf(objs)

	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i).Interface()
		elemType := reflect.TypeOf(elem)
		elemValue := reflect.ValueOf(elem)
		for j := 0; j < elemType.NumField(); j++ {
			field := elemType.Field(j)
			tag := field.Tag.Get("xlsx")
			if strings.TrimSpace(tag) == "" {
				// 会导致空列的出现
				continue
			}
			position := fmt.Sprintf("%s%d", string(rune(basicColIndex+j)), i+2)
			err := xlsx.SetCellValue(sheetName, position, elemValue.Field(j).Interface())
			if err != nil {
				slog.Error("set cell value error", "err", err)
			}
		}
	}
}
