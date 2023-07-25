package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero, //代表本机所有网卡地址
		Port: 51234,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer listen.Close()

	for {
		var buf [1024]byte
		n, addr, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			log.Println(err)
			break
		}
		if n > 0 {
			fmt.Println("udp receive msg:", string(buf[:n]), "地址：", addr.IP.String(), ":", addr.Port)
			log.Println(string(buf[:n]))
		}
		n, err = listen.WriteTo(buf[:n], addr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("write succes ", n)
	}

}
