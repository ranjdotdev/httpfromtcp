package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error){
	data, err := io.ReadAll(reader)
	if err != nil {return nil, err}
	lines := strings.Split(string(data), "\r\n")
	if len(lines)==0 || lines[0] == "" {
		return nil, errors.New("no request line found")
	}

	parts := strings.Split(lines[0], " ")
	// should include method, target, version
	if len(parts) != 3 {
        return nil, errors.New("request line must have 3 parts")
    }

	method := parts[0]
    target := parts[1]
    version := parts[2]

    for _, char := range method {
        if !unicode.IsUpper(char) || !unicode.IsLetter(char){
            return nil, errors.New("method must be uppercase letters")
        }
    }
	if version != "HTTP/1.1" {
        return nil, errors.New("only HTTP/1.1 is supported")
    }
	return &Request{
		RequestLine: RequestLine{
			Method: method,
			RequestTarget: target,
			HttpVersion: "1.1",
		},
	}, nil
}