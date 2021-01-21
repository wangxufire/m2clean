// +build darwin

package cleaner

import (
	"os"
	"syscall"
)

func atime(info os.FileInfo) int64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Atimespec.Sec
	}
	return 0
}
