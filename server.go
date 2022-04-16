package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port string

	Onlinemap map[string]*User //每一个string 对应与一个User的对象

	mapLock sync.RWMutex

	Message chan string //用于消息广播的channel
}

//function : make a object of server
func newServer(ip string, port string) *Server {
	//能够创建一个Server的对象
	server := &Server{
		Ip:        ip,
		Port:      port,
		Onlinemap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (this *Server) BroadCast(user *User, message string) {
	//进行广播的function
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + message

	this.Message <- sendMsg //发送到Server的channel中去

}

//无限监听channel的function， 一旦有了消息就需要传递给用户的channel
func (this *Server) ListenChannel() {
	for {
		//遍历map， 将所得到的string写到用户对应的channel中去

		msg := <-this.Message
		this.mapLock.Lock()
		for _, clientC := range this.Onlinemap {
			clientC.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	//处理接收连接的函数，用户上线成功

	user := newUser(conn)

	//将用户加入到map表中
	this.mapLock.Lock()

	this.Onlinemap[user.Name] = user

	this.mapLock.Unlock()

	//进行广播
	this.BroadCast(user, "已上线")

	select {} //进行阻塞？？？？

	//fmt.Println("连接成功！")
}

//成员函数， Start Server
func (this *Server) Start() {
	fmt.Println("calling start function")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", this.Ip, this.Port))

	if err != nil {
		fmt.Println("ner error: ", err)
		return
	}
	defer listener.Close() //结束的时候自动进行Close

	go this.ListenChannel() //监听Server的channel是否是有数据的

	for {
		//无限循环去accept 新的连接

		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("accept err: ", err)
			return
		}

		go this.Handler(conn)

	}

}
