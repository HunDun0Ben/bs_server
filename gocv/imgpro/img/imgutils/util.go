package imgutils

import (
	"image/color"
	"math/rand"
	"time"

	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/entities/file"
)

var colorR = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandColor() *color.RGBA {
	r := uint8(colorR.Intn(255))
	g := uint8(colorR.Intn(255))
	b := uint8(colorR.Intn(255))
	return &color.RGBA{r, g, b, 0}
}

func GetMaskImg(fileInfo file.ButterflyFile) *gocv.Mat {
	img, _ := gocv.IMDecode(fileInfo.Content, gocv.IMReadColor)
	mask, _ := gocv.IMDecode(fileInfo.MaskContent, gocv.IMReadGrayScale)
	maskimg := gocv.NewMat()
	// 通过 mask 去除无关部分.
	gocv.BitwiseAndWithMask(img, img, &maskimg, mask)
	return &maskimg
}
