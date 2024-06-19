package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func AesEncrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainBytes := []byte(plaintext)
	// 对于CBC模式，需要使用PKCS#7填充plaintext到blocksize的整数倍
	plainBytes = AesPad(plainBytes, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBytes)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AesDecrypt(key []byte, ct string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, err
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	data, err = AesUnpad(data, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// AesPad 使用PKCS#7标准填充数据
func AesPad(buf []byte, blockSize int) []byte {
	padding := blockSize - (len(buf) % blockSize)
	return append(buf, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// AesUnpad 移除PKCS#7标准填充的数据
func AesUnpad(buf []byte, blockSize int) ([]byte, error) {
	if len(buf)%blockSize != 0 {
		return nil, nil
	}
	padding := int(buf[len(buf)-1])
	if len(buf)-padding < 0 {
		return nil, errors.New("aes padding error")
	}
	return buf[:len(buf)-padding], nil
}
