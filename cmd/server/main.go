package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hconn7/vxp/internal/request"
)

const port = "8000"

func main() {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server running on Port: %s\n", port)

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			_, err = request.HandleConn(c)
			if err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
}
