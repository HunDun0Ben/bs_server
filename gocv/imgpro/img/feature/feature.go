package feature

import (
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func GetImgSIFT(mat *gocv.Mat) gocv.Mat {
	if mat.Empty() {
		log.Printf("该图像为空.")
	}
	sift := gocv.NewSIFT()
	keypoints := sift.Detect(*mat)
	sift.Close()
	gocv.DrawKeyPoints(*mat, keypoints, mat, color.RGBA{255, 0, 0, 255}, gocv.DrawDefault)
	return *mat
}
