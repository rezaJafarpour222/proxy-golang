package proxy

import "strings"

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
}

func HttpParser(raw []byte) (Request, error) {
	request := string(raw)

	lines := strings.Split(request, "\r\n")
	requestLine := strings.Split(lines[0], " ")

	headers := make(map[string]string)

	for _, line := range lines[1:] {
		if line == "" {
			break
		}

		key, value, found := strings.Cut(line, ":")
		if found {
			headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}

	return Request{
		Method:  requestLine[0],
		Path:    requestLine[1],
		Version: requestLine[2],
		Headers: headers,
	}, nil
}
