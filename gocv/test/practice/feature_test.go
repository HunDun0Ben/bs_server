package practice_test

import (
	"demo/gocv/imgpro"
	"demo/gocv/imgpro/img/feature"
	"image"
	"image/color"
	"log"
	"math/rand"
	"testing"
	"time"

	"gocv.io/x/gocv"
)

func Test(t *testing.T) {
	filename := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	wrapper := imgpro.NewTrackWindowWrapper("Erosion Demo", filename)
	wrapper.LoadImg()
	wrapper.DstImgMat = feature.GetImgSIFT(&wrapper.SrcImgMat)
	wrapper.Dispaly()
}

func TestDes(t *testing.T) {
	// filename1 := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	// filename2 := `/home/workspace/data/leedsbutterfly/images/0010002.png`
	filename1 := `/home/workspace/data/img/desk001.jpg`
	filename2 := `/home/workspace/data/img/desk002.jpg`
	mat1 := gocv.IMRead(filename1, gocv.IMReadColor)
	mat2 := gocv.IMRead(filename2, gocv.IMReadColor)
	mask := gocv.NewMat()
	size := mat1.Size()
	log.Printf("mat2 = %v", mat2.Size())
	gocv.Resize(mat1, &mat1, image.Pt(size[1], size[1]), 0, 0, gocv.InterpolationNearestNeighbor)
	gocv.Resize(mat2, &mat2, image.Pt(size[1], size[1]), 0, 0, gocv.InterpolationNearestNeighbor)
	log.Printf("mat1 = %v", mat1.Size())
	log.Printf("mat2 = %v", mat2.Size())
	log.Printf("mat1 = %v channel", mat1.Channels())
	log.Printf("mat2 = %v", mat2.Channels())
	sift := gocv.NewORB()
	// orb := gocv.NewORB()
	mat1.ConvertTo(&mat1, gocv.MatTypeCV8U)
	mat2.ConvertTo(&mat2, gocv.MatTypeCV8U)
	kp1, des1 := sift.DetectAndCompute(mat1, mask)
	kp2, des2 := sift.DetectAndCompute(mat2, mask)
	matcher := gocv.NewBFMatcher()
	dm := matcher.KnnMatch(des1, des2, 1)

	imgMatches := gocv.NewMat()
	gocv.Hconcat(mat1, mat2, &imgMatches)
	rand.NewSource(time.Now().UnixNano())
	window := gocv.NewWindow("123")
	// 绘制匹配的特征点
	for _, matches := range dm {
		copy := gocv.NewMat()
		imgMatches.CopyTo(&copy)
		for _, match := range matches {
			r := uint8(rand.Intn(255))
			g := uint8(rand.Intn(255))
			b := uint8(rand.Intn(255))
			pt1 := image.Point{int(kp1[match.QueryIdx].X), int(kp1[match.QueryIdx].Y)}
			pt2 := image.Point{int(kp2[match.TrainIdx].X) + mat1.Cols(), int(kp2[match.TrainIdx].Y)}
			gocv.Circle(&copy, pt1, 5, color.RGBA{r, g, b, 0}, 2)
			gocv.Circle(&copy, pt2, 5, color.RGBA{r, g, b, 0}, 2)
			gocv.Line(&copy, pt1, pt2, color.RGBA{r, g, b, 0}, 2)
			for {
				window.IMShow(copy)
				if window.WaitKey(1) >= 0 || window.GetWindowProperty(gocv.WindowPropertyVisible) == 0 {
					break
				}
			}
		}
	}
	dst := gocv.NewMat()
	for _, d := range dm {

		gocv.DrawMatches(mat1, kp1, mat2, kp2, d, &dst,
			color.RGBA{255, 0, 0, 0}, color.RGBA{0, 255, 0, 0}, nil, gocv.NormconvFilter)
	}

	for {
		window.IMShow(imgMatches)
		if window.WaitKey(1) >= 0 || window.GetWindowProperty(gocv.WindowPropertyVisible) == 0 {
			break
		}
	}
	window.Close()
	// gocv.DrawMatches(mat1, kp1, mat2, kp2, dm[0], &dst,
	// 	color.RGBA{255, 0, 0, 0}, color.RGBA{0, 255, 0, 0}, nil, gocv.NormconvFilter)

}
