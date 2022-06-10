package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

func main() {
	salt := "beyondinfo"
	pwd := "Qq@123456"

	m := md5.New()
	_, err := io.WriteString(m, pwd+salt)
	if err != nil {
		panic(err)
	}
	arr := m.Sum(nil) //已经输出，但是是编码
	// 将编码转换为字符串
	newArr := fmt.Sprintf("%x", arr)
	fmt.Println(newArr)
	//输出字符串字母都是小写，转换为大写
	//sig = strings.ToTitle(newArr)

}
