package imgutils

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
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

func Mat2DBMat(mat *gocv.Mat) (*imongo.DBMat, error) {
	if mat.Empty() {
		return nil, nil
	}
	data := mat.ToBytes()
	return &imongo.DBMat{
		Context: data,
		Col:     mat.Cols(),
		Row:     mat.Rows(),
		MatType: int(mat.Type()),
	}, nil
}

func PrintMat(mat gocv.Mat) {
	rows, cols := mat.Rows(), mat.Cols()
	fmt.Printf("Mat Size: %dx%d, Type: %d\n", rows, cols, mat.Type())

	switch mat.Type() {
	case gocv.MatTypeCV8U:
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				fmt.Printf("%3d ", mat.GetUCharAt(i, j))
			}
			fmt.Println()
		}
	case gocv.MatTypeCV32F:
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				fmt.Printf("%.2f ", mat.GetFloatAt(i, j))
			}
			fmt.Println()
		}
	case gocv.MatTypeCV32S:
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				fmt.Printf("%d ", mat.GetIntAt(i, j))
			}
			fmt.Println()
		}
	default:
		fmt.Println("不支持的 Mat 类型打印")
	}
}
