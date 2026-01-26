package origin

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	// TUNSETIFF is the ioctl command to set the network interface parameters.
	// The value 0x400454ca is derived from _IOW('T', 202, int) on x86_64 Linux.
	// _IOW(type, nr, size) = (write | (size << 16) | (type << 8) | nr)
	// 0x40000000 | (4 << 16) | (0x54 << 8) | 202
	TUNSETIFF = 0x400454ca

	// IFF_TUN indicates a TUN device (IP packets).
	IFF_TUN = 0x0001
	// IFF_TAP indicates a TAP device (Ethernet frames).
	IFF_TAP = 0x0002
	// IFF_NO_PI indicates no packet information header.
	IFF_NO_PI = 0x1000
)

// ifReq is the interface request structure used for ioctl calls.
// It matches the C struct ifreq.
type ifReq struct {
	Name  [16]byte
	Flags uint16
	_     [22]byte // Padding to match the size of struct ifreq (40 bytes on 64-bit Linux)
}

// Open opens a TUN device with the specified name.
// If name is empty, the kernel will pick the next available name (e.g., tun0).
func Open(name string) (*os.File, error) {
	// 1. Open the clone device
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open /dev/net/tun: %w", err)
	}

	// 2. Prepare the interface request
	var ifr ifReq
	copy(ifr.Name[:], name)
	ifr.Flags = IFF_TUN | IFF_NO_PI

	// 3. Perform the ioctl syscall to register the device
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		file.Fd(),
		uintptr(TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr)),
	)

	if errno != 0 {
		file.Close()
		return nil, fmt.Errorf("ioctl TUNSETIFF failed: %w", errno)
	}

	// The actual name might be different if we passed an empty string
	// actualName := string(bytes.Trim(ifr.Name[:], "\x00"))
	// fmt.Printf("TUN device created: %s\n", actualName)

	return file, nil
}
