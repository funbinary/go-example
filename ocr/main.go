package main

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("/mnt/c/Users/xx/Desktop/importBroder.png")
	client.SetLanguage("chi_sim", "eng")
	text, _ := client.Text()
	fmt.Println(text)
	// Hello, World!
}
