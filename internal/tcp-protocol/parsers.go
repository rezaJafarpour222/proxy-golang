package tcpprotocol

import (
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
}

func HttpParser(tcpData []byte) Request {
	request := string(tcpData)
	lines := strings.Split(request, "\r\n")
	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	headers := make(map[string]string)
	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		headers[key] = strings.TrimSpace(value)
	}
	return Request{
		Method:  parts[0],
		Path:    parts[1],
		Version: parts[2],
		Headers: headers,
	}

}
