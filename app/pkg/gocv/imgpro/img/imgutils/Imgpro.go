package imgutils

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func NewSomeMat(src gocv.Mat) *gocv.Mat {
	a := gocv.NewMatWithSizes(src.Size(), src.Type())
	return &a
}

func ResizeWithPadding(src gocv.Mat, targetWidth, targetHeight int, padColor color.RGBA) gocv.Mat {
	srcWidth := src.Cols()
	srcHeight := src.Rows()

	// 计算等比例缩放后的新尺寸
	scale := math.Min(float64(targetWidth)/float64(srcWidth), float64(targetHeight)/float64(srcHeight))
	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	// 缩放图像
	resized := gocv.NewMat()
	gocv.Resize(src, &resized, image.Pt(newWidth, newHeight), 0, 0, gocv.InterpolationLinear)

	// 创建目标图像（填充背景色）
	output := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(float64(padColor.B), float64(padColor.G), float64(padColor.R), 0), targetHeight, targetWidth, src.Type())

	// 计算居中偏移
	xOffset := (targetWidth - newWidth) / 2
	yOffset := (targetHeight - newHeight) / 2

	// 拷贝缩放图像到中心位置
	roi := output.Region(image.Rect(xOffset, yOffset, xOffset+newWidth, yOffset+newHeight))
	resized.CopyTo(&roi)
	roi.Close()
	resized.Close()

	return output
}
