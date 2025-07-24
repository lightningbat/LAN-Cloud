//go:build windows
// +build windows

package fsid

import (
	"fmt"
	"syscall"
)

func GetID(path string) (string, error) {
	pPath, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}

	handle, err := syscall.CreateFile(pPath,
		0, // no access to the file itself
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS, // allow opening dirs
		0)
	if err != nil {
		return "", err
	}
	defer syscall.CloseHandle(handle)

	var info syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(handle, &info)
	if err != nil {
		return "", err
	}

	fileIndex := uint64(info.FileIndexHigh)<<32 | uint64(info.FileIndexLow)
	return fmt.Sprintf("%d-%d", info.VolumeSerialNumber, fileIndex), nil
}