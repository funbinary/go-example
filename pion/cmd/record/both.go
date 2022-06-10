package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bin-work/go-example/pion/cmd/record/mkv"

	"github.com/bin-work/go-example/pion/cmd/record/opus"
	"github.com/bin-work/go-example/pkg/bfile"
)

var opusstream *opus.RtcOgg
var mkvstream *mkv.Stream

func Notify(message chan os.Signal) {

	signal.Notify(message,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGILL,
		syscall.SIGTRAP,
		syscall.SIGABRT,
		syscall.SIGBUS,
		syscall.SIGFPE,
		syscall.SIGKILL,
		syscall.SIGSEGV,
		syscall.SIGPIPE,
		syscall.SIGALRM,
		syscall.SIGTERM, os.Interrupt)
	sig := <-message
	fmt.Println(sig.Signal)
	fmt.Println(sig.String())
	fmt.Println("接受到关闭信号，即将关闭", sig, sig.Signal, sig.String())
	opusstream.Stop()
	mkvstream.Stop()

	time.Sleep(2 * time.Second)

}

func main() {

	host := "192.168.3.249"
	room := "627b5e1a38feb98851a0185c"
	display := "5c718a004d90"
	//display2 := "4eb10b31e915"

	var err error
	savepath1 := bfile.SelfDir() + "/" + room + "_" + display
	mkvstream, err = mkv.NewStream(host, "1985", room, display, savepath1)

	//savepath2 := bfile.SelfDir() + "/" + room + "_" + display2 + ".ogg"
	//opusstream, err = opus.NewOgg(host, room, display2, savepath2)
	fmt.Println(err)
	message := make(chan os.Signal, 1)
	Notify(message)

}
