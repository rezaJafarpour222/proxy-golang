package tcpprotocol

import (
	"fmt"
	"io"
	"net"
)

func Tunneling(clientConn net.Conn) {
	defer clientConn.Close()

	remoteConn, err := net.Dial("tcp", "example.com:80")
	if err != nil {
		return
	}
	fmt.Println("Begin Reading The Stream From Client")

	go func() {
		io.Copy(remoteConn, clientConn)
		defer remoteConn.Close()
	}()

	io.Copy(clientConn, remoteConn)

	fmt.Println("Responded back")
	fmt.Println("--------------------------------------")
}
