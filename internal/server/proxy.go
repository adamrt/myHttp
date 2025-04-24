package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Proxy struct {
	ServerClients map[string]int
	mu            sync.Mutex
}

func ProxyServer() {
	p := NewProxy()
	listener, err := net.Listen("tcp", ":"+proxy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server running on Port: %s\n", proxy)

	defer listener.Close()
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			defer c.Close()
			p.ProxyHandleConn(c)
			if err != nil {
				log.Fatal(err)
			}

		}(conn)
	}
}
func NewProxy() *Proxy {
	return &Proxy{
		ServerClients: make(map[string]int),
	}
}

func (p *Proxy) ProxyInitMap() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ServerClients["8000"] = 0
	p.ServerClients["8001"] = 0
	p.ServerClients["8002"] = 0
}

func (p *Proxy) ProxyHandleConn(conn net.Conn) {
	lowest := p.LowestClients()

	serverAddr := "localhost:" + lowest // Adjust this to use your actual backend server addresses
	backendConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("Failed to connect to backend server %s: %v", lowest, err)
		conn.Close()
		return
	}

	// Step 3: Forward data between client and backend server
	go p.forwardData(conn, backendConn) // Forward from client to server
	go p.forwardData(backendConn, conn) // Forward from server to client
}

func (p *Proxy) forwardData(src, dest net.Conn) {
	// Step 4: Copy data from src to dest
	_, err := io.Copy(dest, src)
	if err != nil {
		log.Printf("Error forwarding data: %v", err)
	}

	// Step 5: Close connections when done
	src.Close()
	dest.Close()
}

func (p *Proxy) LowestClients() string {
	lowest := ""
	lowestCount := int(^uint(0) >> 1)
	for port, clientCount := range p.ServerClients {
		if clientCount < lowestCount {
			lowestCount = clientCount
			lowest = port
		}
	}
	return lowest
}
