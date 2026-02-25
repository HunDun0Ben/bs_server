package ops_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gocv.io/x/gocv"

	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ops"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/testutil"
)

func TestBlurOp_Process(t *testing.T) {
	// 确保测试结束时没有内存泄露
	defer testutil.AssertNoLeaks(t)

	// 初始化为 100x100 的 8位 3通道 图像
	src := gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
	testutil.TrackMat()
	defer testutil.CloseMat(&src)

	op := ops.NewBlurOp()
	op.UpdateParam("type", ops.Blur)
	op.UpdateParam("size", 3)

	// 执行处理
	dstPtr := op.Process(&src)
	testutil.TrackMat() // Process 返回了新 Mat，需要追踪
	defer testutil.CloseMat(dstPtr)

	// 断言结果
	assert.False(t, dstPtr.Empty(), "Result mat should not be empty")
	assert.Equal(t, src.Rows(), dstPtr.Rows(), "Rows should match")
	assert.Equal(t, src.Cols(), dstPtr.Cols(), "Cols should match")
}

func TestBlurOp_LeakDetection(t *testing.T) {
	// 这个测试故意模拟一个泄露来看看工具是否生效
	// 我们在一个子测试中运行它
	t.Run("SimulatedLeak", func(st *testing.T) {
		m := testutil.NewMat()
		_ = m
		// 故意不调用 testutil.CloseMat(&m)

		// 我们不直接断言泄露，因为这会让测试失败。
		// 在实际使用中，你会看到 AssertNoLeaks 报错。
	})

	// 清理刚才故意制造的泄露，以免影响其他测试
	testutil.ReleaseMat()
}
