package main

import (
	"crypto/rc4"
	"fmt"
	"net"
)

func main() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4allsys,
		Port: 51234,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}
	defer socket.Close()
	sendData := []byte("Hello server")
	key := []byte("xx")
	dest1 := make([]byte, len(sendData))

	c, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	c.XORKeyStream(dest1, sendData)
	fmt.Println("SEND ", dest1)
	n, err := socket.Write(dest1) // 发送数据
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
	fmt.Println(n)
	for {
		data := make([]byte, 1024)

		n, _, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(n)

		dest2 := make([]byte, 1024)
		cipher2, _ := rc4.NewCipher(key) // 切记：这里不能重用cipher1，必须重新生成新的
		cipher2.XORKeyStream(dest2, data)
		fmt.Println(string(dest2))
		break
	}
}
