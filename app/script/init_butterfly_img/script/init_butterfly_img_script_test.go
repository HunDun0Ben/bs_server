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

func TestUpdateShiftFeature(t *testing.T) {
	sift := gocv.NewSIFT()
	defer sift.Close()
	svc := butterflysvc.NewButterflyResizedImgSvc()
	resizeList, _ := svc.GetAllList(context.Background(), bson.M{})
	for _, item := range resizeList {
		func() {
			mat, _ := gocv.NewMatFromBytes(200, 200, gocv.MatTypeCV8UC3, item.Content)
			dst := gocv.NewMat()
			defer mat.Close()
			defer dst.Close()
			gocv.CvtColor(mat, &dst, gocv.ColorBGRToGray)
			_, discrib := sift.DetectAndCompute(dst, dst)
			defer discrib.Close()
			// shift 特征的大小通常是128
			dbmat, _ := imgutils.Mat2DBMat(&discrib)
			svc.Update(context.Background(), bson.M{"_id": item.ID}, bson.M{
				"$set": bson.M{
					"describ_mat": dbmat,
				},
			})
		}()
	}
}

var k = 1024
var iterations = 10

func TestKmeans(t *testing.T) {
	buf := make([]byte, 0, 60*1024*1024)
	var rows int
	svc := butterflysvc.NewButterflyResizedImgSvc()
	resizeList, _ := svc.GetAllList(context.Background(), bson.M{})
	// shift 特征的大小通常是128, 拼接所有特征作为 kmeans 所需要的数据集合
	for _, item := range resizeList {
		buf = append(buf, item.DescribMat.Context...)
		rows += item.DescribMat.Row
	}
	// 初始化所有样本的特征点矩阵
	allDescrib, err := gocv.NewMatFromBytes(rows, 128, gocv.MatTypeCV32FC1, buf)
	slog.Info(fmt.Sprintf("Mat allDescrib Size: %dx%d, Type: %d\n",
		allDescrib.Rows(), allDescrib.Cols(), allDescrib.Type()))
	if err != nil {
		slog.Error("cvt all describ is faild")
	}
	defer allDescrib.Close()

	// 特征对应的聚簇的 label. size = feature size
	labels := gocv.NewMat()
	// 每个聚类簇（cluster）的中心位置. size = k
	centers := gocv.NewMat()
	// 终止条件
	criteria := gocv.NewTermCriteria(gocv.Count|gocv.EPS, iterations, 1.0)

	defer labels.Close()
	defer centers.Close()

	// k 表示 聚类的簇（cluster）数
	res := gocv.KMeans(allDescrib, k, &labels, criteria, 3, gocv.KMeansPPCenters, &centers)
	slog.Info(fmt.Sprintf("Mat allDescrib Size: %dx%d, Type: %d\n",
		allDescrib.Rows(), allDescrib.Cols(), allDescrib.Type()))
	slog.Info(fmt.Sprintf("Mat Labels Size: %dx%d, Type: %d\n",
		labels.Rows(), labels.Cols(), labels.Type()))
	slog.Info("结果", "res", res)

	// labels 是一条 features 所分配到 centers 的节点
	storeKmeans(&labels, &centers)

	var start int
	tranningData := gocv.NewMat()

	for _, item := range resizeList {
		imgTag, _ := strconv.Atoi(item.Type)
		bow := buildBowHistogram(&labels, k, start, item.DescribMat.Row, imgTag)
		start += item.DescribMat.Row
		if tranningData.Cols() == 0 {
			tranningData = bow
		} else {
			gocv.Vconcat(tranningData, bow, &tranningData)
		}
	}

	rows, cols := tranningData.Rows(), tranningData.Cols()
	slog.Info(fmt.Sprintf("Mat Size: %dx%d, Type: %d\n", rows, cols, tranningData.Type()))
	// gocv.Normalize(tranningData, &tranningData, 0.0, 1.0, gocv.NormMinMax)
	// for index, item := range resizeList {
	// 	tp, _ := strconv.Atoi(item.Type)
	// 	tranningData.SetFloatAt(index, k, float32(tp))
	// }

	// 对 tranningData 的每一行进行 L2 归一化
	for i := 0; i < tranningData.Rows(); i++ {
		row := tranningData.RowRange(i, i+1)
		// 注意：只对特征部分（前k列）进行归一化，最后一列是标签
		featureRow := row.ColRange(0, k)
		gocv.Normalize(featureRow, &featureRow, 1.0, 0.0, gocv.NormL2)
	}

	defer tranningData.Close()
	slog.Info(fmt.Sprintf("Mat tranningData Size: %dx%d, Type: %d\n",
		tranningData.Rows(), tranningData.Cols(), tranningData.Type()))
	imgutils.SaveMatToCSV(tranningData, "data.csv")
}

func BytesToFloat32sUnsafe(b []byte) []float32 {
	hdr := *(*[]float32)(unsafe.Pointer(&b))
	hdrLen := len(b) / 4
	return hdr[:hdrLen:hdrLen] // 明确 slice 长度
}

// buildBowHistogram creates a histogram of visual word occurrences (Bag of Words) from cluster labels.
// It takes a Mat of cluster labels, number of clusters k, a start index, the number of elements to process (num),
// and a type parameter tp.
// The (k+1)-th column represents the tag type corresponding to the image.
func buildBowHistogram(labels *gocv.Mat, k, start, num, tag int) gocv.Mat {
	bowi := make([]int, k)
	bow := gocv.NewMatWithSize(1, k+1, gocv.MatTypeCV32FC1)
	for i := start; i < start+num; i++ {
		bowi[labels.GetIntAt(i, 0)]++
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
