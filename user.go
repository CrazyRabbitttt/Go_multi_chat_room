package main

import "net"

type Client struct {
	Addr   string
	Name   string
	C      chan string //用于进行通信的channel
	conn   net.Conn
	server *Server //当前Client属于哪一个Server
}

//通过一个conn进行Client端的创建
func newClient(conn net.Conn, server *Server) *Client { //创建一个Client对象

	useradd := conn.RemoteAddr().String()

	client := &Client{
		Addr:   useradd,
		Name:   useradd,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//监听当前的client channel 的消息的goroutne

	go client.ListenMessage()

	return client
}

//clinet就需要无限的进行监听自己的channel， 如果有消息的话就读取出来

func (this *Client) ListenMessage() {
	for {
		message := <-this.C
		this.conn.Write([]byte(message + "\n")) //能够写到客户端的标准输入中去
	}
}

func (this *Client) Online() {
	this.server.mapLock.Lock()

	this.server.OnlineMap[this.Name] = this

	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已上线")

}

func (this *Client) Offline() {
	this.server.mapLock.Lock()

	delete(this.server.OnlineMap, this.Name) //将当前Client的Name进行删除

	this.server.BroadCast(this, "已下线")

	this.server.mapLock.Unlock()
}

//哪一个用户发的什么消息
func (this *Client) DoMessage(client *Client, message string) {
	this.server.BroadCast(this, message)
}
