package main

import (
	"crypto/rc4"
	"encoding/json"
	"fmt"
	"net"
)

type ConfigMsgType int

const (
	ConfigSearch       ConfigMsgType = 1001
	ConfigSingleConfig ConfigMsgType = 1002
	ConfigBatchConfig  ConfigMsgType = 1003
	ConfigMachineInfo  ConfigMsgType = 2001
)

var key = []byte("Bxxx")

func main() {

	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	})
	if err != nil {
		fmt.Println("连接服务端失败，err:", err)
		return
	}
	defer socket.Close()
	send(socket, RespnseMsg{
		Type:        ConfigSearch,
		MachineCode: "123",
		Version:     "123",
	})
	for {
		msg := read(socket)
		if msg == nil {
			continue
		}
		switch msg.Type {
		case ConfigSearch:
			send(socket, RespnseMsg{
				Type:        ConfigMachineInfo,
				MachineCode: "123",
				Version:     "123",
			})
		}
	}
}

type Msg struct {
	Type ConfigMsgType `json:"type"`
}

func read(conn *net.UDPConn) *Msg {
	data := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	dest2 := make([]byte, 1024)
	cipher2, _ := rc4.NewCipher(key) // 切记：这里不能重用cipher1，必须重新生成新的
	cipher2.XORKeyStream(dest2, data)
	fmt.Println(string(dest2[:n]))
	var msg = &Msg{}
	if err := json.Unmarshal(dest2[:n], msg); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return msg
}

type RespnseMsg struct {
	Type        ConfigMsgType `json:"type"`
	MachineCode string        `json:"machineCode"`
	Version     string        `json:"version"`
}

func send(conn *net.UDPConn, msg RespnseMsg) {
	destAddr := &net.UDPAddr{
		IP:   net.ParseIP("192.168.3.255"),
		Port: 51234,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	cryData := make([]byte, len(data))

	c, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	c.XORKeyStream(cryData, data)
	_, err = conn.WriteToUDP(cryData, destAddr) // 发送数据
	if err != nil {
		fmt.Println("发送数据失败，err:", err)
		return
	}
}
