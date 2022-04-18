package main

import (
	"flag"
	"fmt"
	"net"
)

type Client1 struct {
	ServerIp   string
	ServerProt int
	Name       string
	conn       net.Conn
	flag       int //用于菜单的选择
}

func newClient1(Ip string, Port int) *Client1 {

	client1 := &Client1{
		ServerIp:   Ip,
		ServerProt: Port,
		flag:       999,
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

var serverIp string
var serverPort int

//提供一个menu的方法
func (this *Client1) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("=====请输入合法范围内的数字=====")
		return false
	}

}

func (this *Client1) Run() {
	for this.flag != 0 {
		for this.menu() != true {
		} //没有选择就一直挂着
		switch this.flag {
		case 1:
			fmt.Println("公聊模式选择。。。")
		case 2:
			fmt.Println("私聊模式选择。。。")
		case 3:
			fmt.Println("更新用户名选择。。。")
		}
	}
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置Server的Ip，default:127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置Server的Port，default:8888")
}

func main() {

	flag.Parse()

	client := newClient1(serverIp, serverPort)

	if client == nil {
		fmt.Println("=====链接服务器失败=====")
		return
	}

	fmt.Println("=====链接服务器成功=====")

	client.Run()

	//暂时进行阻塞
	select {}
}
