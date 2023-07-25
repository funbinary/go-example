package main

import (
	"fmt"
	"regexp"
)

func main() {
	str := "USB\\VID_0930&PID_6545\\001A4D5F1A5CED41E0000D7D"
	re := regexp.MustCompile(`USB\\VID_(?P<VID>\w+)&PID_(?P<PID>\w+)\\(?P<s>\w+)`)
	match := re.FindStringSubmatch(str)
	for _, m := range match {
		fmt.Println(m)
	}
}
