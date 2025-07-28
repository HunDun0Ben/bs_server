package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/internal/service/butterflysvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo/imongoutil"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

func UploadImg(cxt *gin.Context) {
	file, header, err := cxt.Request.FormFile("file")
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的文件", nil, err))
		return
	}
	defer file.Close()
	slog.Info("Uploading file", "filename", header.Filename)
	fileContent, err := io.ReadAll(file)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "读取文件失败", nil, err))
		return
	}
	fileID, err := imongoutil.StoreFile(
		cxt,
		"updateImg",
		imongo.FileStoreData{FileName: header.Filename, Content: fileContent},
	)
	if err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "存储文件失败", nil, err))
		return
	}
	helper.Success(cxt, gin.H{
		"fileId":   fileID,
		"fileName": header.Filename,
	})
}

func GetImgResult(cxt *gin.Context) {
	var req dto.GetImgResultReq
	if err := cxt.ShouldBindQuery(&req); err != nil {
		cxt.Error(bsvo.NewAppError(http.StatusBadRequest, "无效的请求参数", nil, err))
		return
	}

	_, err := imongoutil.GetFile(cxt, "updateImg", req.ImgID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			cxt.Error(bsvo.NewAppError(http.StatusNotFound, "图片未找到", nil, err))
		} else {
			cxt.Error(bsvo.NewAppError(http.StatusInternalServerError, "获取图片结果失败", nil, err))
		}
		return
	}
	// TODO: 此处应添加实际的图片处理逻辑并返回结果。
	helper.Success(cxt, gin.H{
		"message": "图片结果尚未就绪。",
		"imgId":   req.ImgID,
	})
}

func InsectInfo(cxt *gin.Context) {
}

func ButterflyInfo(cxt *gin.Context) {
	insect_list, err := butterflysvc.NewButterflyTypeSvc().GetAllList(cxt.Request.Context())
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching butterfly info"})
		return
	}
	cxt.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   insect_list,
	})
}
