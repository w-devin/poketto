//go:build !windows

package file

func CopyFileUsedByOtherProcess(srcPath, dstPath string) error {
	// placeholder
	return CopyFile(srcPath, dstPath)
}
