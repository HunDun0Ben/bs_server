package imgpro

import (
	"demo/gocv/imgpro"
	"fmt"
	"log"
	"testing"
	"time"

	"gocv.io/x/gocv"
)

func TestPerformanceBtCloneAndNew(t *testing.T) {
	path := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	mat := gocv.IMRead(path, gocv.IMReadColor)
	var count = 1000
	a := time.Now()
	for i := 0; i < count; i++ {
		mat.Clone()
	}
	fmt.Printf("耗时 %.3fs\n", time.Since(a).Seconds())
	a = time.Now()
	for i := 0; i < count; i++ {
		imgpro.NewSomeMat(mat)
	}
	fmt.Printf("耗时 %.3fs\n", time.Since(a).Seconds())
}

func TestWindowDisplayMat(t *testing.T) {

	filename := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	wrapper := imgpro.NewTrackWindowWrapper("Hello", filename)
	wrapper.SetContext(new(imgpro.TrackWindowCxt))
	wrapper.CreateTrackbar("Blur Type:", 4, 1,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowCxt); ok {
				testCxt.BlurType = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.CreateTrackbar("kernel size:", 20, 1,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowCxt); ok {
				testCxt.Ksize = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.Dispaly()
}
