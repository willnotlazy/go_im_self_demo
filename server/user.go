package server

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn

	Server *Server
}

func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		Server: s,
	}

	go user.ListenMessage()

	return user
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.Conn.Write([]byte(msg + "\n"))
	}
}

func (user *User) Online() {
	user.Server.MapLock.Lock()
	user.Server.OnlineMap[user.Name] = user
	user.Server.MapLock.Unlock()

	user.BroadCast("已上线")
}

func (user *User) Offline() {
	user.Server.MapLock.Lock()
	delete(user.Server.OnlineMap, user.Name)
	user.Server.MapLock.Unlock()

	user.BroadCast("下线")
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		user.Server.MapLock.Lock()
		for _, onlineUser := range user.Server.OnlineMap {
			onlineMsg := "[" + onlineUser.Addr + "]" + onlineUser.Name + ":" + "在线..."
			user.SendMsg(onlineMsg)
		}
		user.Server.MapLock.Unlock()
	} else {
		user.BroadCast(msg)
	}
}

func (user *User) BroadCast(msg string) {
	user.Server.BroadCast(user, msg)
}

func (this *User) SendMsg(msg string) {
	this.C <- msg
}
