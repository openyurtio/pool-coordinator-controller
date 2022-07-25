package fs

import (
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func EnsureDir(dir string) error {
	if DirExists(dir) {
		return nil
	}
	return os.MkdirAll(dir, 0750)
}

func WriteFile(fn string, data []byte) error {
	return os.WriteFile(fn, data, 0660)
}
