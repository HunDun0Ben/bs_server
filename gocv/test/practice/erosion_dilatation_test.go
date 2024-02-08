package practice_test

import (
	"demo/gocv/imgpro"
	"log"
	"testing"
)

const (
	max_elem        = 2
	max_kernel_size = 21
)

func Test(t *testing.T) {
	filename := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	erosionWrapper := imgpro.NewTrackWindowWrapper("Erosion Demo", filename)
	erosionWrapper.SetContext(new(imgpro.TrackWindowErosionCxt))
	erosionWrapper.CreateTrackbar("Element:\n 0: Rect \n 1: Cross \n 2: Ellipse", max_elem, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowErosionCxt); ok {
				testCxt.ErosionElem = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	erosionWrapper.CreateTrackbar("Kernel size:\n 2n +1", max_kernel_size, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowErosionCxt); ok {
				testCxt.ErosionSize = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	go func() { erosionWrapper.Dispaly() }()

	dilatationWrapper := imgpro.NewTrackWindowWrapper("ilatation Demo", filename)
	dilatationWrapper.SetContext(new(imgpro.TrackWindowDilatationCxt))
	dilatationWrapper.CreateTrackbar("Element:\n 0: Rect \n 1: Cross \n 2: Ellipse", max_elem, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowDilatationCxt); ok {
				testCxt.DilationElem = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	dilatationWrapper.CreateTrackbar("Kernel size:\n 2n +1", max_kernel_size, 0,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowDilatationCxt); ok {
				testCxt.DilationSize = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	dilatationWrapper.Dispaly()

}
