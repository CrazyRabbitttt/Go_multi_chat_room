package main

import (
	"fmt"
	"net"
)

type Client1 struct {
	ServerIp   string
	ServerProt int
	Name       string
	conn       net.Conn
}

func newClient1(Ip string, Port int) *Client1 {

	client1 := &Client1{
		ServerIp:   Ip,
		ServerProt: Port,
	}
	//连接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Ip, Port))

	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client1.conn = conn
	return client1
}

func main() {
	client := newClient1("127.0.0.1", 8888)

	if client == nil {
		fmt.Println("=====链接服务器失败=====")
		return
	}

	fmt.Println("=====链接服务器成功=====")

	//暂时进行阻塞
	select {}
}
