package main

import (
	"fmt"

	"github.com/sony/sonyflake"
)

func main() {
	var sf *sonyflake.Sonyflake

	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
	fmt.Println(sonyflake.Decompose(id))

}
