package main

import (
	"fmt"
)

/*
#include <stdio.h>
*/
import "C"

func main() {

	fmt.Println(C.CString("hello world\n"))
}
