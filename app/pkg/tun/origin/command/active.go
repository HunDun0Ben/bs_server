package main

import (
	"log"
	"os/exec"
)

func main() {
	// 激活接口
	cmd := exec.Command("ip", "link", "set", "dev", "tun0", "up")
	if err := cmd.Run(); err != nil {
		log.Fatalf("激活 tun0 失败: %v", err)
	}

	// 分配 IP
	cmd = exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", "tun0")
	if err := cmd.Run(); err != nil {
		log.Fatalf("给 tun0 分配 IP 失败: %v", err)
	}

	log.Println("tun0 已激活并分配 IP")
}
