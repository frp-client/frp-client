package utils

import (
	"os"
	"path/filepath"
)

func AppPath() string {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return p
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	if stat.IsDir() {
		return false
	}
	return true
}
