package initbutterflyimg

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	mcli "github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
)

var (
	basePath      = `/home/workspace/data/leedsbutterfly`
	images        = "images"
	segmentations = "segmentations"
	imgsPath      = filepath.Join(basePath, images)
	segPath       = filepath.Join(basePath, segmentations)
)

func DisplayImg() {
	collection := mcli.FileDatabase().Collection("butterfly_img")
	bf := new(file.ButterflyFile)
	bf.Path = "/home/workspace/data/leedsbutterfly/images/0010001.png"

	err := collection.FindOne(context.TODO(), bf).Decode(bf)
	fmt.Printf("\tbf = %s, %s\n", bf.FileName, bf.Path)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}

	win := ui.NewProcessingWindow("Hello")
	img, _ := gocv.IMDecode(bf.Content, gocv.IMReadColor)
	win.LoadImageFromMat(img)
	win.Display()
}

func InsertImg() {
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
			segPath := filepath.Join(segPath, segFileName)
			content, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			maskContent, err := os.ReadFile(segPath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("name:%s, path:%s \n", info.Name(), path)
			file := file.NewButterflyFile(info.Name(), ext, path, content, maskContent)
			fmt.Printf("file:%s, %s\n", file.FileName, file.Path)
			insertResult, err := collection.InsertOne(context.Background(), file)
			if err != nil {
				log.Fatal(err)
			}

			// 输出插入的文档的 ObjectID
			fmt.Println("Inserted document ID:", insertResult.InsertedID)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func VerifyImgsAndSeg() {
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
			segPath := filepath.Join(segPath, segFileName)
			_, err = os.Stat(segPath)
			if os.IsNotExist(err) {
				count++
				fmt.Printf("File Seg not exist: %s", segPath)
			}
		}
		return nil
	})
	if count == 0 {
		fmt.Println("All images is correctly")
	}
	if err != nil {
		fmt.Println(err)
	}
}
