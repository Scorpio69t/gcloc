package utils

import "os"

// FileExits returns true if the file exists.
func FileExits(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
