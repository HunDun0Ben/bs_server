package initbutterflyinfo_test

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"
)

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
	slog.Info("Mat type", slog.String("type", mat.Type().String()))
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
