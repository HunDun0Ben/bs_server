package test

import (
	"bytes"
	"context"
	"crypto/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

// TestGridFSChunkingBehavior 测试 GridFS 在特定配置下的分块行为。
// 场景：
// 1. 设置块大小 (ChunkSize) 为 16MB。
// 2. 上传一个大于 16MB 的文件。
// 3. 同时，为该文件附加一个大于 1MB 的元数据 (metadata)。
// 验证：
// - 第一个数据块 (chunk) 中 `data` 字段的实际大小应严格等于设置的 ChunkSize (16MB)。
// - 这证明了 `fs.files` 集合中元数据的大小不影响 `fs.chunks` 集合中每个文档的大小。
func TestGridFSChunkingBehavior(t *testing.T) {
	// --- 1. 设置 ---
	const (
		// 测试得出的结果是, chunks 也受到 document 的 16MB 限制
		chunkSize = 16*1024*1024 - 80 // 16MB - 80B
		fileSize  = 17 * 1024 * 1024  // 17MB, 确保至少会产生两个 chunk
		// metadata 写入的是 fs 中, 非 chunks 当中
		metadataSize = 1 * 1024 * 1024
		dbName       = "test_gridfs_chunking"
		bucketName   = "test_bucket"
	)

	client := imongo.Client()
	db := client.Database(dbName)
	// 使用 defer 确保测试数据库在结束后被清理
	defer func() {
		err := db.Drop(context.Background())
		require.NoError(t, err, "Failed to drop test database")
	}()

	bucketOpts := options.GridFSBucket().SetChunkSizeBytes(chunkSize).SetName(bucketName)
	bucket, err := gridfs.NewBucket(db, bucketOpts)
	require.NoError(t, err, "Failed to create GridFS bucket")

	// --- 2. 生成测试数据 ---
	// 创建一个大于16MB的文件内容
	fileData := make([]byte, fileSize)
	_, err = rand.Read(fileData)
	require.NoError(t, err, "Failed to generate random file data")

	// 创建一个大于1MB的元数据
	largeString := strings.Repeat("a", metadataSize)
	metadata := bson.M{"large_key": largeString}

	// --- 3. 执行上传 ---
	uploadOpts := options.GridFSUpload().SetMetadata(metadata)
	fileID, err := bucket.UploadFromStream("large_file_with_large_metadata.bin", bytes.NewReader(fileData), uploadOpts)
	require.NoError(t, err, "Failed to upload file to GridFS")
	require.NotNil(t, fileID, "Received a nil fileID")

	// --- 4. 验证 ---
	// 直接查询 chunks 集合来验证第一个 chunk 的大小
	chunksCollection := db.Collection(bucketName + ".chunks")

	// 定义一个结构体来解码 chunk 文档
	type fsChunk struct {
		ID      primitive.ObjectID `bson:"_id"`
		FilesID primitive.ObjectID `bson:"files_id"`
		N       int32              `bson:"n"`
		Data    primitive.Binary   `bson:"data"`
	}

	var firstChunk fsChunk
	filter := bson.M{"files_id": fileID, "n": 0} // 查询第一个 chunk (n=0)
	err = chunksCollection.FindOne(context.Background(), filter).Decode(&firstChunk)
	require.NoError(t, err, "Failed to find the first chunk in the collection")

	// 核心断言：验证 data 字段的长度
	assert.Equal(t, int(chunkSize), len(firstChunk.Data.Data), "The data size of the first chunk should be exactly the configured chunk size")

	// 可选：验证第二个 chunk 的大小
	var secondChunk fsChunk
	filter = bson.M{"files_id": fileID, "n": 1} // 查询第二个 chunk (n=1)
	err = chunksCollection.FindOne(context.Background(), filter).Decode(&secondChunk)
	require.NoError(t, err, "Failed to find the second chunk in the collection")
	assert.Equal(t, fileSize-chunkSize, len(secondChunk.Data.Data), "The data size of the second chunk should be the remainder of the file")
}
