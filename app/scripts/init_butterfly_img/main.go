package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/internal/repository"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo/imongoutil"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/img/imgutils"

	mcli "github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

var (
	basePath   string
	mode       string
	k          int
	iterations int
)

func init() {
	pflag.StringVarP(&basePath, "path", "p", "/home/workspace/data/leedsbutterfly", "Base path for butterfly data")
	pflag.StringVarP(&mode, "mode", "m", "verify", "Operation mode: insert, verify, display, shift, kmeans")
	pflag.IntVarP(&k, "clusters", "k", 1024, "Number of clusters for KMeans")
	pflag.IntVarP(&iterations, "iterations", "i", 10, "KMeans iterations")
}

func main() {
	pflag.Parse()

	if err := conf.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	repo := repository.NewButterflyRepository(mcli.BizDataBase())
	svc := butterflysvc.NewButterflyService(repo)

	switch mode {
	case "insert":
		InsertImg()
	case "verify":
		VerifyImgsAndSeg()
	case "display":
		DisplayImg(svc)
	case "shift":
		UpdateShiftFeature(svc)
	case "kmeans":
		Kmeans(svc)
	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		pflag.Usage()
	}
}

func DisplayImg(svc butterflysvc.ButterflyService) {
	path := filepath.Join(basePath, "images/0010001.png")
	bf, err := svc.FindImg(context.TODO(), bson.M{"path": path})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Document not found")
			return
		}
		log.Fatalf("Find error: %v", err)
	}
	fmt.Printf("\tbf = %s, %s\n", bf.FileName, bf.Path)

	win := ui.NewProcessingWindow("Butterfly Image")
	img, err := gocv.IMDecode(bf.Content, gocv.IMReadColor)
	if err != nil {
		log.Fatalf("Decode error: %v", err)
	}
	defer img.Close()
	win.LoadImageFromMat(img)
	win.Display()
}

