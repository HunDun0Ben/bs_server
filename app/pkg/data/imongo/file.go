package imongo

import "go.mongodb.org/mongo-driver/bson/primitive"

// FileStoreData supports both embedded content (legacy) and GridFS references.
type FileStoreData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	FileName string             `bson:"file_name,omitempty"`
	TypeName string             `bson:"type_name,omitempty"`
	Path     string             `bson:"path,omitempty"`
	// For new data using GridFS
	FileID primitive.ObjectID `bson:"file_id,omitempty"`
	// For legacy data with embedded content
	Content []byte `bson:"content,omitempty"`
}

type ButterflyFile struct {
	FileStoreData
	Name string `bson:"type,omitempty"` // TypeName is the type of butterfly, e.g., "leedsbutterfly"
}

// UserFile represents a file uploaded by a user.
type UserFile struct {
	FileStoreData
	UserID   primitive.ObjectID `bson:"user_id,omitempty"`
	UserName string             `bson:"user_name,omitempty"`
}

// NewFileStoreWithContent creates a legacy FileStoreData with embedded content.
func NewFileStoreWithContent(fileName, typeName, path string, content []byte) *FileStoreData {
	return &FileStoreData{FileName: fileName, TypeName: typeName, Path: path, Content: content}
}

// NewFileStoreWithGridFS creates a new FileStoreData with a GridFS reference.
func NewFileStoreWithGridFS(fileID primitive.ObjectID, fileName, typeName, path string) *FileStoreData {
	return &FileStoreData{FileID: fileID, FileName: fileName, TypeName: typeName, Path: path}
}

// NewUserFileWithContent creates a legacy UserFile with embedded content.
func NewUserFileWithContent(fileName, typeName, path string, content []byte, userID primitive.ObjectID, userName string) *UserFile {
	return &UserFile{
		FileStoreData: FileStoreData{
			FileName: fileName,
			TypeName: typeName,
			Path:     path,
			Content:  content,
		},
		UserID:   userID,
		UserName: userName,
	}
}

// NewUserFileWithGridFS creates a new UserFile with a GridFS reference.
func NewUserFileWithGridFS(fileID primitive.ObjectID, fileName, typeName, path string, userID primitive.ObjectID, userName string) *UserFile {
	return &UserFile{
		FileStoreData: FileStoreData{
			FileID:   fileID,
			FileName: fileName,
			TypeName: typeName,
			Path:     path,
		},
		UserID:   userID,
		UserName: userName,
	}
}
