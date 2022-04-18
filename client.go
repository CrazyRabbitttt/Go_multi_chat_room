package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

func (this *Client1) SelectUsers() {
	//查询当前在线的用户
	sendMsg := "who\n"

	_, err := this.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (this *Client1) PrivateChat() {
	var remoteName string
	var chatMsg string

	this.SelectUsers()

	fmt.Println("====请输入聊天对象[用户名],exit退出：")

	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("====请输入消息内容，exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不空就发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + "\n\n"
				_, err := this.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("====请输入消息内容，exit退出:")
			fmt.Scanln(&chatMsg)
		}
		this.SelectUsers()
		fmt.Println("====请输入消息内容，exit退出:")
		fmt.Scanln(&chatMsg)

	}

}

func (this *Client1) PublicChat() {
	//提示用户输入信息
	var chatMsg string

	fmt.Println("====请输入聊天内容，exit输出")

	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		sendMsg := chatMsg + "\n"

		_, err := this.conn.Write([]byte(sendMsg))

		if err != nil {
			fmt.Println("conn Write err:", err)
			break
		}

		chatMsg = ""
		fmt.Println("====请输入聊天内容，exit输出")

		fmt.Scanln(&chatMsg)
	}
}

//提供用户名更新的功能
func (this *Client1) UpdateName() bool {
	fmt.Println("====请输入用户名：")
	fmt.Scanln(&this.Name)

	//将我们既定好的格式传输过去

	sendMessage := "rename|" + this.Name + "\n"

	_, err := this.conn.Write([]byte(sendMessage))

	if err != nil {
		fmt.Println("conn.Write error :", err)
		return false
	}

	return true
}

//处理Server回送的消息，输出到标准输出
func (this *Client1) DealResponse() {
	io.Copy(os.Stdout, this.conn)
}

func (this *Client1) Run() {
	for this.flag != 0 {
		for this.menu() != true {
		} //没有选择就一直挂着
		switch this.flag {
		case 1:
			this.PublicChat()
		case 2:
			this.PrivateChat()
		case 3:
			this.UpdateName()
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

	go client.DealResponse()

	client.Run()

	//暂时进行阻塞
	select {}
}
