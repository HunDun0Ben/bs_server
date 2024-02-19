package imgpro

import (
	"image"

	"gocv.io/x/gocv"
)

type TrackWindowErosionCxt struct {
	MorphologyExElem int
	MorphologyExSize int
}

type TrackWindowDilatationCxt struct {
	DilationElem int
	DilationSize int
}

func (cxt *TrackWindowErosionCxt) process(src *gocv.Mat) *gocv.Mat {
	dstp := NewSomeMat(*src)
	if cxt.MorphologyExElem < 0 || cxt.MorphologyExElem > 2 {
		cxt.MorphologyExElem = 0
	}
	kernel := gocv.GetStructuringElement(
		gocv.MorphShape(cxt.MorphologyExElem),
		image.Pt(2*cxt.MorphologyExSize+1, 2*cxt.MorphologyExSize+1))
	gocv.Erode(*src, dstp, kernel)
	return dstp
}

func (cxt *TrackWindowDilatationCxt) process(src *gocv.Mat) *gocv.Mat {
	dstp := NewSomeMat(*src)
	if cxt.DilationElem < 0 || cxt.DilationElem > 2 {
		cxt.DilationElem = 0
	}
	kernel := gocv.GetStructuringElement(
		gocv.MorphShape(cxt.DilationElem),
		image.Pt(2*cxt.DilationSize+1, 2*cxt.DilationSize+1))
	gocv.Dilate(*src, dstp, kernel)
	return dstp
}
