package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) (server *Server) {
	server = &Server{
		Ip:   ip,
		Port: port,
	}
	return
}
func (s *Server) Handler(conn net.Conn) {
	fmt.Println(" create connect success")

}
func (s *Server) Start() {
	//listen
	ip_port := fmt.Sprintf("%s:%d", s.Ip, s.Port)
	listener, err := net.Listen("tcp", ip_port)
	if err != nil {
		log.Fatal(fmt.Sprintf("lister address %s error", ip_port))
	}

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
