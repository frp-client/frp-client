package utils

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func ReadFile(fileName string) (r []byte) {
	f, err := os.Open(fileName)
	if err != nil {
		return make([]byte, 0)
	}
	r, err = io.ReadAll(f)
	if err != nil {
		return make([]byte, 0)
	}
	return r
}

func AppTempFile(elem ...string) string {
	elem = append([]string{os.TempDir(), "frp-client"}, elem...)
	var tmpFile = filepath.Join(elem...)
	_ = os.MkdirAll(filepath.Dir(tmpFile), fs.ModePerm)
	return tmpFile
}

func ReadFileAsByte(filename string) []byte {
	buf, _ := os.ReadFile(filename)
	return buf
}

func ReadFileAsString(filename string) string {
	buf, _ := os.ReadFile(filename)
	return string(buf)
}

func SaveFileAsByte(filename string, data []byte) error {
	return os.WriteFile(filename, data, fs.ModePerm)
}

func SaveFileAsString(filename string, data string) error {
	return os.WriteFile(filename, []byte(data), fs.ModePerm)
}
