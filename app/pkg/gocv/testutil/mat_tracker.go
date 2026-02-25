package testutil

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"
)

var matCount int64

// TrackMat 增加存活 Mat 的计数。
// 在测试代码中，每当创建一个新的 Mat 时调用此函数。
func TrackMat() {
	atomic.AddInt64(&matCount, 1)
}

// ReleaseMat 减少存活 Mat 的计数。
// 在测试代码中，每当调用 Mat.Close() 时调用此函数。
func ReleaseMat() {
	atomic.AddInt64(&matCount, -1)
}

// AssertNoLeaks 断言没有未关闭的 Mat。
func AssertNoLeaks(t *testing.T) {
	count := atomic.LoadInt64(&matCount)
	assert.Equal(t, int64(0), count, "Detected %d unclosed gocv.Mat(s)", count)
}

// NewMat 包装了 gocv.NewMat，并自动进行追踪。
func NewMat() gocv.Mat {
	TrackMat()
	return gocv.NewMat()
}

// CloseMat 包装了 mat.Close()，并自动进行追踪。
func CloseMat(mat *gocv.Mat) {
	if mat != nil {
		mat.Close()
		ReleaseMat()
	}
}
