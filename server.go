package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      string
	OnlineMap map[string]*Client
	mapLock   sync.RWMutex

	Message chan string //用于向Client的channel进行转发数据的channel
}

func newServer(ip string, port string) *Server {

	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*Client),
		Message:   make(chan string),
	}

	return server
}

//哪一个客户， 需要广播的消息是啥
func (this *Server) BroadCast(client *Client, message string) {
	sendMessage := "[" + client.Addr + "]" + client.Name + ":" + message

	this.Message <- sendMessage //将用户对应的消息传送到Server's channel中

}

func (this *Server) ListenChannel() {
	for {
		//遍历map,将所得到的string写到对应的client的channel中

		msg := <-this.Message

		this.mapLock.Lock()

		for _, clientC := range this.OnlineMap {
			clientC.C <- msg
		}

		this.mapLock.Unlock()

	}
}

//Server 的Handler成员函数
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Succefully calling handler function")
	//用户上线成功，需要将用户放到Online Map中去
	user := newClient(conn, this) //将当前的Server进行同client的关联

	user.Online()

	//进行广播的function

	go func() {
		//goroutine 去无限的接收用户传递来的消息，并且进行广播

		buf := make([]byte, 4096)

		for {
			cnt, err := conn.Read(buf) //read from connection

			if cnt == 0 { //如果接受的是空的， 那么就代表是下线了
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:cnt-1]) //从流中提取出数据

			//用户针对与msg进行消息的处理
			user.DoMessage(msg)
		}

	}()

	select {}

}

func (this *Server) Start() {

	fmt.Println("Calling start function")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", this.IP, this.Port))

	if err != nil {
		fmt.Println("Start error: ", err)
		return
	}

	defer listener.Close() //最后退出的时候进行listener 的退出

	go this.ListenChannel() //监听Server的channel是否是有数据的

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Accept error :", err)
			return
		}
		go this.Handler(conn)
	}

}
