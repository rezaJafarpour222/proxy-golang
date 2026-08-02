package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"proxy/internal/proxy"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
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

	fmt.Println("Listening on:", port)
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	// Listen for Ctrl+C or SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nShutdown signal received...")
		_ = listener.Close()
	}()
	go func() {
		fmt.Println()
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			time.Sleep(10 * time.Second)
			fmt.Printf("\rGoroutines: %05d | Go Memory: %.2f MB", runtime.NumGoroutine(), float64(m.Sys)/1024/1024)
		}
	}()
	var wg sync.WaitGroup

	wg.Go(func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Accept error:", err)
				return
			}

			wg.Go(func() {
				proxy.HTTPProxy(conn)
			})
		}
	})

	wg.Wait()

	fmt.Println("Server stopped.")
}
