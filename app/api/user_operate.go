package api

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/HunDun0Ben/bs_server/common/data/imongo"
	"github.com/gin-gonic/gin"
)

func UploadImg(cxt *gin.Context) {
	file, header, err := cxt.Request.FormFile("file")
	if err != nil {
		panic("")
	}
	defer file.Close()
	log.Println(header.Filename)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		panic("")
	}
	imongo.StoreFile(
		cxt,
		"updateImg",
		imongo.FileStoreData{FileName: header.Filename, Content: fileContent},
	)
	cxt.JSON(http.StatusOK, fmt.Sprintf("'%s' uploaded!", header.Filename))
}

func GetImgResult(cxt *gin.Context) {
	var req GetImgResultReq
	if err := cxt.ShouldBindJSON(&req); err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileData, err := imongo.GetFile(cxt, "updateImg", req.ImgID)

	if fileData != nil {
		cxt.JSON(http.StatusOK, gin.H{"GetImgResultReq": req})
	} else {
		cxt.JSON(http.StatusBadRequest, err)
	}
}

func InsectInfo(cxt *gin.Context) {
}
