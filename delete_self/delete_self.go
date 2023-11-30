//go:build !windows
// +build !windows

package delete_self

import (
	"fmt"
	"os"
)

func DeleteSelfDuringRunning() error {
	path, err := os.Executable()
	if err != nil {
		return fmt.Errorf("filed to get current filename, %v\n", err)
	}

	_, err = os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return nil
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete self, %v", err)
	}

	_, err = os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return nil
	}
	return fmt.Errorf("failed to delete self")
}
