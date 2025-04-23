package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {

	t := time.Now()
	var wg sync.WaitGroup
	concurrentClients := 10000
	for i := 0; i < concurrentClients; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", "localhost:8000")
			if err != nil {
				fmt.Printf("Client %d: Error connecting: %v\n", id, err)
				return
			}
			defer conn.Close()

			req := "VXP/1.0 PING localhost\r\nbody: 1\r\n\r\n1"
			_, err = conn.Write([]byte(req))
			if err != nil {
				fmt.Printf("Client %d: Write error: %v\n", id, err)
				return
			}

		}(i)
	}

	wg.Wait()
	e := time.Since(t)
	fmt.Printf("%v clients finished. Time it took %s\n", concurrentClients, e)
}
