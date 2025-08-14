package main

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"testing"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo/imongoutil"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"
)

func TestSift(t *testing.T) {
	sift := gocv.NewSIFT()
	defer sift.Close()
	svc := butterflysvc.NewButterflyResizedImgSvc()
	list, _ := svc.GetAllList(context.Background(), bson.M{})
	for _, item := range list {
		mat, _ := gocv.NewMatFromBytes(200, 200, gocv.MatTypeCV8UC3, item.Content)
		gocv.CvtColor(mat, &mat, gocv.ColorBGRToGray)
		_, discrib := sift.DetectAndCompute(mat, mat)
		dbmat, _ := imgutils.Mat2DBMat(&discrib)
		svc.Update(context.Background(), bson.M{"_id": item.ID}, bson.M{
			"$set": bson.M{
				"describ_mat": dbmat,
			},
		})
	}
}

func TestKmeans(t *testing.T) {
	buf := make([]byte, 0, 60*1024*1024)
	var rows int
	svc := butterflysvc.NewButterflyResizedImgSvc()
	list, _ := svc.GetAllList(context.Background(), bson.M{})
	for _, item := range list {
		buf = append(buf, item.DescribMat.Context...)
		rows += item.DescribMat.Row
	}
	allDescrib, err := gocv.NewMatFromBytes(rows, 128, gocv.MatTypeCV32FC1, buf)
	if err != nil {
		slog.Error("cvt all describ is faild")
	}
	labels := gocv.NewMat()
	centers := gocv.NewMat()
	criteria := gocv.NewTermCriteria(gocv.Count|gocv.EPS, 10, 1.0)

	k := 20
	res := gocv.KMeans(allDescrib, k, &labels, criteria, 3, gocv.KMeansPPCenters, &centers)
	fmt.Printf("Mat allDescrib Size: %dx%d, Type: %d\n", allDescrib.Rows(), allDescrib.Rows(), allDescrib.Type())
	fmt.Printf("Mat Labels Size: %dx%d, Type: %d\n", labels.Rows(), labels.Rows(), labels.Type())
	slog.Info("结果", "res", res)

	// labels 是一条 features 所分配到 centers 的节点
	storeKmeans(&labels, &centers)

	var start int
	tranningData := gocv.NewMat()

	for _, item := range list {
		imgTag, _ := strconv.Atoi(item.Type)
		bow := buildBowHistogram(&labels, k, 0, item.DescribMat.Row, imgTag)
		start += item.DescribMat.Row
		if tranningData.Cols() == 0 {
			tranningData = bow
		} else {
			gocv.Vconcat(tranningData, bow, &tranningData)
		}
	}

	rows, cols := tranningData.Rows(), tranningData.Cols()
	fmt.Printf("Mat Size: %dx%d, Type: %d\n", rows, cols, tranningData.Type())
	// gocv.Normalize(tranningData, &tranningData, 0.0, 1.0, gocv.NormMinMax)
	for index, item := range list {
		tp, _ := strconv.Atoi(item.Type)
		tranningData.SetFloatAt(index, k, float32(tp))
	}

	imgutils.SaveMatToCSV(tranningData, "data.csv")
}

func BytesToFloat32sUnsafe(b []byte) []float32 {
	hdr := *(*[]float32)(unsafe.Pointer(&b))
	hdrLen := len(b) / 4
	return hdr[:hdrLen:hdrLen] // 明确 slice 长度
}

// buildBowHistogram creates a histogram of visual word occurrences (Bag of Words) from cluster labels.
// It takes a Mat of cluster labels, number of clusters k, start and end indices to process,
// and a type parameter tp.
// The (k+1)-th column represents the tag type corresponding to the image.
func buildBowHistogram(labels *gocv.Mat, k, start, end, tag int) gocv.Mat {
	bowi := make([]int, k)
	bow := gocv.NewMatWithSize(1, k+1, gocv.MatTypeCV32FC1)
	for i := start; i < end; i++ {
		bowi[labels.GetIntAt(1, i)]++
	}
	for idx, value := range bowi {
		bow.SetFloatAt(0, idx, float32(value))
	}
	bow.SetFloatAt(0, k, float32(tag))
	return bow
}

func storeKmeans(labels, centers *gocv.Mat) {
	mat1, _ := imgutils.Mat2DBMat(labels)
	mat2, _ := imgutils.Mat2DBMat(centers)
	obj := bson.M{"lables": mat1, "centers": mat2}
	imongoutil.Insert[bson.M](context.Background(), imongo.FileDatabase().Collection("kmeans"), obj)
}

func Test(t *testing.T) {
	a := 123
	slog.Info("", "int 2 to float", float32(a))
}
