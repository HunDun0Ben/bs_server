package imgutils

import (
	"errors"
	"log/slog"

	"gocv.io/x/gocv"
)

func DrawImgSIFT(mat *gocv.Mat) (*gocv.Mat, error) {
	if mat.Empty() {
		return nil, errors.New("image cannot be empty")
	}
	sift := gocv.NewSIFT()
	defer sift.Close()
	// keypoints := sift.Detect(*mat)
	keypoints, discrib := sift.DetectAndCompute(*mat, *mat)
	slog.Info("SIFT 描述符数量: ", "rows", discrib.Rows(), "Cols", discrib.Cols(), "Type", discrib.Type())
	gocv.DrawKeyPoints(*mat, keypoints, mat,
		*RandColor(), gocv.DrawRichKeyPoints)
	return mat, nil
}
