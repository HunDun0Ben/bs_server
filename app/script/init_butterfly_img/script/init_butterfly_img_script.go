package main

import (
	initbutterflyimg "demo/app/script/init_butterfly_img"
	"fmt"
	"os"
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
		fmt.Println(envVar)
	}
}
