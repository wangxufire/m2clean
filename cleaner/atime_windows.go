// +build windows

package cleaner

import (
	"os"
	"syscall"
)

func atime(info os.FileInfo) int64 {
	if stat, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		return stat.LastAccessTime.Nanoseconds() / 1e9
	}
	return 0
}
