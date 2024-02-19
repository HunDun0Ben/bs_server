package service

import (
	"demo/gocv/imgpro/constant"
	"image"

	"gocv.io/x/gocv"
)

func PreMat(src gocv.Mat, preType []int) {
	sum := 0
	for _, v := range preType {
		sum += v
	}
	if sum&constant.GaussianBlur != 0 {
		gocv.GaussianBlur(src, &src, image.Point{3, 3}, 0, 0, gocv.BorderConstant)
	}
}
