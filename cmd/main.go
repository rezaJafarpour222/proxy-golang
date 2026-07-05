package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"proxy/internal/proxy"
	"sync"
	"syscall"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Proxy listening on :", port)

	serverCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-serverCtx.Done():
					return
				default:
					continue
				}
			}
			wg.Go(func() {
				proxy.HTTPProxy(conn)
			})
		}
	})

	<-serverCtx.Done()
	fmt.Println("Shutting down...")
	listener.Close()
	wg.Wait()
}
