package file

import "demo/common/data"

type ButterflyFile struct {
	data.FileStoreData `bson:"inline"`
	MaskContent        []byte `bson:"mask_content,omitempty"`
}

func NewButterflyFile(fileName, typeName, path string, content, maskContent []byte) *ButterflyFile {
	return &ButterflyFile{MaskContent: maskContent,
		FileStoreData: data.FileStoreData{FileName: fileName, TypeName: typeName, Path: path, Content: content},
	}
}
