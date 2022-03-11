package bufio

import (
	"bufio"
	"fmt"
	"strings"
)

func ExampleReadSlice() {
	reader := bufio.NewReader(strings.NewReader("https://golang.google.cn\nIt is the website of golang"))
	line, _ := reader.ReadSlice('\n')
	fmt.Println(string(line))

	line, _ = reader.ReadSlice('\n')
	fmt.Println(string(line))

	// output:
	// https://golang.google.cn
	//
	// It is the website of golang
}

func ExampleReadBytes() {
	reader := bufio.NewReader(strings.NewReader("https://golang.google.cn\nIt is the website of golang"))
	line, _ := reader.ReadBytes('\n')
	fmt.Println(string(line))

	line, _ = reader.ReadBytes('\n')
	fmt.Println(string(line))

	// output:
	// https://golang.google.cn
	//
	// It is the website of golang
}
