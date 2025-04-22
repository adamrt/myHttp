package main

import (
	"fmt"
	"io"
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

		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

	fmt.Printf("Connection accepted at: %s\n", conn.RemoteAddr().String())
	rq := request.RequestLine{}
	buf := make([]byte, 4096)
	readFromIdx := 0
	for {

		n, err := conn.Read(buf[readFromIdx:])
		if err != nil {
			if err == io.EOF {
				conn.Close()
				return
			}
			log.Fatal(err)
		}

		readFromIdx += n
		_, req, err := rq.ParseRequestLine(buf[:readFromIdx])
		if err != nil {
			log.Println(err)
		}
		fmt.Println(req)

	}

}
