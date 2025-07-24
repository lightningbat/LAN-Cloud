//go:build linux || darwin
// +build linux darwin

package fsid

import (
	"fmt"
	"os"
	"syscall"
)

func GetID(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	stat := info.Sys().(*syscall.Stat_t)
	// Inode + device ID
	return fmt.Sprintf("%d-%d", stat.Dev, stat.Ino), nil
}