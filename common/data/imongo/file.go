package imongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type FileStoreData struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	FileName string             `bson:"file_name,omitempty"`
	TypeName string             `bson:"type_name,omitempty"`
	Path     string             `bson:"path,omitempty"`
	Content  []byte             `bson:"content,omitempty"`
}

func NewFileStore(fileName, typeName, path string, content []byte) *FileStoreData {
	return &FileStoreData{FileName: fileName, TypeName: typeName, Path: path, Content: content}
}
