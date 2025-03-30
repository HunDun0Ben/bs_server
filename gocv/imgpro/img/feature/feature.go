package feature

import (
	"log/slog"

	"github.com/HunDun0Ben/bs_server/gocv/imgpro/img/utils"
	"gocv.io/x/gocv"
)

func DrawImgSIFT(mat *gocv.Mat) *gocv.Mat {
	if mat.Empty() {
		slog.Info("该图像为空.")
	}
	sift := gocv.NewSIFT()
	defer sift.Close()
	keypoints := sift.Detect(*mat)
	gocv.DrawKeyPoints(*mat, keypoints, mat,
		*utils.RandColor(), gocv.DrawRichKeyPoints)
	return mat
}
