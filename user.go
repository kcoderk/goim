package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	Ch   chan string
	conn net.Conn
}

func NewUser(conn net.Conn) (user *User) {
	userAddr := conn.RemoteAddr().String()
	user = &User{
		Name: "u" + userAddr,
		Addr: userAddr,
		Ch:   make(chan string),
		conn: conn,
	}
	go user.ListenMessage()
	return
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.Ch
		n, err := u.conn.Write([]byte(msg + "\n"))
		fmt.Println(n, err)

	}
}
