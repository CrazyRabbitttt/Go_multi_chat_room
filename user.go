package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn //进行客户端通信的conn
}

//创建一个User的function
func newUser(conn net.Conn) *User {

	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string), //创建string类型的chann
		conn: conn,
	}

	//启动监听当前的user channel消息的 goroutine
	go user.ListenMessage()
	return user
}

//监听当前的chann，一旦有消息需要发送给客户端

func (this *User) ListenMessage() {
	for {
		message := <-this.C
		this.conn.Write([]byte(message + "\n")) //将message 发送回去给到客户端
	}
}
