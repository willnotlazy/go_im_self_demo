package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	MapLock   sync.Mutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) ListenBroadCast() {
	for {
		msg := <-server.Message
		server.MapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.MapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	//
	u := NewUser(conn, this)

	u.Online()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				u.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn read err:", err)
			}
			msg := string(buf[:n-1])
			u.DoMessage(msg)
		}
	}()

	select {}
}

func (this *Server) Start() {
	// listen
	listener, err := net.Listen("tcp", this.Address())
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close
	defer listener.Close()

	go this.ListenBroadCast()
	// accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		// do handler
		go this.Handler(conn)
	}
}

func (this *Server) Address() string {
	return fmt.Sprintf("%s:%d", this.Ip, this.Port)
}
