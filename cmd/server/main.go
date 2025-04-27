package main

import "github.com/hconn7/vxp/internal/server"

func main() {
	go server.ProxyServer()
	go server.StartServer("8000")
	go server.StartServer("8002")
	go server.StartServer("8001")

	select {}
}
