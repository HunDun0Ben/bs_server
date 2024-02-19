package practice_test

import (
	"demo/gocv/imgpro"
	"log"
	"testing"
)

func TestMorphologyExWindow(t *testing.T) {
	filename := `/home/workspace/data/leedsbutterfly/images/0010046.png`
	windows_name := "Morphology Transformations Demo"
	wrapper := imgpro.NewTrackWindowWrapper(windows_name, filename)
	wrapper.SetContext(new(imgpro.TrackWindowMorphologyExCxt))
	const op_trackbar_name = "Operator:\n 0: Opening - 1: Closing \n 2: Gradient - 3: Top Hat \n 4: Black Hat"
	const ele_trackbar_name = "Element:\n 0: Rect - 1: Cross - 2: Ellipse"
	const kernel_size_trackbar_name = "Kernel size:\n 2n +1"
	wrapper.CreateTrackbar(op_trackbar_name, imgpro.MorphMaxOperator, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if morphCxt, ok := (*cxt).(*imgpro.TrackWindowMorphologyExCxt); ok {
				morphCxt.Operator = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.CreateTrackbar(ele_trackbar_name, imgpro.MorphMaxElem, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if morphCxt, ok := (*cxt).(*imgpro.TrackWindowMorphologyExCxt); ok {
				morphCxt.MorphElem = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.CreateTrackbar(kernel_size_trackbar_name, imgpro.MorphMaxKernelSize, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if morphCxt, ok := (*cxt).(*imgpro.TrackWindowMorphologyExCxt); ok {
				morphCxt.MorphSize = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.Dispaly()
}
