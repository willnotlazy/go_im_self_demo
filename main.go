package main

import "github.com/willnotlazy/go_im_self_demo/server"

func main() {
	s := server.NewServer("192.168.4.22", 8888)
	s.Start()
}
