package main

import (
	"encoding/binary"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		os.Exit(1)
	}

	slog.Info("tun 相关信息",
		"path", f.Name(),
		"fd", f.Fd(),
		"mode", fi.Mode().String(),
		"is_device", fi.Mode()&os.ModeDevice != 0,
		"is_char_device", fi.Mode()&os.ModeCharDevice != 0,
		"size", fi.Size(),
	)
	name := "tun0"
	nameb := make([]byte, 16)
	copy(nameb, name)

	flag := make([]byte, 2)
	binary.LittleEndian.PutUint16(flag, unix.IFF_TUN|unix.IFF_NO_PI)

	ifreq := make([]byte, 18)
	copy(ifreq[0:16], nameb)
	copy(ifreq[16:], flag)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, f.Fd(), uintptr(unix.TUNSETIFF), uintptr((unsafe.Pointer(&ifreq[0]))))
	if errno != 0 {
		log.Fatalf("ioctl 创建 tun 失败: %v", errno)
	}

	log.Println("tun0 创建成功")
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, syscall.SIGINT, syscall.SIGTERM) // 捕获 Ctrl+C 或终止信号

	// 保持程序运行, 并且保持 fd 不被关闭.
	<-sch
	log.Println("收到终止信号，正在清理资源...")
	os.Exit(1)
}
