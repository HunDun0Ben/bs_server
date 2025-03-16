package core

import "gocv.io/x/gocv"

// ImageProcessor 定义基本的图像处理接口
type ImageProcessor interface {
	Process(src *gocv.Mat) *gocv.Mat
	GetName() string
}

// ProcessorContext 定义处理器上下文接口
type ProcessorContext interface {
	ImageProcessor
	UpdateParam(name string, value int) error
	GetParams() map[string]int
}
