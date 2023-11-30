//go:build windows
// +build windows

package delete_self

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"unsafe"
)

const StreamRename = ":tes"

type FileRenameInfo struct {
	ReplaceIfExists bool
	RootDirectory   windows.Handle
	FileNameLength  uint32
	FileName        [2]uint16
}

func openHandle(filePath string) (windows.Handle, error) {
	return windows.CreateFile(
		windows.StringToUTF16Ptr(filePath),
		windows.DELETE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
}

func depositHandle(handle windows.Handle) error {
	yes := byte(1)
	return windows.SetFileInformationByHandle(
		handle,
		windows.FileDispositionInfo,
		&yes,
		1,
	)
}
func DeleteSelfDuringRunning() error {
	path, err := os.Executable()
	if err != nil {
		return fmt.Errorf("filed to get current filename, %v\n", err)
	}

	_, err = os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return nil
	}

	handle, err := openHandle(path)
	if err != nil {
		return fmt.Errorf("filed to open hanle of %s, %v\n", path, err)
	}

	ri := &FileRenameInfo{
		ReplaceIfExists: false,
		RootDirectory:   windows.Handle(0),
		FileNameLength:  uint32(len(StreamRename)),
	}
	copy(ri.FileName[:], windows.StringToUTF16(StreamRename))

	err = windows.SetFileInformationByHandle(
		handle,
		windows.FileRenameInfo,
		(*byte)(unsafe.Pointer(ri)),
		uint32(unsafe.Sizeof(*ri))+ri.FileNameLength,
	)
	if err != nil {
		return fmt.Errorf("failed to rename file, %v\n", err)
	}

	err = windows.CloseHandle(handle)
	if err != nil {
		return fmt.Errorf("failed to close handle, %v\n", err)
	}

	handle, err = openHandle(path)
	if err != nil {
		return fmt.Errorf("filed to open another hanle of %s, %v\n", path, err)
	}

	err = depositHandle(handle)
	if err != nil {
		return fmt.Errorf("failed to delete file, %v", err)
	}

	err = windows.CloseHandle(handle)
	if err != nil {
		return fmt.Errorf("failed to close new handle, %v\n", err)
	}

	_, err = os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return nil
	}
	return fmt.Errorf("failed to delete file")
}
