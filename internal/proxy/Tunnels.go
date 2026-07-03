package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func HTTPTunneling(clientConn net.Conn) {
	defer clientConn.Close()

	_ = clientConn.SetDeadline(time.Now().Add(30 * time.Second))

	for {
		req, err := http.ReadRequest(bufio.NewReader(clientConn))
		if err != nil {
			return
		}

		req.Close = true

		host := req.Host
		if host == "" {
			writeBadRequest(clientConn)
			return
		}

		if !strings.Contains(host, ":") {
			host += ":80"
		}

		fmt.Println("Proxying:", req.Method, "->", host, req.URL.Path)

		remoteConn, err := net.DialTimeout("tcp", host, 10*time.Second)
		if err != nil {
			writeBadGateway(clientConn)
			return
		}

		_ = remoteConn.SetDeadline(time.Now().Add(30 * time.Second))

		err = req.Write(remoteConn)
		if err != nil {
			remoteConn.Close()
			writeBadGateway(clientConn)
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(remoteConn), req)
		if err != nil {
			remoteConn.Close()
			writeBadGateway(clientConn)
			return
		}

		func() {
			defer resp.Body.Close()
			defer remoteConn.Close()

			_ = resp.Write(clientConn)
		}()

		if req.Close || resp.Close {
			return
		}
	}
}

func writeBadRequest(conn net.Conn) {
	_, _ = conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"))
}

func writeBadGateway(conn net.Conn) {
	_, _ = conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\nConnection: close\r\n\r\n"))
}
