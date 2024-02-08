package main

import (
	"demo/gocv/imgpro"
	"log"
)

func main() {
	filename := `/home/workspace/data/leedsbutterfly/images/0010001.png`
	wrapper := imgpro.NewTrackWindowWrapper("Hello", filename)
	wrapper.SetContext(new(imgpro.TrackWindowBlurCxt))
	wrapper.CreateTrackbar("Blur Type:", 4, 1,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowBlurCxt); ok {
				testCxt.BlurType = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.CreateTrackbar("kernel size:", 20, 1,
		func(cxt *imgpro.GenMatCxt, pos int) error {
			if testCxt, ok := (*cxt).(*imgpro.TrackWindowBlurCxt); ok {
				testCxt.Ksize = pos
			} else {
				log.Print("强转类型失败")
			}
			return nil
		})
	wrapper.Dispaly()
}
