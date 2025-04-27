package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type Proxy struct {
	ServerClients map[string]int
	mu            sync.RWMutex
}

var proxy = "9000"

func ProxyServer() {
	p := NewProxy()
	p.ProxyInitMap()
	listener, err := net.Listen("tcp", ":"+proxy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Proxy Server running on Port: %s\n", proxy)

	defer listener.Close()
	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			defer c.Close()
			p.handleConn(c)
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

type CopyJob struct {
	dst net.Conn
	src net.Conn
}

func (p *Proxy) handleConn(clientConn net.Conn) {
	defer clientConn.Close()

	targetAddr := p.LowestClients()

	p.mu.Lock()
	clientCount := p.ServerClients[targetAddr]
	p.mu.Unlock()

	backendConn, err := net.Dial("tcp", "localhost:"+targetAddr)
	if err != nil {
		log.Println("backend connection failed:", err)
		return
	}
	defer backendConn.Close()

	fmt.Printf("Sending to Server at port: %s with %v Clients connected\n", targetAddr, clientCount)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := io.Copy(backendConn, clientConn); err != nil {
			log.Println("Error copying client to backend:", err)
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(clientConn, backendConn); err != nil {
			if !errors.Is(err, net.ErrClosed) && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("Error copying backend to client:", err)
			}
		}
	}()

	wg.Wait()
}

func copyBuffer(dest, src net.Conn, data []byte) error {
	if _, err := io.CopyBuffer(dest, src, data); err != nil {
		return errors.New("Error copying buffer")
	}
	return nil

}
func (p *Proxy) LowestClients() string {
	lowest := ""
	lowestCount := int(^uint(0) >> 1)
	p.mu.RLock()
	for port, clientCount := range p.ServerClients {
		if clientCount < lowestCount {
			lowestCount = clientCount
			lowest = port

		}
	}
	p.mu.RUnlock()

	p.mu.Lock()
	p.ServerClients[lowest]++
	p.mu.Unlock()

	return lowest
}
