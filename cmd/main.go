package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"proxy/internal/proxy"
	"syscall"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}
			go proxy.HTTPTunneling(conn)
		}
	}()
	<-ctx.Done()
	fmt.Printf("\nShutdown signal received.\n")
	time.Sleep(1 * time.Second)
	listener.Close()
	fmt.Printf("Proxy shutdowned gracefully.\n")
}
