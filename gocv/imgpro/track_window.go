package imgpro

import (
	"demo/gocv/imgpro/constant"
	"image"

	"gocv.io/x/gocv"
)

type GenMatCxt interface {
	process(*gocv.Mat) *gocv.Mat
}

type TrackWindowCxt struct {
	BlurType int
	Ksize    int
}

func (cxt TrackWindowCxt) process(src *gocv.Mat) *gocv.Mat {
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

type TrackWindowWrapper struct {
	*gocv.Window
	Path           string
	SrcImgMat      gocv.Mat
	DstImgMat      gocv.Mat
	context        GenMatCxt
	barArray       []*gocv.Trackbar
	oldBarPosArray []int
	onChangeArray  []func(cxt *GenMatCxt, pos int) error
}

func NewTrackWindowWrapper(title, path string) *TrackWindowWrapper {
	win := gocv.NewWindow(title)
	win.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowAutosize)
	return &TrackWindowWrapper{Window: win, Path: path}
}

func (wrap *TrackWindowWrapper) SetContext(cxt GenMatCxt) {
	wrap.context = cxt
}

func (wrap *TrackWindowWrapper) CreateTrackbar(name string, max, posd int, onChange func(cxt *GenMatCxt, pos int) error) {
	bar := wrap.Window.CreateTrackbarWithValue(name, &posd, max)
	wrap.barArray = append(wrap.barArray, bar)
	wrap.onChangeArray = append(wrap.onChangeArray, onChange)
	wrap.oldBarPosArray = append(wrap.oldBarPosArray, 0)
}

func (wrap *TrackWindowWrapper) Dispaly() {
	wrap.SrcImgMat = gocv.IMRead(wrap.Path, gocv.IMReadColor)
	wrap.DstImgMat = wrap.SrcImgMat
	for {
		var update = false
		if len(wrap.barArray) > 0 {
			for i, v := range wrap.barArray {
				if v.GetPos() != wrap.oldBarPosArray[i] {
					wrap.onChangeArray[i](&wrap.context, v.GetPos())
					wrap.oldBarPosArray[i] = v.GetPos()
					update = true
				}
			}
			if update {
				wrap.DstImgMat = *(wrap.context.process(&wrap.SrcImgMat))
			}
		}
		wrap.IMShow(wrap.DstImgMat)
		if wrap.WaitKey(1) >= 0 {
			break
		}
	}
}
