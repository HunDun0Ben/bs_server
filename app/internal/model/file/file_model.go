package file

import (
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

type ButterflyFile struct {
	imongo.FileStoreData `bson:"inline"`
	MaskContent          []byte `bson:"mask_content,omitempty"`
}

type ResizedButteryflyFile struct {
	imongo.FileStoreData `bson:"inline"`
	Col                  int
	Row                  int
	Type                 string       `bson:"type,omitempty"`
	DescribMat           imongo.DBMat `bson:"describ_mat,omitempty"`
}

func NewButterflyFile(fileName, typeName, path string, content, maskContent []byte) *ButterflyFile {
	return &ButterflyFile{
		MaskContent: maskContent,
		FileStoreData: imongo.FileStoreData{
			FileName: fileName,
			TypeName: typeName,
			Path:     path,
			Content:  content,
		},
	}
}
