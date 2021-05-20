package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//online user list
	OnlineMap map[string]*User
	MapLock   sync.RWMutex

	//broadcast channel
	Message chan string
}

func NewServer(ip string, port int) (server *Server) {
	server = &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return
}
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.MapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.Ch <- msg
		}
		s.MapLock.Unlock()

	}

}
func (s *Server) Broadcast(user *User, msg string) {
	message := "[" + user.Addr +  "]:" +user.Name + msg
	s.Message <- message
}
func (s *Server) Handler(conn net.Conn) {
	fmt.Println(" create connect success")
	user := NewUser(conn)
	//用户上线
	s.MapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.MapLock.Unlock()
	//广播消息
	s.Broadcast(user, "已上线")

	go func() {
		buf:=make([]byte,4096)
		for{
			n,err:=conn.Read(buf)
			//约定读到0就下线
			if n==0{
				s.Broadcast(user,"下线")
				return
			}
			if err!=nil && err!=io.EOF{
				fmt.Println("conn read err:",err)
				return
			}
			msg:=string(buf[:n-1])
			s.Broadcast(user,msg)
		}
	}()
	select {
	}

}
func (s *Server) Start() {
	//listen
	ipPort := fmt.Sprintf("%s:%d", s.Ip, s.Port)
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		log.Fatal(fmt.Sprintf("lister address %s error", ipPort))
	}
	go s.ListenMessage()
	//close
	defer func() {
		fmt.Println("connection closed")
		listener.Close()
	}()

	for {
		//accept
		conn, err := listener.Accept()
		if err != err {
			fmt.Println("listener accept error:", err)
			continue
		}
		//do handle
		go s.Handler(conn)
	}

}
