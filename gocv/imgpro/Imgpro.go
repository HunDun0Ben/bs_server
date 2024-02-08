package imgpro

import (
	"gocv.io/x/gocv"
)

func NewSomeMat(src gocv.Mat) *gocv.Mat {
	a := gocv.NewMatWithSizes(src.Size(), src.Type())
	return &a
}
