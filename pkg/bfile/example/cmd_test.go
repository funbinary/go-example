package example

import (
	"beyondinfo.com/baselib/go/base_package.git/bfile"
	"fmt"
)

func ExamplePwd() {
	fmt.Println(bfile.Pwd())
	// output:
	// D:\workspace\kds\base_package\bfile\example
}

func ExampleChdir() {
	fmt.Println(bfile.Pwd())
	bfile.Chdir("/")
	fmt.Println(bfile.Pwd())
	// output:
	// D:\workspace\kds\base_package\bfile\example
	// D:\
}

//func Example(){
//	// output:
//	//
//}
//func Example(){
//	// output:
//	//
//}
//func Example(){
//	// output:
//	//
//}
//func Example(){
//	// output:
//	//
//}
