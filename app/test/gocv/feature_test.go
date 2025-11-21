package gocv_test

import (
	"context"
	"image"
	"image/color"
	"log/slog"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"

	timerutil "github.com/HunDun0Ben/bs_server/app/pkg/util"
)

func Test(t *testing.T) {
	filename := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	win := ui.NewProcessingWindow("Erosion Demo")
	if err := win.LoadImageFromPath(filename); err != nil {
		os.Exit(1)
	}
	win.LoadImageFromPath(filename)
	win.Process(func(src *gocv.Mat) *gocv.Mat {
		pmat, _ := imgutils.DrawImgSIFT(win.GetDstMat())
		return pmat
	})
	win.Display()
}

func TestDes(t *testing.T) {
	filename1 := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	filename2 := `/home/workspace/data/leedsbutterfly/images/0010002.png`
	// filename1 := `/home/workspace/data/img/1.jpg`
	// filename2 := `/home/workspace/data/img/2.jpg`
	mat1 := gocv.IMRead(filename1, gocv.IMReadColor)
	mat2 := gocv.IMRead(filename2, gocv.IMReadColor)

	// 如果有需要将图片的 type 进行转换
	// mat1.ConvertTo(&mat1, gocv.MatTypeCV8U)
	// mat2.ConvertTo(&mat2, gocv.MatTypeCV8U)

	// 因为 Hconcat 需要, 转化为同高的图片大小
	gocv.Resize(mat2, &mat2, image.Pt(mat2.Size()[1], mat1.Size()[0]), 0, 0, gocv.InterpolationNearestNeighbor)
	// 水平合并 mat1 ,mat2
	concatImg := gocv.NewMat()
	gocv.Hconcat(mat1, mat2, &concatImg)

	maskfilename1 := `/home/workspace/data/leedsbutterfly/segmentations/0010001_seg0.png`
	maskfilename2 := `/home/workspace/data/leedsbutterfly/segmentations/0010002_seg0.png`

	// 通过 mask 去除无关部分. 专注需要计算部分的特征点
	mask1 := gocv.IMRead(maskfilename1, gocv.IMReadGrayScale)
	mask2 := gocv.IMRead(maskfilename2, gocv.IMReadGrayScale)

	sift := gocv.NewSIFT()
	kp1, des1 := sift.DetectAndCompute(mat1, mask1)
	kp2, des2 := sift.DetectAndCompute(mat2, mask2)
	matcher := gocv.NewBFMatcher()
	dm := matcher.KnnMatch(des1, des2, 2)

	// 复制体?
	imgCopy := gocv.NewMat()
	concatImg.CopyTo(&imgCopy)

	// 筛选出优秀匹配
	var goodMatches []gocv.DMatch
	for _, m := range dm {
		if len(m) == 2 && m[0].Distance < 0.75*m[1].Distance {
			goodMatches = append(goodMatches, m[0])
		}
	}

	// 绘制匹配的特征点
	for _, match := range goodMatches {
		clr := *imgutils.RandColor()
		pt1 := image.Point{int(kp1[match.QueryIdx].X), int(kp1[match.QueryIdx].Y)}
		pt2 := image.Point{int(kp2[match.TrainIdx].X) + mat1.Cols(), int(kp2[match.TrainIdx].Y)}
		gocv.Circle(&imgCopy, pt1, 5, clr, 2)
		gocv.Circle(&imgCopy, pt2, 5, clr, 2)
		gocv.Line(&imgCopy, pt1, pt2, clr, 2)
	}

	window := gocv.NewWindow("123")
	for {
		window.IMShow(imgCopy)
		if window.WaitKey(100) >= 0 || window.GetWindowProperty(gocv.WindowPropertyVisible) == 0 {
			break
		}
	}

	dst := gocv.NewMat()
	clr := *imgutils.RandColor()
	gocv.DrawMatches(mat1, kp1, mat2, kp2, goodMatches, &dst,
		clr, clr, nil, gocv.NormconvFilter)
	for {
		window.IMShow(dst)
		if window.WaitKey(100) >= 0 || window.GetWindowProperty(gocv.WindowPropertyVisible) == 0 {
			break
		}
	}
	window.Close()
}

func TestLoadMaskImg(t *testing.T) {
	var fileInfo file.ButterflyFile
	err := imongo.FileDatabase().Collection("butterfly_img").FindOne(context.Background(), bson.D{}).Decode(&fileInfo)
	if err != nil {
		t.Fatalf("Failed to find file: %v", err)
		return
	}
	win := ui.NewProcessingWindow("after mask of img demo")
	timerutil.D.StartTimer()
	maskimg := imgutils.GetMaskImg(fileInfo)
	gocv.CvtColor(*maskimg, maskimg, gocv.ColorBGRToGray)
	resized := imgutils.ResizeWithPadding(*maskimg, 200, 200, color.RGBA{0, 0, 0, 0})
	resized.CopyTo(maskimg)
	pmat, _ := imgutils.DrawImgSIFT(win.GetDstMat())
	slog.Info("cost time: ", "d", timerutil.D.GetTimer())
	win.LoadImageFromMat(*pmat)
	win.Display()
}
