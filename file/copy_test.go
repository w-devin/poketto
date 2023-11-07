package file

import (
	"fmt"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	srcPath := "C:\\Users\\jarvis\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Network\\Cookies"
	dstPath := "Cookies"

	_, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("failed to open file, %v\n", err)
	}

	err = CopyFile(srcPath, dstPath)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Printf("file not used by chrome, %v", err)
		return
	}
	fmt.Println("file use by some process, continue...")

	err = CopyFileUsedByOtherProcess(srcPath, dstPath)
	if err != nil {
		fmt.Printf("failed to copy %s to %s, %v", srcPath, dstPath, err)
	}
	fmt.Println("succeed")
}
