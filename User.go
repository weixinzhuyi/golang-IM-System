package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//创建用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听服务
	go user.ListenMessage()

	return user
}

//用户上线的方法
func (u *User) Online() {
	//将用户加入在线用户表中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播上线消息
	u.server.BoardCast(u, "is online")

}

//用户下线的方法
func (u *User) Offline() {
	//将用户加入在线用户表中
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播上线消息
	u.server.BoardCast(u, "is offline")

}

//给当前用户发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//用户发消息的方法
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": is online\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//修改用户名的格式 ：rename|XX
		newName := strings.Split(msg, "|")[1]

		//判断用户名是否已经存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("This name is used!")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("OK!Your new name is :" + newName + "\n")
		}

	} else {
		u.server.BoardCast(u, msg)

	}
}

//监听User channel,一旦收到消息，发送给目标客户端
func (u *User) ListenMessage() {
	for {

		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))

	}
}
