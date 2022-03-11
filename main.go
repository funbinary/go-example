package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	fmt.Println(net.JoinHostPort("127.0.0.1", strconv.Itoa(21)))
}
