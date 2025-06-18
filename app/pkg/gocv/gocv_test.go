package gocv_test

import (
	"testing"

	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ops"
	"github.com/HunDun0Ben/bs_server/app/pkg/gocv/imgpro/core/ui"
)

var filename = `/home/workspace/data/leedsbutterfly/images/0010001.png`

func TestMorphologyOp(t *testing.T) {
	// 创建窗口
	window := ui.NewProcessingWindow("Morphology Demo")
	defer window.Close()

	// 创建处理器
	processor := ops.NewMorphologyOp()

	window.SetProcessor(processor)

	// 添加控制条
	window.AddTrackbar("Operation Type", 3, func(v int) {
		processor.UpdateParam("type", v)
	})
	window.AddTrackbar("Element Type", 2, func(v int) {
		processor.UpdateParam("element", v)
	})
	window.AddTrackbar("Kernel Size", 21, func(v int) {
		processor.UpdateParam("kernel", v)
	})

	// 加载图像
	if err := window.LoadImageFromPath(filename); err != nil {
		panic(err)
	}

	// 显示并处理
	window.Display()
}

func TestBlurWindow(t *testing.T) {
	// 创建窗口
	window := ui.NewProcessingWindow("Blur Demo")
	defer window.Close()

	// 创建模糊处理器
	processor := ops.NewBlurOp()
	window.SetProcessor(processor)

	// 添加控制条
	window.AddTrackbar("Blur Type", 3, func(v int) {
		processor.UpdateParam("type", v)
	})
	window.AddTrackbar("Kernel Size", 21, func(v int) {
		processor.UpdateParam("size", v)
	})

	// 加载图像
	if err := window.LoadImageFromPath(filename); err != nil {
		panic(err)
	}

	// 显示并处理
	window.Display()
}
