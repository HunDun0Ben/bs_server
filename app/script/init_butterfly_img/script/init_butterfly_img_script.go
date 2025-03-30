package main

import (
	"log/slog"
	"os"

	initbutterflyimg "github.com/HunDun0Ben/bs_server/app/script/init_butterfly_img"
)

func main() {
	// initbutterflyimg.InsertImg()
	initbutterflyimg.DisplayImg()
	// initbutterflyimg.VerifyImgsAndSeg()
	// 获取所有环境变量
	envVars := os.Environ()

	// 打印所有环境变量
	for _, envVar := range envVars {
		// 输出格式：KEY=VALUE
		slog.Info(envVar)
	}
}
