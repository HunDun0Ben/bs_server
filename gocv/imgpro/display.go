package imgpro

import (
	"fmt"

	"gocv.io/x/gocv"
)

func Display(title, path string) {
	window := gocv.NewWindow(title)
	bar := window.CreateTrackbar("blur", 4)
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("Error reading image from: %v\n", path)
		return
	}
	for {
		window.IMShow(img)
		fmt.Printf("bar position = %d", bar.GetPos())
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
