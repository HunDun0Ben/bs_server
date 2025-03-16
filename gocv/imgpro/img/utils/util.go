package utils

import (
	"image/color"
	"math/rand"
	"time"
)

var colorR = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandColor() *color.RGBA {

	r := uint8(colorR.Intn(255))
	g := uint8(colorR.Intn(255))
	b := uint8(colorR.Intn(255))

	return &color.RGBA{r, g, b, 0}
}
