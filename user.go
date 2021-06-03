package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Ch     chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) (user *User) {
	userAddr := conn.RemoteAddr().String()
	user = &User{
		Name:   "u" + userAddr,
		Addr:   userAddr,
		Ch:     make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return
}
func (u *User) Online() {
	u.server.MapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.MapLock.Unlock()
	u.server.Broadcast(u, "已上线")
}
func (u *User) Offline() {
	u.server.MapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.MapLock.Unlock()
	u.server.Broadcast(u, "已下线")
}
func (u *User) SendMsgToUser(user *User, msg string) {
	u.server.Broadcast(user, msg)
}
func (u *User) EchoMessage(msg string) {
	u.conn.Write([]byte(msg))

}
func (u *User) SendMessage(msg string) {

	if msg == "who" {
		u.server.MapLock.Lock()
		for _, user := range u.server.OnlineMap {
			msg := "[" + user.Addr + "]" + user.Name + ":" + "在线"
			u.EchoMessage(msg)
		}
		u.server.MapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.EchoMessage("当前用户名已被占用")
		} else {
			u.server.MapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.MapLock.Unlock()
			u.Name = newName
			u.EchoMessage("修改成功")

		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式 to|张三｜内容
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.EchoMessage("消息格式不正确，请使用 \"to|张三｜内容 \" 格式\n")
			return
		}
		remeteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.EchoMessage("用户不存在！！")
			return

		}
		content := strings.Split(msg, "|")[2]
		remeteUser.EchoMessage(u.Name + "对你说::" + content)
	} else {
		u.server.Broadcast(u, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		for msg := range u.Ch {
			n, err := u.conn.Write([]byte(msg + "\n"))
			fmt.Println(n, err)
		}

	}
}
