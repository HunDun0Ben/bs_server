package initbutterflyinfo_test

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"testing"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/internal/model/insect"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"
)

func Test_init_butterfly_type_info(t *testing.T) {
	list, _ := loadTypeInfoFromCSV()
	err := butterflysvc.NewButterflyTypeSvc().InitAll(context.Background(), list)
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

func TestResizeImg(t *testing.T) {
	typeSvc := butterflysvc.NewButterflyTypeSvc()
	svc := butterflysvc.NewButterflyImgSvc()
	resizedSvc := butterflysvc.NewButterflyResizedImgSvc()
	typeList, _ := typeSvc.GetAllList(context.Background())
	for _, v := range typeList {
		res, _ := svc.GetAllList(context.Background(), bson.M{"file_name": bson.M{"$regex": v.Type}})
		slog.Info("length of img list", slog.Int("len", len(res)))
		for _, info := range res {
			mat := imgutils.ResizeWithPadding(*imgutils.GetMaskImg(info), 200, 200, color.RGBA{0, 0, 0, 0})
			resized := file.ResizedButteryflyFile{
				FileStoreData: imongo.FileStoreData{
					Content: mat.ToBytes(),
				},
				Col:  mat.Cols(),
				Row:  mat.Rows(),
				Type: v.Type,
			}
			resizedSvc.Insert(context.Background(), &resized)
		}
	}
}

func TestResizeImgOne(t *testing.T) {
	svc := butterflysvc.NewButterflyImgSvc()
	resizedSvc := butterflysvc.NewButterflyResizedImgSvc()
	res, _ := svc.FindOne(context.Background(), bson.M{})
	mat := imgutils.ResizeWithPadding(*imgutils.GetMaskImg(*res), 200, 200, color.RGBA{0, 0, 0, 0})

	resized := file.ResizedButteryflyFile{
		FileStoreData: imongo.FileStoreData{
			Content: mat.ToBytes(),
		},
		Col: mat.Cols(),
		Row: mat.Rows(),
	}
	resizedSvc.Insert(context.Background(), &resized)
	fmt.Print("mat type", mat.Type())
	slog.Info("mat Type", slog.String("type", fmt.Sprintf("%T", mat.Type())))
	win := ui.NewProcessingWindow("Resize Image")
	win.LoadImageFromMat(mat)
	win.Display()
}

func TestWrite(t *testing.T) {
	resizedSvc := butterflysvc.NewButterflyResizedImgSvc()
	list, _ := resizedSvc.GetAllList(context.Background(), bson.M{})
	for i := 0; i < 10; i++ {
		mat, _ := gocv.NewMatFromBytes(200, 200, gocv.MatTypeCV8UC3, list[i].Content)
		fmt.Print("mat type", mat.Type())
		win := ui.NewProcessingWindow("Resize Image")
		win.LoadImageFromMat(mat)
		win.Display()
	}
}
