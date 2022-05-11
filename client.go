package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, ServerPort int) *Client {
	//创建client对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: ServerPort,
	}

	//建立与服务器的链接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, ServerPort))
	if err != nil {
		fmt.Println("net.Dial err :", err)
		return nil
	}
	client.conn = conn

	//返回对象
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set Server IP (Default: 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "Set Server Port (Default: 8888")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>Link Fail...")
		return
	}
	fmt.Println(">>>>>Link Success!")

}
