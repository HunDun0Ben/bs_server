package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"log/slog"

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
		log.Fatalf("Failed to initialize config: %v", err)
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
		fmt.Printf("Unknown mode: %s\n", mode)
		pflag.Usage()
	}
}

func InitButterflyTypes(svc butterflysvc.ButterflyService) {
	count, err := svc.CountTypes(context.Background())
	if err != nil {
		log.Fatalf("获取数据数量失败: %v", err)
	}
	if count > 0 {
		slog.Info("已经初始化过数据")
		return
	}
	list, err := loadTypeInfoFromExcel(filepath)
	if err != nil {
		log.Fatalf("加载信息失败: %v", err)
	}
	if err := svc.InitTypes(context.Background(), list); err != nil {
		log.Fatalf("初始化蝴蝶信息失败: %v", err)
	}
	fmt.Println("Butterfly types initialized successfully")
}

func ResizeImages(svc butterflysvc.ButterflyService) {
	typeList, err := svc.GetTypes(context.Background())
	if err != nil {
		log.Fatalf("Failed to get types: %v", err)
	}

	for _, v := range typeList {
		res, err := svc.GetImgs(context.Background(), bson.M{"file_name": bson.M{"$regex": v.LatinName}}) // Using LatinName as type marker
		if err != nil {
			log.Printf("Failed to get images for type %s: %v", v.ChineseName, err)
			continue
		}
		slog.Info("length of img list", "type", v.ChineseName, "len", len(res))
		for _, info := range res {
			mat := imgutils.ResizeWithPadding(*imgutils.GetMaskImg(info), 200, 200, color.RGBA{0, 0, 0, 0})
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
				log.Printf("Failed to insert resized image for type %s: %v", v.ChineseName, err)
			}
		}
	}
	fmt.Println("Image resizing completed")
}

func ResizeImgOne(svc butterflysvc.ButterflyService) {
	res, err := svc.FindImg(context.Background(), bson.M{})
	if err != nil || res == nil {
		log.Fatalf("No image found: %v", err)
	}
	mat := imgutils.ResizeWithPadding(*imgutils.GetMaskImg(*res), 200, 200, color.RGBA{0, 0, 0, 0})
	defer mat.Close()

	slog.Info("Mat type", "type", mat.Type().String())
	win := ui.NewProcessingWindow("Resize Image")
	win.LoadImageFromMat(mat)
	win.Display()
}

func DisplayResizedImages(svc butterflysvc.ButterflyService) {
	list, err := svc.GetResizedImgs(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to get resized list: %v", err)
	}

	for i := 0; i < 10 && i < len(list); i++ {
		mat, err := gocv.NewMatFromBytes(200, 200, gocv.MatTypeCV8UC3, list[i].Content)
		if err != nil {
			log.Printf("Error creating mat for index %d: %v", i, err)
			continue
		}
		defer mat.Close()

		fmt.Printf("Displaying image %d, type: %d\n", i, mat.Type())
		win := ui.NewProcessingWindow(fmt.Sprintf("Resize Image %d", i))
		win.LoadImageFromMat(mat)
		win.Display()
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
