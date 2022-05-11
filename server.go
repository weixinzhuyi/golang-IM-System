package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//UserMap
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

//创建一个server接口
func NewServer(ip string, port int) *Server {
	var server = &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
func (s *Server) Handler(conn net.Conn) {
	//处理当前链接的业务
	fmt.Println("正在处理业务")
	user := NewUser(conn, s)

	user.Online()

	//判断活跃
	isLive := make(chan bool)
	//接受用户发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 { //conn 关闭
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Err :", err)
				return
			}

			//提取用户发送的信息
			msg := string(buf[:n-1])

			//广播消息
			user.DoMessage(msg)

			isLive <- true
		}

	}()

	for {
		select {
		case <-isLive: //活跃不做任何处理
		case <-time.After(time.Minute * 5): //超时踢出
			user.SendMsg("You are out!")
			close(user.C)
			conn.Close()
			//delete(s.OnlineMap, user.Name)
			//	不用手动删除在线名单中的数据，因为conn关闭后，监听器会识别长度为0，执行下线行为，删除用户
		}
	}
}

//监听消息的方法
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, client := range s.OnlineMap {
			client.C <- msg
		}

		s.mapLock.Unlock()
	}
}

//广播消息的方法
func (s *Server) BoardCast(user *User, msg string) {
	sendMessage := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMessage
}

//启动服务器接口
func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err :", err)
		return
	}
	defer listener.Close()

	go s.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listen accept err:", err)
			continue
		}

		go s.Handler(conn)
	}
}
