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
	fmt.Printf("Read %d bytes\n", n)
	fmt.Println("-------------------------")
	fmt.Println(string(buffer[:n]))
	fmt.Println("------------------------")
	request := HttpParser(buffer[:n])
	fmt.Print(request)
	return nil
}
