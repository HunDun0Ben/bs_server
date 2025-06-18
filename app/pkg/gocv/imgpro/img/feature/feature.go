package feature

import (
	"log/slog"

	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"
)

func DrawImgSIFT(mat *gocv.Mat) *gocv.Mat {
	if mat.Empty() {
		slog.Info("该图像为空.")
	}
	sift := gocv.NewSIFT()
	defer sift.Close()
	// keypoints := sift.Detect(*mat)
	keypoints, discrib := sift.DetectAndCompute(*mat, *mat)
	slog.Info("SIFT 描述符数量: ", "rows", discrib.Rows(), "Cols", discrib.Cols(), "Type", discrib.Type())
	gocv.DrawKeyPoints(*mat, keypoints, mat,
		*imgutils.RandColor(), gocv.DrawRichKeyPoints)
	return mat
}