func InsertImg() {
	imgsPath := filepath.Join(basePath, "images")
	segPath := filepath.Join(basePath, "segmentations")
	collection := mcli.FileDatabase().Collection("butterfly_img")

	err := filepath.WalkDir(imgsPath, func(path string, d fs.DirEntry, err error) error {
		segSuf := "_seg0"
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				fmt.Println(err)
				return nil
			}
			ext := filepath.Ext(info.Name())
			nameWithoutExt := strings.TrimSuffix(info.Name(), ext)
			segFileName := nameWithoutExt + segSuf + ext
			fullSegPath := filepath.Join(segPath, segFileName)

			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading image %s: %v", path, err)
				return nil
			}
			maskContent, err := os.ReadFile(fullSegPath)
			if err != nil {
				log.Printf("Error reading mask %s: %v", fullSegPath, err)
				return nil
			}

			fmt.Printf("Inserting: %s\n", info.Name())
			bf := file.NewButterflyFileWithContent(info.Name(), ext, path, content, maskContent)
			insertResult, err := collection.InsertOne(context.Background(), bf)
			if err != nil {
				log.Fatalf("Insert error: %v", err)
			}
			fmt.Println("Inserted document ID:", insertResult.InsertedID)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func VerifyImgsAndSeg() {
	imgsPath := filepath.Join(basePath, "images")
	segPath := filepath.Join(basePath, "segmentations")
	var count int
	err := filepath.WalkDir(imgsPath, func(_ string, d fs.DirEntry, err error) error {
		segSuf := "_seg0"
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				fmt.Println(err)
				return nil
			}
			ext := filepath.Ext(info.Name())
			nameWithoutExt := strings.TrimSuffix(info.Name(), ext)
			segFileName := nameWithoutExt + segSuf + ext
			fullSegPath := filepath.Join(segPath, segFileName)
			_, err = os.Stat(fullSegPath)
			if os.IsNotExist(err) {
				count++
				fmt.Printf("File Seg not exist: %s\n", fullSegPath)
			}
		}
		return nil
	})
	if count == 0 {
		fmt.Println("All images and segmentations are verified")
	} else {
		fmt.Printf("Found %d missing segmentations\n", count)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func UpdateShiftFeature(svc butterflysvc.ButterflyService) {
	sift := gocv.NewSIFT()
	defer sift.Close()
	resizeList, err := svc.GetResizedImgs(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to get resized images: %v", err)
	}

	for _, item := range resizeList {
		func() {
			mat, err := gocv.NewMatFromBytes(200, 200, gocv.MatTypeCV8UC3, item.Content)
			if err != nil {
				log.Printf("Error creating mat for item %s: %v", item.ID, err)
				return
			}
			defer mat.Close()

			dst := gocv.NewMat()
			defer dst.Close()
			gocv.CvtColor(mat, &dst, gocv.ColorBGRToGray)

			_, describe := sift.DetectAndCompute(dst, gocv.NewMat())
			defer describe.Close()

			if describe.Empty() {
				log.Printf("No features found for item %s", item.ID)
				return
			}

			dbmat, err := imgutils.Mat2DBMat(&describe)
			if err != nil {
				log.Printf("Error converting mat to DBMat for item %s: %v", item.ID, err)
				return
			}

			err = svc.UpdateResizedImg(context.Background(), bson.M{"_id": item.ID}, bson.M{
				"$set": bson.M{
					"describ_mat": dbmat,
				},
			})
			if err != nil {
				log.Printf("Error updating item %s: %v", item.ID, err)
			}
		}()
	}
	fmt.Println("SIFT features update completed")
}

func Kmeans(svc butterflysvc.ButterflyService) {
	buf := make([]byte, 0, 60*1024*1024)
	var rows int
	resizeList, err := svc.GetResizedImgs(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to get resized images: %v", err)
	}

	// shift 特征的大小通常是128, 拼接所有特征作为 kmeans 所需要的数据集合
	for _, item := range resizeList {
		buf = append(buf, item.DescribMat.Context...)
		rows += item.DescribMat.Row
	}

	if rows == 0 {
		log.Fatal("No features found to run KMeans")
	}

	// 初始化所有样本的特征点矩阵
	allDescrib, err := gocv.NewMatFromBytes(rows, 128, gocv.MatTypeCV32FC1, buf)
	if err != nil {
		log.Fatalf("Failed to create mat from bytes: %v", err)
	}
	defer allDescrib.Close()

	slog.Info("Mat allDescrib Info", "Rows", allDescrib.Rows(), "Cols", allDescrib.Cols(), "Type", allDescrib.Type())

	labels := gocv.NewMat()
	defer labels.Close()
	centers := gocv.NewMat()
	defer centers.Close()

	criteria := gocv.NewTermCriteria(gocv.Count|gocv.EPS, iterations, 1.0)

	res := gocv.KMeans(allDescrib, k, &labels, criteria, 3, gocv.KMeansPPCenters, &centers)
	slog.Info("KMeans completed", "compactness", res)

	storeKmeans(&labels, &centers)

	var start int
	trainingData := gocv.NewMat()
	defer trainingData.Close()

	for _, item := range resizeList {
		imgTag, _ := strconv.Atoi(item.Type)
		bow := buildBowHistogram(&labels, k, start, item.DescribMat.Row, imgTag)
		defer bow.Close()

		start += item.DescribMat.Row
		if trainingData.Cols() == 0 {
			trainingData = bow.Clone()
		} else {
			gocv.Vconcat(trainingData, bow, &trainingData)
		}
	}

	// 对 trainingData 的每一行进行 L2 归一化
	for i := 0; i < trainingData.Rows(); i++ {
		row := trainingData.RowRange(i, i+1)
		featureRow := row.ColRange(0, k)
		gocv.Normalize(featureRow, &featureRow, 1.0, 0.0, gocv.NormL2)
		row.Close()
		featureRow.Close()
	}

	slog.Info("Training data generated", "Rows", trainingData.Rows(), "Cols", trainingData.Cols())
	imgutils.SaveMatToCSV(trainingData, "data.csv")
	fmt.Println("Training data saved to data.csv")
}

func buildBowHistogram(labels *gocv.Mat, k, start, num, tag int) gocv.Mat {
	bowi := make([]int, k)
	bow := gocv.NewMatWithSize(1, k+1, gocv.MatTypeCV32FC1)
	for i := start; i < start+num; i++ {
		clusterIdx := labels.GetIntAt(i, 0)
		if int(clusterIdx) >= 0 && int(clusterIdx) < k {
			bowi[clusterIdx]++
		}
	}
	for idx, value := range bowi {
		bow.SetFloatAt(0, idx, float32(value))
	}
	bow.SetFloatAt(0, k, float32(tag))
	return bow
}

func storeKmeans(labels, centers *gocv.Mat) {
	mat1, _ := imgutils.Mat2DBMat(labels)
	mat2, _ := imgutils.Mat2DBMat(centers)
	obj := bson.M{"labels": mat1, "centers": mat2}
	imongoutil.Insert[bson.M](context.Background(), mcli.FileDatabase().Collection("kmeans"), obj)
}
