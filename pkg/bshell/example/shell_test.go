package example

import (
	"beyondinfo.com/baselib/go/base_package.git/bshell"
	"fmt"
)

func ExampleShell() {
	result, err := bshell.ShellExec(`echo 1`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// output:
	//
}
