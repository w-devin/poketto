//go:build windows

package file

import (
	"fmt"
	"github.com/kbinani/win"
	windigo "github.com/rodrigocfd/windigo/win"
	process "github.com/w-devin/poketto/windows"
	"golang.org/x/sys/windows"
	"strings"
)

var (
	windowsApi = process.WindowsApi{}
)

func CopyFIleUsedByOtherProcess(srcPath, dstPath string) error {
	// 找到目标文件进程及句柄
	pid, fileHandlerNumber, err := FindProcessAndFileHandlerByFileName(srcPath)
	if err != nil {
		fmt.Printf("failed to found process and file handler of %s, %v", srcPath, err)
		return fmt.Errorf("failed to found process and file handler of %s, %v", srcPath, err)
	}
	fmt.Printf("found process, pid: %+v, file: %+v\n", pid, fileHandlerNumber)

	// 复制目标文件句柄
	sourceProcessHandle, err := windowsApi.OpenProcess(windows.PROCESS_DUP_HANDLE, false, pid)
	if err != nil {
		return fmt.Errorf("failed to open process %d, %v", pid, err)
	}
	defer windows.CloseHandle(sourceProcessHandle)

	var duplicatedHandle windows.Handle
	err = windowsApi.DuplicateHandle(sourceProcessHandle, fileHandlerNumber, windows.CurrentProcess(), &duplicatedHandle, 0, false, windows.DUPLICATE_SAME_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to duplicateHandle of source file")
	}

	// 创建文件映射
	mappedPtr, err := windows.CreateFileMapping(duplicatedHandle, nil, windows.PAGE_READONLY, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("failed to create file mapping of source file, %v", err)
	}

	// map the entire file into memory
	mappedViewPtr, err := windows.MapViewOfFile(mappedPtr, windows.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to amp the entire file into memory")
	}

	// copy file
	dstFilePtr, err := windows.CreateFile(windows.StringToUTF16Ptr(dstPath), windows.GENERIC_READ|windows.GENERIC_WRITE, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.CREATE_ALWAYS, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to create file of dstPath: %s, %v", dstPath, err)
	}

	duplicatedFile := windigo.HFILE(duplicatedHandle)
	fileSize, err := duplicatedFile.GetFileSizeEx()
	fmt.Printf("fileSize: %d\n", fileSize)

	var written uint32
	if win.WriteFile(win.HANDLE(dstFilePtr), mappedViewPtr, win.DWORD(fileSize), &written, nil) {
		return nil
	}
	return fmt.Errorf("failed write content to %s", dstPath)
}

func FindProcessAndFileHandlerByFileName(srcPath string) (pid uint32, fileHandler windows.Handle, err error) {
	currentProcessHandle := windows.CurrentProcess()
	defer windows.CloseHandle(currentProcessHandle)

	systemHandlers, err := windowsApi.QuerySystemHandleInformation()
	for _, handler := range systemHandlers {
		sourceProcessHandle, err := windowsApi.OpenProcess(windows.PROCESS_DUP_HANDLE, false, handler.ProcessId)
		if err != nil {
			continue
		}
		defer windows.CloseHandle(sourceProcessHandle)

		err = windowsApi.DuplicateHandle(sourceProcessHandle, windows.Handle(handler.HandleValue), currentProcessHandle, &fileHandler, 0, false, windows.DUPLICATE_SAME_ACCESS)
		if err != nil {
			continue
		}
		fileName, err := windowsApi.GetObjectName(fileHandler)
		if err != nil {
			continue
		}

		if strings.HasSuffix(fileName, srcPath[2:]) {
			return handler.ProcessId, windows.Handle(handler.HandleValue), nil
		}
	}

	return 0, windows.Handle(0), fmt.Errorf("target not found")
}
