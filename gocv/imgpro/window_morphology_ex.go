package imgpro

import (
	"image"

	"gocv.io/x/gocv"
)

const (
	MorphMaxOperator   = 4
	MorphMaxElem       = 2
	MorphMaxKernelSize = 21
)

type TrackWindowMorphologyExCxt struct {
	Operator  int
	MorphElem int
	MorphSize int
}

func (cxt *TrackWindowMorphologyExCxt) process(src *gocv.Mat) *gocv.Mat {
	dstp := NewSomeMat(*src)
	operation := cxt.Operator + 2
	if cxt.MorphElem < 0 || cxt.MorphElem > 2 {
		cxt.MorphElem = 0
	}
	kernel := gocv.GetStructuringElement(
		gocv.MorphShape(cxt.MorphElem),
		image.Pt(2*cxt.MorphSize+1, 2*cxt.MorphSize+1))
	gocv.MorphologyEx(*src, dstp, gocv.MorphType(operation), kernel)
	return dstp
}
