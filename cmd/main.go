package main

import (
	"fmt"
	"net"
	tcpprotocol "proxy/internal/tcp-protocol"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Listen on port 8080")
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		tcpprotocol.Tunneling(connection)
	}

}
