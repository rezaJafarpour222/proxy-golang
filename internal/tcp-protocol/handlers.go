package tcpprotocol

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) error {
	defer conn.Close()
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}
	fmt.Println("Read %d bytes", n)
	fmt.Println("-------------------------")
	fmt.Println(string(buffer[:n]))
	fmt.Println("-------------------------")
	return nil
}
