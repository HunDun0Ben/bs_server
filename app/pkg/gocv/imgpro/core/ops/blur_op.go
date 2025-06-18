package ops

import (
	"image"

	"gocv.io/x/gocv"
)

const (
	Blur = iota
	GaussianBlur
	MedianBlur
	BilateralFilter
)

type BlurOp struct {
	blurType int
	size     int
	params   map[string]int
}

func NewBlurOp() *BlurOp {
	return &BlurOp{
		blurType: 0,
		size:     1,
		params:   make(map[string]int),
	}
}

func (b *BlurOp) GetName() string {
	return "Blur"
}

func (b *BlurOp) Process(src *gocv.Mat) *gocv.Mat {
	dst := gocv.NewMat()
	size := b.size
	if size < 1 {
		size = 1
	}

	switch b.blurType {
	case Blur:
		gocv.Blur(*src, &dst, image.Point{X: size, Y: size})
	case GaussianBlur:
		if size%2 == 0 {
			size++
		}
		gocv.GaussianBlur(*src, &dst, image.Point{X: size, Y: size}, 0, 0, 0)
	case MedianBlur:
		if size%2 == 0 {
			size++
		}
		gocv.MedianBlur(*src, &dst, size)
	case BilateralFilter:
		gocv.BilateralFilter(*src, &dst, size, float64(size*2), float64(size/2))
	}
	return &dst
}

func (b *BlurOp) UpdateParam(name string, value int) error {
	b.params[name] = value
	switch name {
	case "type":
		b.blurType = value
	case "size":
		b.size = value
	}
	return nil
}

func (b *BlurOp) GetParams() map[string]int {
	return b.params
}
