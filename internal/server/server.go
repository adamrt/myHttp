package server

import (
	"fmt"
	"log"
	"net"

	"github.com/hconn7/vxp/internal/handler"
	"github.com/hconn7/vxp/internal/request"
)

type Server struct {
	Handler *handler.HandlerCfg
}

func NewServer() *Server {
	return &Server{
		Handler: handler.NewHandler()}
}

var port = "8000"

func Start() {
	s := NewServer()
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
			defer c.Close()
			req, err := request.HandleConn(c)
			if err != nil {
				log.Fatal(err)
			}

			handler, err := s.Handler.Router.Handle(req.RequestLine.Method, req.RequestLine.RequestTarget)
			if err != nil {
				fmt.Println(err)
			}
			handler(c, req)

		}(conn)
	}
}
