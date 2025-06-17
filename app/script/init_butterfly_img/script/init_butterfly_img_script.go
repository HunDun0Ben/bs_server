package main

import (
	"log/slog"
	"os"
)

func main() {
	// initbutterflyimg.InsertImg()
	// initbutterflyimg.DisplayImg()
	// initbutterflyimg.VerifyImgsAndSeg()
	// 获取所有环境变量
	envVars := os.Environ()

	// 打印所有环境变量
	slog.Info("Environment Variables:", slog.Any("env_vars", envVars))
}
