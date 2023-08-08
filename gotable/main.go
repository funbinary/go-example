package main

import (
	"fmt"
	"github.com/liushuochen/gotable"
)

func main() {
	table, err := gotable.Create("版本号", "包体长度", "业务码", "序列ID")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return
	}

	table.AddRow([]string{"4字节(uint32)", "4字节(uint32)", "4字节(uint32)", "8字节(uint64)"})

	fmt.Println(table)
}
