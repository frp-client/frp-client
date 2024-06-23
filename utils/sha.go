package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func Sha256(data string) string {
	//创建一个基于SHA256算法的hash.Hash接口的对象
	//hash := sha256.New() //sha-256加密
	hash := sha256.New() //SHA-512加密
	//输入数据
	hash.Write([]byte(data))
	//计算哈希值
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(hash.Sum(nil))
	//返回哈希值
	return hashCode
}

func Sha512(data string) string {
	//创建一个基于SHA256算法的hash.Hash接口的对象
	//hash := sha256.New() //sha-256加密
	hash := sha512.New() //SHA-512加密
	//输入数据
	hash.Write([]byte(data))
	//计算哈希值
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(hash.Sum(nil))
	//返回哈希值
	return hashCode
}

func Md5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
