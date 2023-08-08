package main

import (
	"github.com/bin-work/go-example/pkg/bshell"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			bshell.ShellExec("scrot")
		}
	}

}
