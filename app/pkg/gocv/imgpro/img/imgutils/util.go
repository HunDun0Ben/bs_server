package imgutils

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strconv"
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

func SaveMatToCSV(mat gocv.Mat, filename string) {
	if mat.Channels() != 1 {
		fmt.Println("仅支持单通道 Mat")
		return
	}

	rows := mat.Rows()
	cols := mat.Cols()

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("创建文件失败：", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for r := 0; r < rows; r++ {
		var row []string
		for c := 0; c < cols; c++ {
			val := mat.GetFloatAt(r, c)
			row = append(row, strconv.FormatFloat(float64(val), 'f', 8, 64))
		}
		writer.Write(row)
	}
}
