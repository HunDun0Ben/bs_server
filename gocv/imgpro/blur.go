package imgpro

import (
	"image"

	"gocv.io/x/gocv"
)

func Blur(src gocv.Mat) *gocv.Mat {
	dstp := NewSomeMat(src)
	gocv.Blur(src, dstp, image.Point{3, 3})
	return dstp
}

func GaussianBlur(src gocv.Mat) *gocv.Mat {
	dstp := NewSomeMat(src)
	// gocv.GaussianBlur(src, dstp, image.Point{3, 3})
	return dstp
}

func NewSomeMat(src gocv.Mat) *gocv.Mat {
	a := gocv.NewMatWithSizes(src.Size(), src.Type())
	return &a
}
