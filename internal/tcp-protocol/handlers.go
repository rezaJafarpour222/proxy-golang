package tcpprotocol

import (
	"fmt"
	"io"
	"net"
)

func Tunneling(clientConn net.Conn) {
	defer clientConn.Close()

	buffer := make([]byte, 4096)
	n, err := clientConn.Read(buffer)
	if err != nil {
		return
	}

	rawRequest := buffer[:n]
	fmt.Println("--------------------------------------")
	response := HttpParser(rawRequest)

	fmt.Println("Responded back. From ", response.Headers["Host"])
	remoteConn, err := net.Dial("tcp", response.Headers["Host"])
	if err != nil {
		return
	}
	defer remoteConn.Close()
	go io.Copy(remoteConn, clientConn)

	io.Copy(clientConn, remoteConn)

}
