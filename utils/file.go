package utils

import (
	"io"
	"os"
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
