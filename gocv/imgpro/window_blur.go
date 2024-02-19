package imgpro

import (
	"demo/gocv/imgpro/constant"
	"image"

	"gocv.io/x/gocv"
)

type TrackWindowBlurCxt struct {
	BlurType int
	Ksize    int
}

func (cxt *TrackWindowBlurCxt) process(src *gocv.Mat) *gocv.Mat {
	dstp := NewSomeMat(*src)
	size := 1
	if cxt.Ksize > 1 {
		size = cxt.Ksize
	}
	switch cxt.BlurType {
	case constant.Blur:
		gocv.Blur(*src, dstp, image.Point{size, size})
	case constant.GaussianBlur:
		if size%2 == 0 {
			size++
		}
		gocv.GaussianBlur(*src, dstp, image.Point{size, size}, 0, 0, 0)
	case constant.MedianBlur:
		if size%2 == 0 {
			size++
		}
		gocv.MedianBlur(*src, dstp, size)
	case constant.BilateralFiter:
		gocv.BilateralFilter(*src, dstp, size, float64(size*2), float64(size/2))
	}
	return dstp
}
