package ops

import (
	"image"

	"gocv.io/x/gocv"
)

type MorphologyOp struct {
	OpType      int
	ElementType int
	KernelSize  int
	params      map[string]int
}

func NewMorphologyOp() *MorphologyOp {
	return &MorphologyOp{
		OpType:      0,
		ElementType: 0,
		KernelSize:  1,
		params:      make(map[string]int),
	}
}

func (m *MorphologyOp) GetName() string {
	return "Morphology"
}

func (m *MorphologyOp) Process(src *gocv.Mat) *gocv.Mat {
	dst := gocv.NewMat()
	kernel := gocv.GetStructuringElement(
		gocv.MorphShape(m.ElementType),
		image.Pt(2*m.KernelSize+1, 2*m.KernelSize+1),
	)

	switch m.OpType {
	case 0:
		gocv.Erode(*src, &dst, kernel)
	case 1:
		gocv.Dilate(*src, &dst, kernel)
	case 2:
		gocv.MorphologyEx(*src, &dst, gocv.MorphOpen, kernel)
	case 3:
		gocv.MorphologyEx(*src, &dst, gocv.MorphClose, kernel)
	}

	return &dst
}

func (m *MorphologyOp) UpdateParam(name string, value int) error {
	m.params[name] = value
	switch name {
	case "type":
		m.OpType = value
	case "element":
		m.ElementType = value
	case "kernel":
		m.KernelSize = value
	}
	return nil
}

func (m *MorphologyOp) GetParams() map[string]int {
	return m.params
}
