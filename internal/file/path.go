package file

import (
	"path/filepath"
	"strings"
)

// Ext
//  @Description: 获取给定路径的扩展名，包含.
//  @param path 路径
//  @return string给定路径的扩展名，包含.
//
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// ExtName
//  @Description: 获取给定路径的扩展名，不包含.
//  @param path 路径
//  @return string 给定路径的扩展名
//
func ExtName(path string) string {
	return strings.TrimLeft(Ext(path), ".")
}
