package utils

import (
	"io/ioutil"
	"os"
)

// 根据文件名读取文件
func Readfile(name string) []byte {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}

// 判断文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
