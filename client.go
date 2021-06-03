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
	option     int
}

func NewClient(serverIp string, serverPort int) (client *Client) {
	client = &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		option:     99,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error")
		return nil
	}
	client.conn = conn
	return
}
func (client *Client) menu() bool {
	var option int
	fmt.Println("1.聊天室模式")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&option)
	if option <= 3 && option >= 0 {
		client.option = option
		return true
	} else {
		fmt.Println(">>>请输入合法范围的数字")
		return false
	}
}
func (client *Client) selectUsers() {
	msg := "who\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write err", err)

	}
	return
}
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}
func (client *Client) publicChat() {
	var chatMsg string
	fmt.Println(">>>请输入聊天内容,exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {

		if len(chatMsg) > 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write error", err)
				break
			}
		}
		fmt.Println(">>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)

	}
}
func (client *Client) privateChat() {
	var remoteName string
	var chatMsg string
	client.selectUsers()
	fmt.Println(">>>请输入聊天对象「用户名」,exit退出:")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println("请输入消息内容，exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("onn Write err:", err)
				}
			}
			chatMsg = ""
			fmt.Println("请输入消息内容，exit退出")
			fmt.Scanln(&chatMsg)
		}
		client.selectUsers()
		fmt.Println(">>>请输入聊天对象「用户名」,exit退出:")
		fmt.Scanln(&remoteName)

	}

}
func (client *Client) changeName() bool {
	fmt.Println(">>>请输入需要修改的用户名:")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.write error:", err)
		return false
	}
	return true
}
func (client *Client) Run() {
	for client.option != 0 {
		for client.menu() != true {
		}
		switch client.option {
		case 1:
			client.publicChat()
		case 2:
			client.privateChat()
		case 3:
			client.changeName()

		}

	}
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器ip，默认为 localhost即127.0.0.1")
	flag.IntVar(&serverPort, "port", 8080, "设置默认端口为8080")
}
func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println(">>>连接服务器失败")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>连接服务器成功")
	//客户端业务逻辑
	client.Run()
}
