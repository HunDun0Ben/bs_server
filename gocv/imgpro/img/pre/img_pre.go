package pre

import (
	"image"

	"gocv.io/x/gocv"
)

const _unfiicationSize = 500

// 批量统一图片大小至500以下
func UnificationSizeMats(mats ...gocv.Mat) {
	for _, mat := range mats {
		UnificationSizeMat(mat)
	}
}

// 批量统一图片大小至500以下, 每次减半
func UnificationSizeMat(mat gocv.Mat) {
	if mat.Cols() < _unfiicationSize && mat.Rows() < _unfiicationSize {
		return
	}
	gocv.PyrDown(mat, &mat, image.Point{mat.Cols() / 2, mat.Rows() / 2}, gocv.BorderDefault)
}
