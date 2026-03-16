package main

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/internal/model/insect"
	"github.com/HunDun0Ben/bs_server/app/internal/repository"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"
)

var (
	mode     string
	filepath string
)

func init() {
	pflag.StringVarP(&mode, "mode", "m", "init", "Operation mode: init (load type info), resize (resize images), display-one, display-batch")
	pflag.StringVarP(&filepath, "file", "f", "./scripts/initbutterflytypeinfo/蝴蝶信息.xlsx", "Path to butterfly info Excel/CSV file")
}

func main() {
	pflag.Parse()

	if err := conf.InitConfig(); err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	repo := repository.NewButterflyRepository(imongo.BizDataBase())
	svc := butterflysvc.NewButterflyService(repo)

	switch mode {
	case "init":
		InitButterflyTypes(svc)
	case "resize":
		ResizeImages(svc)
	case "display-one":
		ResizeImgOne(svc)
	case "display-batch":
		DisplayResizedImages(svc)
	default:
		slog.Warn("Unknown mode", "mode", mode)
		pflag.Usage()
	}
}

func InitButterflyTypes(svc butterflysvc.ButterflyService) {
	count, err := svc.CountTypes(context.Background())
	if err != nil {
		slog.Error("Failed to count types", "error", err)
		os.Exit(1)
	}
	if count > 0 {
		slog.Info("Already initialized data")
		return
	}
	list, err := loadTypeInfoFromExcel(filepath)
	if err != nil {
		slog.Error("Failed to load information", "error", err)
		os.Exit(1)
	}
	if err := svc.InitTypes(context.Background(), list); err != nil {
		slog.Error("Failed to initialize butterfly information", "error", err)
		os.Exit(1)
	}
	slog.Info("Butterfly types initialized successfully")
}
func ResizeImages(svc butterflysvc.ButterflyService) {
	typeList, err := svc.GetTypes(context.Background())
	if err != nil {
		slog.Error("Failed to get types", "error", err)
		os.Exit(1)
	}

	for _, v := range typeList {
		res, err := svc.GetImgs(context.Background(), bson.M{"file_name": bson.M{"$regex": v.LatinName}}) // Using LatinName as type marker
		if err != nil {
			slog.Error("Failed to get images for type", "type", v.ChineseName, "error", err)
			continue
		}
		slog.Info("length of img list", "type", v.ChineseName, "len", len(res))
		for _, info := range res {
			func() {
				src := imgutils.GetMaskImg(info)
				defer src.Close()

				mat := imgutils.ResizeWithPadding(*src, 200, 200, color.RGBA{0, 0, 0, 0})
				defer mat.Close()

				resized := file.ResizedButteryflyFile{
					FileStoreData: imongo.FileStoreData{
						Content: mat.ToBytes(),
					},
					Col:  mat.Cols(),
					Row:  mat.Rows(),
					Type: v.LatinName,
				}
				err := svc.InsertResizedImg(context.Background(), &resized)
				if err != nil {
					slog.Error("Failed to insert resized image", "type", v.ChineseName, "error", err)
				}
			}()
		}
	}
	slog.Info("Image resizing completed")
}

func ResizeImgOne(svc butterflysvc.ButterflyService) {
	res, err := svc.FindImg(context.Background(), bson.M{})
	if err != nil || res == nil {
		slog.Error("No image found", "error", err)
		os.Exit(1)
	}

	src := imgutils.GetMaskImg(*res)
	defer src.Close()

	mat := imgutils.ResizeWithPadding(*src, 200, 200, color.RGBA{0, 0, 0, 0})
	defer mat.Close()

	slog.Info("Mat type", "type", mat.Type().String())
	win := ui.NewProcessingWindow("Resize Image")
	win.LoadImageFromMat(mat)
	win.Display()
}

func DisplayResizedImages(svc butterflysvc.ButterflyService) {
	list, err := svc.GetResizedImgs(context.Background(), bson.M{})
	if err != nil {
		slog.Error("Failed to get resized list", "error", err)
		os.Exit(1)
	}

	for i := 0; i < 10 && i < len(list); i++ {
		func() {
			mat, err := gocv.NewMatFromBytes(200, 200, gocv.MatTypeCV8UC3, list[i].Content)
			if err != nil {
				slog.Error("Error creating mat for index", "index", i, "error", err)
				return
			}
			defer mat.Close()

			slog.Info("Displaying image", "index", i, "type", mat.Type())
			win := ui.NewProcessingWindow(fmt.Sprintf("Resize Image %d", i))
			defer win.Close()
			win.LoadImageFromMat(mat)
			win.Display()
		}()
	}
}

func loadTypeInfoFromExcel(path string) ([]insect.Insect, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %s, 错误: %v", path, err)
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return nil, fmt.Errorf("无法获取工作表的行数据: %v", err)
	}

	list := make([]insect.Insect, 0)
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}
		item := insect.Insect{
			ChineseName:        row[0],
			LatinName:          row[1],
			EnglishName:        row[2],
			FeatureDescription: row[3],
		}
		list = append(list, item)
	}
	return list, nil
}
