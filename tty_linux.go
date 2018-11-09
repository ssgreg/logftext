// +build linux
// +build !appengine

package logftext

import (
	"syscall"
	"unsafe"
)

func enableSeqTTY(fd uintptr, flag bool) error {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if errno != 0 {
		return errno
	}

	return nil
}
