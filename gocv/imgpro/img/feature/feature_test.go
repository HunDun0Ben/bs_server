package feature_test

import (
	"demo/gocv/imgpro/core/ui"
	"demo/gocv/imgpro/img/feature"
	"demo/gocv/imgpro/img/utils"
	"image"
	"testing"

	"gocv.io/x/gocv"
)

func Test(t *testing.T) {
	filename := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	win := ui.NewProcessingWindow("Erosion Demo")
	win.LoadImageFromPath(filename)
	win.Process(func(src *gocv.Mat) *gocv.Mat {
		return feature.DrawImgSIFT(win.GetDstMat())
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
	copy := gocv.NewMat()
	concatImg.CopyTo(&copy)

	// 筛选出优秀匹配
	var goodMatches []gocv.DMatch
	for _, m := range dm {
		if len(m) == 2 && m[0].Distance < 0.75*m[1].Distance {
			goodMatches = append(goodMatches, m[0])
		}
	}

	// 绘制匹配的特征点
	for _, match := range goodMatches {
		clr := *utils.RandColor()
		pt1 := image.Point{int(kp1[match.QueryIdx].X), int(kp1[match.QueryIdx].Y)}
		pt2 := image.Point{int(kp2[match.TrainIdx].X) + mat1.Cols(), int(kp2[match.TrainIdx].Y)}
		gocv.Circle(&copy, pt1, 5, clr, 2)
		gocv.Circle(&copy, pt2, 5, clr, 2)
		gocv.Line(&copy, pt1, pt2, clr, 2)
	}

	window := gocv.NewWindow("123")
	for {
		window.IMShow(copy)
		if window.WaitKey(100) >= 0 || window.GetWindowProperty(gocv.WindowPropertyVisible) == 0 {
			break
		}
	}

	dst := gocv.NewMat()
	clr := *utils.RandColor()
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
