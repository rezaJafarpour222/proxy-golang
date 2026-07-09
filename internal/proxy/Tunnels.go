package proxy

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var dialer = &net.Dialer{
	Timeout:   10 * time.Second,
	KeepAlive: 30 * time.Second,
}

func HTTPProxy(clientConn net.Conn) {
	defer clientConn.Close()
	httpContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	reader := bufio.NewReader(clientConn)
	for {
		clientConn.SetReadDeadline(time.Now().Add(30 * time.Second))
		req, err := http.ReadRequest(reader)
		if err != nil {
			return
		}
		if req.Method == http.MethodConnect {
			_ = clientConn.SetDeadline(time.Time{})
			handleHTTPS(httpContext, clientConn, req, cancel)
			return
		}
		handleHTTP(clientConn, req)
		if req.Close {
			return
		}
	}
}

func handleHTTP(clientConn net.Conn, req *http.Request) {
	host := req.Host
	if host == "" {
		writeBadRequest(clientConn)
		return
	}

	if !strings.Contains(host, ":") {
		host += ":80"
	}

	remoteConn, err := dialer.Dial("tcp", host)
	if err != nil {
		writeBadGateway(clientConn)
		return
	}
	defer remoteConn.Close()

	req.RequestURI = ""

	if err := req.Write(remoteConn); err != nil {
		writeBadGateway(clientConn)
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(remoteConn), req)
	if err != nil {
		writeBadGateway(clientConn)
		return
	}
	defer resp.Body.Close()

	_ = resp.Write(clientConn)
}

func handleHTTPS(ctx context.Context, clientConn net.Conn, req *http.Request, cancel context.CancelFunc) {

	remoteConn, err := dialer.DialContext(ctx, "tcp", req.Host)
	if err != nil {
		writeBadGateway(clientConn)
		return
	}
	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		remoteConn.Close()
		return
	}

	var once sync.Once

	closeAll := func() {
		once.Do(func() {
			_ = remoteConn.Close()
			_ = clientConn.Close()
			cancel()
		})
	}
	go func() {
		defer closeAll()
		io.Copy(remoteConn, clientConn)
	}()

	go func() {
		defer closeAll()
		io.Copy(clientConn, remoteConn)

	}()

	<-ctx.Done()
}

func writeBadRequest(conn net.Conn) {
	_, _ = conn.Write([]byte(
		"HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n",
	))
}

func writeBadGateway(conn net.Conn) {
	_, _ = conn.Write([]byte(
		"HTTP/1.1 502 Bad Gateway\r\nConnection: close\r\n\r\n",
	))
}
