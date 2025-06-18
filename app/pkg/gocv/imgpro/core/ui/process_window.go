package ui

import (
	"fmt"

	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"
)

type ProcessingWindow struct {
	window    *gocv.Window
	src       *gocv.Mat
	dst       *gocv.Mat
	processor core.ProcessorContext
	trackbars []*Trackbar
}

func NewProcessingWindow(title string) *ProcessingWindow {
	win := gocv.NewWindow(title)
	win.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowAutosize)

	return &ProcessingWindow{
		window:    win,
		trackbars: make([]*Trackbar, 0),
	}
}

func (w *ProcessingWindow) Process(fun func(src *gocv.Mat) *gocv.Mat) {
	w.dst = fun(w.dst)
}

func (w *ProcessingWindow) GetDstMat() *gocv.Mat {
	return w.dst
}

func (w *ProcessingWindow) SetProcessor(p core.ProcessorContext) {
	w.processor = p
}

func (w *ProcessingWindow) LoadImageFromPath(path string) error {
	src := gocv.IMRead(path, gocv.IMReadColor)
	if src.Empty() {
		return fmt.Errorf("failed to load image: %s", path)
	}
	w.src = &src
	w.dst = imgutils.NewSomeMat(*w.src)
	src.CopyTo(w.dst)
	return nil
}

// LoadImageFromMat loads an image from a gocv.Mat and initializes the src and dst fields.
func (w *ProcessingWindow) LoadImageFromMat(img gocv.Mat) {
	if img.Empty() {
		return
	}
	w.src = &img
	w.dst = imgutils.NewSomeMat(*w.src)
	img.CopyTo(w.dst)
}

func (w *ProcessingWindow) Display() {
	for {
		if w.updateTrackbars() {
			w.dst = w.processor.Process(w.src)
		}
		w.window.IMShow(*w.dst)
		if w.window.WaitKey(100) >= 0 {
			break
		}
	}
}

func (w *ProcessingWindow) Close() {
	if w.src != nil {
		w.src.Close()
	}
	if w.dst != nil {
		w.dst.Close()
	}
	w.window.Close()
}
