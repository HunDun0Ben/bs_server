package file

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

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

// NewButterflyFileWithContent creates a legacy ButterflyFile with embedded content.
func NewButterflyFileWithContent(fileName, typeName, path string, content, maskContent []byte) *ButterflyFile {
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

// NewButterflyFileWithGridFS creates a new ButterflyFile with a GridFS reference.
func NewButterflyFileWithGridFS(fileID primitive.ObjectID, fileName, typeName, path string, maskContent []byte) *ButterflyFile {
	return &ButterflyFile{
		MaskContent: maskContent,
		FileStoreData: imongo.FileStoreData{
			FileID:   fileID,
			FileName: fileName,
			TypeName: typeName,
			Path:     path,
		},
	}
}
