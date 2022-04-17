package main

import (
	"net"
)

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

//发送个特定的Client， 谁查的就传送给谁
func (this *Client) sendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

//哪一个用户发的什么消息
func (this *Client) DoMessage(message string) {

	if message == "who" { //如果发送的是"who", 就是查询所有在线的Client
		this.server.mapLock.Lock()

		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":在线。。。。。\n"
			this.sendMsg(onlineMsg)
		}

		this.server.mapLock.Unlock()
	} else if len(message) > 7 && message[:7] == "rename|" {
		//消息格式：rename|张三

		newName := message[7:] //获取到新起的名字

		//judge if the new_name is there
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.sendMsg("当前用户名已经被使用了！")
		} else {
			this.server.mapLock.Lock()
			//将映射当前对象的Name的字段首先删除，再添加
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.sendMsg("您已经成功更新用户名：" + newName + "\n")

		}

	} else {
		this.server.BroadCast(this, message)
	}

}
