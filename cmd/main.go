package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"proxy/internal/proxy"
	"runtime"
	"sync"
)

func main() {
	cert, err := tls.LoadX509KeyPair(
		"cert.pem",
		"key.pem",
	)
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listener, err := tls.Listen(
		"tcp",
		":"+port,
		&tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Listen on :", port)
	fmt.Println("GOMAXPROCS: ", runtime.GOMAXPROCS(0))
	var wg sync.WaitGroup

	wg.Go(func() {

		for {
			fmt.Println("goroutines:", runtime.NumGoroutine())

			conn, err := listener.Accept()
			if err != nil {
				return
			}

			wg.Go(func() {
				proxy.HTTPProxy(conn)
			})
		}
	})

	wg.Wait()
}
