package service

import (
	"image"

	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/constant"
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
