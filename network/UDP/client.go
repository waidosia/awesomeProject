package main

import (
	"fmt"
	"net"
)

func main() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("服务器连接失败,Err : ", err)
		return
	}
	defer socket.Close()
	sendData := []byte("Hello server")
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("数据发送失败,err:", err)
		return
	}
	data := make([]byte, 4096)
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("数据接受失败,err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data), remoteAddr, n)
}
