package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, ServerPort int) *Client {
	//创建client对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: ServerPort,
		flag:       999,
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

func (c *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入正确的选项")
		return false
	}
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
	go client.DealResponse()
	client.run()
}

func (c *Client) run() {
	for c.flag != 0 {
		for c.menu() != true {
		}
		switch c.flag {
		case 1:
			c.PublicChat()
			break
		case 2:
			fmt.Println("选择私聊模式")
		case 3:
			c.ChangeName()

		}
	}

}

func (c *Client) DealResponse() {

	//监听服务器返回的消息并执行标准输出,永久阻塞监听
	io.Copy(os.Stdout, c.conn)
	//等价于
	//for {
	//	buf := make([]byte,4096)
	//	c.conn.Read(buf)
	//	fmt.Println(buf)
	//}
}

func (c *Client) ChangeName() bool {
	fmt.Println("请输入新用户名：")
	newName := ""
	fmt.Scanln(&newName)
	sendMsg := "rename|" + newName + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err :", err)
		return false
	}

	c.Name = newName
	return true
}

func (c *Client) PublicChat() {
	//提示用户输入信息
	fmt.Println("请输入聊天内容，想要退出请输入exit")
	var msg string
	fmt.Scanln(&msg)

	for msg != "exit" {
		//消息不为空则发送给服务器
		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}

		}
		msg = ""
		fmt.Println("请输入聊天内容，退出请输入exit")
		fmt.Scanln(&msg)
	}
}
