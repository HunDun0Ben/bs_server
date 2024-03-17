package imgpro

import (
	"log"

	"gocv.io/x/gocv"
)

type GenMatCxt interface {
	process(*gocv.Mat) *gocv.Mat
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
	loadImg        bool
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

func (wrap *TrackWindowWrapper) LoadImg() {
	wrap.SrcImgMat = gocv.IMRead(wrap.Path, gocv.IMReadColor)
	wrap.DstImgMat = wrap.SrcImgMat
	wrap.DstImgMat = wrap.SrcImgMat
	wrap.loadImg = true
}

func (wrap *TrackWindowWrapper) Dispaly() {
	if !wrap.loadImg {
		wrap.LoadImg()
	}
	for {
		var update = false
		if len(wrap.barArray) > 0 {
			if wrap.context == nil {
				log.Fatal("Context can't be nil")
				return
			}
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
